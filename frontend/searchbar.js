const backendApiAddress = "http://localhost:8080"
searchBarVersions = [
    {
        id: "only-frontend",
        searchFunc: localExecution
    },
    {
        id: "dumb-backend",
        searchFunc: slowApi
    },
    {
        id: "suggested-backend",
        searchFunc: indexedTable
    },
    {
        id: "async-backend",
        searchFunc: indexedTableAsyncCalls
    },
    {
        id: "improved-backend",
        searchFunc: prefixTable
    }
]


searchBarVersions.forEach(function (element) {
    let searchbar = document.querySelector("#" + element.id + " .searchbar");
    let suggestionsContainer = document.querySelector("#" + element.id + " .suggestions");
    let metricsContainer = document.querySelector("#" + element.id + " .performance");

    let timerId;
    element["count"] = 0
    element["fast"] = 0
    element["average"] = 0
    element["slow"] = 0

    searchbar.addEventListener("keyup", function (event) {
        // cancel previous request
        clearTimeout(timerId);

        let searchValue = event.target.value;

        timerId = setTimeout(async function () {
            let start = performance.now();
            element.searchFunc(searchValue, (suggestions) => {
                let end = performance.now();
                displaySearchResults(suggestionsContainer, suggestions)
                // calculate and display metrics
                let timeInMs = Math.ceil(end - start)
                element["count"]++
                if (timeInMs > element["slow"]) {
                    element["slow"] = timeInMs
                }
                if (timeInMs != 0 && (timeInMs < element["fast"] || element["count"] == 1)) {
                    element["fast"] = timeInMs
                }
                let sum = element["average"] * (element["count"] - 1) + timeInMs;
                element["average"] = Math.ceil(sum / element["count"]);
                console.log(`impl: ${element.id}, search: ${searchValue}, time: ${timeInMs}`)
                displayPerformanceMetrics(metricsContainer, element["average"], element["fast"], element["slow"])
            });


        }, 300);
    });
});


function displaySearchResults(suggestionsDiv, suggestions) {
    suggestionsDiv.innerHTML = "";
    suggestions.forEach(suggestion => {
        let suggestionEl = document.createElement("div");
        suggestionEl.innerText = suggestion;
        suggestionsDiv.appendChild(suggestionEl);
    });
}

function displayPerformanceMetrics(metricsDiv, average, fast, slow) {
    metrics = [
        { text: "Fastest", time: fast },
        { text: "Slowest", time: slow },
        { text: "Average", time: average }
    ]
    metricsDiv.innerHTML = "";
    metrics.forEach(metric => {
        let metricEl = document.createElement("div");
        metricEl.innerText = `${metric.text}: ${metric.time}ms`;
        if (metric.time > 300) {
            metricEl.className = "slow"
        }
        if (metric.time < 50) {
            metricEl.className = "fast"
        }
        metricsDiv.appendChild(metricEl);
    });
}

///////////////////////////////////
// Different search implementations 
///////////////////////////////////

// from local file
function localExecution(searchValue, callback) {
    callback(getFromWordsJSFile(searchValue, 5))
}

// from a slow API
function slowApi(searchValue, callback) {
    getListFromApi(searchValue, "/mysql/get-words-no-index", callback)
}

// from an indexed table
function indexedTable(searchValue, callback) {
    getListFromApi(searchValue, "/mysql/get-words", callback)
}

// from an indexed table but queries performed in async manner
function indexedTableAsyncCalls(searchValue, callback) {
    getListFromApi(searchValue, "/mysql/get-words-async", callback)
}

// from an indexed table + indexed prefix table
function prefixTable(searchValue, callback) {
    getListFromApi(searchValue, "/mysql/get-words-v2", callback)
}


///////////////////////////////
// Helper functions for search
///////////////////////////////

function getListFromApi(searchValue, subPath, callback) {
    if (searchValue == "") {
        callback([])
    } else {
        fetch(`${backendApiAddress}${subPath}?search=${searchValue}`)
            .then(response => response.json())
            .then(data => callback(data))
            .catch(error => console.error(error));
    }
}

// Search from the words array in words.js file
function getFromWordsJSFile(searchValue, results) {
    if (searchValue == "") {
        return []
    }
    let expressions = []
    let aggregatedResult = []
    // expressions based on most desirable matches
    expressions.push(new RegExp("^" + searchValue + ".*", 'i')) // joh*
    expressions.push(new RegExp(searchValue + ".*", 'i')) // *joh*
    expressions.push(new RegExp(".*" + searchValue.split('').join('.*') + ".*", 'i')) // *j*o*h*

    let i = 0
    while (aggregatedResult.length <= results && i < expressions.length) {
        let filteredWords = words.filter(word => {
            return expressions[i].test(word);
        });
        if (filteredWords.length > 0) {
            let leftSlots = Math.min(filteredWords.length, results - aggregatedResult.length)
            aggregatedResult = aggregatedResult.concat(filteredWords.slice(0, leftSlots));
        }
        i++
    }

    return aggregatedResult
}