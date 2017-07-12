from glob import glob
import sys
import os
from subprocess import call


def move(path, dest):
    base = os.path.basename(path)
    base = os.path.splitext(base)[0]
    return os.path.join(dest, base + '.txt')


pattern = sys.argv[2:]
dest = sys.argv[1]

print('pattern="%s"' % pattern)
print('dest="%s"' % dest)

try:
    os.mkdir(dest)
except FileExistsError:
    pass

files = []
for p in pattern:
    print('$$', p)
    if os.path.isdir(p):
        for f in glob(p):
            print('!!', f)
            files.append(f)
    else:
        files.append(p)
print('%d files' % len(files))

for i, path in enumerate(files):
    dst = move(path, dest)
    print('%3d: %s' % (i, dst))
    call(["pdftotext", path, dst])
