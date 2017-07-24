# -*- coding: utf-8 -*-
"""
    Analyze test results from csv files produced by pdf_transform_content_streams

    Shows encoding filters and colorspace spaces used in
        - all files
        - color files
        - color files that pdf_transform_content_streams can't convert to grayscale
"""
from __future__ import division, print_function
import csv
from collections import defaultdict


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


def num_values(a_dict):
    return sum(len(v) for v in a_dict.values())


def encoding_colorspace_counts(files):
    """
        Returns
            {encoding: [files in `files` that contains image XObjects with this encoding]},
            {colorspace: [files in `files` that contains image XObjects with this colorspace]}
            Names of all files with this
    """
    encoding_name = defaultdict(set)
    colorspace_name = defaultdict(set)
    for name in all_files:
        for encoding in name_encoding[name]:
            encoding_name[encoding].add(name)
        for colorspace in name_colorspace[name]:
            colorspace_name[colorspace].add(name)
    return encoding_name, colorspace_name


def summarize(name, files):
    """Prints
        Total number of files in `files`
        Encodings used in those files
        Color spaces used in those files
    """
    encoding_name, colorspace_name = encoding_colorspace_counts(files)
    print('-' * 80)
    print('%s: %d total' % (name, len(files)))
    print('%d encodings' % len(encoding_name))
    for i, encoding in enumerate(sorted(encoding_name, key=lambda k: -len(encoding_name[k]))):
        print('%4d: %-20s %s' % (i, encoding, len(encoding_name[encoding])))
    print('%d colorspace spaces' % len(colorspace_name))
    for i, colorspace in enumerate(sorted(colorspace_name, key=lambda k: -len(colorspace_name[k]))):
        print('%4d: %-20s %s' % (i, colorspace, len(colorspace_name[colorspace])))

    assert num_values(encoding_name) >= len(files), (num_values(encoding_name), len(files))
    assert num_values(colorspace_name) >= len(files), (num_values(colorspace_name), len(files))

#
# Print the summary data
#
summarize('all files', all_files)      # all files tested
summarize('color_files', color_files)  # color files in all_files
summarize('fail_files', fail_files) # files in color_files that couldn't be converted to grayscale

