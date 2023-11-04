#include <iostream>
#include <string>
#include <unordered_set>
using namespace std;

vector<string> loops;

void dfs(unordered_map<string, vector<string>>& graph, vector<string> cur_path, unordered_set<string>& visited, vector<vector<string>>& loops) {
    for (auto& e: graph[cur_path.back()]) {
        if (e == cur_path[0] && cur_path.size() > 2) {
            loops.push_back(cur_path);
            loops.back().push_back(e);
            continue;
        }
        if (visited.contains(e)) {
            continue;
        }
        cur_path.push_back(e);
        visited.insert(e);
        dfs(graph, cur_path, visited, loops);
        cur_path.pop_back();
        visited.erase(e);
    }
}

vector<vector<string>> find_loops(vector<pair<string, string>> paths) {
    unordered_map<string, vector<string>> graph;

    // Build graph from paths
    for(auto& p: paths) {
        graph[p.first].push_back(p.second);
        graph[p.second].push_back(p.first);
    }

    vector<vector<string>> loops;

    for(const auto& [elem, routes]: graph) {
        // cout << elem << ": ";
        // for (auto& node: routes) {
        //     cout << node << ",";
        // }
        // cout << endl;
        vector<string> cur_path = {elem};
        unordered_set<string> visited = {elem};
        dfs(graph, cur_path, visited, loops);
    }

    return loops;
}

int main() {
    int pairs_cnt = 0;
    std::cin >> pairs_cnt;

    vector<pair<string, string>> paths;

    for (int i = 0; i < pairs_cnt; ++i) {
        string from, to;
        std::cin >> from >> to;
        paths.push_back({from, to});
    }

    auto loops = find_loops(paths);

    for(auto& loop: loops) {
        for (auto& node: loop) {
            cout << node << " ";
        }
        cout << endl;
    }

    return 0;
}
