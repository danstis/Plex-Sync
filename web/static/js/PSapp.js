/* global angular */
var PSapp = angular.module("Plex-Sync", []);

PSapp.controller("indexImages", function($scope, $http) {
    $http.get("/api/shows").then(function(response) {
        $scope.records = response.data;
    });
});

PSapp.controller("settings", function($scope, $http) {
    $http.get("/api/settings").then(function(response) {
        $scope.settings = response.data;
    });

    $scope.saveSettings = function() {
        $http.post("/api/settings", $scope.settings).then(
            // Success
            function() {
                $scope.statusMessage = "Settings Updated";
                $scope.statusLevel = "alert-success";
            },
            // Failure
            function(errMsg) {
                $scope.statusMessage = "Error updating settings: " + errMsg;
                $scope.statusLevel = "alert-warning";
            }
        );
    };
});
