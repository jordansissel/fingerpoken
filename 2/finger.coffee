distance = (x1, y1, x2, y2) ->
  # distance between two points using pythagorean theorem
  # a^2 + b^2 = c^2, solving for c, where a == (x1 - x2) and b == (y1 - y2)
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

window.Finger = Finger

if typeof(window.test) == "function"
  touch = (x, y) ->
    return {
      clientX: x,
      clientY: y,
      identifier: 0,
      pageX: x,
      pageY: y,
      screenX: x,
      screenY: y
    }

  test("movement", () => 
    f = new Finger(touch(0, 0))
    ok(f.touch.clientX == 0, "clientX is correct")
    ok(f.touch.clientY == 0, "clientY is correct")
    ok(f.travel == 0, "travel should be zero")
    f.move(touch(0, 1))

    ok(f.touch.clientX == 0, "clientX is correct")
    ok(f.touch.clientY == 1, "clientY is correct")
    ok(f.travel == 1, "travel should be 1")

    f.move(touch(1, 1))
    ok(f.touch.clientX == 1, "clientX is correct")
    ok(f.touch.clientY == 1, "clientY is correct")
    ok(f.travel == 2, "travel should be 2")
  )
