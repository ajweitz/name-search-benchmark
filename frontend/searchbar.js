const backendApiAddress = "http://localhost:8080"
searchBarVersions = [
    {
        id: "only-frontend",
        searchFunc: localExecution
    },
    {
        id: "dumb-backend",
        searchFunc: slowApi
    }
]


searchBarVersions.forEach(function(element) {
    let searchbar = document.querySelector("#" + element.id + " .searchbar");
    let suggestionsContainer = document.querySelector("#" + element.id + " .suggestions");
    let metricsContainer = document.querySelector("#" + element.id + " .performance");

    let fast = slow = average = count = 0
    let timerId;

    searchbar.addEventListener("keyup", function(event) {
        // cancel previous request
        clearTimeout(timerId);

        let searchValue = event.target.value;

        timerId = setTimeout(async function() {
            let start = performance.now();
            element.searchFunc(searchValue,(suggestions)=>{
                let end = performance.now();
                displaySearchResults(suggestionsContainer, suggestions)
                // calculate and display metrics
                let timeInMs = Math.ceil(end-start)
                count++
                if (timeInMs > slow){
                    slow = timeInMs
                } 
                if (timeInMs != 0 && (timeInMs < fast || count == 1)){
                    fast = timeInMs
                } 
                let sum = average * (count - 1) + timeInMs;
                average = Math.ceil(sum / count);
                console.log(count)
                displayPerformanceMetrics(metricsContainer, average, fast, slow)
            });

  
        },300);
    });
});


function displaySearchResults(suggestionsDiv, suggestions ){
    suggestionsDiv.innerHTML = "";
        suggestions.forEach(suggestion => {
            let suggestionEl = document.createElement("div");
            suggestionEl.innerText = suggestion;
            suggestionsDiv.appendChild(suggestionEl);
        });
}

function displayPerformanceMetrics(metricsDiv, average, fast, slow){
    metrics = [
        `Average: ${average}ms`,
        `Fastest: ${fast}ms`,
        `Slowest: ${slow}ms`
    ]
    metricsDiv.innerHTML = "";
    metrics.forEach(metric => {
        let metricEl = document.createElement("div");
        metricEl.innerText = metric;
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
function slowApi(searchValue, callback){
    getListFromApi(searchValue,"/mysql/get-words-no-index",callback)
}


///////////////////////////////
// Helper functions for search
///////////////////////////////

function getListFromApi(searchValue, subPath, callback){
    fetch(`${backendApiAddress}${subPath}?search=${searchValue}`)
    .then(response => response.json())
    .then(data => callback(data))
    .catch(error => console.error(error));
}

// Search from the words array in words.js file
function getFromWordsJSFile(searchValue, results){
    if (searchValue == "") {
        return []
    }
    let expressions = []
    let aggregatedResult = []
    // expressions based on most desirable matches
    expressions.push(new RegExp("^"+searchValue+".*",'i')) // joh*
    expressions.push(new RegExp(searchValue+".*",'i')) // *joh*
    expressions.push(new RegExp(".*"+searchValue.split('').join('.*')+".*",'i')) // *j*o*h*

    let i = 0
    while (aggregatedResult.length <= results && i < expressions.length) {
        let filteredWords = words.filter(word => {
            return expressions[i].test(word);
        });
        if(filteredWords.length > 0){
            let leftSlots = Math.min(filteredWords.length,results-aggregatedResult.length)
            aggregatedResult = aggregatedResult.concat(filteredWords.slice(0,leftSlots));
        }
        i++
    }

    return aggregatedResult
}