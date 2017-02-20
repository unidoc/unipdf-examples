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
from collections import defaultdict, namedtuple
from glob import glob


def dict_counts(a_dict):
    return {k: len(v) for k, v in a_dict.items()}


def num_values(a_dict):
    return sum(len(v) for v in a_dict.values())


def toBool(s):
    return s.strip().lower() == 'true'


def makedir(a_dir):
    try:
        os.makedirs(a_dir)
    except:
        pass


def read_csv(path):
    with open(path, 'r') as f:
        r = csv.reader(f)
        header = next(r)
        body = list(r)
    return header, body



testResultPath = "xform.test.results.csv"
imageInfoPath = "xform.image.info.csv"
trHeader, trBody = read_csv(testResultPath)

key_type_list = [
    ('name', str),
    ('colorIn', toBool),
    ('colorOut', toBool),
    ('numPages', int),
    ('duration', float),
    ('image_xobj', int),
    ('form_xobj', int)
]
types = dict(key_type_list)

key_list = [key for key, _ in key_type_list]
header = [key.replace(' ', '_') for key in trHeader]
column_key = {header.index(key): key for key in key_list}
print('column_key=%s' % column_key)
column_type = {col: types[key] for col, key in column_key.items()}
print('column_type=%s' % column_type)
Result = namedtuple('Result', key_list)

print(trHeader)
print('%s: %d %dx%d %s' % (testResultPath, len(trHeader), len(trBody), len(trBody[0]), trHeader))

trBody = [Result(*[column_type[i](x) for i, x in enumerate(row)]) for row in trBody]

all_files = {row.name for row in trBody}
color_files = {row.name for row in trBody if row.colorIn}
fail_files = {row.name for row in trBody if (row.colorIn and row.colorOut)}
success_files = all_files - fail_files
img_xobj_files = {row.name for row in trBody if row.image_xobj > 0}


basedir = 'test.corpus.new'
colordirs = {
    'color': set(color_files),
    'gray': set(all_files) - set(color_files)
}
xobjdirs = {
    'xobj': set(img_xobj_files),
    'no.xobj': set(all_files) - set(img_xobj_files)
}

print('colordirs=%s %d' % (dict_counts(colordirs), num_values(colordirs)))
print('xobjdirs=%s %d' % (dict_counts(xobjdirs), num_values(xobjdirs)))

name_dir = {}
print('Results match')
for dc in colordirs:
    for dx in xobjdirs:
        dcx = os.path.join(basedir, dc, dx)
        match_names = colordirs[dc] & xobjdirs[dx]
        for name in match_names:
            name_dir[name] = dcx
        print('%4d [%d pass + %d fail] "%s"' % (len(match_names), len(match_names & success_files),
                                      len(match_names & fail_files), dcx))
        makedir(dcx)


dir_other = os.path.join(basedir, 'other')
makedir(dir_other)


path_list = [path for a in sys.argv[1:] for path in glob(a)]
path_list = sorted(set(path_list))

dest_count = defaultdict(int)
for path in path_list:
    name = os.path.basename(path)
    dest_dir = name_dir.get(name, dir_other)
    dest = os.path.join(dest_dir, name)
    assert dest.lower() != path.lower()
    dest_count[dest_dir] += 1
    # print('%50s => %s' % (name, dest))
    shutil.copyfile(path, dest)

print('Files copied. Total = %d' % len(path_list))
for dest in sorted(dest_count):
     print('%4d "%s"' % (dest_count[dest], dest))
