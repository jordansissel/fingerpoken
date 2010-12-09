#!/usr/bin/env ruby

require "rubygems"
require "em-websocket"
require "json"
require "rack"
require "sinatra/async"
require "optparse"

class FingerPoken < Sinatra::Base
  register Sinatra::Async
  set :haml, :format => :html5
  set :logging, true
  set :public, "#{File.dirname(__FILE__)}/../public"
  set :views, "#{File.dirname(__FILE__)}/../views"

  aget '/' do
    headers "Content-Type" => "text/html"
    body haml :index
  end # GET /

  aget '/style.css' do
    headers "Content-Type" => "text/css; charset=utf8"
    body sass :style
  end # GET /style.css
end

def main(args)
  opts = OptionParser.new do |opts|
  end
  EventMachine::run do
    $:.unshift(File.dirname(__FILE__) + "/lib")
    channel = EventMachine::Channel.new

    # TODO(sissel): Pick up here and make command flags to choose the 
    # target (vnc, xdo, etc)
    #require "fingerpoken/xdo"
    #target = FingerPoken::Target::Xdo.new :channel => channel
    
    require "fingerpoken/vnc"
    target = FingerPoken::Target::VNC.new :channel => channel

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
end

exit(main(ARGV))
