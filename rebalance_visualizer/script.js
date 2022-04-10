
THROTTLING_DATA = 10

source_data = new Map()
var MY_CHART = null

function readFiles(e)
{
    for (let f_idx = 0; f_idx < e.target.files.length; f_idx++)
    {
        var file = e.target.files[f_idx];
        var file_name = file.name
        if (!file) {
            console.error("fail to read file " + file_name)
            return;
        }
        console.log("reeding file " + file_name)

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
            source_data.set(this.FileName, ticker_data)
            console.log("reeding file " + this.FileName + " done.")
        };
        reader.FileName = file_name
        
        reader.readAsText(file);
    }
}

function visualize(e)
{
    console.log(source_data)
    console.log("visualization started!")
    var series_data = []
    for (const [key, value] of source_data.entries()) {
        randomColor = "#" + Math.floor(Math.random()*16777215).toString(16);
        console.log(randomColor)
        series_data.push({
            label: key,
            // backgroundColor: "rgba(255,99,132,0.2)",
            borderColor: randomColor,
            borderWidth: 2,
            hoverBackgroundColor: randomColor,
            hoverBorderColor:randomColor,
            data: value,
            // xAxisID: "time",
            // yAxisId: "amount"
        })
    }
    console.log(series_data)
    var chart_data = {
        datasets: series_data
    };
    if (MY_CHART)
        MY_CHART.destroy()

    MY_CHART = new Chart("chart", {
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

document.getElementById('file-input').addEventListener('change', readFiles, false);
document.getElementById('visualize').addEventListener('click', visualize, false);
  
