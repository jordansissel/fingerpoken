task :default => [:package]

task :test do
  system("cd test; ruby alltests.rb")
end

task :package => [:test, :package_real]  do
end

task :package_real do
  system("gem build fingerpoken.gemspec")
end

task :publish do
  latest_gem = %x{ls -t fingerpoken*.gem}.split("\n").first
  system("gem push #{latest_gem}")
end
