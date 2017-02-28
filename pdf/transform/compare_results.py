
TEXT_PETER = '''

 0 `/Users/pcadmin/go-work/src/trust-me/testdata/unlexicalized-parsing.pdf`
  1 `/Users/pcadmin/go-work/src/trust-me/testdata/pcfgs.pdf`
  2 `/Users/pcadmin/go-work/src/trust-me/testdata/multi_string_search.pdf`
  3 `/Users/pcadmin/go-work/src/trust-me/testdata/lexpcfgs.pdf`
  4 `/Users/pcadmin/go-work/src/trust-me/testdata/lcs.pdf`
  5 `/Users/pcadmin/go-work/src/trust-me/testdata/SA_implementations.pdf`
  6 `/Users/pcadmin/go-work/src/trust-me/testdata/fawkes.pdf`
  7 `/Users/pcadmin/go-work/src/trust-me/testdata/lec03-sp01.pdf`
  8 `/Users/pcadmin/go-work/src/trust-me/testdata/se.pdf`
  9 `/Users/pcadmin/go-work/src/trust-me/testdata/Suffix.pdf`
 10 `/Users/pcadmin/go-work/src/trust-me/testdata/nips01-discriminativegenerative.pdf`
 11 `/Users/pcadmin/go-work/src/trust-me/testdata/1603.06186v2.pdf`
 12 `/Users/pcadmin/go-work/src/trust-me/testdata/1206.6382.pdf`
 13 `/Users/pcadmin/go-work/src/trust-me/testdata/1409.1556v6.pdf`
 14 `/Users/pcadmin/go-work/src/trust-me/testdata/CaiHof-CIKM2004.pdf`
 15 `/Users/pcadmin/go-work/src/trust-me/testdata/geoff_hinton_dark14.pdf`
 16 `/Users/pcadmin/go-work/src/trust-me/testdata/ESCP-R reference_151008.pdf`
 17 `/Users/pcadmin/go-work/src/trust-me/testdata/extreme-days-and-nights-daylight-variation-in-the-arctic-reykjavik-murmansk-and-alert_e634.pdf`
 18 `/Users/pcadmin/go-work/src/trust-me/testdata/jkraaijeveld_thesis.pdf`
 19 `/Users/pcadmin/go-work/src/trust-me/testdata/Pages from Wolfcub Vision 01-Jan-16-6749.pdf`
 20 `/Users/pcadmin/go-work/src/trust-me/testdata/gv05-joc.pdf`
 21 `/Users/pcadmin/go-work/src/trust-me/testdata/lda.pdf`
 22 `/Users/pcadmin/go-work/src/trust-me/testdata/power_law_bins.pdf`
 23 `/Users/pcadmin/go-work/src/trust-me/testdata/rules_of_ml.pdf`
 24 `/Users/pcadmin/go-work/src/trust-me/testdata/compression.kdd06(1).pdf`
 25 `/Users/pcadmin/go-work/src/trust-me/testdata/twitter.pdf`
 26 `/Users/pcadmin/go-work/src/trust-me/testdata/pearson_science_8_sb_chapter_5_unit_5.2.pdf`
 27 `/Users/pcadmin/go-work/src/trust-me/testdata/1512.03547v2.pdf`
 28 `/Users/pcadmin/go-work/src/trust-me/testdata/BLUEBOOK.pdf`
 29 `/Users/pcadmin/go-work/src/trust-me/testdata/cvxopt_1306.0057v1.pdf`
 30 `/Users/pcadmin/go-work/src/trust-me/testdata/scan_alan_2016-03-30-10-38-15.pdf`
 31 `/Users/pcadmin/go-work/src/trust-me/testdata/Parsing-Probabilistic.pdf`
 32 `/Users/pcadmin/go-work/src/trust-me/testdata/2015-09-16-T23-39-51_ec2-user_ip-172-31-6-72_jim.pdf`
 33 `/Users/pcadmin/go-work/src/trust-me/testdata/Hierarchical Detection of Hard Exudates.pdf`
 34 `/Users/pcadmin/go-work/src/trust-me/testdata/WhatIsEnergy.pdf`
 35 `/Users/pcadmin/go-work/src/trust-me/testdata/Physics_Sample_Chapter_3.pdf`
 36 `/Users/pcadmin/go-work/src/trust-me/testdata/day3_TemporalImageProcessing.pdf`
 37 `/Users/pcadmin/go-work/src/trust-me/testdata/Lesson_054_handout.pdf`
 38 `/Users/pcadmin/go-work/src/trust-me/testdata/talk_Simons_part1_pdf.pdf`
 39 `/Users/pcadmin/go-work/src/trust-me/testdata/nips-tutorial-policy-optimization-Schulman-Abbeel.pdf`
peteraah:pdf pcadmin$ ./x_peter -a -o out.peter /Users/pcadmin/go-work/src/trust-me/testdata/*.pdf > blah.pete
'''

