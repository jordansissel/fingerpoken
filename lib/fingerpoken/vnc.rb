require "rubygems"
require "eventmachine-vnc"
require "fingerpoken/target"

class FingerPoken::Target::VNC < FingerPoken::Target
  attr_accessor :x
  attr_accessor :y
  attr_accessor :buttonmask

  def initialize(config)
    super(config)
    # TODO(sissel): Make this configurable
    @password = ENV["VNCPASS"]
    @host = "sadness"
    @port = 5900

    @vnc = EventMachine::connect(@host, @port, VNCClient, self)
    @x = 0
    @y = 0
    @buttonmask = 0
  end

  def update
    p [@x, @y, @buttonmask]
    @vnc.pointerevent(@x, @y, @buttonmask)

    # TODO(sissel): Hack to make it work in TF2.
    # Mouse movement is always "from center"
    # So after each move, center the cursor.
    if ENV["RECENTER"]
      @x = (@vnc.screen_width / 2).to_i
      @y = (@vnc.screen_height / 2).to_i
    end
  end

  def mousemove_relative(x, y)
    @x += x
    @y += y
    update
  end

  def mousedown(button)
    button = (1 << (button.to_i - 1))
    return if @buttonmask & button != 0
    @buttonmask |= button
    update
  end

  def mouseup(button)
    button = (1 << (button.to_i - 1))
    return if @buttonmask & button == 0
    @buttonmask &= (~button)
    update
  end

  class VNCClient < EventMachine::Connection
    include EventMachine::Protocols::VNC::Client

    def initialize(target)
      @target = target
    end
    
    def ready
      @target.register
      @target.x = (@screen_width / 2).to_i
      @target.y = (@screen_height / 2).to_i
      @target.buttonmask = 0
    end
  end # class VNCClient
end
