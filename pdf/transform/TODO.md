colorspace.go
image.go
xobject.go
check all !@#$
TODO
====
p297 10.3 Conversions among Device Colour Spaces
gray = .3r + .59g + .11b
     = 1.0 -min(1.0, k + .3c + .59m + .11y)


p182 8.7.4.3Shading Dictionaries

Convert colors to grayscale in shading dicts. See extreme-days-and-nights-daylight-variation-in-the-arctic-reykjavik-murmansk-and-alert_e634.pdf

ahs-scn-bjh-quick-facts.pdf


34 0 obj
<</BBox [350.782000 452.200000 626.534000 51.347300]/Group 45 0 R/Length 3917/Matrix [1.000000 0.000000 0.000000 1.000000 0.000000 0.000000]/Resources <</ColorSpace <</CS1 31 0 R/CS0 35 0 R>>/ExtGState <</GS0 29 0 R>>/Shading <</Sh0 40 0 R>>>>/Subtype /Form/Type /XObject>>
stream

40 0 obj
<</Domain [0.000000 1.000000]/Extend [true true]/Function 41 0 R/ShadingType 2/AntiAlias false/ColorSpace 35 0 R/Coords [0.000000 0.000000 1.000000 0.000000]>>
endobj

41 0 obj
<</Domain [0.000000 1.000000]/Encode [0.000000 1.000000 0.000000 1.000000 1.000000 0.000000]/FunctionType 3/Functions [42 0 R 43 0 R 44 0 R]/Bounds [0.269226 0.994507]>>
endobj
42 0 obj
<</FunctionType 2/N 1.000000/C0 [0.000000 0.474808]/C1 [0.000000 1.000000]/Domain [0.000000 1.000000]>>
endobj
43 0 obj
<</FunctionType 2/N 1.000000/C0 [0.000000 1.000000]/C1 [1.000000 0.000000]/Domain [0.000000 1.000000]>>
endobj
44 0 obj
<</C0 [1.000000 0.000000]/C1 [1.000000 0.000000]/Domain [0.000000 1.000000]/FunctionType 2/N 13.513400>>
endobj

