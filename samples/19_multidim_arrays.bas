10 REM Multi-dimensional Arrays Test
20 PRINT "=== Multi-dimensional Arrays Test ==="
30 PRINT
40 REM Test 1D array
50 PRINT "1. One-dimensional Array"
60 DIM A(5)
70 A(0) = 10
80 A(1) = 20
90 A(2) = 30
100 A(3) = 40
110 A(4) = 50
120 FOR I = 0 TO 4
130   PRINT "A("; I; ") = "; A(I)
140 NEXT I
150 PRINT
160 REM Test 2D array
170 PRINT "2. Two-dimensional Array"
180 DIM B(3, 4)
190 REM Initialize the 2D array
200 FOR I = 0 TO 2
210   FOR J = 0 TO 3
220     B(I, J) = I * 10 + J
230   NEXT J
240 NEXT I
250 REM Print the 2D array
260 FOR I = 0 TO 2
270   FOR J = 0 TO 3
280     PRINT "B("; I; ", "; J; ") = "; B(I, J)
290   NEXT J
300 NEXT I
310 PRINT
320 REM Test 3D array
330 PRINT "3. Three-dimensional Array"
340 DIM C(2, 2, 2)
350 REM Initialize and print 3D array
360 FOR I = 0 TO 1
370   FOR J = 0 TO 1
380     FOR K = 0 TO 1
390       C(I, J, K) = I * 100 + J * 10 + K
400       PRINT "C("; I; ", "; J; ", "; K; ") = "; C(I, J, K)
410     NEXT K
420   NEXT J
430 NEXT I
440 PRINT
450 REM Test array access with variables
460 PRINT "4. Array Access with Variables"
470 DIM D(4, 4)
480 X = 2
490 Y = 3
500 D(X, Y) = 999
510 PRINT "D(2, 3) = "; D(2, 3)
520 PRINT "D(X, Y) = "; D(X, Y)
530 PRINT
540 REM Test large 2D array
550 PRINT "5. Larger 2D Array (Matrix Multiplication Example)"
560 DIM M1(2, 3)
570 DIM M2(3, 2)
580 DIM RESULT(2, 2)
590 REM Initialize M1
600 M1(0, 0) = 1: M1(0, 1) = 2: M1(0, 2) = 3
610 M1(1, 0) = 4: M1(1, 1) = 5: M1(1, 2) = 6
620 REM Initialize M2
630 M2(0, 0) = 7: M2(0, 1) = 8
640 M2(1, 0) = 9: M2(1, 1) = 10
650 M2(2, 0) = 11: M2(2, 1) = 12
660 REM Matrix multiplication
670 FOR I = 0 TO 1
680   FOR J = 0 TO 1
690     SUM = 0
700     FOR K = 0 TO 2
710       SUM = SUM + M1(I, K) * M2(K, J)
720     NEXT K
730     RESULT(I, J) = SUM
740   NEXT J
750 NEXT I
760 REM Print result
770 PRINT "Matrix Multiplication Result:"
780 FOR I = 0 TO 1
790   FOR J = 0 TO 1
800     PRINT "RESULT("; I; ", "; J; ") = "; RESULT(I, J)
810   NEXT J
820 NEXT I
830 PRINT
840 END
