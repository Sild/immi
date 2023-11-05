use std::{
    fs::{self, read_to_string},
    io::Error,
};

use crate::objects::{self, SymbolLoops};

pub fn read_pairs(fname: &str) -> Result<Vec<objects::SymbolPair>, Error> {
    let lines: Vec<String> = read_to_string(fname)?.lines().map(String::from).collect();

    let mut res = Vec::default();
    for l in lines {
        let split: Vec<&str> = l.split(' ').collect();
        res.push(objects::SymbolPair {
            first: split[0].to_string(),
            second: split[1].to_string(),
        })
    }
    Ok(res)
}

pub fn write_pairs(pairs: &[objects::SymbolPair], fname: &str) -> Result<(), Error> {
    let mut data = String::default();
    for p in pairs.iter() {
        data.push_str(format!("{},{}\n", p.first, p.second).as_str());
    }
    fs::write(fname, data)
}

pub fn read_loops(fname: &str) -> Result<SymbolLoops, Error> {
    let lines: Vec<String> = read_to_string(fname)?.lines().map(String::from).collect();

    let mut res = SymbolLoops::default();
    for l in lines {
        res.add_from_str(l.as_str());
    }
    Ok(res)
}
