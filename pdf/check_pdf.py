import re
import sys


path = sys.argv[1]
text = open(path, 'rb').read()

'''
    <</BitsPerComponent 1/ColorSpace /DeviceGray/Decode [1 0]/DecodeParams <</Columns 35/K -1>>
    /Filter /CCITTFaxDecode/Height 46/ImageMask true/Length 51/Subtype /Image/Type
    /XObject/Width 35>>
    stream
    &?P$?-?????????????????#A????????x;Xp?@
    endstream
'''


def find_dict(i1, level):
    print('##', i1, level)
    i = i1
    d0, d1 = -1, -1
    while i >= max(1, i1 - 1000):
        if text[i:i + 2] == b'>>':
            # print("  ->> ", text[i:i + 2])
            if d1 < 0:
                d1 = i + 2
            else:
                i, _ = find_dict(i, level + 1)
        elif text[i:i + 2] == b'<<':
            # print("  ->> ", text[i:i + 2])
            assert d1 >= 0
            d0 = i
            return d0, d1
        i -= 1
    print('@@@ text[%d:%d]="%s"' % (i, i1, text[i:i1]))
    raise ValueError


RE_LENGTH = re.compile(b'/Length\s+(\d+)', re.MULTILINE | re.DOTALL)
RE_STREAM = re.compile(b'stream\s+(.+?)\s+endstream\s+', re.MULTILINE | re.DOTALL)


def find_length(text):
    m = RE_LENGTH.search(text)
    assert m
    return int(m.group(1))


def find_stream(text):
    # b = 'stream'
    # e = 'endstream'
    # while isspace(text[i]):
    #     i += 1
    m = RE_STREAM.search(text)
    assert m, 'find_stream text=%d "%s"' % (len(text), text[:100])
    return m.group(0), m.group(1)


obj_list = []
i0 = 0
while True:
    i = text.find(b'stream', i0)
    if i < 0:
        break
    d0, d1 = find_dict(i, 0)
    params = text[d0:d1]
    length = find_length(params)
    print('%3d: %6d: length=%3d params=%s' %
        (0, i, length, params))
    whole, contents = find_stream(text[i-2:i + 2 * length])
    obj = (i, params, length, contents)

    # print('%3d: %6d: length=%3d params=%s contents="%s"' %
    #     (0, i, length, params, contents))
    obj_list.append(obj)
    i0 = i + len(whole)

for n, (i, params, length, contents) in enumerate(obj_list):
    print('%3d: %6d: length=%d=%d params=%s contents="%s"' %
        (n, i, length, len(contents), params, contents[:10]))
