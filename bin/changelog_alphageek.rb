#!/home/vbatts/opt/bin/ruby

#require 'fileutils'
require 'logger'
require 'tempfile'
require 'stringio'

require 'rubygems'
require 'slackware'
require 'slackware/changelog/rss'

#include FileUtils

$LOG = Logger.new(STDERR)
$LOG.level = Logger::WARN

FEEDS_BASE_DIR = "/home/vbatts/public_html/feeds/"
#url = 'http://alphageek.dyndns.org/linux/slackware-packages/slack-13.1/ChangeLog.txt'
# Sun Feb 13 08:44:35 PST 2011
# new url
URL = 'http://alphageek.dyndns.org/mirrors/alphageek/slackware-%s/ChangeLog.txt'

VERSIONS = %w{ 14.0 14.1 }

def url(ver)
  URL % ver
end

if ARGV.include?('-v')
  $LOG.level = Logger::DEBUG
end

VERSIONS.each {|ver|
  begin
    #tmp_file = File.open("/tmp/vbatts/alpha_log-#{(rand*1000).to_i}.xxx", "w+")
    tmp_file = Tempfile.new("alpha_log")
    $LOG.debug('tmp_file') { tmp_file }

    strio = StringIO.new()
    $LOG.debug('created ') { strio }

    buffer = `lynx -source #{url(ver)}`
    $LOG.debug('buffer length') { buffer.length }

    tmp_file.write(buffer)
    tmp_file.flush

    changelog = Slackware::ChangeLog.new(tmp_file.path)
    changelog.parse
    strio.write(changelog.to_rss(
      :noimage => true,
      :title => "alphageek's #{ver} ChangeLog",
      :url => url(ver)))
  ensure
    strio.seek(0)
    tmp_file.close
  end
    feed_file = File.open(FEEDS_BASE_DIR + "alphageek-#{ver}_ChangeLog.rss", "w+")
    $LOG.debug('feed_file') { feed_file }
    feed_file.write(strio.read())
    feed_file.close
}
