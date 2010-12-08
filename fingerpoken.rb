#!/usr/bin/env ruby

require "rubygems"
require "em-websocket"
require "json"
require "rack"
require "sinatra/async"
require "eventmachine-vnc"

class FingerPoken < Sinatra::Base
  register Sinatra::Async
  set :haml, :format => :html5
  set :logging, true
  set :public, "#{File.dirname(__FILE__)}/public"
  set :views, "#{File.dirname(__FILE__)}/views"

  aget '/' do
    headers "Content-Type" => "text/html"
    body haml :index
  end # GET /

  aget '/style.css' do
    headers "Content-Type" => "text/css; charset=utf8"
    body sass :style
  end # GET /style.css
end

class Mouse
  attr_accessor :x
  attr_accessor :y
  attr_accessor :buttonmask
end

class VNCClient < EventMachine::Connection
  include EventMachine::Protocols::VNC::Client

  def initialize(channel)
    @channel = channel
    @mouse = Mouse.new
  end

  def ready
    puts "vnc ready"
    # Start in the center of the screen.
    # VNC doesn't have a 'pointer motion relative' feature, so
    # we have to track relative motion internally.
    @mouse.x = (@screen_width / 2).to_i
    @mouse.y = (@screen_height / 2).to_i
    @mouse.buttonmask = 0

    @channel.subscribe do |request|
      handle(request)
    end
  end

  def handle(request)
    p request["action"]

    button = (1 << request["button"])

    case request["action"]
    when "move"
      rel_x = request["rel_x"]
      rel_y = request["rel_y"]
      @mouse.x += rel_x
      @mouse.y += rel_y
      pointerevent(@mouse.x, @mouse.y, @mouse.buttonmask)
    when "click"
      @mouse.buttonmask |= button
      pointerevent(@mouse.x, @mouse.y, @mouse.buttonmask)
      # Bad form to sleep, but we want to block the EM reactor until we finish
      # clicking.
      #sleep(0.020)
      @mouse.buttonmask &= ~(button)
      pointerevent(@mouse.x, @mouse.y, @mouse.buttonmask)
    when "mousedown"
      if @mouse.buttonmask & button != 0
        puts "already down"
        return
      end
      @mouse.buttonmask |= button
      pointerevent(@mouse.x, @mouse.y, @mouse.buttonmask)
    when "mouseup"
      return if @mouse.buttonmask & button == 0
      @mouse.buttonmask &= ~(button)
      pointerevent(@mouse.x, @mouse.y, @mouse.buttonmask)
    else
      puts "Message type '#{request["action"]}' not supported"
    end
    puts "Mouse button: #{@mouse.buttonmask}"
    @mouse.x = (@screen_width / 2).to_i
    @mouse.y = (@screen_height / 2).to_i
  end
end

EventMachine::run do
  channel = EventMachine::Channel.new
  EventMachine::connect("sadness", 5900, VNCClient, channel)

  EventMachine::WebSocket.start(:host => "0.0.0.0", :port => 5001) do |ws|
    ws.onmessage do |message|
      request = JSON.parse(message)
      p request
      channel.push(request)
    end # ws.onmessage
  end # WebSocket
  
  Rack::Handler::Thin.run(
    Rack::CommonLogger.new( \
        Rack::ShowExceptions.new( \
              FingerPoken.new)), :Port => 5000)
end # EventMachine::run
