
TEXT = '''

  {"\"",  3, {tchkNum,    tchkNum,    tchkString},
          &Gfx::opMoveSetShowText},
  {"'",   1, {tchkString},
          &Gfx::opMoveShowText},
  {"B",   0, {tchkNone},
          &Gfx::opFillStroke},
  {"B*",  0, {tchkNone},
          &Gfx::opEOFillStroke},
  {"BDC", 2, {tchkName,   tchkProps},
          &Gfx::opBeginMarkedContent},
  {"BI",  0, {tchkNone},
          &Gfx::opBeginImage},
  {"BMC", 1, {tchkName},
          &Gfx::opBeginMarkedContent},
  {"BT",  0, {tchkNone},
          &Gfx::opBeginText},
  {"BX",  0, {tchkNone},
          &Gfx::opBeginIgnoreUndef},
  {"CS",  1, {tchkName},
          &Gfx::opSetStrokeColorSpace},
  {"DP",  2, {tchkName,   tchkProps},
          &Gfx::opMarkPoint},
  {"Do",  1, {tchkName},
          &Gfx::opXObject},
  {"EI",  0, {tchkNone},
          &Gfx::opEndImage},
  {"EMC", 0, {tchkNone},
          &Gfx::opEndMarkedContent},
  {"ET",  0, {tchkNone},
          &Gfx::opEndText},
  {"EX",  0, {tchkNone},
          &Gfx::opEndIgnoreUndef},
  {"F",   0, {tchkNone},
          &Gfx::opFill},
  {"G",   1, {tchkNum},
          &Gfx::opSetStrokeGray},
  {"ID",  0, {tchkNone},
          &Gfx::opImageData},
  {"J",   1, {tchkInt},
          &Gfx::opSetLineCap},
  {"K",   4, {tchkNum,    tchkNum,    tchkNum,    tchkNum},
          &Gfx::opSetStrokeCMYKColor},
  {"M",   1, {tchkNum},
          &Gfx::opSetMiterLimit},
  {"MP",  1, {tchkName},
          &Gfx::opMarkPoint},
  {"Q",   0, {tchkNone},
          &Gfx::opRestore},
  {"RG",  3, {tchkNum,    tchkNum,    tchkNum},
          &Gfx::opSetStrokeRGBColor},
  {"S",   0, {tchkNone},
          &Gfx::opStroke},
  {"SC",  -4, {tchkNum,   tchkNum,    tchkNum,    tchkNum},
          &Gfx::opSetStrokeColor},
  {"SCN", -33, {tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN},
          &Gfx::opSetStrokeColorN},
  {"T*",  0, {tchkNone},
          &Gfx::opTextNextLine},
  {"TD",  2, {tchkNum,    tchkNum},
          &Gfx::opTextMoveSet},
  {"TJ",  1, {tchkArray},
          &Gfx::opShowSpaceText},
  {"TL",  1, {tchkNum},
          &Gfx::opSetTextLeading},
  {"Tc",  1, {tchkNum},
          &Gfx::opSetCharSpacing},
  {"Td",  2, {tchkNum,    tchkNum},
          &Gfx::opTextMove},
  {"Tf",  2, {tchkName,   tchkNum},
          &Gfx::opSetFont},
  {"Tj",  1, {tchkString},
          &Gfx::opShowText},
  {"Tm",  6, {tchkNum,    tchkNum,    tchkNum,    tchkNum,
        tchkNum,    tchkNum},
          &Gfx::opSetTextMatrix},
  {"Tr",  1, {tchkInt},
          &Gfx::opSetTextRender},
  {"Ts",  1, {tchkNum},
          &Gfx::opSetTextRise},
  {"Tw",  1, {tchkNum},
          &Gfx::opSetWordSpacing},
  {"Tz",  1, {tchkNum},
          &Gfx::opSetHorizScaling},
  {"W",   0, {tchkNone},
          &Gfx::opClip},
  {"W*",  0, {tchkNone},
          &Gfx::opEOClip},
  {"b",   0, {tchkNone},
          &Gfx::opCloseFillStroke},
  {"b*",  0, {tchkNone},
          &Gfx::opCloseEOFillStroke},
  {"c",   6, {tchkNum,    tchkNum,    tchkNum,    tchkNum,
        tchkNum,    tchkNum},
          &Gfx::opCurveTo},
  {"cm",  6, {tchkNum,    tchkNum,    tchkNum,    tchkNum,
        tchkNum,    tchkNum},
          &Gfx::opConcat},
  {"cs",  1, {tchkName},
          &Gfx::opSetFillColorSpace},
  {"d",   2, {tchkArray,  tchkNum},
          &Gfx::opSetDash},
  {"d0",  2, {tchkNum,    tchkNum},
          &Gfx::opSetCharWidth},
  {"d1",  6, {tchkNum,    tchkNum,    tchkNum,    tchkNum,
        tchkNum,    tchkNum},
          &Gfx::opSetCacheDevice},
  {"f",   0, {tchkNone},
          &Gfx::opFill},
  {"f*",  0, {tchkNone},
          &Gfx::opEOFill},
  {"g",   1, {tchkNum},
          &Gfx::opSetFillGray},
  {"gs",  1, {tchkName},
          &Gfx::opSetExtGState},
  {"h",   0, {tchkNone},
          &Gfx::opClosePath},
  {"i",   1, {tchkNum},
          &Gfx::opSetFlat},
  {"j",   1, {tchkInt},
          &Gfx::opSetLineJoin},
  {"k",   4, {tchkNum,    tchkNum,    tchkNum,    tchkNum},
          &Gfx::opSetFillCMYKColor},
  {"l",   2, {tchkNum,    tchkNum},
          &Gfx::opLineTo},
  {"m",   2, {tchkNum,    tchkNum},
          &Gfx::opMoveTo},
  {"n",   0, {tchkNone},
          &Gfx::opEndPath},
  {"q",   0, {tchkNone},
          &Gfx::opSave},
  {"re",  4, {tchkNum,    tchkNum,    tchkNum,    tchkNum},
          &Gfx::opRectangle},
  {"rg",  3, {tchkNum,    tchkNum,    tchkNum},
          &Gfx::opSetFillRGBColor},
  {"ri",  1, {tchkName},
          &Gfx::opSetRenderingIntent},
  {"s",   0, {tchkNone},
          &Gfx::opCloseStroke},
  {"sc",  -4, {tchkNum,   tchkNum,    tchkNum,    tchkNum},
          &Gfx::opSetFillColor},
  {"scn", -33, {tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN,   tchkSCN,    tchkSCN,    tchkSCN,
          tchkSCN},
          &Gfx::opSetFillColorN},
  {"sh",  1, {tchkName},
          &Gfx::opShFill},
  {"v",   4, {tchkNum,    tchkNum,    tchkNum,    tchkNum},
          &Gfx::opCurveTo1},
  {"w",   1, {tchkNum},
          &Gfx::opSetLineWidth},
  {"y",   4, {tchkNum,    tchkNum,    tchkNum,    tchkNum},
          &Gfx::opCurveTo2},
};
'''

