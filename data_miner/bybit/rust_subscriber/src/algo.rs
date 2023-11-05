use std::{collections::HashMap, collections::HashSet};

use crate::objects::{self, LoopRoute};

fn find_loop_dfs<'a>(
    graph: &HashMap<&'a String, Vec<&'a String>>,
    cur_path: &mut Vec<&'a String>,
    visited: &mut HashSet<&'a String>,
    result: &mut Vec<LoopRoute>,
) {
    for e in graph.get(*cur_path.last().unwrap()).unwrap() {
        if e == cur_path.first().unwrap() && cur_path.len() > 2 {
            cur_path.push(e);
            let r = LoopRoute {
                route: cur_path
                    .iter()
                    .map(|x| String::from(*x))
                    .collect::<Vec<String>>(),
            };
            result.push(r);
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

pub fn find_loops(pairs: &[objects::SymbolPair]) -> Vec<objects::LoopRoute> {
    let mut graph: HashMap<&String, Vec<&String>> = HashMap::default();
    for p in pairs.iter() {
        graph.entry(&p.first).or_default().push(&p.second);
        graph.entry(&p.second).or_default().push(&p.first);
    }

    let mut result = Vec::<objects::LoopRoute>::default();

    for elem in graph.keys() {
        let mut cur_path = vec![*elem];
        let mut visited = HashSet::from([*elem]);
        find_loop_dfs(&graph, &mut cur_path, &mut visited, &mut result);
    }
    result
}
