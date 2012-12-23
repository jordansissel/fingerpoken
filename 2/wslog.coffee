# A websocket logger.
#
# * events sent while not connected will be queued until connection is healthy
# * on disconnect/error, we will automatically reconnect.
#
# Example:
#
#   logger = new WSLogger("ws://somehost:12345/")
#   logger.log({ some: "object" })
# 
# Logs go to 'console.log' as well as the websocket.

class WSLogger
  constructor: (url) ->
    @url = url
    @connect()
    @queue = []
    @is_connected = false

  connected: () ->
    return @is_connected

  connect: () ->
    @socket = new WebSocket(@url)
    @socket.onopen = (event) => 
      @is_connected = true
      console.log("Logger connected")
      for obj in @queue
        @socket.send(JSON.stringify(obj))
      @queue = []
    @socket.onerror = (event) => 
      console.log("websocket error: " + event)
      @socket.onclose(event)
    @socket.onclose = (event) => 
      @is_connected = false
      @socket.close()
      retry = () => @connect()
      setTimeout(retry, 1000)

  log: (obj) ->
    event = @prepare(obj)
    if !@connected()
      console.log("Not connected")
      @queue.push(event)
    else
      @socket.send(JSON.stringify(event))

  # Prepare an object for logging. This will try to do the most correct
  # thing to convert any object into something that can be turned into
  # JSON.
  prepare: (obj) ->
    #now = (new Date()).toISOString()
    event = {}
    safetypes = ["number", "string", "boolean"]
    if typeof(obj) == "string"
      event.message = obj
    else if typeof(obj) == "object"
      for key, value of obj
        if typeof(value) in safetypes
          event[key] = value 
        else if value == undefined
          # skip
        else if typeof(value) == "function"
          # skip
        else
          try
            # See if we can convert to json, use the object if this succeeds
            JSON.stringify(value)
            event[key] = value
          catch TypeError
            # JSON conversion failed, use the toString value
            event[key] = value.toString()
         
# expose this to the browser.
window.WSLogger = WSLogger
