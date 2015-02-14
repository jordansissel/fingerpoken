'use strict';

angular.module('myApp.onkyo', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/onkyo', {
    templateUrl: 'onkyo/onkyo.html',
    controllerAs: "onkyo",
    controller: 'OnkyoCtrl'
  });
}])
.directive("fancy", function() {
  return {
    template: "fancy",
    link: function($scope, $element) {
      $element.on("click", function() {
        $scope.onkyo.fancy()
      });
    }
  };
})
.controller('OnkyoCtrl', ['$http', function($http) {
  this.fancy = function() {
    $http.post("/lirc/once/onkyo/GAME", {})
  };
}]);
