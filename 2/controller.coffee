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

window.Controller = Controller
