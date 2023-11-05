use std::{
    cmp,
    collections::HashMap,
    sync::{Arc, Mutex},
    thread::{self, JoinHandle},
};

use bybit::{ws::response::SpotPublicResponse, WebSocketApiClient};

use crate::{
    market_data,
    objects::{self, SymbolPair},
};
use crypto_rest_client::BybitRestClient;

const SUBSCRIBE_BATCH_SIZE: usize = 5;

pub fn populate_trades(
    market_data: Arc<Mutex<market_data::MarketData>>,
    sym_to_pair: Arc<HashMap<String, objects::SymbolPair>>,
    symbols: Vec<objects::SymbolPair>,
) {
    let callback = |res: SpotPublicResponse| match res {
        // SpotPublicResponse::Orderbook(res) => println!("Orderbook: {:?}", res),
        SpotPublicResponse::Trade(res) => {
            for d in res.data.iter() {
                market_data.lock().unwrap().add(d, sym_to_pair.as_ref());
            }
            // println!("Trade: {:?}", res);
        }
        // SpotPublicResponse::Ticker(res) => println!("Ticker: {:?}", res),
        // SpotPublicResponse::Kline(res) => println!("Kline: {:?}", res),
        // SpotPublicResponse::LtTicker(res) => println!("LtTicker: {:?}", res),
        // SpotPublicResponse::LtNav(res) => println!("LtNav: {:?}", res),
        SpotPublicResponse::Op(_res) => {
            // println!("Op: {:?}", res)
        }
        _ => println!("Unexpected event"),
    };

    let mut client = WebSocketApiClient::spot().build();

    for p in symbols {
        println!("Subscribing to symbol {}", p.to_bybit_symbol());
        client.subscribe_trade(p.to_bybit_symbol());
    }

    match client.run(callback) {
        Ok(_) => {}
        Err(e) => println!("{}", e),
    }
}

pub fn run_trades_population(
    market_data: &Arc<Mutex<market_data::MarketData>>,
    sym_to_pair: &Arc<HashMap<String, objects::SymbolPair>>,
    pairs: &Vec<objects::SymbolPair>,
) -> Vec<JoinHandle<()>> {
    let mut threads = Vec::default();
    for i in (0..pairs.len()).step_by(SUBSCRIBE_BATCH_SIZE) {
        let after_last = cmp::min(i + SUBSCRIBE_BATCH_SIZE, pairs.len());
        let batch = pairs[i..after_last].to_vec();
        println!(
            "starting thread for symbols: {:?} {}/{})",
            batch,
            after_last,
            pairs.len()
        );
        let md_copy = market_data.clone();
        let s_t_p_copy = sym_to_pair.clone();
        let th = thread::spawn(move || {
            populate_trades(md_copy, s_t_p_copy, batch);
        });
        threads.push(th);
    }
    threads
}

pub fn get_recent_symbol_pairs() -> Vec<SymbolPair> {
    // let mut result = Vec::default();

    // mo meed credentials for public call
    let _client = BybitRestClient::new(None, None);
    todo!();
}
