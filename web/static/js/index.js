/* global angular */
var app = angular.module("Plex-Sync", []);

app.controller("indexImages", function($scope, $http){
    $http.get("/api/shows")
    .then(function(response) {
        $scope.records = response.data;
    });
});