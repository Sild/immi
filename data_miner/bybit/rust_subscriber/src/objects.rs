use std::{collections::HashMap, error::Error, fs, fmt::{Display, self}};

use crate::algo;

#[derive(Debug, Clone, Eq, Ord, PartialEq, PartialOrd)]

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
        self.pairs.sort();
        self.pairs.dedup();

        let mut pairs_str = String::default();
        for p in self.pairs.iter() {
            pairs_str.push_str(format!("{}->{}\n", p.first, p.second).as_str());
        }
        log::trace!("run recalc based on pairs:\n{}", pairs_str);

        self.routes = algo::find_loops(&self.pairs).routes;
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
        let mut repr = String::from(format!("===SymbolLoops State===\npairs ({}):\n", self.pairs.len()).as_str());

        for p in self.pairs.iter() {
            repr.push_str(format!("{}->{}\n", p.first, p.second).as_str());
        }
        repr.push_str(format!("loops ({}):\n", self.routes.len()).as_str());
        for r in self.routes.values() {
            repr.push_str(format!("{}\n", r.join("->")).as_str());
        }
        repr.push_str("===/SymbolLoops State===");
        write!(f, "{}", repr)
    }
}