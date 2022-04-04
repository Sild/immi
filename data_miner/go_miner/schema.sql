CREATE TABLE instruments (
    id INTEGER PRIMARY KEY,
    figi TEXT NOT NULL,
    ticker TEXT,
    isin TEXT NOT NULL,
    name TEXT NOT NULL,
    min_price_increment REAL NOT NULL,
    lot INTEGER NOT NULL,
    currency TEXT,
    type TEXT
);

CREATE TABLE candles (
    id INTEGER PRIMARY KEY ,
    instrument_id INTEGER NOT NULL,
    interval TEXT NOT NULL,
    open REAL NOT NULL,
    close REAL NOT NULL,
    hight REAL NOT NULL,
    low REAL NOT NULL,
    volume INTEGER,
    ts INTEGER NOT NULL
);