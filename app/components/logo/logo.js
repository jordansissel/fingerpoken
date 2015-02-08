'use strict';
angular.module('myApp.logo', [])
.directive('logo', function(version) {
  return {
    template: "<i class='fa fa-hand-o-up'></i>"
  };
});

