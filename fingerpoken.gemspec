Gem::Specification.new do |spec|
  files = []
  dirs = %w{lib public views samples test bin}
  dirs.each do |dir|
    files += Dir["#{dir}/**/*"]
  end

  #svnrev = %x{svn info}.split("\n").grep(/Revision:/).first.split(" ").last.to_i
  spec.name = "fingerpoken"
  spec.version = "0.3.0"
  spec.summary = "fingerpoken - turns your ipad/itouch/iphone into a remote touchpad, keyboard, etc"
  spec.description = "fingerpoken - turns your ipad/itouch/iphone into a remote touchpad, keyboard, etc"
  spec.add_dependency("ffi")
  spec.add_dependency("ruby-hmac")
  spec.add_dependency("eventmachine")
  spec.add_dependency("em-websocket")
  spec.add_dependency("async_sinatra")
  spec.add_dependency("json")
  spec.add_dependency("thin")
  spec.add_dependency("haml")
  spec.add_dependency("eventmachine-vnc")
  spec.files = files
  spec.require_paths << "lib"
  spec.bindir = "bin"
  spec.executables << "fingerpoken.rb"

  spec.author = "Jordan Sissel"
  spec.email = "jls@semicomplete.com"
  spec.homepage = "https://github.com/jordansissel/eventmachine-vnc"
end

