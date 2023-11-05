use std::{collections::HashMap};

use bybit::ws::response::Trade;

use crate::objects;

#[derive(Default)]
pub struct TradeDetail {
    pub ts: u64,
    pub price: f64,
    pub quantity: f64,
}

impl TradeDetail {
    fn new(ts: u64, price: f64, quantity: f64) -> Self {
        TradeDetail {
            ts,
            price,
            quantity,
        }
    }
}

pub struct MarketData {
    pub trades: HashMap<String, HashMap<String, TradeDetail>>, // sell->buy->data
    pub trades_zero_cnt: usize,
    pub trades_total_cnt: usize,
}

impl MarketData {
    pub fn new(pairs: &[objects::SymbolPair]) -> Self {
        let mut trades: HashMap<String, HashMap<String, TradeDetail>> = HashMap::new();
        let mut size: usize = 0;
        for p in pairs.iter() {
            trades
                .entry(p.first.clone())
                .or_default()
                .insert(p.second.clone(), TradeDetail::default());
            trades
                .entry(p.second.clone())
                .or_default()
                .insert(p.first.clone(), TradeDetail::default());
            size += 2;
        }

        MarketData {
            trades,
            trades_zero_cnt: size,
            trades_total_cnt: size,
        }
    }
    pub fn add(&mut self, trade: &Trade, sym_to_pair: &HashMap<String, objects::SymbolPair>) {
        let pair = sym_to_pair.get(&trade.s.to_string()).unwrap();
        let sell_symb = match trade.S {
            "Buy" => pair.first.clone(),
            _ => pair.second.clone(),
        };
        let buy_symb = match trade.S {
            "Buy" => pair.second.clone(),
            _ => pair.first.clone(),
        };
        let price = trade.p.parse::<f64>().unwrap();
        let quantity = trade.v.parse::<f64>().unwrap();

        let details = TradeDetail::new(trade.T, price, quantity);
        let tmp = self.trades.entry(sell_symb.to_string()).or_default();
        let old = tmp.insert(buy_symb.to_string(), details).unwrap();
        if old.ts == 0 {
            self.trades_zero_cnt -= 1;
        }
    }
}
