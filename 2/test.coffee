window.onerror = (message, url, line) =>
  window.$logger ||= new WSLogger("ws://10.0.0.3:8081/")
  window.$logger.log(message: message, url: url, line: line)

distance = (x1, y1, x2, y2) ->
  return Math.sqrt(Math.pow(x1 - x2, 2.0) + Math.pow(y1 - y2, 2.0))

angle = (x1, y1, x2, y2) ->
  # SOHCAHTOA. we compute O and A and use arctan to get the angle.
  return Math.atan((y1 - y2) / (x1 - x2))

# This is required because iOS reuses TouchEvent objects, it seems,
# so we copy each value we care about.
copyTouch = (touch) ->
  return {
    clientX: touch.clientX,
    clientY: touch.clientY,
    identifier: touch.identifier,
    pageX: touch.pageX,
    pageY: touch.pageY,
    screenX: touch.screenX,
    screenY: touch.screenY,
    target: touch.target
  }

# A finger!
class Finger
  constructor: (touch) ->
    @callbacks = {}
    @origin = @touch = copyTouch(touch)

    # How far has this finger moved
    @travel = 0

  trigger: (name, t) ->
    if @callbacks[name]
      for callback in @callbacks[name] 
        callback(this, t)

  move: (t) ->
    t.distance = distance(t.pageX, t.pageY, @touch.pageX, @touch.pageY)
    t.angle = angle(t.pageX, t.pageY, @touch.pageX, @touch.pageY)
    @travel += t.distance
    t.travel = @travel
    @touch = copyTouch(t)
    @trigger("move", t)

  down: (t) -> 
    @touch = t
    @trigger("down", t)

  up: (t) -> 
    @trigger("up", t || @touch)
    delete @touch

  bind: (event, callback) ->
    @callbacks[event] ||= []
    @callbacks[event].push(callback)

  # Give the distance from the original finger-down.
  origin_distance: () ->
    return distance(@origin.pageX, @origin.pageY, @touch.pageX, @touch.pageY)

  origin_angle: () ->
    return angle(@origin.pageX, @origin.pageY, @touch.pageX, @touch.pageY)

class Controller
  constructor: (@element) ->
    @logger = new WSLogger("ws://10.0.0.3:8081/")
    @fingers = {}

    # Disable scrolling by dragging
    $("#controller").bind("touchmove", false)

    @canvas = d3.select(@element).append("svg").node()

    d3.select(@canvas)
      .attr("id", "canvasthing")
      .on("touchstart", () => @touchstart())
      .on("touchend", () => @touchend())
      .on("touchmove", () => @touchmove())
    @log("ready")

  touchstart: () -> # Got a new finger! :)
    d3.event.preventDefault()
    for touch in d3.event.changedTouches
      #@log("finger: " + touch.identifier)
      finger = @fingers[touch.identifier] = new Finger(touch)
      @circle_cursor(finger)

  touchmove: () -> # Moved a finger (or more)
    d3.event.preventDefault()
    for touch in d3.event.changedTouches
      @fingers[touch.identifier].move(touch)

  touchend: () -> # Lost this finger :(
    d3.event.preventDefault()
    for touch in d3.event.changedTouches
      #@log("finger: " + touch.identifier)
      @fingers[touch.identifier].up(touch)
      delete @fingers[touch.identifier]

  log: (obj) ->
    @logger.log(obj)

  circle_cursor: (finger) ->
    @palette ||= d3.scale.category10()
    @palette_i ||= 0
    @palette_i++
    color = @palette(@palette_i)

    circle = d3.select(@canvas).append("circle")
    circle.attr("r", 50)
      .attr("cx", finger.touch.pageX)
      .attr("cy", finger.touch.pageY)
      .attr("stroke", "#000")
      .attr("fill", color)

    finger.bind("move", (finger, touch) => 
      @log(distance: finger.origin_distance(touch), \
           angle: finger.origin_angle(touch), \
           travel: finger.travel)
      circle.attr("cx", touch.pageX).attr("cy", touch.pageY)

      # drop a tracer
      d3.select(@canvas).append("circle")
        .attr("r", "30")
        .attr("cx", finger.touch.pageX)
        .attr("cy", finger.touch.pageY)
        .attr("fill", color)
        .transition()
          .duration(500)
          .attr("r", "0")
          .remove()
    )

    finger.bind("up", (finger, touch) => 
      circle.style("opacity", 1)
      circle.transition().duration(500)
        .style("opacity", 0)
        .attr("r", circle.attr("r") * 0.50)
        .remove()
    )

    
window.addEventListener("load", () -> 
  new Controller(document.querySelector("#controller"))
)
