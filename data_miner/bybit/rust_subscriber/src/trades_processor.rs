use crate::fs_cache;
use crate::helpers::{cur_ts_ms, cur_ts_sec};
use crate::market_data::Trades;
use crate::objects::SymbolPair;
use std::sync::{Arc, RwLock};
use std::thread;
use std::thread::{sleep, JoinHandle};
use std::time::Duration;

const LOOP_PATH: &str = "../symbol_loops.csv";

pub fn run(trades: &Arc<RwLock<Trades>>, recalc_interval: &Duration) -> JoinHandle<()> {
    let trades = trades.clone();
    let recalc_interval = recalc_interval.clone();
    return thread::spawn(move || process(trades, recalc_interval));
}

fn need_recalc_loops(last_loop_recalc_ts: u64, recalc_interval: &Duration) -> bool {
    cur_ts_sec() - last_loop_recalc_ts > recalc_interval.as_secs()
}

fn process(trades: Arc<RwLock<Trades>>, recalc_interval: Duration) {
    let mut last_recalc_ts = 0;

    let mut loops = fs_cache::read_loops(LOOP_PATH).expect("fail to read loop file");
    loop {
        let loop_start_ts = cur_ts_sec();
        {
            let mut locked = trades.write().unwrap();
            log::trace!(
                "cur_ts_ms: {}, {}",
                cur_ts_ms(),
                locked.recent_trades_stat()
            );
        }
        if need_recalc_loops(last_recalc_ts, &recalc_interval) {
            let trades_data = trades.read().unwrap().data.clone();
            let recent_trades = trades.read().unwrap().recent_trades.clone();

            let recent_pairs = recent_trades
                .iter()
                .map(|x| SymbolPair {
                    first: x.sell_symb.clone(),
                    second: x.buy_symb.clone(),
                })
                .collect::<Vec<SymbolPair>>();

            loops.recalc(&recent_pairs);
            log::debug!("loops recalculated");

            last_recalc_ts = loop_start_ts;
            log::info!("===profitable loops===");
            for l in loops.routes.values() {
                let mut report_line = String::default();
                let base_symbol = l[0].clone();
                let init_cnt = 100f64;
                let mut counts = vec![init_cnt];
                report_line.push_str(format!("{}(100)->", base_symbol).as_str());
                for i in 0..(l.len() - 1) {
                    let cur = &l[i];
                    let next = &l[i + 1];
                    let price = trades_data.get(cur).unwrap().get(next).unwrap().price;
                    counts.push(counts.last().unwrap() * price);
                    report_line.push_str(
                        format!("{}({:.4},p={:.6})->", next, counts.last().unwrap(), price)
                            .as_str(),
                    );
                }

                if *counts.last().unwrap() > init_cnt {
                    log::info!("{}", report_line);
                }
            }
        }

        sleep(Duration::from_secs(2));
    }
}
