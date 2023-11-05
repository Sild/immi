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

#[derive(Debug)]

pub struct LoopRoute {
    pub route: Vec<String>,
}

impl LoopRoute {
    pub fn from_str(loop_str: &str) -> Self {
        let route = loop_str
            .split(' ')
            .map(|f| f.to_string())
            .collect::<Vec<String>>();
        LoopRoute { route }
    }
}
