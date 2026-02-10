10 PRINT "Hello World from VM"
20 LET A = 10
30 LET B = 20
40 PRINT "A + B =", A + B
50 IF A < B THEN GOSUB 200
60 IF A >= B THEN GOTO 80
70 PRINT "Skipping this line"
80 PRINT "Back from IF Check"
90 END
200 PRINT "Inside Subroutine: A is less than B"
210 RETURN