# TEXT = '''{"'",   1, {tchkString},   &Gfx::opMoveShowText},'''


name_map = {
    'Array': 'pdfTypeArray',
    'Int': 'pdfTypeInteger',
    'Name': 'pdfTypeName',
    'Num': 'pdfTypeNumber',
    'Props': 'pdfTypeNameDict',
    'SCN': 'pdfTypeNameNumber',
    'String': 'pdfTypeString',
    'None' : '',
}


import re

RE_OP = re.compile(r'\{\s*"(.{,4})",\s*(-?\d+),\s*\{([\w\s,]+)\},\s*&Gfx::(\w+)\},', re.MULTILINE | re.DOTALL)

all_ops = []
all_args = set()
for i, m in enumerate(RE_OP.finditer(TEXT)):
    op = m.group(1)
    n = int(m.group(2))
    args = m.group(3)
    args = args.split(',')
    args = [a.strip() for a in args]
    args = [a[4:] for a in args]
    args = [name_map[a] for a in args]
    op = op.replace('"', '\\"')

    if len(args) == 1 and args[0] == '':
        args = []
    all_args = all_args.union(args)
    func = m.group(4)[2:]

    replacement = r'\1 \2'
    func = re.sub(r'(.)([A-Z][a-z]+)', replacement, func, count=20)
    # func = re.sub(r'(.)([A-Z][a-z]+)', replacement, func)

    # print("%2d: %4s %d %s %s" % (i, op, n, args, func[2:]))
    assert abs(n) == len(args), (m.groups(), n, args)
    all_ops.append((op, n, args, func))


def case_key(op, n, args, func):
   return op.lower(), op

all_ops.sort(key=lambda x: case_key(*x))

print('var allValidOperations = map[string] operationSpec {')
for i, (op, n, args, func) in enumerate(all_ops):
    arg_str = '[]PdfObjectType{%s}' % ', '.join( args)
    if n < 0:
      arg_str = 'nil'
    # print('\t"%s": {"%s", %s, "%s"}, // %3d ' % (op, op, arg_str, func, i))
    print('\t"%s": %s, // %s ' % ( op, arg_str, func))
print('}')

print('// all_args: %s' % sorted(all_args))
