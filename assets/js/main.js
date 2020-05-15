// Entry point
let hostList = [];
let sourceList = [];
var hashParams = {
    "host" : "",
    "source" : "",
    "severity" : ""
};
var lastUpdate = Math.floor(Date.now() / 1000);

fetchResults("");

/*
    Refresh every 5 seconds
*/
window.setInterval(function() {
    fetchResultsSince($("#searchbox").val(), lastUpdate);
}, 5000);

/*
    Handlers
*/
$("#searchform").submit(function(e){
    e.preventDefault();
    fetchResults($("#searchbox").val());
});

$("#hostfilter").change(function() {
    hashParams["host"] = $("#hostfilter").val();
    fetchResults($("#searchbox").val());
});

$("#sourcefilter").change(function() {
    hashParams["source"] = $("#sourcefilter").val();
    fetchResults($("#searchbox").val());
});

$("#severityfilter").change(function() {
    hashParams["severity"] = $("#severityfilter").val();
    fetchResults($("#searchbox").val());
});

$("#search-button").click(function() {
    fetchResults($("#searchbox").val());
});

$("#additional-options-toggle").click(function() {
    $("#additional-options").slideToggle(200);
});