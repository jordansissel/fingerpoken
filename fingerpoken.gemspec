Gem::Specification.new do |spec|
  files = []
  dirs = %w{lib samples test bin}
  dirs.each do |dir|
    files += Dir["#{dir}/**/*"]
  end

  #svnrev = %x{svn info}.split("\n").grep(/Revision:/).first.split(" ").last.to_i
  rev = Time.now.strftime("%Y%m%d%H%M%S")
  spec.name = "fingerpoken"
  spec.version = "0.2.#{rev}"
  spec.summary = "fingerpoken - turns your ipad/itouch/iphone into a remote touchpad, keyboard, etc"
  spec.description = "fingerpoken - turns your ipad/itouch/iphone into a remote touchpad, keyboard, etc"
  spec.add_dependency("eventmachine-vnc")
  spec.add_dependency("ffi")
  spec.add_dependency("eventmachine")
  spec.files = files
  spec.require_paths << "lib"
  spec.bindir = "bin"
  spec.executables << "fingerpoken.rb"

  spec.author = "Jordan Sissel"
  spec.email = "jls@semicomplete.com"
  spec.homepage = "https://github.com/jordansissel/eventmachine-vnc"
end

