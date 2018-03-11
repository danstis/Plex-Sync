$(document).ready(function() {
    $.ajax({
        type: "GET",
        url: "/api/version",
        dataType: "json",
        success: function(data) {
            $("#shortVersion").html("Plex-Sync v" + data.shortVersion);
            $("#fullVersion").html("Plex-Sync v" + data.fullVersion);
        },
        failure: function(errMsg) {
            console.log("error getting version from API: " + errMsg);
        }
    });
});
