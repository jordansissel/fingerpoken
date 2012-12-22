line_length = (start, end) ->
  return Math.sqrt(
    Math.pow(start[0] - end[0], 2.0) + Math.pow(start[1] - end[1], 2.0)
  )

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
    console.log(obj)
    if !@connected()
      console.log("Not connected")
      @queue.push(obj)
    else
      @socket.send(JSON.stringify(obj))

class Controller
  constructor: (@element) ->
    @logger = new WSLogger("ws://10.0.0.3:8081/")
    cancel = (e) => 
      e.preventDefault()
      e.cancelBubble = true
      e.returnValue = false
      e.stopPropagation()
      @log("Cancelling " + e.constructor.toString())
      return false

    # Disable scrolling by dragging
    $("#controller").bind("touchmove", false)

    @canvas = d3.select(@element).append("svg").node()

    d3.select(@canvas)
      .attr("id", "canvasthing")
      .on("touchstart", () => @touchstart())
      .on("touchmove", () => @touchmove())
      .append("circle").attr("cx", "30").attr("cy", "30").attr("r", 50).attr("id", "finger")
    @log("ready")

  touchstart: () ->
    touches = d3.touches(@canvas)
    d3.select(@canvas).select("#finger")
      .attr("cx", touches[0][0]).attr("cy", touches[0][1])
    @lastf1 = touches[0]

  touchmove: () ->
    touches = d3.touches(@canvas)
    f1 = touches[0]
    distance = line_length(f1, @lastf1)
    d3.select(@canvas).select("#finger")
      .attr("cx", touches[0][0]).attr("cy", touches[0][1])
    @log("Distance: " + distance + " ---> " + f1 + " <=> " + @lastf1)
    @lastf1 = f1

  click: (event) ->
    @logger.log("CLICK")

  log: (obj) ->
    @logger.log(obj)
    
window.addEventListener("load", () -> 
  new Controller(document.querySelector("#controller"))
)
