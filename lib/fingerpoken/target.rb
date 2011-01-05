require "logger"

class FingerPoken::Target
  def initialize(config)
    @channel = config[:channel]
    @logger = Logger.new(STDERR)
    @logger.level = ($DEBUG ? Logger::DEBUG: Logger::WARN)
  end

  def register
    if @registered
      @logger.warn("Ignoring extra call to #{self.class.name}#register. Trace:\n#{caller[0..3].join("\n")}")
      return
    end

    @registered = true
    @logger.debug(:register => self.class.name)
    @channel.subscribe do |obj|
      request = obj[:request]
      callback = obj[:callback]
      @logger.debug(:request => request)
      response = case request["action"]
        when "mousemove_relative"
          mousemove_relative(request["rel_x"], request["rel_y"])
        when "mousemove_absolute"
          mousemove_absolute(request["percent_x"], request["percent_y"])
        when "move_end"
          move_end()
        when "click"
          click(request["button"])
        when "mousedown"
          mousedown(request["button"])
        when "mouseup"
          mouseup(request["button"])
        when "type"
          type(request["string"])
        when "keypress"
          keypress(request["key"])
        else
          p ["Unsupported action", request]
      end

      if response.is_a?(Hash)
        callback.call(response)
      end
    end # @channel.subscribe
  end # def register

  # Subclasses should implement this.
  def mousemove_relative(x, y)
    @logger.info("mousemove not supported")
  end

  def mousedown(button)
    @logger.info("mousedown not supported")
  end

  def mouseup(button)
    @logger.info("mouseup not supported")
  end

  def click(button)
    mousedown(button)
    mouseup(button)
  end

  def type(string)
    @logger.info("typing not supported")
  end

  def keypress(key)
    @logger.info("keypress not supported")
  end

  def keydown(key)
    @logger.info("keydown not supported")
  end

  def keyup(key)
    @logger.info("keyup not supported")
  end

  def move_end()
    @logger.info("move_end not supported")
  end
end # class FingerPoken::Target
