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
  targets = []
  opts = OptionParser.new do |opts|
    opts.banner = "Usage: #{$0} [options]"

    opts.on("-t TARGET", "--target TARGET",
            "Target a url. Can be given multiple times to target multiple things.") do |url|
      target = URI.parse(url)
      case target.scheme
      when "xdo"
        require "fingerpoken/#{target.scheme}"
        targets << [:Xdo, {}]
      when "vnc"
        require "fingerpoken/#{target.scheme}"
        targets << [:VNC, {}]
      when "tivo"
        require "fingerpoken/#{target.scheme}"
        targets << [:Tivo, { :host => "192.168.0.134" }]
      end
    end
  end
  opts.parse(args)

  puts targets

  EventMachine::run do
    $:.unshift(File.dirname(__FILE__) + "/lib")
    channel = EventMachine::Channel.new

    targets.each do |klass, args|
      args.merge!({ :channel => channel })
      puts FingerPoken::Target.const_get(klass).new(args)
    end

    EventMachine::WebSocket.start(:host => "0.0.0.0", :port => 5001) do |ws|
      ws.onmessage do |message|
        request = JSON.parse(message)
        puts "Request: #{request.inspect}"
        channel.push(
          :request => request,
          :callback => proc { |message| ws.send(JSON.dump(message)) }
        )
      end # ws.onmessage
    end # WebSocket
    
    Rack::Handler::Thin.run(
      Rack::CommonLogger.new( \
          Rack::ShowExceptions.new( \
                FingerPoken.new)), :Port => 5000)
  end # EventMachine::run
end

exit(main(ARGV))
