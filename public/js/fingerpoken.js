(function() {
  $(document).ready(function() {
    var state = {
      x: -1,
      y: -1,
      moving: false,
      dragging: false,
      width: window.innerWidth,
      height: window.innerHeight,
    }

    var connect = function(state) {
      var websocket = new WebSocket("ws://" + document.location.hostname + ":5001");
      websocket.onopen = function(event) {
        console.log("websocket ready");
      }

      websocket.onclose = function(event) {
        status.html("Closed, trying to reopen.");
        setTimeout(1000, function() {
          connect(state);
        });
      }

      state.websocket = websocket;
    }

    var status = $("#status");
    status.html("connecting...");

    connect(state);

    $("#area").bind("touchstart", function(event) {
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
        state.websocket.send(JSON.stringify({
          action: "mousedown",
          button: state.button,
        }))
        state.dragging = true;
      }
    });

    $("#area").bind("touchend", function(event) {
      var e = event.originalEvent;
      var touches = e.touches;
      if (state.dragging) {
        state.websocket.send(JSON.stringify({
          action: "mouseup",
          button: state.button,
        }));
        state.dragging = false;
      } else {
        if (state.moving) {
          var e = state.last_move;
          if (e.rotation < 0) {
            e.rotation += 360;
          }

          status.html(e.rotation);
          if (e.rotation > 80 && e.rotation < 100) {
            /* Activate the keyboard when there's a 90-dgree rotation*/
            var keyboard = $("<textarea id='keyboard' rows='10'></textarea>");
            keyboard.css("width", "100%");
            keyboard.css("height", "100%");
            status.html("");
            keyboard.appendTo(status).focus();
            keyboard.bind("keypress", function(event) {
              var e = event.originalEvent;
              var key = e.charCode;
              console.log(key);
              if (!key) {
                key = (e.keyCode ? e.keyCode : e.which);
              }
              state.websocket.send(JSON.stringify({ 
                action: "log",
                shift: e.shiftKey,
                char: e.charCode,
                ctrl: e.ctrlKey,
                meta: e.ctrlKey,
              }));
              state.websocket.send(JSON.stringify({ 
                action: "keypress",
                key: key,
                shift: e.shiftKey,
              }));

              e.preventDefault();
            });

            keyboard.bind("keyup", function(event) {
              var e = event.originalEvent;
              state.websocket.send(JSON.stringify({ 
                action: "log",
                shift: e.shiftKey,
                char: e.charCode,
                key: e.which,
                ctrl: e.ctrlKey,
                meta: e.ctrlKey,
              }));
              if (e.charCode != 0) {
                /* non-symbol keys zero-charcode */
                return;
              }
              console.log(key);
              if (!key) {
                key = (e.keyCode ? e.keyCode : e.which);
              }
              state.websocket.send(JSON.stringify({ 
                action: "keypress",
                key: key,
                shift: e.shiftKey,
              }));

              e.preventDefault();
            });
            //status.html("<textarea id='keyboard'></textarea>");
            //$("#keyboard").focus();
          }
        } else {
          /* No movement, click! */
          status.html("Click!");
          state.websocket.send(JSON.stringify({ 
            action: "click",
            button: state.button,
          }));
          state.last_click = (new Date()).getTime();
        }
      }
      state.moving = false;
      event.preventDefault();
    });

    $("#area").bind("touchmove", function(event) {
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

      var r = e.rotation;
      if (r < 0) {
        r += 360;
      }
      output += "rotation: " + r + "\n";
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
      state.websocket.send(JSON.stringify({ 
        action: "move",
        rel_x: delta_x,
        rel_y: delta_y
      }));
    });
  });
})();
