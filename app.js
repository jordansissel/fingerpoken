(function() {
  var app = angular.module("fingerpoken", []);

  app.directive("logo", function() {
    return {
      restrict: "E",
      template: "<i class='fa fa-hand-o-up'></i> fingerpoken"
    };
  });

  app.directive("navbar", function() {
    return {
      restrict: "E",
      templateUrl: "navbar.html"
    };
  });

  app.directive("touchpad", function() {
    return {
      restrict: "E",
      template: "<div class='full'></div>
    };
  });
})();
