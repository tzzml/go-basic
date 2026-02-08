10 REM ============================================================
20 REM Comment Style Examples
30 REM Demonstrating both REM and single-quote comments in zork-basic
40 REM ============================================================
50
60 REM ------------------------------------------------------------
70 REM 1. REM Comments (Traditional BASIC style)
80 REM ------------------------------------------------------------
90
100 REM Full line REM comment
110 PRINT "1. REM Comments"
120
130 REM Variable assignment with REM comment after colon
140 A = 10: REM Assign value to A
150 B = 20: REM Assign value to B
160
170 PRINT "   A ="; A
180 PRINT "   B ="; B
190 PRINT
200
210 REM ------------------------------------------------------------
220 REM 2. Single-Quote Comments (GW-BASIC/QuickBASIC style)
230 REM ------------------------------------------------------------
240
250 ' Full line single-quote comment
260 PRINT "2. Single-Quote Comments"
270
280 ' Variable assignment with single-quote comment after colon
290 C = 30: ' Assign value to C
300 D = 40: ' Assign value to D
310
320 PRINT "   C ="; C
330 PRINT "   D ="; D
340 PRINT
350
360 REM ------------------------------------------------------------
370 REM 3. Mixed Styles
380 REM ------------------------------------------------------------
390
400 PRINT "3. Mixed Comment Styles"
410 TOTAL = A + B + C + D: ' Calculate total
420 PRINT "   Total ="; TOTAL: REM Display result
430 PRINT
440
450 REM ------------------------------------------------------------
460 REM 4. Commenting Out Code
470 REM ------------------------------------------------------------
480
490 PRINT "4. Commenting Out Code"
500
510 ' The following lines are commented out:
520 ' PRINT "This will not be printed"
530 REM PRINT "Neither will this"
540
550 PRINT "   (The above lines were commented out)"
560 PRINT
570
580 REM ------------------------------------------------------------
590 REM Best Practices
600 REM ------------------------------------------------------------
610
620 PRINT "5. Best Practices:"
630 PRINT "   - Use REM for traditional BASIC compatibility"
640 PRINT "   - Use ' for shorter, cleaner comments"
650 PRINT "   - Both styles work identically"
660 PRINT "   - Use colon : to separate comments from code"
670 PRINT
680
690 REM Program complete
700 ' END statement is optional at the end
710 END
