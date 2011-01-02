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

  def update
    @vnc.pointerevent(@x, @y, @buttonmask)

    # TODO(sissel): Hack to make it work in TF2.
    # Mouse movement is always "from center"
    # So after each move, center the cursor.
    if @recenter
      @x = (@vnc.screen_width / 2).to_i
      @y = (@vnc.screen_height / 2).to_i
    end
  end

  def mousemove_relative(x, y)
    @x += x
    @y += y
    update
    return nil
  end

  def mousemove_absolute(px, py)
    # Edges may be hard to hit on some devices, so inflate things a bit.
    xbuf = @screen_x * 0.1
    ybuf = @screen_y * 0.1
    @x = (((@screen_x + xbuf) * px) - (xbuf / 2)).to_i
    @y = (((@screen_y + ybuf) * py) - (ybuf / 2)).to_i
    update
    return nil
  end

  def mousedown(button)
    button = (1 << (button.to_i - 1))
    return if @buttonmask & button != 0
    @buttonmask |= button
    update
    return nil
  end

  def mouseup(button)
    button = (1 << (button.to_i - 1))
    return if @buttonmask & button == 0
    @buttonmask &= (~button)
    update
    return nil
  end

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
      @target.register
    end
  end # class VNCClient
end
