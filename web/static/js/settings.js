$(document).ready(function() {
    getSettings();

    $("#generalSave").click(function() {
        saveSettings();
    });
    $("#localserverSave").click(function() {
        saveSettings();
    });
    $("#remoteserverSave").click(function() {
        saveSettings();
    });
});

function saveSettings() {
    var settings = JSON.stringify(createSettingsJSON());
    $.ajax({
        type: "POST",
        url: "/api/settings",
        data: settings,
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function(data) {
            loadSettings(data);
        },
        failure: function(errMsg) {
            console.log("error writing settings to DB: " + errMsg);
        }
    });
}

function loadSettings(data) {
    // Populate settings in form
    $("#port").val(data.WebserverPort);
    $("#syncInterval").val(data.SyncInterval);
    $("#logSize").val(data.MaxLogSize);
    $("#logCount").val(data.MaxLogCount);
    $("#logAge").val(data.MaxLogAge);
    $("#cacheLifetime").val(data.CacheLifetime);

    $("#localserverName").val(data.LocalServer.Name);
    $("#localserverHostname").val(data.LocalServer.Hostname);
    $("#localserverPort").val(data.LocalServer.Port);
    $("#localserverSSL").val(data.LocalServer.Ssl); //Fix this

    $("#remoteserverName").val(data.RemoteServer.Name);
    $("#remoteserverHostname").val(data.RemoteServer.Hostname);
    $("#remoteserverPort").val(data.RemoteServer.Port);
    $("#remoteserverSSL").val(data.RemoteServer.Ssl); //Fix this
}

function getSettings() {
    $.ajax({
        type: "GET",
        url: "/api/settings",
        dataType: "json",
        success: function(data) {
            loadSettings(data);
        },
        failure: function(errMsg) {
            alert("error getting settings from DB: " + errMsg);
        }
    });
}

function createSettingsJSON() {
    return (settings = {
        SyncInterval: JSON.parse($("#syncInterval").val()),
        WebserverPort: JSON.parse($("#port").val()),
        MaxLogSize: JSON.parse($("#logSize").val()),
        MaxLogCount: JSON.parse($("#logCount").val()),
        MaxLogAge: JSON.parse($("#logAge").val()),
        CacheLifetime: JSON.parse($("#cacheLifetime").val()),
        LocalServer: {
            Name: $("#localserverName").val(),
            Hostname: $("#localserverHostname").val(),
            Port: JSON.parse($("#localserverPort").val()),
            Ssl: JSON.parse($("#localserverSSL").val())
        },
        RemoteServer: {
            Name: $("#remoteserverName").val(),
            Hostname: $("#remoteserverHostname").val(),
            Port: JSON.parse($("#remoteserverPort").val()),
            Ssl: JSON.parse($("#remoteserverSSL").val())
        }
    });
}
