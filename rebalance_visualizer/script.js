
THROTTLING_DATA = 10

SOURCE_DATA = new Map() // ticker -> [(x1, y1), (x2, y2)]
TIME_POINTS = new Array() // [x1, x2, x3...] из значений мапки, описанной выше
CHARTS = new Map() // xpath_id -> Chart obj

function fillMock() {
    console.log("filling mock data...")
    for (i = 0; i < 2; ++i) {
        label = "mock" + i
        ticker_data = []
        value = 5 + i * 1000
        var now = new Date();
        for (var d = new Date(2020, 0, 1); d <= now; d.setDate(d.getDate() + 5)) {
            switch (d.getDate()) {
                case 1:
                    value += 10
                    break
                case 2:
                    value /= 1.5
                    break
                case 3:
                    value -= 30
                    break
                case 4:
                    value *=1.1
                    break
                default:
                    value += 20
                
            }
            date = new Date(d).getTime()
            ticker_data.push({x: date, y: value});
            TIME_POINTS.push(date)
        }
        SOURCE_DATA.set(label, ticker_data)
    }
}

// заполним пробелы в данных последними существующими данными
// это облегчит итерацию по данным при перебалансировке
function fillGaps() {
    console.log("filling gaps...")
    TIME_POINTS = [...new Set(TIME_POINTS)]
    TIME_POINTS.sort()
    console.log("totally", TIME_POINTS.length, "uniq points received")
    for (const [key, series] of SOURCE_DATA.entries()) {
        filled_series = []
        i = 0, j = 0
        for (;i < TIME_POINTS.length; ++i) {
            ts = TIME_POINTS[i]
            if (series[j].x != ts) {
                if (j != 0)
                    filled_series.push({x: ts, y: series[j - 1].y})
                else
                    filled_series.push({x: ts, y: 0})
            } else {
                filled_series.push(series[j])
                if (++j == series.length)
                    --j // оператор подергивания. Если данные кончились - всегда берем последние
            }
        }
        SOURCE_DATA.set(key, filled_series)
    }
}

function readFiles(e, mock=null)
{
    console.log("cleanup data ...")
    SOURCE_DATA = new Map()
    TIME_POINTS = new Array()
    if (mock) {
        fillMock()
        return
    }
    for (let f_idx = 0; f_idx < e.target.files.length; f_idx++) {
        var file = e.target.files[f_idx];
        var file_name = file.name
        if (!file) {
            console.error("fail to read file " + file_name)
            return;
        }
        // console.log("reeding file " + file_name)

        var reader = new FileReader();
        reader.onload = function(e) {
            var contents = e.target.result
            var ticker_data = []
            var splitted = contents.split("\n")
            for (let i = 1; i < splitted.length; i++) {
                // if (i % THROTTLING_DATA != 0)
                //     continue
                var line = splitted[i].split(",")
                var value = parseFloat(line[5])
                if (this.FileName.includes("dollar"))
                    value *= 80
                ticker_data.push({x: line[0], y: value})
                TIME_POINTS.push(line[0])
            }
            ticker_data.sort()
            SOURCE_DATA.set(this.FileName, ticker_data)
            console.log("reeding file " + this.FileName + " done:", ticker_data.length, "points received")
        };
        reader.FileName = file_name
        reader.readAsText(file);
    }
}

function draw(xpath_chart_id, data, title_text) {
    // data contains label -> [(x1, y1), (x2, y2)]
    console.log("drawing", xpath_chart_id)
    var series_data = []
    for (const [key, value] of data.entries()) {
        randomColor = "#" + Math.floor(Math.random()*16777215).toString(16);
        // console.log(randomColor)
        series_data.push({
            label: key,
            backgroundColor:randomColor,
            borderWidth: 0,
            hoverBackgroundColor: randomColor,
            hoverBorderColor:randomColor,
            data: value,
            pointRadius: 1,
            pointHoverRadius: 5,
        })
    }
    // console.log(series_data)
    var chart_data = {
        datasets: series_data
    };

    chart = CHARTS[xpath_chart_id]
    if (chart)
        chart.destroy()

    CHARTS[xpath_chart_id] = new Chart(xpath_chart_id, {
        type: "line",
        data: chart_data,
        options: {
            scales: {
                x: {
                    type: 'time',
                }
            },
            plugins: {
                title: {
                    display: true,
                    text: title_text,
                },
                legend: {
                    position: "right",
                    // align: "start",
                }
            }
        }
    });

}

// cчитает какое количество бумаг каждого типа нужно иметь, чтобы портфель был сбалансированным
// начиная с переданной даты и до конца доступных дней
function calcBalancedTickersCount(source_data, timestamp, previous_count) {
    // проверим что входные данные валидны для алгоритма
    for (const [_, data] of source_data) {
        keys = data.map(a => a.x);
        console.assert(keys.length == TIME_POINTS.length)
        for (i = 0; i < keys.length; ++i) {
            console.assert(keys[i] == TIME_POINTS[i])
        }
    }
    startAmount = document.getElementById('start-amount').value
    if (previous_count)
    {
        startAmount = 0
        // посчитаем сколько будет денег, когда продадим то что есть
        for (const [key, value] of previous_count) {
            for (i = 0; i < value.length; ++i) {
                serie = value[i]
                if (serie.x >= timestamp) {  // попали в запрошенную дату
                    // количество умножаем на стоимость в этот день
                    // уж постараемся чтобы таймстемпы у количества и сырых данных совпадали
                    // возможным расхождением из-за дырок в данных пренебрегаем
                    startAmount += serie.y * source_data[key][i]
                    break
                }
            }
        }
    }
    amountForEach = startAmount / source_data.length

    countsPerTicker = new Map()
    if (!previous_count) {
        // простой кейс - берем первую стоимость которая есть
        for (const [key, value] of source_data.entries()) {
            countsPerTicker.set(key, amountForEach / value[0].y)
        }
    } else {
        // кейс сложнее - нужно посчитать сколько можно купить в правильный день
        for (const [key, value] of data.entries()) {
            for (i = 0; i < value.length; ++i) {
                serie = value[i]
                if (serie.x >= timestamp) {  // попали в запрошенную дату
                    countsPerTicker[key] = amountForEach / value[0].y
                    break
                }
            }
            
        }
    }
    return countsPerTicker
}

function visualize(e) {
    fillGaps()
    // console.log(SOURCE_DATA)
    console.log("visualization started!")
    draw('chart-row-data', SOURCE_DATA, 'price')
    countsPerTicker = calcBalancedTickersCount(SOURCE_DATA, null, null)
    console.log(countsPerTicker)
    draw('chart-row-data2', countsPerTicker, 'count')

}

document.getElementById('file-input').addEventListener('change', readFiles, false);
document.getElementById('visualize').addEventListener('click', visualize, false);
  
readFiles(null, true)
visualize(null)