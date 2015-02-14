'use strict';

angular.module('myApp.navigation', [])
.controller("NavigationCtrl", function() {
  this.links = {
    "Example2": "view1",
    "Foo": "view2",
    "Onkyo": "onkyo"
  };
  this.selected = "Example2";

  this.select = function(text) {
    this.selected = text;
  }

  this.isSelected = function(text) {
    return this.selected === text;
  }
})
.directive("navigation", function() {
  return {
    restrict: "E",
    controllerAs: "navigation",
    controller: "NavigationCtrl",
    templateUrl: "navigation/navigation.html"
  };
})
