

use std::io::Error;
use std::fs::read_to_string;

use bybit::ws::response::SpotPublicResponse;


use bybit::WebSocketApiClient;

#[derive(Debug)]

struct SymbolPair {
    first: String,
    second: String,

}

impl SymbolPair {
    fn to_bybit_symbol(&self) -> String {
        return format!("{}{}", self.first, self.second)
    }
}

#[derive(Debug)]

struct LoopRoute {
    route: Vec<String>
}

impl LoopRoute {
    fn new(loop_str: &str) -> Self {
        let route = loop_str.split(' ').map(|f| f.to_string()).collect::<Vec<String>>();
        return LoopRoute {
            route
         }
    }
}

fn read_pairs(fname: &str) -> Result<Vec<SymbolPair>, Error> {
    let lines: Vec<String> = read_to_string(fname)?
    .lines()
    .map(String::from)
    .collect();

    let mut res = Vec::default();
    for l in lines {
        let split: Vec<&str> = l.split(" ").collect();
        res.push(SymbolPair{
            first: split[0].to_string(),
            second: split[1].to_string()
        })
    }
    Ok(res)
}

fn read_loops(fname: &str) -> Result<Vec<LoopRoute>, Error> {
    let lines: Vec<String> = read_to_string(fname)?
    .lines()
    .map(String::from)
    .collect();

    let mut res = Vec::default();
    for l in lines {
        res.push(LoopRoute::new(&l));
    }
    Ok(res)
}

fn main()-> Result<(), Error> {
    let pairs = read_pairs("../routes.txt")?;
    let _loops = read_loops("../symbol_loops.csv")?;


    let mut client = WebSocketApiClient::spot().build();

    for p in pairs {
        println!("Subscribing to symbol {}", p.to_bybit_symbol());
        client.subscribe_trade("ETHUSDT");
        break;
    }

    // for p in pairs {
    //     println!("{:?}", p);
    // }

    // for l in loops {
    //     println!("{:?}", l)
    // }

    // client.subscribe_trade(symbol);

    let callback = |res: SpotPublicResponse| match res {
        // SpotPublicResponse::Orderbook(res) => println!("Orderbook: {:?}", res),
        SpotPublicResponse::Trade(res) => println!("Trade: {:?}", res),
        // SpotPublicResponse::Ticker(res) => println!("Ticker: {:?}", res),
        // SpotPublicResponse::Kline(res) => println!("Kline: {:?}", res),
        // SpotPublicResponse::LtTicker(res) => println!("LtTicker: {:?}", res),
        // SpotPublicResponse::LtNav(res) => println!("LtNav: {:?}", res),
        SpotPublicResponse::Op(res) => println!("Op: {:?}", res),
        _ => println!("Unexpected event")
    };
    
    match client.run(callback) {
        Ok(_) => {}
        Err(e) => println!("{}", e),
    }
    Ok(())
}
