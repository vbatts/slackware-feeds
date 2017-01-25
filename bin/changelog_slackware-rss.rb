#!/usr/bin/env ruby
# Sun Jan 23 11:30:53 PST 2011
# Created by vbatts, vbatts@hashbangbash.com

$PROGRAM_NAME = File.basename(__FILE__)

require 'find'

require 'rubygems'
require 'ruby-prof'
require 'slackware'
require 'slackware/changelog/rss'
require 'rb-inotify'


BASE_URL	= "http://slackware.osuosl.org/"
MIRROR_BASE_DIR = "/mirrors/ftp.slackware.com/pub/slackware/"
FEEDS_BASE_DIR	= "/home/vbatts/public_html/feeds/"
RE_REPO_NAME	= Regexp.new(/slackware(\d{2})?-(\d+\.\d+|current)\/(patches)?\/?.*/)

def generate_new_if_none
	files = []

	Find.find(MIRROR_BASE_DIR) {|file|
		relative_name = file.sub(MIRROR_BASE_DIR, "")
		if File.basename(file) == "ChangeLog.txt"
			if not(relative_name.include?("zipslack"))
				files << relative_name
				Find.prune
			end
		end
		# putting this check *after* the one above,
		# lets us get the patches directories too
		# while still getting a bit of speed (1.5s)
		if relative_name.split("/").count > 2
			Find.prune
		end
	}
	puts "%f: watching %d changelogs" % [Time.now.to_f, files.count]
	files.each {|file|
		m = RE_REPO_NAME.match file
		if m[3].nil?
			file_name = "%sslackware%s-%s_ChangeLog.rss" % [FEEDS_BASE_DIR, m[1], m[2]]
		else
			file_name = "%sslackware%s-%s_%s_ChangeLog.rss" % [FEEDS_BASE_DIR, m[1], m[2], m[3]]
		end
		unless File.exist?(file_name)
			c_file = MIRROR_BASE_DIR + file
			changelog = Slackware::ChangeLog.new(c_file, :version => m[2])
			changelog.opts[:arch] = m[1] unless m[1].nil?
			if m[3].nil?
				changelog.opts[:url] = "%sslackware%s-%s/ChangeLog.txt" % [BASE_URL, m[1], m[2]]
				feed = File.open( "%sslackware%s-%s_ChangeLog.rss" % [FEEDS_BASE_DIR, m[1], m[2]], "w+")
			else
				changelog.opts[:url] = "%sslackware%s-%s/%s/ChangeLog.txt" % [BASE_URL, m[1], m[2], m[3]]
				feed = File.open( "%sslackware%s-%s_%s_ChangeLog.rss" % [FEEDS_BASE_DIR, m[1], m[2], m[3]], "w+")
			end
			changelog.parse
			puts "%f: Making a first feed: %s" % [Time.now.to_f, feed.path]
			feed << changelog.to_rss
			feed.close
			changelog = nil
		end
	}
end

def run_notifier
	n = INotify::Notifier.new
	dirs = Dir.glob(MIRROR_BASE_DIR + "*")
	dirs.concat(Dir.glob(MIRROR_BASE_DIR + "*/patches/"))
	dirs.each {|dir|
		next unless File.exist?(File.join(dir, "ChangeLog.txt"))
		puts "%f: working with %s" % [Time.now.to_f, dir]
		n.watch(dir, :moved_to) {|mfile|
			file_name = mfile.absolute_name
			if File.basename(file_name) == "ChangeLog.txt"
				puts "%f: looking into %s" % [Time.now.to_f, file_name]
				match_data = RE_REPO_NAME.match(file_name)

				unless match_data.nil?
					changelog = Slackware::ChangeLog.new(file_name, :version => match_data[2])
					changelog.opts[:arch] = match_data[1] unless match_data[1].nil?

					if match_data[3].nil?
						changelog.opts[:url] = "%sslackware%s-%s/ChangeLog.txt" % [
							BASE_URL,
							match_data[1],
							match_data[2]
						]
						feed = File.open( "%sslackware%s-%s_ChangeLog.rss" % [
								 FEEDS_BASE_DIR,
								 match_data[1],
								 match_data[2]
						], "w+")
					else
						changelog.opts[:url] = "%sslackware%s-%s/%s/ChangeLog.txt" % [
							BASE_URL,
							match_data[1],
							match_data[2],
							match_data[3]
						]
						feed = File.open( "%sslackware%s-%s_%s_ChangeLog.rss" % [
								 FEEDS_BASE_DIR,
								 match_data[1],
								 match_data[2],
								 match_data[3]
						], "w+")
					end
					begin
						changelog.parse
					rescue StandardError => ex
						puts "%f: %s" % [Time.now.to_f, ex.message]
						puts "%f: %s" % [Time.now.to_f, file_name]
						next
					end

					puts "%f: parsed %s to %s" % [Time.now.to_f, file_name, feed.path]

					feed << changelog.to_rss
					feed.close
					changelog = nil
				end
			end
		}
	}
	begin
		n.run
	rescue Interrupt
	end
end

## Main

#generate_new_if_none()
begin
	RubyProf.start
	run_notifier()
ensure
	result = RubyProf.stop

	RubyProf.measure_mode = RubyProf::PROCESS_TIME
	RubyProf.measure_mode = RubyProf::WALL_TIME
	RubyProf.measure_mode = RubyProf::CPU_TIME
	#RubyProf.measure_mode = RubyProf::ALLOCATIONS
	#RubyProf.measure_mode = RubyProf::MEMORY
	#RubyProf.measure_mode = RubyProf::GC_RUNS
	#RubyProf.measure_mode = RubyProf::GC_TIME

	output_file_name = File.join(ENV["HOME"],"%s-%s%s" % [Time.now.to_i.to_s,File.basename(__FILE__),".log"])
	output_file = File.open(output_file_name, "w+")
	printer = RubyProf::FlatPrinter.new(result)
	printer.print(output_file,0)
	puts "%f: %s written" % [Time.now.to_f, output_file_name]
	output_file.close
end