35 0 obj
[/DeviceN [/PANTONE#20300#20C /Yellow] /DeviceCMYK 36 0 R 37 0 R]
endobj
36 0 obj
<</Domain [0.000000 1.000000 0.000000 1.000000]/Filter /FlateDecode/FunctionType 4/Length 126/Range [0.000000 1.000000 0.000000 1.000000 0.000000 1.000000 0.000000 1.000000]>>
stream
H‰ª6Ô3....

37 0 obj
<</Colorants 38 0 R/Process 39 0 R/Subtype /NChannel>>
endobj

38 0 obj
<</PANTONE#20300#20C 31 0 R>>
endobj
39 0 obj
<</ColorSpace /DeviceCMYK/Components [/Cyan /Magenta /Yellow /Black]>>
endobj

31 0 obj
[/Separation /PANTONE#20300#20C
/DeviceCMYK
<</Range [0.000000 1.000000 0.000000 1.000000 0.000000 1.000000 0.000000 1.000000]
  /C0 [0.000000 0.000000 0.000000 0.000000]
  /C1 [1.000000 0.440002 0.000000 0.000000]
  /N 1.000000
  /FunctionType 2
  /Domain [0.000000 1.000000]>>]
endobj

p186 Table 80 – Additional Entries Specific to a Type 2 Shading Dictionary
Type
Value
Coords
array
(Required) An array of four numbers [x0 y0 x1 y1] specifying the starting and ending coordinates
of the axis, expressed in the shading’s target coordinate space.
Domain
array
(Optional) An array of two numbers [t0 t1] specifying the limiting values of a parametric variable t.
The variable is considered to vary linearly between these two values as the colour gradient varies
between the starting and ending points of the axis. The variable t becomes the input argument to the
colour function(s). Default value: [0.0 1.0].
Function
function
(Required) A 1-in, n-out function or an array of n 1-in, 1-out functions (where n is the number of
colour components in the shading dictionary’s colour space). The function(s) shall be called with
values of the parametric variable t in the domain defined by the Domain entry. Each function’s domain
shall be a superset of that of the shading dictionary. If the value returned by the function for a
given colour component is out of range, it shall be adjusted to the nearest valid value.
Extend
array
(Optional) An array of two boolean values specifying whether to extend the shading beyond the
starting and ending points of the axis, respectively. Default value: [false false].

p92 7.10 Functions

HansRosling.pdf grayscale response  looks wrong.

pre-algebra_translate.pdf QR code is all black after conversion

explosion-database-choice.pdf  encoding.go:270:DecodeStream: Invalid row length...(39195/291=134+201/291

Cirrato Whitepaper.pdf     encoding.go:762:DecodeStream: jpeg.Decode failed. err=invalid JPEG format: bad Th value


SemanticStates201.pdf PdfColorspaceSpecialPattern

PdfColorspaceCalRGB Has 1 component per pixel 2012060560369753.pdf



upload1419017431.pdf Im1 1320 0 obj
<</BitsPerComponent 8/ColorSpace/DeviceCMYK/Filter/DCTDecode/Height 182/Length 41827/Subtype/Image/Type/XObject/Width 166>>
Decompressed jpeg has # pixels = 3 * width * height

327 of 333 `Gradle_Beyond_the_Basics.pdf`  (9892515->10321238 104%) 80 pages 0.396 sec => `output/Gradle_Beyond_the_Basics.pdf`
[ERROR]  pdf_transform_content_streams.go:712:runGhostscript runGhostscript: Could not process pdf="output/Gradle_Beyond_the_Basics.pdf" err=exit status 1
stdout=
{GPL Ghostscript 9.16 (2015-03-30)
Copyright (C) 2015 Artifex Software, Inc.  All rights reserved.
This software comes with NO WARRANTY: see the file PUBLIC for details.
Processing pages 1 through 80.
Page 1
Page 2
Page 3
Page 4
Page 5
Page 6
Page 7
Page 8
Error: /unregistered in --run--
Operand stack:

Execution stack:
   %interp_exit   .runexec2   --nostringval--   --nostringval--   --nostringval--   2   %stopped_push   --nostringval--   --nostringval--   --nostringval--   false   1   %stopped_push   1951   1   3   %oparray_pop   1950   1   3   %oparray_pop   1934   1   3   %oparray_pop   --nostringval--   --nostringval--   9   1   80   --nostringval--   %for_pos_int_continue   --nostringval--   --nostringval--
Dictionary stack:
   --dict:1187/1684(ro)(G)--   --dict:1/20(G)--   --dict:83/200(L)--   --dict:83/200(L)--   --dict:117/127(ro)(G)--   --dict:281/300(ro)(G)--   --dict:28/32(L)--   --dict:6/8(L)--   --dict:23/40(L)--   --dict:1/1(ro)(G)--   --dict:1/1(ro)(G)--   --dict:10/15(L)--   --dict:1/1(ro)(G)--   --dict:1/1(ro)(G)--   --dict:1/1(ro)(G)--   --dict:1/1(ro)(G)--   --dict:4/5(L)--
Current allocation mode is local
Last OS error: Invalid argument
 %!s(int=0)   %!s(bytes.readOp=0)}
stderr=
{   **** Error reading a content stream. The page may be incomplete.
   **** File did not complete the page properly and may be damaged.

   **** File has unbalanced q/Q operators (too many Q's) ****



[ERROR]  pdf_transform_content_streams.go:210:main isPdfColor: 51 Color pages
347 files 0 bad 224 failed
0 bad
224 fail
  0 `/Users/pcadmin/go-work/src/trust-me/testdata/pre-algebra_translate.pdf`
  1 `/Users/pcadmin/go-work/src/trust-me/testdata/Jedi.pdf`
  2 `/Users/pcadmin/go-work/src/trust-me/testdata/49-7441.pdf`
  3 `/Users/pcadmin/go-work/src/trust-me/testdata/A Sunshine - IEP.pdf`
  4 `/Users/pcadmin/go-work/src/trust-me/testdata/android.pdf`
  5 `/Users/pcadmin/go-work/src/trust-me/testdata/10.235.1.2_SMTP_via_LDAP_05-16-2016_09-43-48.pdf.pdf`
  6 `/Users/pcadmin/go-work/src/trust-me/testdata/Whitepaper-XMP-metadata-in-PDFlib-products.pdf`
  7 `/Users/pcadmin/go-work/src/trust-me/testdata/docs.pdf`
  8 `/Users/pcadmin/go-work/src/trust-me/testdata/mh-guide-matrix.pdf`
  9 `/Users/pcadmin/go-work/src/trust-me/testdata/letter to Parents re Telstra Phone Line.pdf`
 10 `/Users/pcadmin/go-work/src/trust-me/testdata/Doc_Test_papercut.pdf`
 11 `/Users/pcadmin/go-work/src/trust-me/testdata/invoice88786339.pdf`
 12 `/Users/pcadmin/go-work/src/trust-me/testdata/HH factsheet.pdf`
 13 `/Users/pcadmin/go-work/src/trust-me/testdata/2017 Arts Festival Letter.pdf`
 14 `/Users/pcadmin/go-work/src/trust-me/testdata/ICT Parent letter 7 2016.pdf`
 15 `/Users/pcadmin/go-work/src/trust-me/testdata/designofasearchablesymmetrickeyciphersystem-120529001108-phpapp02.pdf`
 16 `/Users/pcadmin/go-work/src/trust-me/testdata/twitter_sentimenti.pdf`
 17 `/Users/pcadmin/go-work/src/trust-me/testdata/Zendesk FAQ.pdf`
 18 `/Users/pcadmin/go-work/src/trust-me/testdata/RAND_PE198.pdf`
 19 `/Users/pcadmin/go-work/src/trust-me/testdata/p.pdf`
 20 `/Users/pcadmin/go-work/src/trust-me/testdata/SemanticStates201.pdf`
 21 `/Users/pcadmin/go-work/src/trust-me/testdata/KEI_UT_Raessler.pdf`
 22 `/Users/pcadmin/go-work/src/trust-me/testdata/Term One Sport Junior.pdf`
 23 `/Users/pcadmin/go-work/src/trust-me/testdata/explosion-database-choice.pdf`
 24 `/Users/pcadmin/go-work/src/trust-me/testdata/expressions.pdf`
 25 `/Users/pcadmin/go-work/src/trust-me/testdata/wald2002-2.pdf`
 26 `/Users/pcadmin/go-work/src/trust-me/testdata/epson_pages3_color_pages1.pdf`
 27 `/Users/pcadmin/go-work/src/trust-me/testdata/MFP Scan.pdf`
 28 `/Users/pcadmin/go-work/src/trust-me/testdata/cure53.pdf`
 29 `/Users/pcadmin/go-work/src/trust-me/testdata/frege.pdf`
 30 `/Users/pcadmin/go-work/src/trust-me/testdata/250a5c08cb131c12125c916530830b6d-4836-RJASET-DOI.pdf`
 31 `/Users/pcadmin/go-work/src/trust-me/testdata/manber.pdf`
 32 `/Users/pcadmin/go-work/src/trust-me/testdata/Implementing Gentry’s Fully-Homomorphic Encryption Scheme.pdf`
 33 `/Users/pcadmin/go-work/src/trust-me/testdata/word_jpeg_color.pdf`
 34 `/Users/pcadmin/go-work/src/trust-me/testdata/school.pdf`
 35 `/Users/pcadmin/go-work/src/trust-me/testdata/word_color_png_24bit.pdf`
 36 `/Users/pcadmin/go-work/src/trust-me/testdata/word_color_png_32bit.pdf`
 37 `/Users/pcadmin/go-work/src/trust-me/testdata/PaperCutv15.3internalupdate.pdf`
 38 `/Users/pcadmin/go-work/src/trust-me/testdata/Andrew Fuller flyer.pdf`
 39 `/Users/pcadmin/go-work/src/trust-me/testdata/AcrobatDC_acrobat_digital_signature_appearances.pdf`
 40 `/Users/pcadmin/go-work/src/trust-me/testdata/Cirrato Whitepaper.pdf`
 41 `/Users/pcadmin/go-work/src/trust-me/testdata/0912.0779v1.pdf`
 42 `/Users/pcadmin/go-work/src/trust-me/testdata/Little_niss_oct1014.pdf`
 43 `/Users/pcadmin/go-work/src/trust-me/testdata/508_icmlpaper.pdf`
 44 `/Users/pcadmin/go-work/src/trust-me/testdata/age_text.pdf`
 45 `/Users/pcadmin/go-work/src/trust-me/testdata/geoff_hinton_dark14.pdf`
 46 `/Users/pcadmin/go-work/src/trust-me/testdata/s.pdf`
 47 `/Users/pcadmin/go-work/src/trust-me/testdata/spool-snapshot (2).spl.pdf`
 48 `/Users/pcadmin/go-work/src/trust-me/testdata/3page_test.pdf`
 49 `/Users/pcadmin/go-work/src/trust-me/testdata/2012060560369753.pdf`
 50 `/Users/pcadmin/go-work/src/trust-me/testdata/extreme-days-and-nights-daylight-variation-in-the-arctic-reykjavik-murmansk-and-alert_e634.pdf`
 51 `/Users/pcadmin/go-work/src/trust-me/testdata/Parallel and Dynamic Searchable Symmetric Encryption.pdf`
 52 `/Users/pcadmin/go-work/src/trust-me/testdata/1306.1840v1.pdf`
 53 `/Users/pcadmin/go-work/src/trust-me/testdata/2017 Apps List.pdf`
 54 `/Users/pcadmin/go-work/src/trust-me/testdata/w13229.pdf`
 55 `/Users/pcadmin/go-work/src/trust-me/testdata/pca_text.pdf`
 56 `/Users/pcadmin/go-work/src/trust-me/testdata/raluca-cryptdb.pdf`
 57 `/Users/pcadmin/go-work/src/trust-me/testdata/Presto_UserGuide.pdf`
 58 `/Users/pcadmin/go-work/src/trust-me/testdata/address-change.pdf`
 59 `/Users/pcadmin/go-work/src/trust-me/testdata/ibm12.pdf`
 60 `/Users/pcadmin/go-work/src/trust-me/testdata/science.pdf`
 61 `/Users/pcadmin/go-work/src/trust-me/testdata/integrity.pdf`
 62 `/Users/pcadmin/go-work/src/trust-me/testdata/upload1419017431.pdf`
 63 `/Users/pcadmin/go-work/src/trust-me/testdata/Principl Letter - start of year - 2017.pdf`
 64 `/Users/pcadmin/go-work/src/trust-me/testdata/Artificiell Intelligens och machine learing för sjukvård och life science.pdf`
 65 `/Users/pcadmin/go-work/src/trust-me/testdata/signature.pdf`
 66 `/Users/pcadmin/go-work/src/trust-me/testdata/384.full.pdf`
 67 `/Users/pcadmin/go-work/src/trust-me/testdata/PauloCarrilloPenaCameraReadyJuly132016.pdf`
 68 `/Users/pcadmin/go-work/src/trust-me/testdata/dna - Wolfram_Alpha.pdf`
 69 `/Users/pcadmin/go-work/src/trust-me/testdata/1702.06959.pdf`
 70 `/Users/pcadmin/go-work/src/trust-me/testdata/journal.pone.0048386.pdf`
 71 `/Users/pcadmin/go-work/src/trust-me/testdata/shattered-1.pdf`
 72 `/Users/pcadmin/go-work/src/trust-me/testdata/shattered-2.pdf`
 73 `/Users/pcadmin/go-work/src/trust-me/testdata/DeBruijnPrimer.pdf`
 74 `/Users/pcadmin/go-work/src/trust-me/testdata/Blackburn%2B2016TOPLAS.pdf`
 75 `/Users/pcadmin/go-work/src/trust-me/testdata/cs-pwgmsn20-20130328-5101.1.pdf`
 76 `/Users/pcadmin/go-work/src/trust-me/testdata/cs.pdf`
 77 `/Users/pcadmin/go-work/src/trust-me/testdata/20122845.full.pdf`
 78 `/Users/pcadmin/go-work/src/trust-me/testdata/spanner-osdi2012.pdf`
 79 `/Users/pcadmin/go-work/src/trust-me/testdata/power_law_bins.pdf`
 80 `/Users/pcadmin/go-work/src/trust-me/testdata/VerifiableDataStructures.pdf`
 81 `/Users/pcadmin/go-work/src/trust-me/testdata/348_icmlpaper.pdf`
 82 `/Users/pcadmin/go-work/src/trust-me/testdata/2568.pdf`
 83 `/Users/pcadmin/go-work/src/trust-me/testdata/Politics_of_Voter_Fraud_Final.pdf`
 84 `/Users/pcadmin/go-work/src/trust-me/testdata/pasted-image-5398.pdf`
 85 `/Users/pcadmin/go-work/src/trust-me/testdata/climate_science_words.pdf`
 86 `/Users/pcadmin/go-work/src/trust-me/testdata/1409.6070v1.pdf`
 87 `/Users/pcadmin/go-work/src/trust-me/testdata/SK Hynix Case Study November 16 2015 - Final.pdf`
 88 `/Users/pcadmin/go-work/src/trust-me/testdata/dev_plan.pdf`
 89 `/Users/pcadmin/go-work/src/trust-me/testdata/Pfeffer_Nonprofit-Management-Institute-2016.pdf`
 90 `/Users/pcadmin/go-work/src/trust-me/testdata/1611.01989.pdf`
 91 `/Users/pcadmin/go-work/src/trust-me/testdata/fernbach13.pdf`
 92 `/Users/pcadmin/go-work/src/trust-me/testdata/l09r01.pdf`
 93 `/Users/pcadmin/go-work/src/trust-me/testdata/sa.pdf`
 94 `/Users/pcadmin/go-work/src/trust-me/testdata/Bez_test.pdf`
 95 `/Users/pcadmin/go-work/src/trust-me/testdata/HyperLogLog.pdf`
 96 `/Users/pcadmin/go-work/src/trust-me/testdata/IJETAE_1012_47.pdf`
 97 `/Users/pcadmin/go-work/src/trust-me/testdata/PaperCut-Professional-Services-Guide.pdf`
 98 `/Users/pcadmin/go-work/src/trust-me/testdata/shattered.pdf`
 99 `/Users/pcadmin/go-work/src/trust-me/testdata/legislative-process-poster-espanol.pdf`
100 `/Users/pcadmin/go-work/src/trust-me/testdata/searchable-symmetric-encryption-elcrypt.pdf`
101 `/Users/pcadmin/go-work/src/trust-me/testdata/14(1).pdf`
102 `/Users/pcadmin/go-work/src/trust-me/testdata/Deep learning for image denoising.pdf`
103 `/Users/pcadmin/go-work/src/trust-me/testdata/PaperCut MF - MFD Integration Matrix.pdf`
104 `/Users/pcadmin/go-work/src/trust-me/testdata/1206.4637.pdf`
105 `/Users/pcadmin/go-work/src/trust-me/testdata/How-to-sign-a-PDF-form-using-a-digital-signature.pdf`
106 `/Users/pcadmin/go-work/src/trust-me/testdata/tout.pdf`
107 `/Users/pcadmin/go-work/src/trust-me/testdata/Updated Bring Your Own Technology 2017 Letter.pdf`
108 `/Users/pcadmin/go-work/src/trust-me/testdata/1410.4615v3.pdf`
109 `/Users/pcadmin/go-work/src/trust-me/testdata/nips2010.pdf`
110 `/Users/pcadmin/go-work/src/trust-me/testdata/joi160132.pdf`
111 `/Users/pcadmin/go-work/src/trust-me/testdata/Lise_GenomeLett2002.pdf`
112 `/Users/pcadmin/go-work/src/trust-me/testdata/CRG Leaders or Leaderfulness Digital 6_16_16.pdf`
113 `/Users/pcadmin/go-work/src/trust-me/testdata/7Vol30No1.pdf`
114 `/Users/pcadmin/go-work/src/trust-me/testdata/RootstockWhitePaperv9-Overview.pdf`
115 `/Users/pcadmin/go-work/src/trust-me/testdata/cosy_small.pdf`
116 `/Users/pcadmin/go-work/src/trust-me/testdata/1610.02306v1.pdf`
117 `/Users/pcadmin/go-work/src/trust-me/testdata/1306.0239v1.pdf`
118 `/Users/pcadmin/go-work/src/trust-me/testdata/BSA Fiery Direct Print Issue.pdf`
119 `/Users/pcadmin/go-work/src/trust-me/testdata/004635142b2089f675000000.pdf`
120 `/Users/pcadmin/go-work/src/trust-me/testdata/conditonal_random_fields.pdf`
121 `/Users/pcadmin/go-work/src/trust-me/testdata/gpads.pdf`
122 `/Users/pcadmin/go-work/src/trust-me/testdata/privacy-fact-sheet-17-australian-privacy-principles_2.pdf`
123 `/Users/pcadmin/go-work/src/trust-me/testdata/LATTICE-BASED CRYPTOGRAPHY.pdf`
124 `/Users/pcadmin/go-work/src/trust-me/testdata/sciencetopics.pdf`
125 `/Users/pcadmin/go-work/src/trust-me/testdata/2016095pap.pdf`
126 `/Users/pcadmin/go-work/src/trust-me/testdata/sugihara-causality-science-2012.pdf`
127 `/Users/pcadmin/go-work/src/trust-me/testdata/twitter.pdf`
128 `/Users/pcadmin/go-work/src/trust-me/testdata/Meissen_MQP2.pdf`
129 `/Users/pcadmin/go-work/src/trust-me/testdata/WeeklyBulletin (1).pdf`
130 `/Users/pcadmin/go-work/src/trust-me/testdata/p253-porter.pdf`
131 `/Users/pcadmin/go-work/src/trust-me/testdata/sharp_brochure.pdf`
132 `/Users/pcadmin/go-work/src/trust-me/testdata/pearson_science_8_sb_chapter_5_unit_5.2.pdf`
133 `/Users/pcadmin/go-work/src/trust-me/testdata/CCM sandbox issue reproduction.pdf`
134 `/Users/pcadmin/go-work/src/trust-me/testdata/Ball_Drop_activity.pdf`
135 `/Users/pcadmin/go-work/src/trust-me/testdata/Privet201.pdf`
136 `/Users/pcadmin/go-work/src/trust-me/testdata/threshold_sigs.pdf`
137 `/Users/pcadmin/go-work/src/trust-me/testdata/PrinterOn Enterprise & PaperCut Integration.pdf`
138 `/Users/pcadmin/go-work/src/trust-me/testdata/zerocash-extended-20140518.pdf`
139 `/Users/pcadmin/go-work/src/trust-me/testdata/srep40678.pdf`
140 `/Users/pcadmin/go-work/src/trust-me/testdata/protocol.pdf`
141 `/Users/pcadmin/go-work/src/trust-me/testdata/whitepaper2015.pdf`
142 `/Users/pcadmin/go-work/src/trust-me/testdata/1507.02672v1.pdf`
143 `/Users/pcadmin/go-work/src/trust-me/testdata/cryptdb.pdf`
144 `/Users/pcadmin/go-work/src/trust-me/testdata/PNAS-2016-Battiston-10031-6.pdf`
145 `/Users/pcadmin/go-work/src/trust-me/testdata/mueller.pdf`
146 `/Users/pcadmin/go-work/src/trust-me/testdata/a0w20000000dikuAAA.pdf`
147 `/Users/pcadmin/go-work/src/trust-me/testdata/2013-12-12_rg_final_report.pdf`
148 `/Users/pcadmin/go-work/src/trust-me/testdata/Implementing Homomorphic Encryption.pdf`
149 `/Users/pcadmin/go-work/src/trust-me/testdata/HansRosling.pdf`
150 `/Users/pcadmin/go-work/src/trust-me/testdata/Read-the-draft-of-the-executive-order-on-CIA.pdf`
151 `/Users/pcadmin/go-work/src/trust-me/testdata/Longevity Bulletin Issue 9.pdf`
152 `/Users/pcadmin/go-work/src/trust-me/testdata/stick_breaking.pdf`
153 `/Users/pcadmin/go-work/src/trust-me/testdata/unsupervised.pdf`
154 `/Users/pcadmin/go-work/src/trust-me/testdata/weakLongrangeForcesTheory_01.pdf`
155 `/Users/pcadmin/go-work/src/trust-me/testdata/pcl_xl_2_0_technical_reference_rev2_2.pdf`
156 `/Users/pcadmin/go-work/src/trust-me/testdata/pnas.201308477.pdf`
157 `/Users/pcadmin/go-work/src/trust-me/testdata/Lattice Based Cryptography for Beginners.pdf`
158 `/Users/pcadmin/go-work/src/trust-me/testdata/text_13pages.pdf`
159 `/Users/pcadmin/go-work/src/trust-me/testdata/gap.pdf`
160 `/Users/pcadmin/go-work/src/trust-me/testdata/Parsing-Probabilistic.pdf`
161 `/Users/pcadmin/go-work/src/trust-me/testdata/cosy.pdf`
162 `/Users/pcadmin/go-work/src/trust-me/testdata/12 Image reconstruction.pdf`
163 `/Users/pcadmin/go-work/src/trust-me/testdata/The_State_of_Late_Payment_MarketInvoice_2016.pdf`
164 `/Users/pcadmin/go-work/src/trust-me/testdata/Framingham Heart Study.pdf`
165 `/Users/pcadmin/go-work/src/trust-me/testdata/paper.pdf`
166 `/Users/pcadmin/go-work/src/trust-me/testdata/word_color_gif.pdf`
167 `/Users/pcadmin/go-work/src/trust-me/testdata/Redmon_You_Only_Look_CVPR_2016_paper.pdf`
168 `/Users/pcadmin/go-work/src/trust-me/testdata/Szegedy_Going_Deeper_With_2015_CVPR_paper.pdf`
169 `/Users/pcadmin/go-work/src/trust-me/testdata/RMG_Factsheet_2017.pdf`
170 `/Users/pcadmin/go-work/src/trust-me/testdata/cs-ippprodprint10-20010212-5100.3.pdf`
171 `/Users/pcadmin/go-work/src/trust-me/testdata/imagenet.pdf`
172 `/Users/pcadmin/go-work/src/trust-me/testdata/4824-imagenet-classification-with-deep-convolutional-neural-networks.pdf`
173 `/Users/pcadmin/go-work/src/trust-me/testdata/02-Lexing.pdf`
174 `/Users/pcadmin/go-work/src/trust-me/testdata/PhysRevLett.118.060401.pdf`
175 `/Users/pcadmin/go-work/src/trust-me/testdata/fgvc-2015-fast-bird-part.pdf`
176 `/Users/pcadmin/go-work/src/trust-me/testdata/1603.09056.pdf`
177 `/Users/pcadmin/go-work/src/trust-me/testdata/ICA_2017_01.pdf`
178 `/Users/pcadmin/go-work/src/trust-me/testdata/p329-galle.pdf`
179 `/Users/pcadmin/go-work/src/trust-me/testdata/1207.0580.pdf`
180 `/Users/pcadmin/go-work/src/trust-me/testdata/1701.02434v1.pdf`
181 `/Users/pcadmin/go-work/src/trust-me/testdata/art%3A10.1186%2Fs13673-015-0039-9.pdf`
182 `/Users/pcadmin/go-work/src/trust-me/testdata/sparse3d.pdf`
183 `/Users/pcadmin/go-work/src/trust-me/testdata/BALU_Raghavendran.pdf`
184 `/Users/pcadmin/go-work/src/trust-me/testdata/US5553145.pdf`
185 `/Users/pcadmin/go-work/src/trust-me/testdata/PaperCut Case Study - Bakers Delight.pdf`
186 `/Users/pcadmin/go-work/src/trust-me/testdata/1502.04623.pdf`
187 `/Users/pcadmin/go-work/src/trust-me/testdata/MondayAM.pdf`
188 `/Users/pcadmin/go-work/src/trust-me/testdata/BookClubGuide16.pdf`
189 `/Users/pcadmin/go-work/src/trust-me/testdata/PaperCut Device REST Web Service API.pdf`
190 `/Users/pcadmin/go-work/src/trust-me/testdata/The_Block_Cipher_Companion.pdf`
191 `/Users/pcadmin/go-work/src/trust-me/testdata/fireflyforparents.pdf`
192 `/Users/pcadmin/go-work/src/trust-me/testdata/94-487-1-PB.pdf`
193 `/Users/pcadmin/go-work/src/trust-me/testdata/Zero_Configuration_Networking_The_Definitive_Guide.pdf`
194 `/Users/pcadmin/go-work/src/trust-me/testdata/Afternoon Need States Report_FINAL 6.25.15.pdf`
195 `/Users/pcadmin/go-work/src/trust-me/testdata/dark-internet-mail-environment-march-2015.pdf`
196 `/Users/pcadmin/go-work/src/trust-me/testdata/day2_ListeningToHypercolors.pdf`
197 `/Users/pcadmin/go-work/src/trust-me/testdata/IPC3661.pdf`
198 `/Users/pcadmin/go-work/src/trust-me/testdata/PaperCutMF-Top-10-Reasons.pdf`
199 `/Users/pcadmin/go-work/src/trust-me/testdata/1601.06759v2.pdf`
200 `/Users/pcadmin/go-work/src/trust-me/testdata/2015-09-16-T23-39-51_ec2-user_ip-172-31-6-72_jim.pdf`
201 `/Users/pcadmin/go-work/src/trust-me/testdata/Hierarchical Detection of Hard Exudates.pdf`
202 `/Users/pcadmin/go-work/src/trust-me/testdata/rigamonti_cvpr13.pdf`
203 `/Users/pcadmin/go-work/src/trust-me/testdata/05-20.pdf`
204 `/Users/pcadmin/go-work/src/trust-me/testdata/OpenCV Computer Vision with Python.pdf`
205 `/Users/pcadmin/go-work/src/trust-me/testdata/The great chain of being sure about things _ The Economist.pdf`
206 `/Users/pcadmin/go-work/src/trust-me/testdata/t1_week_6_3rd_march_2016.pdf`
207 `/Users/pcadmin/go-work/src/trust-me/testdata/1701.07164v1.pdf`
208 `/Users/pcadmin/go-work/src/trust-me/testdata/WhatIsEnergy.pdf`
209 `/Users/pcadmin/go-work/src/trust-me/testdata/Physics_Sample_Chapter_3.pdf`
210 `/Users/pcadmin/go-work/src/trust-me/testdata/hypercolumn.pdf`
211 `/Users/pcadmin/go-work/src/trust-me/testdata/day3_TemporalImageProcessing.pdf`
212 `/Users/pcadmin/go-work/src/trust-me/testdata/PracticalPythonAndOpenCV_Chapter10.pdf`
213 `/Users/pcadmin/go-work/src/trust-me/testdata/augmentingReality_salon01.pdf`
214 `/Users/pcadmin/go-work/src/trust-me/testdata/Spatial Transformer Networks.pdf`
215 `/Users/pcadmin/go-work/src/trust-me/testdata/Lesson_054_handout.pdf`
216 `/Users/pcadmin/go-work/src/trust-me/testdata/1701.07875v1.pdf`
217 `/Users/pcadmin/go-work/src/trust-me/testdata/Gradle_Beyond_the_Basics.pdf`
218 `/Users/pcadmin/go-work/src/trust-me/testdata/CAFR_2012.pdf`
219 `/Users/pcadmin/go-work/src/trust-me/testdata/Debug.Hacks中文版_深入调试的技术和工具.pdf`
220 `/Users/pcadmin/go-work/src/trust-me/testdata/faa4bd35a36296305ce7e4cbbbb9d7f0a909de18.pdf`
221 `/Users/pcadmin/go-work/src/trust-me/testdata/Trump-Intelligence-Allegations.pdf`
222 `/Users/pcadmin/go-work/src/trust-me/testdata/talk_Simons_part1_pdf.pdf`
223 `/Users/pcadmin/go-work/src/trust-me/testdata/nips-tutorial-policy-optimization-Schulman-Abbeel.pdf`



67 fail
  0 `/Users/pcadmin/go-work/src/trust-me/testdata/pearson_science_8_sb_chapter_5_unit_5.2.pdf`
  1 `/Users/pcadmin/go-work/src/trust-me/testdata/Privet201.pdf`
  2 `/Users/pcadmin/go-work/src/trust-me/testdata/zerocash-extended-20140518.pdf`
  3 `/Users/pcadmin/go-work/src/trust-me/testdata/srep40678.pdf`
  4 `/Users/pcadmin/go-work/src/trust-me/testdata/protocol.pdf`
  5 `/Users/pcadmin/go-work/src/trust-me/testdata/whitepaper2015.pdf`
  6 `/Users/pcadmin/go-work/src/trust-me/testdata/1507.02672v1.pdf`
  7 `/Users/pcadmin/go-work/src/trust-me/testdata/cryptdb.pdf`
  8 `/Users/pcadmin/go-work/src/trust-me/testdata/mueller.pdf`
  9 `/Users/pcadmin/go-work/src/trust-me/testdata/a0w20000000dikuAAA.pdf`
 10 `/Users/pcadmin/go-work/src/trust-me/testdata/2013-12-12_rg_final_report.pdf`
 11 `/Users/pcadmin/go-work/src/trust-me/testdata/Implementing Homomorphic Encryption.pdf`
 12 `/Users/pcadmin/go-work/src/trust-me/testdata/HansRosling.pdf`
 13 `/Users/pcadmin/go-work/src/trust-me/testdata/Longevity Bulletin Issue 9.pdf`
 14 `/Users/pcadmin/go-work/src/trust-me/testdata/stick_breaking.pdf`
 15 `/Users/pcadmin/go-work/src/trust-me/testdata/unsupervised.pdf`
 16 `/Users/pcadmin/go-work/src/trust-me/testdata/pcl_xl_2_0_technical_reference_rev2_2.pdf`
 17 `/Users/pcadmin/go-work/src/trust-me/testdata/pnas.201308477.pdf`
 18 `/Users/pcadmin/go-work/src/trust-me/testdata/Lattice Based Cryptography for Beginners.pdf`
 19 `/Users/pcadmin/go-work/src/trust-me/testdata/12 Image reconstruction.pdf`
 20 `/Users/pcadmin/go-work/src/trust-me/testdata/The_State_of_Late_Payment_MarketInvoice_2016.pdf`
 21 `/Users/pcadmin/go-work/src/trust-me/testdata/paper.pdf`
 22 `/Users/pcadmin/go-work/src/trust-me/testdata/word_color_gif.pdf`
 23 `/Users/pcadmin/go-work/src/trust-me/testdata/Redmon_You_Only_Look_CVPR_2016_paper.pdf`
 24 `/Users/pcadmin/go-work/src/trust-me/testdata/Szegedy_Going_Deeper_With_2015_CVPR_paper.pdf`
 25 `/Users/pcadmin/go-work/src/trust-me/testdata/RMG_Factsheet_2017.pdf`
 26 `/Users/pcadmin/go-work/src/trust-me/testdata/imagenet.pdf`
 27 `/Users/pcadmin/go-work/src/trust-me/testdata/4824-imagenet-classification-with-deep-convolutional-neural-networks.pdf`
 28 `/Users/pcadmin/go-work/src/trust-me/testdata/02-Lexing.pdf`
 29 `/Users/pcadmin/go-work/src/trust-me/testdata/PhysRevLett.118.060401.pdf`
 30 `/Users/pcadmin/go-work/src/trust-me/testdata/1603.09056.pdf`
 31 `/Users/pcadmin/go-work/src/trust-me/testdata/ICA_2017_01.pdf`
 32 `/Users/pcadmin/go-work/src/trust-me/testdata/1701.02434v1.pdf`
 33 `/Users/pcadmin/go-work/src/trust-me/testdata/sparse3d.pdf`
 34 `/Users/pcadmin/go-work/src/trust-me/testdata/3725b275.pdf`
 35 `/Users/pcadmin/go-work/src/trust-me/testdata/BALU_Raghavendran.pdf`
 36 `/Users/pcadmin/go-work/src/trust-me/testdata/PaperCut Case Study - Bakers Delight.pdf`
 37 `/Users/pcadmin/go-work/src/trust-me/testdata/1502.04623.pdf`
 38 `/Users/pcadmin/go-work/src/trust-me/testdata/MondayAM.pdf`
 39 `/Users/pcadmin/go-work/src/trust-me/testdata/BookClubGuide16.pdf`
 40 `/Users/pcadmin/go-work/src/trust-me/testdata/The_Block_Cipher_Companion.pdf`
 41 `/Users/pcadmin/go-work/src/trust-me/testdata/94-487-1-PB.pdf`
 42 `/Users/pcadmin/go-work/src/trust-me/testdata/Zero_Configuration_Networking_The_Definitive_Guide.pdf`
 43 `/Users/pcadmin/go-work/src/trust-me/testdata/Afternoon Need States Report_FINAL 6.25.15.pdf`
 44 `/Users/pcadmin/go-work/src/trust-me/testdata/dark-internet-mail-environment-march-2015.pdf`
 45 `/Users/pcadmin/go-work/src/trust-me/testdata/IPC3661.pdf`
 46 `/Users/pcadmin/go-work/src/trust-me/testdata/PaperCutMF-Top-10-Reasons.pdf`
 47 `/Users/pcadmin/go-work/src/trust-me/testdata/1601.06759v2.pdf`
 48 `/Users/pcadmin/go-work/src/trust-me/testdata/2015-09-16-T23-39-51_ec2-user_ip-172-31-6-72_jim.pdf`
 49 `/Users/pcadmin/go-work/src/trust-me/testdata/Hierarchical Detection of Hard Exudates.pdf`
 50 `/Users/pcadmin/go-work/src/trust-me/testdata/rigamonti_cvpr13.pdf`
 51 `/Users/pcadmin/go-work/src/trust-me/testdata/05-20.pdf`
 52 `/Users/pcadmin/go-work/src/trust-me/testdata/The great chain of being sure about things _ The Economist.pdf`
 53 `/Users/pcadmin/go-work/src/trust-me/testdata/t1_week_6_3rd_march_2016.pdf`
 54 `/Users/pcadmin/go-work/src/trust-me/testdata/1701.07164v1.pdf`
 55 `/Users/pcadmin/go-work/src/trust-me/testdata/WhatIsEnergy.pdf`
 56 `/Users/pcadmin/go-work/src/trust-me/testdata/Physics_Sample_Chapter_3.pdf`
 57 `/Users/pcadmin/go-work/src/trust-me/testdata/hypercolumn.pdf`
 58 `/Users/pcadmin/go-work/src/trust-me/testdata/day3_TemporalImageProcessing.pdf`
 59 `/Users/pcadmin/go-work/src/trust-me/testdata/Spatial Transformer Networks.pdf`
 60 `/Users/pcadmin/go-work/src/trust-me/testdata/Lesson_054_handout.pdf`
 61 `/Users/pcadmin/go-work/src/trust-me/testdata/1701.07875v1.pdf`
 62 `/Users/pcadmin/go-work/src/trust-me/testdata/Gradle_Beyond_the_Basics.pdf`
 63 `/Users/pcadmin/go-work/src/trust-me/testdata/CAFR_2012.pdf`
 64 `/Users/pcadmin/go-work/src/trust-me/testdata/faa4bd35a36296305ce7e4cbbbb9d7f0a909de18.pdf`
 65 `/Users/pcadmin/go-work/src/trust-me/testdata/talk_Simons_part1_pdf.pdf`
 66 `/Users/pcadmin/go-work/src/trust-me/testdata/nips-tutorial-policy-optimization-Schulman-Abbeel.pdf`
