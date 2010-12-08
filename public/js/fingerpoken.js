(function() {
  $(document).ready(function() {
    var state = {
      x: -1,
      y: -1,
      moving: false,
      dragging: false,
      width: window.innerWidth,
      height: window.innerHeight,
      key: undefined,
    }
    var status = $("#status");

    var connect = function(state) {
      status.html("connecting...");
      var websocket = new WebSocket("ws://" + document.location.hostname + ":5001");
      websocket.onopen = function(event) {
        status.html("websocket ready");
      }

      websocket.onclose = function(event) {
        status.html("Closed, trying to reopen.");
        setTimeout(function() {
          connect(state);
        }, 1000);
      }

      state.websocket = websocket;
    }


    connect(state);

    /* This will track orientation/motion changes with the accelerometer and
     * gyroscope. Not sure how useful this would be... */
    //$(window).bind("devicemotion", function(event) {
      //var e = event.originalEvent;
      //state.accel = e.accelerationIncludingGravity;

      /* Trim shakes */
      //if (Math.abs(state.accel.x) < 0.22 && Math.abs(state.accel.y) < 0.22) {
        //return;
      //}
      //status.html("Motion: \nx: " + state.accel.x + "\ny: " + state.accel.y + "\nz: " + state.accel.z);
      //state.websocket.send(JSON.stringify({
        //action: "move",
        //rel_x: Math.ceil(state.accel.x) * -1,
        //rel_y: Math.ceil(state.accel.y) * -1,
      //}));
    //});
    
    $("#menu").bind("touchmove", function(event) { 
      event.preventDefault();
    });

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
    }).bind("touchend", function(event) { /* $("#area").bind("touchend" ...  */
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
          var r = e.rotation;
          if (r < 0) {
            r += 360;
          }

          status.html(r);
          if (r > 75 && r < 105) {
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
            }).bind("change", function(event) {
              /* Skip empty changes */
              if (keyboard.val() == "") {
                return;
              }

              state.websocket.send(JSON.stringify({ 
                action: "type",
                string: keyboard.val(),
              }));

              /* Clear the field */
              keyboard.val("");
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

              console.log(key);
              var key = (e.keyCode ? e.keyCode : e.which);
              if (key >= 32 && key <= 127) {
                /* skip printable keys (a-z, etc) */
                return;
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
    }).bind("touchmove", function(event) { /* $("#area").bind("touchmove" ... */
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
      delta_x = (x - state.x);
      delta_y = (y - state.y);

      /* Apply acceleration */
      sign_x = (delta_x < 0 ? -1 : 1);
      sign_y = (delta_y < 0 ? -1 : 1);
      delta_x = Math.ceil(Math.pow(Math.abs(delta_x), 1.5) * sign_x);
      delta_y = Math.ceil(Math.pow(Math.abs(delta_y), 1.5) * sign_y);


      output += "Delta: " + delta_x + ", " + delta_y + "\n";
      status.html(output);

      state.x = x;
      state.y = y;

      if (e.rotation < -10 || e.rotation > 10) {
        /* Skip rotations that are probably not mouse-cursor-wanting movements */
        return;
      }
      if (e.scale < 0.9 || e.scale > 1.1) {
        /* Skip scales that are probably not mouse-cursor-wanting movements */
        return;
      }

      if (touches.length > 1 && !state.dragging) {
        /* Multifinger movement, probably should scroll? */
        if (delta_y < 0 || delta_y > 0) {
          /* Scroll */
          state.websocket.send(JSON.stringify({
            action: "click",
            button: (delta_y < 0) ? 4 : 5,
          }))
        }
        
      } else {
        state.websocket.send(JSON.stringify({
          action: "move",
          rel_x: delta_x,
          rel_y: delta_y
        }));
      }
    }); /*  $("#area").bind( ... )*/

    $("#leftarrow").bind("touchstart", function(event) {
      event.preventDefault();
      state.key = "Left";
    }).bind("touchmove", function(event) {
      event.preventDefault();
    }).bind("touchend", function(event) {
      event.preventDefault();
      if (state.key == "Left") {
        state.websocket.send(JSON.stringify({ 
          action: "keypress",
          key: "Left",
        }));
      }
    });

    $("#rightarrow").bind("touchstart", function(event) {
      event.preventDefault();
      state.key = "Right";
    }).bind("touchmove", function(event) {
      event.preventDefault();
    }).bind("touchend", function(event) {
      event.preventDefault();
      if (state.key == "Right") {
        state.websocket.send(JSON.stringify({ 
          action: "keypress",
          key: "Right",
        }));
      }
    });

    
  }); /* $(document).ready */
})();
