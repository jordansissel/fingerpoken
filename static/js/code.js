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
    this.rpc_id = 0;

    this.flush_rate = 16; // milliseconds
    this.setup_flush();
  }

  Fingerpoken.prototype.flush = function() {
    if (this.next_rpc === undefined) {
      this.stop_flush();
    } else {
      fp.send(JSON.stringify(this.next_rpc));
      this.next_rpc = undefined;
    }
  };

  Fingerpoken.prototype.stop_flush = function() {
    if (this.live_interval !== undefined) {
      //log("Stopping flush")
      clearInterval(this.live_interval);
      this.live_interval = undefined;
    }
  }

  Fingerpoken.prototype.setup_flush = function() {
    if (this.live_interval !== undefined) {
      this.stop_flush();
    }
    //log("Starting flush")
    var self = this;
    this.live_interval = setInterval(function() {
      self.flush()
    }, this.flush_rate)
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

    this.websocket.onmessage = function(event) {
      console.log(event.data);
      var result = JSON.parse(event.data)
      if (result.id === null) {
        // Notification
        log(result.params.line);
      } else {
        // Reactive congestion control.
        // If latency is larger than flush rate, slow
        // down flush_rate (increase iterval time).
        // If latency is less than flush rate, speed
        // up flush_rate (decreate interval time)
        // This should ad some rough congestion correction
        // where we can absorb latency spikes and recover
        // to whatever the average round-trip latency is.
        var latency = (Date.now() - result.id.ts)
        //log(latency + "|" + self.flush_rate);
        if (latency > self.flush_rate * 2) {
          self.rate(Math.floor(self.flush_rate * 1.1));
        } else if (latency < self.flush_rate / 1.5) {
          self.rate(Math.floor(self.flush_rate * 0.9));
        }
        log("Message[" + latency + "]: " + event.data);
      }
    }
  };

  Fingerpoken.prototype.rate = function(val) {
    if (val < 16) {
      return self.rate(16);
    }
    if (val == self.flush_rate) {
      return;
    }
    log("Setting frame rate: " + val);
    this.flush_rate = val;
    this.setup_flush();
  }

  Fingerpoken.prototype.ready = function() {
    this.onready();
  }

  Fingerpoken.prototype.onready = function() {
    // Nothing, users can override per-instance.
  }

  Fingerpoken.prototype.send = function(message) {
    // TODO(sissel): Add proactive congestion control.
    // This should track how many messages are in-flight/delays/etc
    // and slow the flush_rate if there's too many unanswered calls.
    this.websocket.send(message);
  }

  Fingerpoken.prototype.nextRPCId = function() {
    this.rpc_id += 1;
    return this.rpc_id;
  }

  Fingerpoken.prototype.set_next_rpc = function(rpc) {
    this.next_rpc = rpc;
    if (this.live_interval === undefined) {
      this.setup_flush()
    }
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
      id: { id: fp.nextRPCId(), ts: Date.now() },
    }
    fp.set_next_rpc(rpc)
    //fp.send(JSON.stringify(rpc));
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
