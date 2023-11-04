package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/hirokisan/bybit/v2"
)

// type Pair struct {
// 	first string
// 	second struct
// }

func sys_exec(cmd_str []string) (string, error) {
	// fmt.Printf(strings.Join(cmd_str, " "))
	cmd := exec.Command("/bin/bash", "-c", strings.Join(cmd_str, " "))

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func build_loops(pairs []string, f_path string, rebuild bool) error {
	if _, err := os.Stat(f_path); err == nil && !rebuild {
		return nil
	}

	loops := make([][]string, 0)

	stdin := fmt.Sprintf("%d\n", len(pairs))
	for _, d := range pairs {
		stdin += fmt.Sprintf("%s\n", d)
	}
	cmd := []string{"echo", fmt.Sprintf("'%s'", stdin), "|", "./a.out"}
	out, err := sys_exec(cmd)
	if err != nil {
		return err
	}

	out_split := strings.Split(out, "\n")
	for _, s := range out_split {
		s = strings.TrimSpace(s)
		if s != " " && s != "" {
			loops = append(loops, strings.Split(s, " "))
		}
	}

	os.Remove(f_path)
	file, err := os.OpenFile(f_path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range loops {
		if _, err := file.WriteString(fmt.Sprintf("%s\n", strings.Join(line, ","))); err != nil {
			return err
		}
	}
	return nil
}

func read_loops(fname string) ([][]string, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	result := [][]string{}
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		result = append(result, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func format_trade_response_v2(rsp *bybit.SpotWebsocketV1PublicV2TradeResponse) string {
	// data := []string{}
	action := "buy"
	r := rsp.Data
	if !r.IsBuySideTaker {
		action = "sell"
	}
	return fmt.Sprintf("ts: %d, %s, a=%s, p=%s, q=%s", r.Timestamp, rsp.Params.Symbol, action, r.Price, r.Quantity)
	// return strings.Join(data, "\n")
}

func write_routes(path string, routes *[]string) error {
	os.Remove(path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range *routes {
		if _, err := file.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
			return err
		}
	}
	return nil
}

func format_trade_response_v1(rsp *bybit.SpotWebsocketV1PublicV1TradeResponse) string {
	data := []string{}
	action := "buy"
	for _, r := range rsp.Data {
		if !r.IsBuySideTaker {
			action = "sell"
		}
		line := fmt.Sprintf("ts: %d, %s, a=%s, p=%s, q=%s", r.Timestamp, rsp.Symbol, action, r.Price, r.Quantity)
		data = append(data, line)
	}
	return strings.Join(data, "\n")
}

type Pair struct {
	first  string
	second string
}

func populate_rates(pairs *[]string, rates *map[string]map[string]float64, rates_mtx *sync.Mutex) error {
	mtx := sync.Mutex{}
	symbol_to_pair := map[string]Pair{}

	handler := func(response bybit.SpotWebsocketV1PublicV2TradeResponse) error {
		rsp_str := format_trade_response_v2(&response)
		mtx.Lock()
		fmt.Println(rsp_str)
		mtx.Unlock()
		p1, p2 := "", ""
		if response.Data.IsBuySideTaker {
			p1 = symbol_to_pair[response.Params.SymbolName].second
			p2 = symbol_to_pair[response.Params.SymbolName].first
		} else {
			p2 = symbol_to_pair[response.Params.SymbolName].second
			p1 = symbol_to_pair[response.Params.SymbolName].first
		}
		price, _ := strconv.ParseFloat(response.Data.Price, 32)
		rates_mtx.Lock()
		if _, found := (*rates)[p1]; !found {
			(*rates)[p1] = map[string]float64{}
		}
		(*rates)[p1][p2] = price
		rates_mtx.Unlock()
		return nil
	}

	executors := []bybit.WebsocketExecutor{}
	wsClient := bybit.NewWebsocketClient()
	svcRoot := wsClient.Spot().V1()

	for pos, p := range *pairs {
		symbol := strings.ReplaceAll(p, " ", "")
		arr := strings.Split(p, " ")
		p1, p2 := arr[0], arr[1]

		symbol_to_pair[symbol] = Pair{p1, p2}
		svc, err := svcRoot.PublicV2()
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("subscribing for symbol=%s (%d/%d)", symbol, pos+1, len(*pairs)))
		_, err = svc.SubscribeTrade(bybit.SymbolSpot(symbol), handler)
		if err != nil {
			return err
		}
		executors = append(executors, svc)
		if pos == 10 {
			break
		}
	}

	wsClient.Start(context.Background(), executors)
	return nil
}

func main() {
	client := bybit.NewClient().WithAuth("your api key", "your api secret")
	rsp, err := client.Spot().V1().SpotSymbols()
	if err != nil {
		fmt.Printf("Fail to get symbols: %v", err)
	}

	symbols_routes := []string{}
	for _, r := range rsp.Result {
		line := fmt.Sprintf("%s %s", r.BaseCurrency, r.QuoteCurrency)
		symbols_routes = append(symbols_routes, line)
	}

	// sort.SliceStable(symbols_routes, func(i, j int) bool {
	// 	return symbols_routes[i] < symbols_routes[j]
	// })

	if err := write_routes("routes.txt", &symbols_routes); err != nil {
		fmt.Printf("Fail to write routes file")
		return
	}

	loops_fname := "symbol_loops.csv"
	if err := build_loops(symbols_routes, loops_fname, false); err != nil {
		fmt.Printf("Fail to build loops in file %s: %w", loops_fname, err)
		return
	}

	// loops, err := read_loops(loops_fname)
	// if err != nil {
	// 	fmt.Printf("Fail to read loops from file %s : %w", loops_fname, err)
	// 	return
	// }
	rates := map[string]map[string]float64{}
	rates_mtx := sync.Mutex{}
	if err := populate_rates(&symbols_routes, &rates, &rates_mtx); err != nil {
		fmt.Println(err)
	}

	// for _, l := range loops {
	// 	fmt.Println(l)
	// }
}
