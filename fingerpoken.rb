#!/usr/bin/env ruby

require "rubygems"
require "em-websocket"
require "json"
require "rack"
require "sinatra/async"
require "ffi"

class FingerPoken < Sinatra::Base
  register Sinatra::Async

  aget '/' do
    headers "Content-Type" => "text/html"
    body haml :index
  end
end

module Xdotool
  extend FFI::Library
  ffi_lib "libxdo.so"

  attach_function :xdo_new, [:string], :pointer
  attach_function :xdo_mousemove, [:pointer, :int, :int, :int], :int
  attach_function :xdo_mousemove_relative, [:pointer, :int, :int], :int
  attach_function :xdo_click, [:pointer, :long, :int], :int
end

EventMachine::run do
  xdo = Xdotool.xdo_new(nil)

  EventMachine::WebSocket.start(:host => "0.0.0.0", :port => 5001) do |ws|
    ws.onmessage do |message|
      request = JSON.parse(message)
      case request["action"]
        when "move"
          Xdotool.xdo_mousemove_relative(xdo, request["rel_x"], request["rel_y"])
        when "click"
          Xdotool.xdo_click(xdo, 0, request["button"]);
      end
    end # ws.onmessage
  end # WebSocket
  
  Rack::Handler::Thin.run(
    Rack::CommonLogger.new( \
        Rack::ShowExceptions.new( \
              FingerPoken.new)), :Port => 5000)
end # EventMachine::run
