# -*- coding: utf-8 -*-
"""
    Separate files in test corpus into this directory structure
        test.corpus/
            color/
                xobj/
                no.xobj/
            gray/
                xobj/
                no.xobj/

    Usage:
        python corpus_separate.py <corpus dir>*.pdf
    e.g.
        python corpus_separate.py /Users/peter/testdata/*.pdf

"""
from __future__ import division, print_function
import sys
import os
import csv
import shutil
from collections import defaultdict
from glob import glob


testResultPath = "xform.test.results.csv"
imageInfoPath = "xform.image.info.csv"


def read_csv(path):
    with open(path, 'r') as f:
        r = csv.reader(f)
        header = next(r)
        body = list(r)
    return header, body


trHeader, trBody = read_csv(testResultPath)
iiHeader, iiBody = read_csv(imageInfoPath)

print('%s: %d %dx%d %s' % (testResultPath, len(trHeader), len(trBody), len(trBody[0]), trHeader))
print('%s: %d %dx%d %s' % (imageInfoPath, len(iiHeader), len(iiBody), len(iiBody[0]), iiHeader))

all_files = [row[0] for row in trBody]
color_files = [row[0] for row in trBody if row[1] == 'true']
fail_files = [row[0] for row in trBody if row[1] == 'true' and row[2] == 'true']


# name_encoding = {name: [encodings]} : name ∈ all files in corpus
# name_colorspace = {name: [color spaces]} : name ∈ all files in corpus
name_encoding = defaultdict(set)
name_colorspace = defaultdict(set)
for name, _, _, _, encoding1, encoding2, colorspace in iiBody:
    if encoding1:
        name_encoding[name].add(encoding1)
    if encoding2:
        name_encoding[name].add(encoding2)
    if not encoding1 and not encoding2:
        name_encoding[name].add('[No encoding]')
    if colorspace == '':
        colorspace = '[No colorspace]'
    name_colorspace[name].add(colorspace)


NO_XOBJ = '[No XObjects]'

# Mark files that don't contain image XObjects
for name in all_files:
    e = name in name_encoding
    c = name in name_colorspace
    assert e == c, (name, e, c)
    if not e:
        name_encoding[name].add(NO_XOBJ)
        name_colorspace[name].add(NO_XOBJ)

# Files that contain image XObjects
xobj_files = [name for name in all_files if NO_XOBJ not in name_encoding[name]]


path_list = [path for a in sys.argv[1:] for path in glob(a)]
path_list = sorted(set(path_list))

basedir = 'test.corpus'
colordirs = {
    'color': set(color_files),
    'gray': set(all_files) - set(color_files)
}
xobjdirs = {
    'xobj': set(xobj_files),
    'no.xobj': set(all_files) - set(xobj_files)
}

for dc in colordirs:
    for dx in xobjdirs:
        dcx = os.path.join(basedir, dc, dx)
        match_names = colordirs[dc] & xobjdirs[dx]
        print('Copying %d files to %s' % (len(match_names), dcx))
        try:
            os.makedirs(dcx)
        except:
            pass

        for path in path_list:
            name = os.path.basename(path)
            if name not in match_names:
                continue
            dest = os.path.join(dcx, name)
            shutil.copyfile(path, dest)
