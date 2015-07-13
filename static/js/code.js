(function() {
  var logs = []
  function log(message) {
    var el = document.querySelector("log");
    var entry = document.createElement("pre")
    entry.innerHTML = message
    el.appendChild(entry)
    el.scrollTop = el.scrollHeight;
  }

  var Fingerpoken = function() {
    log("New fingerpoken");
    this.connect();
  }

  Fingerpoken.prototype.connect = function() { 
    if (this.websocket !== undefined) {
      this.websocket.close();
    }

    var self = this;
    var url = "ws://" + document.location.host + "/ws";
    log(url);
    this.websocket = new WebSocket(url);
    log("New websocket" + this.websocket);
    this.websocket.onopen = function(event) {
      log("Socket ready");
      self.ready()
    }

    this.websocket.onclose = function(event) {
      log("Socket closed, will reopen")
      this.websocket = undefined;
      setTimeout(function() { self.connect(); }, 250);
    }
  }

  Fingerpoken.prototype.ready = function() {
    this.onready();
  }

  Fingerpoken.prototype.onready = function() {
    // Nothing, users can override per-instance.
  }

  Fingerpoken.prototype.send = function(message) {
    this.websocket.send(message);
  }

  var touch_el = document.querySelector("touch");
  function handleMove(fp, event) {
    event.preventDefault();
    var text = "";
    touches = event.touches
    var x = touches[0].clientX;
    var y = touches[0].clientY;
    touch_el.innerHTML = "x: " + x + ", y: " + y;

    var rpc = {
      method: "Mouse.Move",
      params: [ { x: x, y: y } ],
      id: 1,
    }
    fp.send(JSON.stringify(rpc));
  };

  log("Starting up...");
  try {
    var fp = new Fingerpoken();
    fp.onready = function() {

    }

    var body = document.body;
    body.addEventListener("touchmove", function(e) { handleMove(fp, e); }, false);
    body.addEventListener("touchstart", function(e) { e.preventDefault() }, false);
    body.addEventListener("touchend", function(e) { e.preventDefault() }, false);
    body.addEventListener("touchcancel", function(e) { e.preventDefault() }, false);
  } catch (err) {
    log(err);
  }

})();
