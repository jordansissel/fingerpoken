'use strict';

function ForeverWebSocket(url) {
  this.url = url;
  this.connect();
}

ForeverWebSocket.connect = function() {
  this.ws = new WebSocket(url);
  this.setup_handlers(this.ws);
};

ForeverWebSocket.connect_if_necessary = function() {
  if (this.ws !== undefined) {
    this.connect();
  }
}

ForeverWebSocket.setup_handlers = function(ws) {
  var _this = this;
  if (this.onmessage !== undefined) {
    ws.onmessage(this.onmessage);
  }

  ws.onclose = function(event) {
    _this.connect();
  };
};

ForeverWebSocket.onmessage = function(callback) {
  this.callback = callback;
  if (this.ws !== undefined) {
    this.ws.onmessage(callback);
  }
}

ForeverWebSocket.send = function(arg) {
  this.connect_if_necessary();
  return this.ws.send(arg);
}

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
  this.ws = new ForeverWebSocket("ws://" + document.location.hostname + ":" + document.location.port + "/websocket");
  this.fancy = function() {
    this.ws.send(JSON.stringify({ "lirc": { "command": "once", "remote": "onkyo", "code": "GAME" }});
  };
}]);