TEXT_UNIDOC = '''

 0 `/Users/pcadmin/go-work/src/trust-me/testdata/lec03-sp01.pdf`
  1 `/Users/pcadmin/go-work/src/trust-me/testdata/geoff_hinton_dark14.pdf`
  2 `/Users/pcadmin/go-work/src/trust-me/testdata/ESCP-R reference_151008.pdf`
  3 `/Users/pcadmin/go-work/src/trust-me/testdata/extreme-days-and-nights-daylight-variation-in-the-arctic-reykjavik-murmansk-and-alert_e634.pdf`
  4 `/Users/pcadmin/go-work/src/trust-me/testdata/Pages from Wolfcub Vision 01-Jan-16-6749.pdf`
  5 `/Users/pcadmin/go-work/src/trust-me/testdata/lda.pdf`
  6 `/Users/pcadmin/go-work/src/trust-me/testdata/rules_of_ml.pdf`
  7 `/Users/pcadmin/go-work/src/trust-me/testdata/pearson_science_8_sb_chapter_5_unit_5.2.pdf`
  8 `/Users/pcadmin/go-work/src/trust-me/testdata/BLUEBOOK.pdf`
  9 `/Users/pcadmin/go-work/src/trust-me/testdata/scan_alan_2016-03-30-10-38-15.pdf`
 10 `/Users/pcadmin/go-work/src/trust-me/testdata/Parsing-Probabilistic.pdf`
 11 `/Users/pcadmin/go-work/src/trust-me/testdata/2015-09-16-T23-39-51_ec2-user_ip-172-31-6-72_jim.pdf`
 12 `/Users/pcadmin/go-work/src/trust-me/testdata/Hierarchical Detection of Hard Exudates.pdf`
 13 `/Users/pcadmin/go-work/src/trust-me/testdata/WhatIsEnergy.pdf`
 14 `/Users/pcadmin/go-work/src/trust-me/testdata/Physics_Sample_Chapter_3.pdf`
 15 `/Users/pcadmin/go-work/src/trust-me/testdata/day3_TemporalImageProcessing.pdf`
 16 `/Users/pcadmin/go-work/src/trust-me/testdata/Lesson_054_handout.pdf`
 17 `/Users/pcadmin/go-work/src/trust-me/testdata/talk_Simons_part1_pdf.pdf`
 18 `/Users/pcadmin/go-work/src/trust-me/testdata/nips-tutorial-policy-optimization-Schulman-Abbeel.pdf`
'''



import re

RE_NAME = re.compile(r'\d+\s+`(.*?)`', re.MULTILINE | re.DOTALL)

def get_names(text):
    names = []
    for m in RE_NAME.finditer(text):
        n = m.group(1)
        n = n.split('/')[-1]
        names.append(n)
    return names

names_peter = set(get_names(TEXT_PETER))
names_unidoc = set(get_names(TEXT_UNIDOC))

peter_unidoc = sorted(names_peter - names_unidoc)
unidoc_peter = sorted(names_unidoc - names_peter)

print('peter - unidoc: %d' % len(peter_unidoc))
print('unidoc - peter: %d' % len(unidoc_peter))
for i, n in enumerate(sorted(names_unidoc)):
    print('%2d: %s' % (i, n))

