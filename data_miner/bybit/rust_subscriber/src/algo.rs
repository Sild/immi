use std::{collections::HashMap, collections::HashSet, ops::Add};

use crate::objects::{SymbolLoops, SymbolPair};
use logging_timer::stime;

fn find_loop_dfs<'a>(
    graph: &HashMap<&'a String, Vec<&'a String>>,
    cur_path: &mut Vec<&'a String>,
    visited: &mut HashSet<&'a String>,
    result: &mut SymbolLoops,
) {
    for e in graph.get(*cur_path.last().unwrap()).unwrap() {
        log::trace!("cur path: {:?}", cur_path);
        if **e == **cur_path.first().unwrap() && cur_path.len() > 2 {
            cur_path.push(e);
            let r = cur_path
            .iter()
            .map(|x| String::from(*x))
            .collect::<Vec<String>>();
            result.add(&r);
            log::trace!("added route: {:?}", r);
            cur_path.pop();
            continue;
        }
        if visited.contains(e) {
            continue;
        }
        cur_path.push(e);
        visited.insert(e);
        find_loop_dfs(graph, cur_path, visited, result);
        cur_path.pop();
        visited.remove(e);
    }
}

#[stime("info")]
pub fn find_loops(pairs: &[SymbolPair]) -> SymbolLoops {
    let mut graph: HashMap<&String, Vec<&String>> = HashMap::default();
    for p in pairs.iter() {
        graph.entry(&p.first).or_default().push(&p.second);
        graph.entry(&p.second).or_default().push(&p.first);
    }

    let mut result = SymbolLoops::default();

    let mut counter = 1usize;
    for elem in graph.keys() {
        log::debug!("find_loops: start with elem: {} ({}/{})", elem, counter, graph.len());
        counter += 1;
        let mut cur_path = vec![*elem];
        let mut visited = HashSet::from([*elem]);
        find_loop_dfs(&graph, &mut cur_path, &mut visited, &mut result);
    }
    result
}
