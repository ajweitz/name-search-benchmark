

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
    let searchRegex = new RegExp(".*"+searchValue.split('').join('.*')+".*",'i')
    console.log(searchRegex)
    let filteredNames = names.filter(name => {
        return searchRegex.test(name);
      });

    return rate(filteredNames)
}

// Ideally, we would have some rating logic
function rate(names){
    return names.slice(0,5)
}