use std::{collections::HashMap, error::Error, fs, fmt::{Display, self}};

use crate::algo;

#[derive(Debug, Clone)]

pub struct SymbolPair {
    pub first: String,
    pub second: String,
}

impl SymbolPair {
    pub fn to_bybit_symbol(&self) -> String {
        format!("{}{}", self.first, self.second)
    }
}

#[derive(Default, Debug, Clone)]
pub struct SymbolLoops {
    pub routes: HashMap<String, Vec<String>>,
    pairs: Vec<SymbolPair>
}

impl SymbolLoops {
    pub fn recalc(&mut self, pairs: &[SymbolPair]) {
        self.pairs = pairs.to_vec();
        self.routes = algo::find_loops(pairs).routes;
    }

    pub fn add_from_str(&mut self, loop_str: &str) {
        let route = loop_str
            .split(' ')
            .map(|f| f.to_string())
            .collect::<Vec<String>>();
        self.add(&route);
        self.routes.insert(route.first().unwrap().clone(), route);
    }

    pub fn add(&mut self, route: &[String]) {
        self.routes.insert(route.first().unwrap().clone(), route.to_vec());
    }

    pub fn to_file(&self, file_path: &str) -> Result<(), std::io::Error> {
        let mut data = String::default();
        for r in self.routes.values() {
            data.push_str(format!("{}\n", r.join(",")).as_str());
        }
        fs::write(file_path, data)
    }
}

impl Display for SymbolLoops {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        let mut repr = String::from("===SymbolLoops State===\npairs:\n");

        for p in self.pairs.iter() {
            repr.push_str(format!("{}->{}\n", p.first, p.second).as_str());
        }
        repr.push_str("loops:\n");
        for r in self.routes.values() {
            repr.push_str(format!("{}\n", r.join("->")).as_str());
        }
        repr.push_str("===/SymbolLoops State===");
        write!(f, "{}", repr)
    }
}