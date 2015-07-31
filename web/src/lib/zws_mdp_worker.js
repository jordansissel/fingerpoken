var ForeverSocket = require("lib/forever_socket")
const FINAL_FRAME = '0';
const MORE_FRAME = '1';
const CLIENT_HEADER = "MDPC01";
const COMMAND_REQ = "\x01";
const COMMAND_HEARTBEAT = "\x04"
const COMMAND_REQUEST = "\x02"

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
  this.buffer = []
  this.computeBrokerExpiration()
}


MajordomoWorker.prototype.computeBrokerExpiration = function() {
  this.brokerExpiration = Date.now() + (this.heartbeatInterval * this.maxMissedHeartbeats);
}

MajordomoWorker.prototype.heartbeat = function() {
  if (Date.now() > this.brokerExpiration) {
    console.log("Broker hasn't been seen for a while. Closing connection.")
    this.close();
    return;
  }
  console.log("Sending heartbeat.");
  this.socket.send([
    MORE_FRAME + "",
    MORE_FRAME + "MDPW01",
    FINAL_FRAME + "\x04",
  ]);
}

MajordomoWorker.prototype.handleMessage = function(request) {
  this.computeBrokerExpiration()
  data = request.data
  flag = data.slice(0, 1)
  response = data.slice(1)
  this.buffer.push(response)

  if (flag != FINAL_FRAME) {
    return
  }

  if (validateWorkerMessage(this.buffer)) {
    switch (this.buffer[2]) {
      case COMMAND_HEARTBEAT:
        // Nothing to do...
        break;
      case COMMAND_REQUEST:
        reply = this.handleRPC(this.buffer);
        this.socket.send(reply)
        break;
      default:
        console.log("Unknown command: " + this.buffer[2].charCodeAt(0));
        break;
    }
    // TODO(sissel): Call callback
  }

  //console.log("Resetting buffer");
  this.buffer = [];
}

MajordomoWorker.prototype.handleRPC = function(buffer) {
  var client = buffer[3];
  console.log(JSON.stringify(buffer));
  if (client == "") {
    console.log("Invalid request command (fourth frame must contain the client id)");
    return
  }
  if (buffer[4] != "") {
    console.log("Invalid request command (fifth frame must be empty)");
    return
  }

  var obj = JSON.parse(buffer[5])
  console.log("Method: " + obj.method);

  // TODO(sissel): Implement JSONRPC for this.
  // buffer.slice(5....) is the request payload. Parse as JSON, handle as JSONRPC.
  var reply = {
    id: obj.id,
    result: null,
    error: null
  };
  if (this.rpcHandler === undefined) {
    reply.error = "This worker has no RPC handler set."
  } else if (this.rpcHandler[obj.method] === undefined) {
    reply.error = "No such method named '" + obj.method + "'";
  } else {
    this.rpcHandler[obj.method](obj.params, reply);
  }

  reply_frames = [
    MORE_FRAME + "",
    MORE_FRAME + "MDPW01",
    MORE_FRAME + "\x03",
    MORE_FRAME + client,
    MORE_FRAME + "",
    FINAL_FRAME + JSON.stringify(reply)
  ]
  console.log(JSON.stringify(reply));
  return reply_frames
}

var validateWorkerMessage = function(buffer) {
  if (buffer[0] != "") {
    console.log("First frame must be empty");
    return false;
  }

  if (buffer[1] != "MDPW01") {
    console.log("Second frame must be MDPW01");
    return false;
  }

  if (buffer[2].length != 1) {
    console.log("Invalid command frame");
  }

  if (buffer[2] == COMMAND_HEARTBEAT) {
    if (buffer.length != 3) {
      console.log("Got invalid frame count for heartbeat command");
      return false;
    }
  }

  return true
}

MajordomoWorker.prototype.handleConnected = function() {
  // Send READY message
  this.socket.send([
    MORE_FRAME + "",
    MORE_FRAME + "MDPW01",
    MORE_FRAME + "\x01",
    FINAL_FRAME + this.service
  ]);
  this.scheduleHeartbeat();
}

MajordomoWorker.prototype.scheduleHeartbeat = function() {
  var self = this;
  this.cancelHeartbeat();
  this.heartbeatTimer = setInterval(function() { self.heartbeat() }, this.heartbeatInterval);
}

MajordomoWorker.prototype.cancelHeartbeat = function() {
  if (this.heartbeatTimer !== undefined) {
    clearInterval(this.heartbeatTimer);
    this.heartbeatTimer = undefined;
  }
}

MajordomoWorker.prototype.handleClose = function() {
  this.cancelHeartbeat();
}

MajordomoWorker.prototype.close = function() {
  this.socket.close();
}

module.exports = MajordomoWorker;
