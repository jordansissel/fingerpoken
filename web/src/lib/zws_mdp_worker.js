var ForeverSocket = require("lib/forever_socket")
const FINAL_FRAME = '0';
const MORE_FRAME = '1';
const CLIENT_HEADER = "MDPC01";
const COMMAND_REQ = "\x01";

function MajordomoWorker(url, service) {
  // TODO(sissel): Parse the url and add 'type=req'
  this.url = url + "?type=dealer";
  this.socket = new ForeverSocket(this.url, "ZWS1.0");
  this.service = service

  this.heartbeatInterval = 5000; // milliseconds
  this.maxMissedHeartbeats = 3;

  var self = this;
  this.socket.setMessageHandler(function(m) { self.handleMessage(m) });
  this.socket.setConnectedHandler(function() { self.handleConnected() });

}

MajordomoWorker.prototype.heartbeat = function() {
  console.log("Sending heartbeat.");
  this.socket.send([
    MORE_FRAME + "",
    MORE_FRAME + "MDPW01",
    FINAL_FRAME + "\x04",
  ]);
}

MajordomoWorker.prototype.handleMessage = function(request) {
  this.lastSeen = Date.now();
  console.log(request.data.length + ":" + request.data);
}

MajordomoWorker.prototype.handleConnected = function() {
  // Send READY message
  this.socket.send([
    MORE_FRAME + "",
    MORE_FRAME + "MDPW01",
    MORE_FRAME + "\x01",
    FINAL_FRAME + this.service
  ]);
  var self = this;
  this.heartbeatTimer = setInterval(function() { self.heartbeat() }, this.heartbeatInterval);
}

MajordomoWorker.prototype.handleClose = function() {
  if (this.heartbeatTimer !== undefined) {
    clearInterval(this.heartbeatTimer);
    this.heartbeatTimer = undefined;
  }
}

MajordomoWorker.prototype.close = function() {
  this.socket.close();
}

module.exports = MajordomoWorker;
