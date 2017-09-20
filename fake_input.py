#!/usr/bin/python

import sys, datetime, time

def main():
    while True:
        now = datetime.datetime.now()
        sys.stdout.write("The time is now %s\n" % now)
        sys.stdout.flush()
        time.sleep(1)

if __name__ == '__main__':
    main()
