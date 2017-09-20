#!/usr/bin/python

import sys, datetime, time

def main():
    while True:
        now = datetime.datetime.now()
        print "The time is now %s" % now
        time.sleep(1)

if __name__ == '__main__':
    main()
