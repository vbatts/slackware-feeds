#!/home/vbatts/opt/bin/ruby

require 'logger'

$log = Logger.new(STDERR)
$log.level = Logger::DEBUG

# put this in a loader function, because the 
# rss library is SOO SLOW to load. we don't want to load it, 
# if the script is going to fail early.
def load_libs()
  require 'rubygems'
  require 'slackware'
  require 'slackware/changelog/rss'
  require 'rb-inotify'
end


BASE_URL  = "http://slackware.osuosl.org/"
FEEDS_BASE_DIR  = "/home/vbatts/public_html/feeds/"
RE_REPO_NAME  = Regexp.new(/slackware(\d{2})?-(\d+\.\d+|current)\/(patches)?\/?.*/)

def gen_file(file)
  m = RE_REPO_NAME.match file
  if m[3].nil?
    file_name = "%sslackware%s-%s_ChangeLog.rss" % [FEEDS_BASE_DIR, m[1], m[2]]
  else
    file_name = "%sslackware%s-%s_%s_ChangeLog.rss" % [FEEDS_BASE_DIR, m[1], m[2], m[3]]
  end

  if File.exist?(file_name)
    if File.mtime(file) < File.mtime(file_name)
      printf("%f: INFO: %s is newer than %s\n", Time.now, file, file_name)
    end
  end

  changelog = Slackware::ChangeLog.new(file) #, :version => m[2])
  opts = Hash.new
  opts[:arch] = m[1] unless m[1].nil?
  if m[3].nil?
    opts[:url] = "%sslackware%s-%s/ChangeLog.txt" % [BASE_URL, m[1], m[2]]
    feed = File.open( "%sslackware%s-%s_ChangeLog.rss" % [FEEDS_BASE_DIR, m[1], m[2]], "w+")
  else
    opts[:url] = "%sslackware%s-%s/%s/ChangeLog.txt" % [BASE_URL, m[1], m[2], m[3]]
    feed = File.open( "%sslackware%s-%s_%s_ChangeLog.rss" % [FEEDS_BASE_DIR, m[1], m[2], m[3]], "w+")
  end
  changelog.parse
  printf("%f: INFO: generating feed: %s\n", Time.now.to_f, feed.path)
  feed << changelog.to_rss(opts)
  feed.close
  changelog = nil
end

if ARGV.count == 0
  $log.error("#{Time.now}: ERROR: ChangeLog.txt files must be passed\n")
  exit(2)
else
  load_libs()
  for file in ARGV
    if File.exist?(file)
      gen_file(file)
    else
      $log.warn("#{Time.now}: WARN: #{file} does not exist\n")
    end
  end
end

# vim: set sts=2 sw=2 et ai:
