
data = {}

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
                var line = splitted[i].split(",")
                var value = parseFloat(line[5])
                if (this.FileName.includes("dollar"))
                    value *= 80
                ticker_data.push({date: line[0], close_cost: value})
            }
            data[this.FileName] = ticker_data
            console.log("reeding file " + this.FileName + " done.")
        };
        reader.FileName = file_name
        
        reader.readAsText(file);
    }
}

function visualize(e)
{
    console.log(data)
    console.log("visualization started!")
    var chart_data = {
        labels: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul"],
        datasets: [
          {
            label: "Dataset #1",
            backgroundColor: "rgba(255,99,132,0.2)",
            borderColor: "rgba(255,99,132,1)",
            borderWidth: 2,
            hoverBackgroundColor: "rgba(255,99,132,0.4)",
            hoverBorderColor: "rgba(255,99,132,1)",
            data: [65, 59, 20, 81, 56, 55, 40]
          }
        ]
      };
      
    // var option = {
    //     responsive: false,
    //     scales: {
    //         y: {
    //             stacked: true,
    //             grid: {
    //                 display: true,
    //                 color: "rgba(255,99,132,0.2)"
    //             }
    //         },
    //         x: {
    //             grid: {
    //                 display: false
    //             }
    //         }
    //     }
    //   };
      
    new Chart("chart", {
        type: "bar",
        // options: option,
        data: chart_data
    });

}

document.getElementById('file-input').addEventListener('change', readFiles, false);
document.getElementById('visualize').addEventListener('click', visualize, false);
  
