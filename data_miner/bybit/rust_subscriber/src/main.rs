#![feature(test)]
extern crate test;

use std::collections::HashMap;
use std::io::Error;
use std::sync::Arc;
use std::time::Duration;

mod algo;
mod fs_cache;
mod helpers;
mod market_data;
mod objects;
mod trades_processor;
mod trades_updater;

use log::LevelFilter;

extern crate bybit;
extern crate log;

const PAIRS_PATH: &str = "../routes.txt";
const LOOP_RECALC_INTERVAL: Duration = Duration::from_secs(30);

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

    let market_data = market_data::MarketData::new(&possible_pairs, Duration::from_secs(30));

    let mut threads =
        trades_updater::run_trades_population(&market_data.trades, &sym_to_pair, &possible_pairs);
    threads.push(trades_processor::run(
        &market_data.trades,
        &LOOP_RECALC_INTERVAL,
    ));

    for t in threads {
        let _ = t.join();
    }

    Ok(())
}
