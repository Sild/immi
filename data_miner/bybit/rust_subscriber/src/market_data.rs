use std::{collections::HashMap, time::Duration};
use drain_while::*;
use bybit::ws::response::Trade;

use crate::{objects, helpers::{cur_ts_sec, cur_ts_ms}};

#[derive(Default, Clone)]
pub struct TradeDetail {
    pub sell_symb: String,
    pub buy_symb: String,
    pub ts_ms: u128,
    pub price: f64,
    pub quantity: f64,
}

impl TradeDetail {
    fn new(sell_symb: &str, buy_symb: &str, ts_ms: u128, price: f64, quantity: f64) -> Self {
        TradeDetail {
            sell_symb: sell_symb.to_string(),
            buy_symb: buy_symb.to_string(),
            ts_ms,
            price,
            quantity,
        }
    }
}

pub struct MarketData {
    pub trades: HashMap<String, HashMap<String, TradeDetail>>, // sell->buy->data
    pub recent_trades: Vec<TradeDetail>,
    pub trades_zero_cnt: usize,
    pub trades_total_cnt: usize,
    recent_trades_window: Duration,
}

impl MarketData {
    pub fn new(pairs: &[objects::SymbolPair], recent_trades_window: Duration) -> Self {
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
            recent_trades: Vec::default(),
            trades_zero_cnt: size,
            trades_total_cnt: size,
            recent_trades_window,
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

        let details = TradeDetail::new(&sell_symb, &buy_symb, trade.T as u128, price, quantity);
        let tmp = self.trades.entry(sell_symb.to_string()).or_default();
        let old = tmp.insert(buy_symb.to_string(), details.clone()).unwrap();
        if old.ts_ms == 0 {
            self.trades_zero_cnt -= 1;
        }
        self.recent_trades.push(details);
        self.remove_old_trades();
    }

    pub fn recent_trades_stat(&mut self) -> String {
        self.remove_old_trades();
        let default_trade = TradeDetail::default();
        format!(
            "recent_trades size: total={}, filled={}, empty={}, oldest_trade_ts={}, recent_trade_ts={},",
            self.trades_zero_cnt,
            self.trades_total_cnt - self.trades_zero_cnt,
            self.trades_zero_cnt,
            self.recent_trades.first().unwrap_or(&default_trade).ts_ms,
            self.recent_trades.last().unwrap_or(&default_trade).ts_ms,
        )
    }

    fn remove_old_trades(&mut self) {
        let oldest_ts = cur_ts_ms() - self.recent_trades_window.as_millis();
        self.recent_trades.drain_while(|x| x.ts_ms < oldest_ts);
    }
}
