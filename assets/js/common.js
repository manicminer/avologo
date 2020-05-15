function fetchEntry(id) {
    $.getJSON("/log/" + id, function(data) {
        $("#field_timestamp").html("<td>Timestamp</td><td>" + data.timeformatted + " (unix: " + data.timestamp +  ")</td>");
        $("#field_host").html("<td>Host</td><td>" + data.host + "</td>");
        $("#field_source").html("<td>Source</td><td>" + data.source + "</td>");
        $("#field_message").html("<td>Message</td><td>" + data.message + "</td>");
        $("#field_severity").html("<td>Severity</td>" + formatSeverity(data.severity));
    });
}

function fetchClients() {
    $.getJSON("/queryClients", function(data) {
        data.forEach(function(entry) {
            $("#clienttable > tbody:last-child").append("<tr>" + formatClientRow(entry) + "</tr>");
        });
    });
}

function fetchResults(query) {
    clearTable();
    var queryString = buildQueryString(query);

    $.getJSON(queryString, function(data) {
        $.each(data, function(index) {
            addToTable(data[index]);
        });
    });

    hashParams["search"] = query;

    lastUpdate = Math.floor(Date.now() / 1000);
}

function fetchResultsSince(query, since) {
    var queryString = buildQueryString(query) + '&lower=' + since;

    $.getJSON(queryString, function(data) {
        $.each(data, function(index) {
            addToTableLive(data[index]);
        });
    });

    lastUpdate = Math.floor(Date.now() / 1000);
}

function buildQueryString(query) {
    var queryString = '/query?search=' + query;
    if (hashParams["host"] != "") {
        queryString = queryString + '&host=' + hashParams["host"];
    }
    if (hashParams["source"] != "") {
        queryString = queryString + '&source=' + hashParams["source"];
    }
    if (hashParams["severity"] != "") {
        queryString = queryString + '&severity=' + hashParams["severity"];
    }

    return queryString;
}

function parseHash() {
    var hashArray = location.hash.split("&");
    hashArray.forEach(function(entry) {
        console.log(entry);
    });
}

function formatSeverity(severity) {
    if (severity == 2) {
        return "<td class='severity_warn'>WARN</td>";
    }
    else if (severity == 3) {
        return "<td class='severity_error'>ERROR</td>";
    }
    else {
        return "<td class='severity_info'>INFO</td>";
    }
}

function round(value, decimals) {
    return Number(Math.round(value+'e'+decimals)+'e-'+decimals);
}
  
  

function formatRow(entry) {
    return "<td><a title='View Entry' href='/view/" + entry.id + "'>" + entry.timeformatted +
            "</a></td><td>" + entry.message +
            "</td><td>" + entry.source + " from " + entry.host +
            "</td>" + formatSeverity(entry.severity);
}

function formatClientRow(entry) {
    return "<td>" + entry.host +
            "</td><td>" + entry.count +
            "</td><td>" + entry.rate.toFixed(5) +
            " / hour</td><td>" + entry.lastmessage +
            "</td>";
}

function clearTable() {
    $("#logresults").empty();
}

function addToTable(entry) {
    addHost(entry.host);
    addSource(entry.source);

    $("#maintable > tbody:last-child").append("<tr>" + formatRow(entry) + "</tr>");
    refreshSelects();
}

function addToTableLive(entry) {
    addHost(entry.host);
    addSource(entry.source);

    $("<tr id='"+ entry.id +"' style='display:none;'>" + formatRow(entry) + "</tr>").prependTo("#maintable > tbody");
    $('#' + entry.id).fadeIn(1000);
    refreshSelects();
}

function addHost(host) {
    if (!hostList.includes(host)) {
        hostList.push(host);
    }
}

function addSource(source) {
    if (!sourceList.includes(source)) {
        sourceList.push(source);
    }
}

function refreshSelects() {
    refreshHostList();
    refreshSourceList();
}

function refreshHostList() {
    $("#hostfilter").find('option').remove();
    $('#hostfilter').append($('<option>', {value:'', text:'Filter by Host'}));
    hostList.forEach(element => $('#hostfilter').append($('<option>', {value:element, text:element})));
    
    if (hostList.includes(hashParams["host"])) {
        $("#hostfilter").val(hashParams["host"]);
    }
    else {
        $("#hostfilter").val('');
    }
}

function refreshSourceList() {
    $("#sourcefilter").find('option').remove();
    $('#sourcefilter').append($('<option>', {value:'', text:'Filter by Source'}));
    sourceList.forEach(element => $('#sourcefilter').append($('<option>', {value:element, text:element})));
    
    if (sourceList.includes(hashParams["source"])) {
        $("#sourcefilter").val(hashParams["source"]);
    }
    else {
        $("#sourcefilter").val('');
    }
}
