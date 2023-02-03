

let searchbar = document.getElementById("searchbar");
let suggestionsContainer = document.getElementById("search-suggestions");

searchbar.addEventListener("keyup", function(event) {
  let searchValue = event.target.value;
      let suggestions = getSuggestions(searchValue)
      suggestionsContainer.innerHTML = "";
      suggestions.forEach(suggestion => {
        let suggestionEl = document.createElement("div");
        suggestionEl.innerText = suggestion;
        suggestionsContainer.appendChild(suggestionEl);
      });
});


function getSuggestions(searchValue) {
    return rank(searchValue, 5)
}

// Search for best match
function rank(searchValue, results){
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