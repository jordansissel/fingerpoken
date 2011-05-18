(function() {
  /* TODO(sissel): This could use some serious refactoring. */

    /* From http://alastairc.ac/2010/03/detecting-touch-based-browsing/ */
  function isTouchDevice() {
    var el = document.createElement('div');
    el.setAttribute('ongesturestart', 'return;');
    if (typeof el.ongesturestart == "function"){
      return true;
    } else {
      return false;
    }
  };
  var TOUCH_SUPPORTED = isTouchDevice;

  var Fingerpoken = function() {
    this.x = -1;
    this.y = -1;
    this.moving = false;
    this.dragging = false;
    this.width = window.innerWidth;
    this.height = window.innerHeight;
    this.key = undefined; /* TODO(sissel) = unused? */
    this.keyboard = false;
    this.touchpad_active = false;
    this.mouse = { };
    this.scroll = {
      y: 0
    };
    this.sequence = Math.floor(Math.random() * 10000);
  };

  Fingerpoken.prototype.send = function(obj) {
    if (this.config("fingerpoken/passphrase")) {
      this.sign(obj);
    }
    this.websocket.send(JSON.stringify(obj))
  };

  Fingerpoken.prototype.sign = function(obj) {
    var passphrase = this.config("fingerpoken/passphrase");
    obj.sequence = this.sequence;
    this.sequence += 1;

    console.log("Passphrase: " + passphrase);
    console.log("Sequence: " + this.sequence);
    obj.signature = Crypto.HMAC(Crypto.MD5, "" + obj.sequence, passphrase, { asBytes: true });
  };

  Fingerpoken.prototype.connect = function() {
    if (this.status === undefined) {
      this.status = $("#status");
    }
    this.status.html("connecting...");
    var websocket = new WebSocket("ws://" + document.location.hostname + ":5001");
    var self = this;
    websocket.onopen = function(event) {
      self.status.html("websocket ready");
    }

    websocket.onclose = function(event) {
      self.status.html("Closed, trying to reopen.");
      setTimeout(function() {
        self.connect();
      }, 1000);
    }

    websocket.onmessage = function(event) {
      var request = JSON.parse(event.data);
      fingerpoken.message_callback(request)
    }

    this.websocket = websocket;
  };

  Fingerpoken.prototype.config = function(key, value, default_value) {
    if (value) {
      this.status.html("config[" + key + "] = " + value);
      //alert(key + " => " + value);
      
      window.localStorage[key] = value
      return value
    } else {
      return window.localStorage[key] || default_value
    }
  };

  $(document).ready(function() {
    var status = $("#status");
    var keyboard = $('#keyboard');
    var keyboard_button = keyboard.prev('a');
    var fingerpoken = new Fingerpoken();
    fingerpoken.status = status;

    keyboard.width(keyboard_button.width());
    keyboard.height(keyboard_button.height());
    /* TODO(sissel): get the computed width (margin, padding, width) */
    keyboard.css('margin-left', '-' + keyboard_button.width() + 'px');
    keyboard.show();

    keyboard.bind("focus", function() {
      /* move the textarea away so we don't see the caret */
      keyboard.css('margin-left', '-10000px');
      fingerpoken.keyboard = true;
      $(window).triggerHandler("resize");
    });
    keyboard.bind("blur", function(){
      keyboard.css('margin-left', '-' + keyboard_button.width() + 'px');
      fingerpoken.keyboard = false;
      $(window).triggerHandler("resize");
    });
    keyboard.bind("keypress", function(event) {
      var e = event.originalEvent;
      var key = e.charCode;
      if (!key) {
        key = (e.keyCode ? e.keyCode : e.which);
      }
      fingerpoken.send(JSON.stringify({ 
        action: "log",
        shift: e.shiftKey,
        char: e.charCode,
        ctrl: e.ctrlKey,
        meta: e.ctrlKey,
      }));
      fingerpoken.send(JSON.stringify({ 
        action: "keypress",
        key: key,
        shift: e.shiftKey,
      }));

      /* Only prevent default if we're not backspace,
       * this lets 'backspace' do keyrepeat. */
      if (key != 8) { 
        event.preventDefault();
      }
    }).bind("change", function(event) {
      /* Skip empty changes */
      if (keyboard.val() == "") {
        return;
      }

      fingerpoken.send(JSON.stringify({ 
        action: "type",
        string: keyboard.val(),
      }));

      /* Clear the field */
      keyboard.val("");
    });

    keyboard.bind("keyup", function(event) {
      var e = event.originalEvent;
      fingerpoken.send(JSON.stringify({ 
        action: "log",
        shift: e.shiftKey,
        char: e.charCode,
        key: e.which,
        ctrl: e.ctrlKey,
        meta: e.ctrlKey,
      }));

      var key = (e.keyCode ? e.keyCode : e.which);
      if (key >= 32 && key <= 127) {
        /* skip printable keys (a-z, etc) */
        return;
      }

      fingerpoken.send(JSON.stringify({ 
        action: "keypress",
        key: key,
        shift: e.shiftKey,
      }));

      event.preventDefault();
    });

    fingerpoken.message_callback = function(request) {
      action = request["action"];
      switch (action) {
        case "status":
          /* Use eval to do unicode escaping */
          var status = eval("\"" + request["status"] + "\"");
          var el = $("<h1 class='status'>" + status + "</h1>");
          $("#area").empty().append(el);
          el.delay(500).fadeOut(500, function() { $(this).remove() });
          break;
      }
    };

    /* Sync configuration elements */

    /* Mouse movement */
    console.log(fingerpoken.config("fingerpoken/mouse/movement"));
    $("input[name = \"mouse-config\"]")
      .bind("change", function(event) {
        fingerpoken.config("fingerpoken/mouse/movement", event.target.value);
      }).filter("[value = \"" + fingerpoken.config("fingerpoken/mouse/movement") + "\"]")
      .attr("checked", "checked").click()
    
    $("input[name = \"mouse-acceleration\"]")
      .bind("change", function(event) {
        fingerpoken.config("fingerpoken/mouse/acceleration", parseInt(event.target.value));
      }).val(fingerpoken.config("fingerpoken/mouse/acceleration")).change();

    $("input[name = \"fingerpoken-passphrase\"]")
      .bind("change", function(event) {
        fingerpoken.config("fingerpoken/passphrase", event.target.value);
      }).val(fingerpoken.config("fingerpoken/passphrase")).change();

    /* Changing orientation sometimes leaves the viewport
     * not starting at 0,0. Fix it with this hack.
     * Also, we want to make the content size full height. */
    $(window).bind("orientationchange resize pageshow", function(event) {
      scroll(0, 0);

      var header = $(".header:visible");
      var footer = $(".footer:visible");
      var content = $(".content:visible");
      var viewport_height = $(window).height();

      var content_height = viewport_height - header.outerHeight() - footer.outerHeight();

      /* Trim margin/border/padding height */
      content_height -= (content.outerHeight() - content.height());

      /* TODO(sissel): Make this special handling only for iphones.
       * http://developer.apple.com/library/safari/#documentation/appleapplications/reference/safariwebcontent/UsingtheViewport/UsingtheViewport.html
       */
      if (fingerpoken.keyboard) {
        if (window.orientation == 90 || window.orientation == -90) {
          content_height -= 162; /* landscape orientation keyboard */
          content_height -= 32; /* "form assistant" aka FormFill, this height is undocumented. */
        } else {
          content_height -= 216; /* portrait orientation keyboard */
          content_height -= 44; /* "form assistant" aka FormFill */
        }
      }
      status.html("Resize / " + window.orientation + " / " + fingerpoken.keyboard + " / " + content_height);
      content.height(content_height);
    });


    fingerpoken.connect();

    /* This will track orientation/motion changes with the accelerometer and
     * gyroscope. Not sure how useful this would be... */
    //$(window).bind("devicemotion", function(event) {
      //var e = event.originalEvent;
      //fingerpoken.accel = e.accelerationIncludingGravity;

      /* Trim shakes */
      //if (Math.abs(fingerpoken.accel.x) < 0.22 && Math.abs(fingerpoken.accel.y) < 0.22) {
        //return;
      //}
      //status.html("Motion: \nx: " + fingerpoken.accel.x + "\ny: " + fingerpoken.accel.y + "\nz: " + fingerpoken.accel.z);
      //fingerpoken.send({
        //action: "move",
        //rel_x: Math.ceil(fingerpoken.accel.x) * -1,
        //rel_y: Math.ceil(fingerpoken.accel.y) * -1,
      //});
    //});
    

    /* TODO(sissel): add mousedown/mousemove/mouseup support */
    $("#area").bind("touchstart mousedown", function(event) {
      var e = event.originalEvent;
      fingerpoken.touchpad_active = true;
      /* if no 'touches', use the event itself, one finger/mouse */
      var touches = e.touches || [ e ];
      var output = "Start: " + touches[0].clientX + "," + touches[0].clientY + "\n";
      output += "Fingers: " + touches.length + "\n";
      status.html(output);

      /* number of fingers == mouse button */
      fingerpoken.fingers = touches.length;
      switch (fingerpoken.fingers) {
        case 1: fingerpoken.button = 1; break;
        case 2: fingerpoken.button = 3; break;
        case 3: fingerpoken.button = 2; break;
      }

      var now = (new Date()).getTime();
      if ((now - fingerpoken.last_click) < 170) {
        /* Start dragging */
        fingerpoken.send({
          action: "mousedown",
          button: fingerpoken.button,
        })
        fingerpoken.dragging = true;
      }
      event.preventDefault();
    }).bind("touchend mouseup", function(event) { /* $("#touchpadsurface").bind("touchend" ...  */
      var e = event.originalEvent;
      var touches = e.touches || [ e ];

      if (fingerpoken.mouse.vectorTimer) {
        clearInterval(fingerpoken.mouse.vectorTimer);
        fingerpoken.mouse.vectorTimer = null;
      }

      if (fingerpoken.dragging) {
        fingerpoken.send({
          action: "mouseup",
          button: fingerpoken.button,
        });
        fingerpoken.dragging = false;
      } else {
        if (fingerpoken.moving && !fingerpoken.scrolling) {
          var e = fingerpoken.last_move;
          var r = e.rotation;
          if (r < 0) {
            r += 360;
          }

          status.html(r);
        } else if (fingerpoken.scrolling) {
          /* nothing for now */
        } else {
          /* No movement, click! */
          status.html("Click!");
          console.log("click");
          fingerpoken.send({
            action: "click",
            button: fingerpoken.button,
          });
          fingerpoken.last_click = (new Date()).getTime();
        }
      }
      if (touches.length == 0 || !e.touches) {
        fingerpoken.moving = false;
        fingerpoken.scrolling = false;
        fingerpoken.touchpad_active = false;
      }
      event.preventDefault();
    }).bind("touchmove mousemove", function(event) { /* $("#touchpadsurface").bind("touchmove" ... */
      var e = event.originalEvent;
      var touches = e.touches || [ e ];

      //if (!fingerpoken.touchpad_active) {
        //event.preventDefault();
        //return;
      //}

      if (!fingerpoken.moving) {
        /* Start calculating delta offsets now */
        fingerpoken.moving = true;
        fingerpoken.start_x = fingerpoken.x = touches[0].clientX;
        fingerpoken.start_y = fingerpoken.y = touches[0].clientY;
        /* Skip this event */
        return;
      }

      fingerpoken.last_move = e;

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

      var x = touches[0].clientX;
      var y = touches[0].clientY;
      var delta_x = (x - fingerpoken.x);
      var delta_y = (y - fingerpoken.y);

      /* Apply acceleration */
      var sign_x = (delta_x < 0 ? -1 : 1);
      var sign_y = (delta_y < 0 ? -1 : 1);

      /* jQuery Mobile or HTML 'range' inputs don't support floating point.
       * Hack around it by using larger numbers and compensating. */
      var accel = fingerpoken.config("fingerpoken/mouse/acceleration", null, 150) / 100.0;
      output += "Accel: " + accel + "\n";

      var delta_x = Math.ceil(Math.pow(Math.abs(delta_x), accel) * sign_x);
      var delta_y = Math.ceil(Math.pow(Math.abs(delta_y), accel) * sign_y);
      output += "Delta: " + delta_x + ", " + delta_y + "\n";

      fingerpoken.x = x;
      fingerpoken.y = y;

      /* TODO(sissel): Make this a config option */
      if (e.rotation < -10 || e.rotation > 10) {
        /* Skip rotations that are probably not mouse-cursor-wanting movements */
        return;
      }
      /* TODO(sissel): Make this a config option */
      if (e.scale < 0.9 || e.scale > 1.1) {
        /* Skip scales that are probably not mouse-cursor-wanting movements */
        return;
      }

      if (touches.length > 1 && !fingerpoken.dragging) {
        /* Multifinger movement, probably should scroll? */
        if (Math.abs(delta_y) > 0) {
          /* Scroll */
          fingerpoken.scroll.y += delta_y;

          /* Don't scroll every time we move, wait until we move enough
           * that it is more than 10 pixels. */
          /* TODO(sissel): Make this a config option */
          if (Math.abs(fingerpoken.scroll.y) > 10) {
            fingerpoken.scrolling = true;
            fingerpoken.moving = false;
            /* TODO(sissel): Support horizontal scrolling (buttons 6 and 7) */
            fingerpoken.scroll.y = 0;
            fingerpoken.send({
              action: "click",
              button: (delta_y < 0) ? 4 : 5,
            });
          }
        } /* if (Math.abs(delta_y) > 0) */
      } else {
        /* Only 1 finger, and we aren't dragging. So let's move! */
        /* TODO(sissel): Refactor these in to fumctions */
        var movement = fingerpoken.config("fingerpoken/mouse/movement");
        if (movement == "relative") {
          fingerpoken.send({
            action: "mousemove_relative",
            rel_x: delta_x,
            rel_y: delta_y
          });
        } else if (movement == "absolute") {
          /* Send absolute in terms of percentages. */
          var content = $(".content:visible");
          fingerpoken.send({
            action: "mousemove_absolute",
            percent_x: x / content.innerWidth(),
            percent_y: y / content.innerHeight(),
          });
        } else if (movement == "vector") {
          if (!fingerpoken.mouse.vectorTimer) {
            fingerpoken.mouse.vectorTimer = setInterval(function() {
              var rx = fingerpoken.x - fingerpoken.start_x;
              var ry = fingerpoken.y - fingerpoken.start_y;
              if (rx == 0 || ry == 0) {
                return;
              }

              var sign_rx = (rx < 0 ? -1 : 1);
              var sign_ry = (ry < 0 ? -1 : 1);
              var vector_accel = accel / 1.7 /* feels like the right ratio */

              output = "rx,ry = " + rx + ", " + ry + "\n";
              rx = Math.ceil(Math.pow(Math.abs(rx), vector_accel) * sign_rx);
              ry = Math.ceil(Math.pow(Math.abs(ry), vector_accel) * sign_ry);
              output += "rx2,ry2 = " + rx + ", " + ry + "\n"; 

              fingerpoken.send({
                action: "mousemove_relative",
                rel_x: rx,
                rel_y: ry
              });
            }, 15);
          } /* if (!fingerpoken.mouse.vectorTimer) */
        } /* mouse vector movement */
        status.html(output)
      } /* finger movement */
    }); /*  $("#touchpadsurface").bind( ... )*/


    /* Take commands like this:
     * 
     * Key press:
     * <a class="command" data-action="keypress" data-key="key to press">
     *
     * Mouse click
     * <a class="command" data-action="click" data-button="button to click">
     */

    var cmd = $("a.command")
    cmd.bind("touchstart", function(event) {
      fingerpoken.touchelement = this;
    }).bind("touchmove", function(event) {
      event.preventDefault();
    }).bind("touchend", function(event) {
      if (fingerpoken.touchelement == this) {
        fingerpoken.send({
          action: $(this).attr("data-action"),
          key: $(this).attr("data-key"),
          button: parseInt($(this).attr("data-button")),
        });
      }
    });

    if (!TOUCH_SUPPORTED) {
      cmd.bind("mousedown", function(event) {
        fingerpoken.touchelement = this;
        event.preventDefault();
      }).bind("mousemove", function(event) {
        event.preventDefault();
      }).bind("touchend", function(event) {
        if (fingerpoken.touchelement == this) {
          fingerpoken.send({
            action: $(this).attr("data-action"),
            key: $(this).attr("data-key"),
            button: parseInt($(this).attr("data-button")),
          });
        }
      });
    } /* if !TOUCH_SUPPORTED */

  }); /* $(document).ready */
})();
