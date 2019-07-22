#pragma once
#include <map>
#include <csv.h>
#include <optional>
#include <filesystem>

namespace NStockAnalyzer {

struct TTickerGeneralData {
    std::string Symbol;
    std::string Name;
    std::optional<size_t> IpoYear;
    std::string Sector;
    std::string Industry;
    std::string SummaryLink;

};

struct TTickerDayData {
    std::string Date;
    double Open;
    double High;
    double Low;
    double Close;
    size_t Volume;
};

auto ParseTickersGeneralData() {
    const auto& path = "/Users/d.korchagin/Projects/Personal/imt/data/companylist.csv";
    std::map<std::string, TTickerGeneralData> result;
    io::CSVReader<6, io::trim_chars<' ', '\t'>, io::double_quote_escape<',', '"'>> in(path);
    in.read_header(io::ignore_extra_column, "Symbol", "Name", "IPOyear", "Sector", "Industry", "Summary Quote");
    TTickerGeneralData entity;
    std::string ipoYear;
    try {
        while (in.read_row(entity.Symbol, entity.Name, ipoYear, entity.Sector, entity.Industry, entity.SummaryLink)) {
            if (ipoYear == "n/a") {
                entity.IpoYear.reset();
            } else {
                entity.IpoYear = static_cast<size_t>(std::atoi(ipoYear.c_str()));
            }
            result.insert(std::make_pair(entity.Symbol, entity));
        }
    } catch (const io::error::base& e) {
        std::cerr << e.what() << std::endl;
    }
    
    return result;
}

auto ParseTickerDaysData() {
    using TTickerDaysData = std::map<std::string, TTickerDayData>;
    std::map<std::string, TTickerDaysData> result;

    for (auto& p: std::filesystem::directory_iterator("/Users/d.korchagin/Projects/Personal/imt/data")) {
        if (!p.is_directory()) {
            continue;
        }
        auto historyDataPath = p.path().string() + "/history_data.csv"; 
        std::cout << historyDataPath << std::endl;
        TTickerDaysData daysData;
        io::CSVReader<6> in(historyDataPath);
        in.read_header(io::ignore_extra_column, "Date", "Open", "High", "Low", "Close", "Volume");
        TTickerDayData entity;
        try {
            while (in.read_row(entity.Date, entity.Open, entity.High, entity.Low, entity.Close, entity.Volume)) {
            daysData.insert(std::make_pair(entity.Date, entity));
        }
        } catch (const io::error::base& e) {
            std::cerr << e.what() << std::endl;
        }
        result.insert(std::make_pair(p.path(), daysData));
    }
    
    
    return result;
}
}