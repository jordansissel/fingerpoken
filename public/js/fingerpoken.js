(function() {
  $(document).ready(function() {
    var status = $("#status");
    status.html("Ready!");

    var state = { 
      x: -1,
      y: -1,
      moving: false,
      dragging: false,
      width: window.innerWidth,
      height: window.innerHeight,
    }

    var websocket = new WebSocket("ws://" + document.location.hostname + ":5001");
    websocket.onopen = function(event) {
      console.log("websocket ready");
    }

    $(document).bind("touchstart", function(event) {
      event.preventDefault();
      var e = event.originalEvent;
      var touches = e.touches;
      var output = "Start: " + touches[0].clientX + "," + touches[0].clientY + "\n";
      output += "Fingers: " + touches.length + "\n";
      status.html(output);

      /* number of fingers == mouse button */
      state.fingers = touches.length;
      switch (state.fingers) {
        case 1: state.button = 1; break;
        case 2: state.button = 3; break;
        case 3: state.button = 2; break;
      }

      var now = (new Date()).getTime();
      if ((now - state.last_click) < 170) {
        /* Start dragging */
        websocket.send(JSON.stringify({
          action: "mousedown",
          button: state.button,
        }))
        state.dragging = true;
      }
    });

    $(document).bind("touchend", function(event) {
      var e = event.originalEvent;
      var touches = e.touches;
      if (state.dragging) {
        websocket.send(JSON.stringify({
          action: "mouseup",
          button: state.button,
        }));
        state.dragging = false;
      } else {
        if (state.moving) {
          var e = state.last_move;
          status.html(e.rotation);
          if (e.rotation > 80 && e.rotation < 100) {
            /* Activate the keyboard */
            var keyboard = $("<textarea id='keyboard'></textarea>");
            status.html("");
            keyboard.appendTo(status).focus();
            keyboard.bind("keyup", function(event) {
              var e = event.originalEvent;
              var code = (e.keyCode ? e.keyCode : e.which);
              websocket.send(JSON.stringify({ 
                action: "keypress",
                key: code
              }));

              e.preventDefault();
            });
            //status.html("<textarea id='keyboard'></textarea>");
            //$("#keyboard").focus();
          }
        } else {
          /* No movement, click! */
          status.html("Click!");
          websocket.send(JSON.stringify({ 
            action: "click",
            button: state.button,
          }));
          state.last_click = (new Date()).getTime();
        }
      }
      state.moving = false;
      event.preventDefault();
    });

    $(document).bind("touchmove", function(event) {
      var e = event.originalEvent;
      var touches = e.touches;
      event.preventDefault();
      if (!state.moving) {
        /* Start calculating delta offsets now */
        state.moving = true;
        state.x = touches[0].clientX;
        state.y = touches[0].clientY;
        /* Skip this event */
        return;
      }

      state.last_move = e;

      var output = "";
      for (var i in touches) {
        output += i + ": " + touches[i].clientX + "," + touches[i].clientY + "\n";
      }
      output += "rotation: " + e.rotation + "\n";
      output += "scale: " + e.scale + "\n";

      x = touches[0].clientX;
      y = touches[0].clientY;
      delta_x = (x - state.x) * 3;
      delta_y = (y - state.y) * 3;
      output += delta_x + ", " + delta_y + "\n";
      status.html(output);

      state.x = x;
      state.y = y;

      if (e.rotation < -10 || e.rotation > 10) {
        /* Skip rotations that are probably not mouse-cursor-wanting movements */
        return;
      }
      websocket.send(JSON.stringify({ 
        action: "move",
        rel_x: delta_x,
        rel_y: delta_y
      }));
    });
  });
})();
