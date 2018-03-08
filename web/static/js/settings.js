$(document).ready(function() {
    // Get settings data from API
    $.getJSON("/api/settings", function(data, status) {
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
        $("#localserverSSL").val(data.LocalServer.Ssl);

        $("#remoteserverName").val(data.RemoteServer.Name);
        $("#remoteserverHostname").val(data.RemoteServer.Hostname);
        $("#remoteserverPort").val(data.RemoteServer.Port);
        $("#remoteserverSSL").val(data.RemoteServer.Ssl);
    });
});
