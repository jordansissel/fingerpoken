#!/usr/bin/env ruby

require "rubygems"

# Hack to skip 'openssl' if we don't have it, since we don't use it.
# https://github.com/jordansissel/fingerpoken/issues/#issue/1
begin
  require "openssl"
rescue LoadError => e
  # Lie, and say we loaded "openssl"
  # 'thin' usees 'rack' which requires 'openssl' so it can compute hashes.
  # We don't need that feature, anyway.
  $LOADED_FEATURES << "openssl.rb"
end

require "em-websocket"
require "json"
require "rack"
require "sinatra/async"
require "ffi"

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

module Xdotool
  extend FFI::Library
  ffi_lib "libxdo.so"

  attach_function :xdo_new, [:string], :pointer
  attach_function :xdo_mousemove, [:pointer, :int, :int, :int], :int
  attach_function :xdo_mousemove_relative, [:pointer, :int, :int], :int
  attach_function :xdo_click, [:pointer, :long, :int], :int
  attach_function :xdo_mousedown, [:pointer, :long, :int], :int
  attach_function :xdo_mouseup, [:pointer, :long, :int], :int
  attach_function :xdo_type, [:pointer, :long, :string, :long], :int
  attach_function :xdo_keysequence, [:pointer, :long, :string, :long], :int
end

EventMachine::run do
  xdo = Xdotool.xdo_new(nil)
  #vnc = Net::VNC.new("sadness:0", :shared => true, :password => ENV["VNCPASS"])

  EventMachine::WebSocket.start(:host => "0.0.0.0", :port => 5001) do |ws|
    ws.onmessage do |message|
      request = JSON.parse(message)
      p request
      case request["action"]
        when "move"
          Xdotool.xdo_mousemove_relative(xdo, request["rel_x"], request["rel_y"])
        when "click"
          Xdotool.xdo_click(xdo, 0, request["button"])
        when "mousedown"
          Xdotool.xdo_mousedown(xdo, 0, request["button"])
        when "mouseup"
          Xdotool.xdo_mouseup(xdo, 0, request["button"])
        when "type"
          Xdotool.xdo_type(xdo, 0, request["string"], 12000)
        when "keypress"
          key = request["key"]
          if key.is_a?(String)
            if key.length == 1
              # Assume letter
              Xdotool.xdo_type(xdo, 0, key, 12000)
            else
              # Assume keysym
              Xdotool.xdo_keysequence(xdo, 0, key, 12000)
            end
          else
            # type printables, key others.
            if 32.upto(127).include?(key)
              Xdotool.xdo_type(xdo, 0, request["key"].chr, 12000)
            else
              case key
                when 8 
                  Xdotool.xdo_keysequence(xdo, 0, "BackSpace", 12000)
                when 13
                  Xdotool.xdo_keysequence(xdo, 0, "Return", 12000)
                else
                  puts "I don't know how to type web keycode '#{key}'"
                end # case key
            end # if 32.upto(127).include?(key)
          end # if key.is_a?String
      end # case request["action"]
    end # ws.onmessage
  end # WebSocket
  
  Rack::Handler::Thin.run(
    Rack::CommonLogger.new( \
        Rack::ShowExceptions.new( \
              FingerPoken.new)), :Port => 5000)
end # EventMachine::run
