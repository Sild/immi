
THROTTLING_DATA = 10

SOURCE_DATA = new Map() // ticker -> [(x1, y1), (x2, y2)]
CHARTS = new Map() // xpath_id -> Chart obj

function fillMock()
{
    for (i = 0; i < 2; ++i)
    {
        label = "mock" + i
        ticker_data = []
        value = 5 + i * 1000
        var now = new Date();
        for (var d = new Date(2020, 0, 1); d <= now; d.setDate(d.getDate() + 5)) {
            switch (d.getDate())
            {
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
            ticker_data.push({x: new Date(d).getTime(), y: value});
        }
        SOURCE_DATA.set(label, ticker_data)
    }
}

function readFiles(e, mock=null)
{
    if (mock)
    {
        fillMock()
        return
    }
    for (let f_idx = 0; f_idx < e.target.files.length; f_idx++)
    {
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
                if (i % THROTTLING_DATA != 0)
                    continue
                var line = splitted[i].split(",")
                var value = parseFloat(line[5])
                if (this.FileName.includes("dollar"))
                    value *= 80
                ticker_data.push({x: line[0], y: value})
            }
            ticker_data.sort()
            SOURCE_DATA.set(this.FileName, ticker_data)
            console.log("reeding file " + this.FileName + " done.")
        };
        reader.FileName = file_name
        
        reader.readAsText(file);
    }
}

function draw(xpath_chart_id, data)
{
    // data contains label -> [(x1, y1), (x2, y2)]
    console.log("drawing", xpath_chart_id)
    var series_data = []
    for (const [key, value] of data.entries()) {
        randomColor = "#" + Math.floor(Math.random()*16777215).toString(16);
        // console.log(randomColor)
        series_data.push({
            label: key,
            backgroundColor: "rgba(255,99,132,0.2)",
            borderColor: randomColor,
            borderWidth: 2,
            hoverBackgroundColor: randomColor,
            hoverBorderColor:randomColor,
            data: value,
            pointRadius: 1,
            pointHoverRadius: 5,
            // xAxisID: "time",
            // yAxisId: "amount"
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
                    // time: {
                    //     unit: 'day'
                    // }
                }
            }
        }
    });

}

function visualize(e)
{
    // console.log(SOURCE_DATA)
    console.log("visualization started!")
    draw('chart-row-data', SOURCE_DATA)
    draw('chart-row-data2', SOURCE_DATA)

}

document.getElementById('file-input').addEventListener('change', readFiles, false);
document.getElementById('visualize').addEventListener('click', visualize, false);
  
readFiles(null, true)
visualize(null)