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

class VNCClient < EventMachine::Connection
  include EventMachine::Protocols::VNC::Client

  def initialize(channel)
    @channel = channel

  end

  def mousemove(x, y)
    message = [ POINTER_EVENT, 0, x, y ].pack("CCnn")
    send_data(message)
  end

  def ready
    puts "vnc ready"
    @x = (@screen_width / 2).to_i
    @y = (@screen_height / 2).to_i

    @channel.subscribe do |request|
      rel_x = request["rel_x"]
      rel_y = request["rel_y"]
      @x += rel_x
      @y += rel_y
      mousemove(@x, @y)
    end
  end
end

EventMachine::run do
  channel = EventMachine::Channel.new
  EventMachine::connect("sadness", 5900, VNCClient, channel)

  EventMachine::WebSocket.start(:host => "0.0.0.0", :port => 5001) do |ws|
    ws.onmessage do |message|
      request = JSON.parse(message)
      p request
      case request["action"]
        when "move"
          #Xdotool.xdo_mousemove_relative(xdo, request["rel_x"], request["rel_y"])
          channel.push(request)
      end # case request["action"]
    end # ws.onmessage
  end # WebSocket
  
  Rack::Handler::Thin.run(
    Rack::CommonLogger.new( \
        Rack::ShowExceptions.new( \
              FingerPoken.new)), :Port => 5000)
end # EventMachine::run
