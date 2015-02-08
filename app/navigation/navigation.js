'use strict';

angular.module('myApp.navigation', [])
.directive("navigation", function() {
  return {
    restrict: "E",
    templateUrl: "navigation/navigation.html"
  };
});
