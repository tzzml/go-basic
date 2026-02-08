10 REM BASIC 运算符示例
20 REM 演示 MOD 和 NOT 运算符的使用
30 PRINT
40 PRINT "=== BASIC 运算符示例 ==="
50 PRINT
60
70 REM ============================================================
80 REM MOD 取模运算符
90 REM ============================================================
100 PRINT "1. MOD 取模运算符"
110 PRINT
120 PRINT "基本用法:"
130 PRINT "  10 MOD 3 ="; 10 MOD 3
140 PRINT "  15 MOD 4 ="; 15 MOD 4
150 PRINT "  100 MOD 7 ="; 100 MOD 7
160 PRINT
170 PRINT "浮点数取模:"
180 PRINT "  20.5 MOD 3 ="; 20.5 MOD 3
190 PRINT
200 PRINT "应用：判断奇偶数"
210 FOR I = 1 TO 10
220   IF (I MOD 2) = 0 THEN PRINT "  "; I; " 是偶数" ELSE PRINT "  "; I; " 是奇数"
230 NEXT I
240 PRINT
250 PRINT
260
270 REM ============================================================
280 REM NOT 逻辑运算符
290 REM ============================================================
300 PRINT "2. NOT 逻辑运算符"
310 PRINT
320 A = 5
330 B = 15
340 C = 10
350 PRINT "A ="; A; ", B ="; B; ", C ="; C
360 PRINT
370 PRINT "基本用法:"
380 PRINT "  NOT (A > 10) ="; NOT (A > 10)
390 PRINT "  NOT (B > 10) ="; NOT (B > 10)
400 PRINT
410 PRINT "条件语句中的使用:"
420 IF NOT (A > 10) THEN PRINT "  A 不大于 10"
430 IF NOT (B > 10) THEN PRINT "  B 不大于 10" ELSE PRINT "  B 大于 10"
440 PRINT
450 PRINT "复杂逻辑:"
460 PRINT "  NOT (A > C AND B > 20) ="; NOT (A > C AND B > 20)
470 PRINT "  NOT (A > B OR C > 20) ="; NOT (A > B OR C > 20)
480 PRINT
490 PRINT
500
510 REM ============================================================
520 REM 综合应用
530 REM ============================================================
540 PRINT "3. 综合应用：闰年判断"
550 PRINT
560 REM 判断闰年的规则：
570 REM 能被4整除但不能被100整除，或者能被400整除
580 FOR YEAR = 2000 TO 2024 STEP 4
590   ISLEAP = (YEAR MOD 4 = 0 AND YEAR MOD 100 <> 0) OR (YEAR MOD 400 = 0)
600   IF ISLEAP THEN PRINT "  "; YEAR; " 是闰年"
610 NEXT YEAR
620 PRINT
630 PRINT
640 END
