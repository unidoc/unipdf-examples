import re


def get_lines(path):
    lines = []
    with open(path, 'rt') as f:
        for ln in f:
            lines.append(ln.rstrip('\n'))
    return lines


split_re = re.compile(r'\s*[,;:]\s*')
if True:
    s1 = '12.0,0.0,0.0,12.0:0.0,0.0'
    s2 = '12.0,0.0,0.0,12.0:0.0,0.0'
    m1 = split_re.search(s1)
    assert m1, s1
    m2 = split_re.search(s2)
    assert m2, s2

# [DEBUG]  text.go:648 "A" stateMatrix=[12.0,0.0,0.0,12.0:0.0,0.0] CTM=[1.0,0.0,0.0,1.0:0.0,0.0] Tm=[1.0,0.0,0.0,1.0:219.6,633.0]
uni_re = re.compile(r'text.go:\d+\s*"(.*?)".*Tm=\[(.*)\]$')

if True:
    ln = '[DEBUG]  text.go:654 "A" stateMatrix=[12.0,0.0,0.0,12.0:0.0,0.0] CTM=[1.0,0.0,0.0,1.0:0.0,0.0] Tm=[1.0,0.0,0.0,1.0:219.6,633.0]'
    m = uni_re.search(ln)
    assert m

# ###~ unicode=A
box_re1 = re.compile(r'^###~ unicode=\s*(.*?)\s*$')
box_re2 = re.compile(r'^@2\s*textMatrix=\[(.*)\]$')

if True:
    ln = '###~ unicode=A]'
    m = box_re1.search(ln)
    assert m

substitutions = {
    "fl": "ﬂ",
    "fi": "ﬁ",
}


def normalize(s):
    if s in substitutions:
        return substitutions[s]
    return s


def to_matrix(s):
    parts = split_re.split(s)
    return [float(v) for v in parts]


DELTA = 0.2


def same_matrix(m1, m2):
    assert len(m1) == 6, m1
    assert len(m2) == 6, m2
    return all(abs(x1-x2) <= DELTA for x1, x2 in zip(m1, m2))


def get_tm_uni(lines):
    tm = []
    for i, ln in enumerate(lines):
        m = uni_re.search(ln)
        if m is None:
            continue
        tm.append((m.group(1), to_matrix(m.group(2)), i + 1))
    return tm


def get_tm_box(lines):
    tm = []
    i = 0
    s = ''
    for j, ln in enumerate(lines):
        m = box_re1.search(ln)
        if m is not None:
            # print('---', j, ln)
            # assert s == '', (j, ln)
            i = j
            s = m.group(1)
            if s == '':
                s = ' '
            continue
        m = box_re2.search(ln)
        if m is not None:
            # print('+++', j, ln)
            # assert s != '', (j, ln)
            tm.append((s, to_matrix(m.group(1)), i + 1))
            s = ''
    return tm


# go run pdf_render_text.go -d  file.pdf > blah
uni_path = "/Users/pcadmin/go-work/src/github.com/unidoc/unidoc-examples/pdf/text/blah"
# java -jar ./app/target/pdfbox-app-2.0.9-SNAPSHOT.jar ExtractText -sort filt.pdf out.box.26 >blah
box_path = "/Users/pcadmin/pdf/pdfbox.orig/blah"

uni_lines = get_lines(uni_path)
box_lines = get_lines(box_path)
print('uni="%s" %d' % (uni_path, len(uni_lines)))
print('box="%s" %d' % (box_path, len(box_lines)))

tm_uni = get_tm_uni(uni_lines)
print('tm_uni=%d' % len(tm_uni))
assert tm_uni

tm_box = get_tm_box(box_lines)
print('tm_box=%d' % len(tm_box))
assert tm_box

print('-' * 80)
for i, tm in enumerate(tm_uni[:5]):
    print('%5d: %s' % (i, tm))
print('-' * 80)

for i, tm in enumerate(tm_box[:5]):
    print('%5d: %s' % (i, tm))

print('-' * 80)
for i in range(3000):
    print('%5d: %s %s' % (i, tm_uni[i], tm_box[i]))
    if normalize(tm_uni[i][0]) != normalize(tm_box[i][0]):
        break
    if not same_matrix(tm_uni[i][1], tm_box[i][1]):
        break


def show(tm, i0, i1):
    return ''.join(normalize(v[0]) for v in tm[i0:i1])


i0 = max(0, i - 40)
i1 = max(0, i + 10)
print('uni:', show(tm_uni, i0, i+1), '--', show(tm_uni, i+1, i1))
print('box:', show(tm_box, i0, i+1), '--', show(tm_box, i+1, i1))

# print(''.join(v[0] for v in tm_box[i0:i+1]))
