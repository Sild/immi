#include <iostream>
#include "entities.h"
int main() {
    std::cout << "hello analyzer" << std::endl;
    // auto tickerGeneralInfo = NStockAnalyzer::ParseTickersGeneralData();
    // for (auto& [k, v]: tickerGeneralInfo) {
    //     if (v.IpoYear) {
    //         std::cout << "ticker=" << k << ", ipoYear=" << *(v.IpoYear) << ", sector=" << v.Sector << std::endl;
    //     } else {
    //         std::cout << "ticker=" << k << ", ipoYear=" << "n/a" << ", sector=" << v.Sector << std::endl;

    //     }
    // }

    auto tickerDaysData = NStockAnalyzer::ParseTickerDaysData();
    for (auto& [ticker, data]: tickerDaysData) {
        for (auto& [day, v]: data) {
        std::cout << "ticker=" << ticker << ", date=" << day << ", open=" << v.Open << ", close=" << v.Close << std::endl;

        }
    }
    
}