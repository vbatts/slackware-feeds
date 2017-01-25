#!/usr/bin/env python

import os
import sys
import glob
import time

sys.path.insert(0, "/home/vbatts/opt/lib/python2.5/site-packages")
sys.path.insert(0, "/home/vbatts/opt/lib/python2.5")
import pyinotify

dir_path	= "/mirrors/ftp.slackware.com/pub/slackware"

def process_changelog_rss(event):
	if os.path.basename(event.pathname) == "ChangeLog.txt":
		print "%f: proccessing %s" % (time.time(), event)
		os.system("/home/vbatts/opt/bin/ruby /home/vbatts/bin/gen_changlog_rss.rb %s" % event.pathname)

def main(args):
	wm = pyinotify.WatchManager()

	notifier = pyinotify.Notifier(wm)

	for dir in glob.glob(dir_path + "/*/"):
		if os.path.exists(dir + "ChangeLog.txt"):
			print "%f: Adding watch for %s" % (time.time(), dir)
			wm.add_watch(dir, pyinotify.IN_MOVED_TO, rec=False, proc_fun=process_changelog_rss)

	for dir in glob.glob(dir_path + "/*/patches/"):
		print "%f: Adding watch for %s" % (time.time(), dir)
		wm.add_watch(dir, pyinotify.IN_MOVED_TO, rec=False, proc_fun=process_changelog_rss)

	#wm.add_watch("/home/vbatts/", pyinotify.IN_MOVED_TO, rec=False, proc_fun=process_changelog_rss)

	notifier.loop()


if __name__ == "__main__": main(sys.argv[1:])

