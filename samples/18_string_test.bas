10 REM String Functions Test
20 PRINT "=== String Functions Test ==="
30 PRINT
40 REM Test string concatenation
50 PRINT "1. String Concatenation (+)"
60 A$ = "Hello"
70 B$ = "World"
80 C$ = A$ + " " + B$
90 PRINT "A$ = "; A$
100 PRINT "B$ = "; B$
110 PRINT "A$ + B$ = "; C$
120 PRINT
130 REM Test number to string concatenation
140 PRINT "2. Number + String Concatenation"
150 D$ = "Number: " + 42
160 PRINT "Result: "; D$
170 PRINT
180 REM Test string comparison
190 PRINT "3. String Comparison"
200 X$ = "Apple"
210 Y$ = "Banana"
220 IF X$ = Y$ THEN PRINT "X$ = Y$ is True" ELSE PRINT "X$ = Y$ is False"
230 IF X$ <> Y$ THEN PRINT "X$ <> Y$ is True" ELSE PRINT "X$ <> Y$ is False"
240 IF X$ < Y$ THEN PRINT "X$ < Y$ is True" ELSE PRINT "X$ < Y$ is False"
250 IF X$ > Y$ THEN PRINT "X$ > Y$ is True" ELSE PRINT "X$ > Y$ is False"
260 PRINT
270 REM Test LEN function
280 PRINT "4. LEN Function"
290 S$ = "Hello"
300 PRINT "LEN(S$) = "; LEN(S$)
310 S2$ = ""
320 PRINT "LEN(empty) = "; LEN(S2$)
330 S3$ = "BASIC Programming"
340 PRINT "LEN(S3$) = "; LEN(S3$)
350 PRINT
360 REM Test LEFT$ function
370 PRINT "5. LEFT$ Function"
380 T$ = "Hello World"
390 PRINT "Text = "; T$
400 PRINT "LEFT$(T$, 5) = "; LEFT$(T$, 5)
410 PRINT "LEFT$(T$, 10) = "; LEFT$(T$, 10)
420 PRINT
430 REM Test RIGHT$ function
440 PRINT "6. RIGHT$ Function"
450 PRINT "RIGHT$(T$, 5) = "; RIGHT$(T$, 5)
460 PRINT "RIGHT$(T$, 10) = "; RIGHT$(T$, 10)
470 PRINT
480 REM Test MID$ function
490 PRINT "7. MID$ Function"
500 PRINT "MID$(T$, 7, 5) = "; MID$(T$, 7, 5)
510 PRINT "MID$(T$, 1, 5) = "; MID$(T$, 1, 5)
520 PRINT "MID$(T$, 7) = "; MID$(T$, 7)
530 PRINT "MID$(T$, 1) = "; MID$(T$, 1)
540 PRINT
550 REM Test INSTR function
560 PRINT "8. INSTR Function"
570 STR$ = "Hello World"
580 SEARCH1$ = "World"
590 SEARCH2$ = "Hello"
600 SEARCH3$ = "xyz"
610 PRINT "STR$ = "; STR$
620 PRINT "INSTR(STR$, SEARCH1$) = "; INSTR(STR$, SEARCH1$)
630 PRINT "INSTR(STR$, SEARCH2$) = "; INSTR(STR$, SEARCH2$)
640 PRINT "INSTR(STR$, SEARCH3$) = "; INSTR(STR$, SEARCH3$)
650 PRINT "INSTR(7, STR$, SEARCH2$) = "; INSTR(7, STR$, SEARCH2$)
660 PRINT
670 REM Test UCASE$ and LCASE$ functions
680 PRINT "9. UCASE$ and LCASE$ Functions"
690 M$ = "Hello World"
700 PRINT "Mixed = "; M$
710 PRINT "UCASE$(M$) = "; UCASE$(M$)
720 PRINT "LCASE$(M$) = "; LCASE$(M$)
730 PRINT
740 REM Test SPACE$ function
750 PRINT "10. SPACE$ Function"
760 PRINT "SPACE$(5) test: "; "A" + SPACE$(3) + "B"
770 PRINT
780 REM Test CHR$ and ASC functions
790 PRINT "11. CHR$ and ASC Functions"
800 PRINT "CHR$(65) = "; CHR$(65)
810 PRINT "CHR$(97) = "; CHR$(97)
820 C1$ = "A"
830 C2$ = "a"
840 PRINT "ASC(C1$) = "; ASC(C1$)
850 PRINT "ASC(C2$) = "; ASC(C2$)
860 PRINT "ASC(CHR$(72)) = "; ASC(CHR$(72))
870 PRINT
880 REM Complex example
890 PRINT "12. Complex Example"
900 FULLNAME$ = "John Doe"
910 PRINT "Full Name: "; FULLNAME$
920 POS = INSTR(FULLNAME$, " ")
930 FIRST$ = LEFT$(FULLNAME$, POS - 1)
940 LAST$ = MID$(FULLNAME$, POS + 1)
950 PRINT "First Name: "; FIRST$
960 PRINT "Last Name: "; LAST$
970 PRINT
980 END
