# fingerpoken - use your iphone/ipad as a mouse and remote

![fingerpoken main screen](http://farm6.static.flickr.com/5209/5315764134_fd91969a2c_m.jpg "fingerpoken main touch screen")

## Video demos: 

* Mouse control: [http://www.youtube.com/watch?v=39PtZoxW_fM](http://www.youtube.com/watch?v=39PtZoxW_fM)
* TiVo control: [http://www.youtube.com/watch?v=2GUkVDrAFbY](http://www.youtube.com/watch?v=2GUkVDrAFbY)

## What works:

* dragging a finger moves the mouse.
* tap clicks, 2-finger tap right-clicks, 3 finger-tap middle-clicks
* two-finger drag up/down scrolls
* double tap clicks twice, etc.
* tap-drag works like you might expect (select text, etc).
* arbitrary keyboard input (press the 'keyboard' button)

## Supported targets:

* X11 (via libxdo)
* VNC (via rubygem eventmachine-vnc)
* TiVo (no external dependencies)

## Configuring the UI:

* The 'config' button lets you change a few things:
  * 'mouse movement' can be relative (normal touchpad), absolute positioning,
    and vector (start-point + current-point == mouse direction and speed)
  * You can change the mouse sensitivity.

Options set in the UI will persist across sessions using HTML5 localStorage.

## What's planned:

* security (https, authentication)
* many other things.
* Better PC-as-a-client support (mousemovement, clicking, etc)
* Got suggestions? 

## TODO:

* special keystroke input (control/shift/alt, function keys, page up/down, etc)

## What you need to run it:

client:
  * client: an iphone or ipad running iOS >=4.2 (requires websocket support in safari)

server:
  * ruby
  * rubygems: em-websocket, eventmachine, ffi, async_sinatra, json, rack
  * For the xdo target: libxdo (from the xdotool project)
  * For the vnc target: rubygem eventmachine-vnc
  * For the tivo target: nothing.
  * Linux, and X server. OS X and Windows support is probably easy.

## Run it:

  * xdo (X11): fingerpoken.rb -t xdo:///
  * vnc: fingerpoken.rb -t vnc:///password@host:port/
    * password is optional
    * port defaults to 5900
  * tivo: fingerpoken.rb -t tivo://yourtivoip:port
    * default port is 31339

Notes:

  * TiVo control requires the mouse mode be 'vector' (tap on 'config' to change this)
  * TiVo support also requires you enable the network remote control on your TiVo.


1) Run fingerpoken 
2) Point your iphone browser at http://yourmachine:5000/
3) Use your phone as a touchpad.

Optional:
* Bookmark to home screen. Works from there, too.
