var baseConfig = `# Server configuration
# Only required on avologo master node
server:
    host: "{server_host}"
    port: {server_port}

# Database credentials
# Only required on avologo master node
database:
    host: "{db_host}"
    port: {db_port}
    user: "{db_user}"
    password: "{db_pass}"
    dbname: "{db_name}"

# Client configuration
client:
    destination: "127.0.0.1:{server_port}"
    friendly_name: ""
    watch:
    error_keywords:
        - "error"
        - "fatal"
    warning_keywords:
        - "warning"
        - "warn"
`;

$("#test-button").click(function() {
    connectionParams = new Object();
    connectionParams['host'] = $("#db_host").val();
    connectionParams['port'] = $("#db_port").val();
    connectionParams['user'] = $("#db_user").val();
    connectionParams['password'] = $("#db_pass").val();
    connectionParams['dbname'] = $("#db_name").val();

    $.post("/testConnection", connectionParams)
        .done(function(response) {
            if (response['valid']) {
                // Valid connection
                $("#test-results").html("<i class='fa fa-check'></i>&nbsp;&nbsp;Your connection was successful");
                $("#test-results").css("color", "#00d412");
                $("#test-results").fadeIn(200);
            }
            else {
                // Invalid connection
                $("#test-results").html("<i class='fa fa-times'></i>&nbsp;&nbsp;Your connection was unsuccessful");
                $("#test-results").css("color", "#ff0000");
                $("#test-results").fadeIn(200);
            }
    });
});

$("#generate-button").click(function() {

    // Generate config file
    var generatedConfig = baseConfig;
    generatedConfig = generatedConfig.replace(/{db_host}/g, $("#db_host").val());
    generatedConfig = generatedConfig.replace(/{db_port}/g, $("#db_port").val());
    generatedConfig = generatedConfig.replace(/{db_user}/g, $("#db_user").val());
    generatedConfig = generatedConfig.replace(/{db_pass}/g, $("#db_pass").val());
    generatedConfig = generatedConfig.replace(/{db_name}/g, $("#db_name").val());
    generatedConfig = generatedConfig.replace(/{server_host}/g, $("#server_host").val());
    generatedConfig = generatedConfig.replace(/{server_port}/g, $("#server_port").val());

    $("#config-text").val(generatedConfig);
    $('#configmodal').modal('toggle');
});