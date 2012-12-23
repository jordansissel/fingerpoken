# A websocket logger.
#
# * events sent while not connected will be queued until connection is healthy
# * on disconnect/error, we will automatically reconnect.
#
# Example:
#
#   logger = new WSLogger("ws://somehost:12345/")
#   logger.log({ some: "object", another: "value", hello: [1,2,3] })
# 
# Logs go to 'console.log' as well as the websocket.
#
# You should probably also use this to capture javascript exceptions/errors as
# well. To do so, you will need to use 'window.onerror':
#
#     window.onerror = function(message, url, line) {
#       if (window.$logger === undefined) {
#         window.$logger = new WSLogger("ws://somehost:12345/")
#       }
#       window.$logger.log({ message: message, url: url, line: line })
#     }

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
          event[key] = null # pretend null
        else if typeof(value) == "function"
          # skip functions
        else
          try
            # See if we can convert to json, use the object if this succeeds
            JSON.stringify(value)
            event[key] = value
          catch e
            # JSON conversion failed, use the toString value
            event[key] = value.toString()

    return event
         
# expose this to the browser.
window.WSLogger = WSLogger
