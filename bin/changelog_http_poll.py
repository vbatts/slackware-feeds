#!/usr/bin/env python
# Mon Oct 17 08:25:29 PDT 2011
# copyright 2011  Vincent Batts, Vienna, VA, USA

# switching from an inotify watcher, to an http poll
# since what lands on connie.slackware.com usually doesn't go public 
# immediately


import os
import sys
import glob
import time
from datetime import datetime
from datetime import timedelta
from time import mktime
import urllib2
import anydbm

DEFAULT_DB = os.path.join(os.getenv('HOME'), '.slackware_changelog.db')
DEFAULT_URL = "http://slackware.osuosl.org/"
SLACKWARE_DIR_PATH = "/mirrors/ftp.slackware.com/pub/slackware"
RSS_DIR_PATH = "/home/vbatts/public_html/feeds"

'''
slackware-12.2_ChangeLog.rss
/home/vbatts/public_html/feeds/slackware-10.1_patches_ChangeLog.rss
/home/vbatts/public_html/feeds/slackware-8.1_patches_ChangeLog.rss
>>> for i in c.slackware_versions(): print i
...
/mirrors/ftp.slackware.com/pub/slackware/slackware64-13.0/ChangeLog.txt
/mirrors/ftp.slackware.com/pub/slackware/slackware-8.1/ChangeLog.txt
/mirrors/ftp.slackware.com/pub/slackware/slackware64-13.37/ChangeLog.txt
/mirrors/ftp.slackware.com/pub/slackware/slackware-13.0/ChangeLog.txt
/mirrors/ftp.slackware.com/pub/sla
'''

def rss_files():
	for item in glob.glob(RSS_DIR_PATH + "/*.rss"):
		yield item

def rss_files_format(str):
	if str.startswith(RSS_DIR_PATH + "/"):
		str = str[len(RSS_DIR_PATH + "/"):]
	if str.endswith(".rss"):
		str = str[:-4]
		str = str + '.txt'
	return str.replace('_','/')

def rss_files_cleaned():
	for i in rss_files():
		yield rss_files_format(i)

def slackware_versions():
	changes = glob.glob(SLACKWARE_DIR_PATH + "/*/ChangeLog.txt")
	patches = glob.glob(SLACKWARE_DIR_PATH + "/*/patches/ChangeLog.txt")
	for item in changes + patches:
		yield item

def slackware_versions_format(str):
	if str.startswith(SLACKWARE_DIR_PATH + "/"):
		str = str[len(SLACKWARE_DIR_PATH + "/"):]
	if str.endswith("/"):
		str = str[:-1]
	if str.startswith("/"):
		str = str[1:]
	if str.endswith(".txt"):
		str = str[:-4]
	return str.replace('/','_')

def slackware_versions_strip():
	for i in slackware_versions():
		yield i[len(SLACKWARE_DIR_PATH + "/"):]

def slackware_versions_rss():
	for i in slackware_versions():
		yield slackware_versions_format(i)

def process_changelog_rss(pathname):
	if os.path.basename(pathname) == "ChangeLog.txt":
		print "%f: proccessing %s" % (time.time(), pathname)
		# XXX REPLACE ME!!
		cmd = "/home/vbatts/opt/bin/ruby /home/vbatts/bin/gen_changlog_rss.rb %s" % pathname
		print cmd
		print os.system(cmd)
	else:
		print '[WARN] "%s" is not a ChangeLog.txt file' % pathname

def db_setup(name = DEFAULT_DB):
	try:
		return anydbm.open(name, 'c')
	except:
		return None

def db_teardown(db):
	try:
		return db.close()
	except:
		return None

def db_add_ts(db, key, val):
	if type(val) == float:
		db[key] = str(val)
	if type(val) == datetime:
		db[key] = str(unix_time(val))
	return db[key]

def db_get_ts(db, key):
	try:
		return datetime.fromtimestamp(float(db[key]))
	except KeyError:
		return None

def unix_time(dt):
	return mktime(dt.timetuple())+1e-6*dt.microsecond

def time_from_header(str):
	return datetime.strptime(str, "%a, %d %b %Y %H:%M:%S %Z")

def get_remote_header(url, header):
	try:
		req = urllib2.Request(url)
		resp = urllib2.urlopen(req)
		return resp.headers.getheader(header)
	except:
		return None

def get_remote_time_str(url):
	return get_remote_header(url,"last-modified")

def get_remote_time(url):
	time_str = get_remote_time_str(url)
	if time_str:
		return time_from_header(time_str)
	else:
		return None

def get_local_time(path):
	try:
		time_flt = os.stat(path).st_mtime
		return datetime.fromtimestamp(time_flt)
	except:
		return None

def main(args):
	try:
		db = db_setup()
		if db == None:
			print "ERROR: could not setup database at %s" % DEFAULT_DB
			return 1

		for i in slackware_versions_strip():
			# i'm not going to worry about this file, right now
			if i == 'slackware/ChangeLog.txt':
				continue

			rss_file_name = os.path.join(RSS_DIR_PATH,
					slackware_versions_format(i) + ".rss")
			rss_ts = get_local_time(rss_file_name)
			curr_ts = get_local_time(os.path.join(SLACKWARE_DIR_PATH, i))
			prev_ts = db_get_ts( db, "local_" + i)

			# Go no further for this file
			if curr_ts == prev_ts and os.path.exists(rss_file_name) and rss_ts > prev_ts:
				print '[INFO] Local time of "%s" is same as the database has' % i
				continue

			db_add_ts( db, "local_" + i, curr_ts)

			remote_ts = get_remote_time(DEFAULT_URL + i)
			print '[INFO] inserting remote_%s: %s' % (i,remote_ts)
			db_add_ts( db, "remote_" + i, remote_ts)

			if prev_ts == None or (remote_ts - prev_ts) == timedelta(hours=7):
				print '[INFO] local and remote ChangeLog times match'
				if rss_ts == None:
					print '[INFO] RSS file (%s) does not exist' % (rss_ts)
					print '[INFO] Processing "%s"' % rss_file_name
					process_changelog_rss(os.path.join(SLACKWARE_DIR_PATH, i))
				elif prev_ts == None or rss_ts < prev_ts:
					print '[INFO] RSS file (%s) is older than the ChangeLog (%s)' % (rss_ts, prev_ts)
					print '[INFO] Processing "%s"' % rss_file_name
					process_changelog_rss(os.path.join(SLACKWARE_DIR_PATH, i))
				else:
					print '[INFO] RSS seems current'
	finally:
		try:
			os.wait()
		except:
			pass
		db_teardown(db)

if __name__ == "__main__": sys.exit(main(sys.argv[1:]))

