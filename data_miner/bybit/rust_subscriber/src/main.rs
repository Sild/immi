use std::collections::HashMap;
use std::io::Error;
use std::sync::{Arc, Mutex};
use std::thread::sleep;
use std::time::Duration;

use helpers::cur_ts_sec;
use objects::SymbolPair;

mod algo;
mod fs_cache;
mod helpers;
mod market_data;
mod objects;
mod updater;
use log::LevelFilter;

extern crate bybit;
extern crate log;

const PAIRS_PATH: &str = "../routes.txt";
const LOOP_PATH: &str = "../symbol_loops.csv";
const LOOP_RECALC_INTERVAL: Duration = Duration::from_secs(30);

fn need_recalc_loops(last_loop_recalc_ts: u64) -> bool {
    cur_ts_sec() - last_loop_recalc_ts > LOOP_RECALC_INTERVAL.as_secs()
}

fn init_logger() {
    let mut builder = env_logger::Builder::from_default_env();
    if std::env::var("RUST_LOG").is_err() {
        // override default 'error'
        builder.filter_level(LevelFilter::Debug);
    }
    builder.init();
}

fn main() -> Result<(), Error> {
    init_logger();

    let possible_pairs = fs_cache::read_pairs(PAIRS_PATH)?;
    let sym_to_pair = Arc::new(
        possible_pairs
            .iter()
            .map(|p| (p.to_bybit_symbol(), p.clone()))
            .collect::<HashMap<_, _>>(),
    );

    let mut loops = fs_cache::read_loops(LOOP_PATH)?;

    let market_data = Arc::new(Mutex::new(market_data::MarketData::new(
        &possible_pairs,
        Duration::from_secs(30),
    )));

    let _threads = updater::run_trades_population(&market_data, &sym_to_pair, &possible_pairs);

    let mut last_loop_recalc_ts = cur_ts_sec()-20;
    loop {
        let loop_start_ts = cur_ts_sec();
        {
            let mut locked = market_data.lock().unwrap();
            log::debug!(
                "cur_ts_ms: {}, {}",
                helpers::cur_ts_ms(),
                locked.recent_trades_stat()
            );
        }
        if need_recalc_loops(last_loop_recalc_ts) {
            let recent_pairs = market_data
                .lock()
                .unwrap()
                .recent_trades
                .clone()
                .into_iter()
                .map(|x| SymbolPair {
                    first: x.sell_symb,
                    second: x.buy_symb,
                })
                .collect::<Vec<SymbolPair>>();

            loops.recalc(&recent_pairs);
            log::info!("loops recalculated:\n{}", loops);
            last_loop_recalc_ts = cur_ts_sec(); //loop_start_ts;
        }

        sleep(Duration::from_secs(2));
        if false {
            break;
        }
    }

    // for t in threads {
    //     let _ = t.join();
    // }

    Ok(())
}
