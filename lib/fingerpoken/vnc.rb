require "rubygems"
require "eventmachine-vnc"
require "fingerpoken/target"

class FingerPoken::Target::VNC < FingerPoken::Target
  attr_accessor :x
  attr_accessor :y
  attr_accessor :screen_x
  attr_accessor :screen_y
  attr_accessor :buttonmask

  def initialize(config)
    super(config)
    # TODO(sissel): eventmachine-vnc needs to suppore more auth mechanisms
    @user = config[:user]
    @password = (config[:password] or config[:user])
    @host = config[:host]
    @port = (config[:port] or 5900)
    @ready = false
    @recenter = config[:recenter]

    # For eventmachine-vnc
    ENV["VNCPASS"] = @password

    if @host == nil
      raise "#{self.class.name}: No host given to connect to"
    end

    @vnc = EventMachine::connect(@host, @port, VNCClient, self)
    @x = 0
    @y = 0
    @buttonmask = 0
  end

  def update_mouse
    if !@ready
      @logger.warn("VNC connection is not ready. Ignoring update.")
      return { "action" => "status", "status" => "VNC connection not ready, yet" }
    end
    @vnc.pointerevent(@x, @y, @buttonmask)

    # TODO(sissel): Hack to make it work in TF2.
    # Mouse movement is always "from center"
    # So after each move, center the cursor.
    if @recenter
      @x = (@vnc.screen_width / 2).to_i
      @y = (@vnc.screen_height / 2).to_i
    end
  end

  def ready
    @ready = true
    return { "action" => "status", "status" => "VNC READY!" }
  end

  def mousemove_relative(x, y)
    @x += x
    @y += y
    update_mouse
    return nil
  end

  def mousemove_absolute(px, py)
    # Edges may be hard to hit on some devices, so inflate things a bit.
    xbuf = @screen_x * 0.1
    ybuf = @screen_y * 0.1
    @x = (((@screen_x + xbuf) * px) - (xbuf / 2)).to_i
    @y = (((@screen_y + ybuf) * py) - (ybuf / 2)).to_i
    update_mouse
    return nil
  end

  def mousedown(button)
    button = (1 << (button.to_i - 1))
    return if @buttonmask & button != 0
    @buttonmask |= button
    update_mouse
    return nil
  end

  def mouseup(button)
    button = (1 << (button.to_i - 1))
    return if @buttonmask & button == 0
    @buttonmask &= (~button)
    update_mouse
    return nil
  end

  # TODO(sissel): Add keyboard support.
  # VNC uses the same keysym values as X11, so that's a win. We can likely
  # leverage xdo's char-to-keysym magic with VNC.
  def keypress(key)
    puts "Got key: #{key} (#{key.class})"
    if key.is_a?(String)
      if key.length == 1
        # Assume letter
        @vnc.keyevent(key.chr, true)
        @vnc.keyevent(key.chr, false)
      else
        # Assume keysym
        puts "I don't know how to type '#{key}'"
        return { :action => "status", :status => "I don't know how to type '#{key}'" }
      end
    else
      # type printables, key others.
      if 32.upto(127).include?(key)
        @vnc.keyevent(key, true)
        @vnc.keyevent(key, false)
      else
        case key
          when 8 
            @vnc.keyevent(0xff08, true)
            @vnc.keyevent(0xff08, false)
          when 13
            @vnc.keyevent(0xff0D, true)
            @vnc.keyevent(0xff0D, false)
          else
            puts "I don't know how to type web keycode '#{key}'"
            return { :action => "status", :status => "I don't know how to type '#{key}'" }
          end # case key
      end # if 32.upto(127).include?(key)
    end # if key.is_a?String
    return nil
  end # def keypress

  class VNCClient < EventMachine::Connection
    include EventMachine::Protocols::VNC::Client

    def initialize(target)
      @target = target
    end
    
    def ready
      @target.screen_x = @screen_width
      @target.screen_y = @screen_height
      @target.buttonmask = 0
      @target.x = (@screen_width / 2).to_i
      @target.y = (@screen_height / 2).to_i
      @target.ready()
    end
  end # class VNCClient
end
