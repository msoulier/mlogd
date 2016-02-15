#!/usr/bin/python

import sys
from random_words import LoremIpsum

def main():
    nlines = 1000
    if len(sys.argv) > 1:
        nlines = int(sys.argv[1])
    li = LoremIpsum()
    for sentence in li.get_sentences_list(nlines):
        print sentence

if __name__ == '__main__':
    main()
