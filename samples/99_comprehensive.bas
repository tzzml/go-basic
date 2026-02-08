10 REM ============================================================
20 REM 综合示例 - 演示所有 BASIC 语言特性
30 REM ============================================================
40 PRINT "=== BASIC 语言综合示例 ==="
50 PRINT
60
70 REM ===== 变量赋值 =====
80 PRINT "1. 变量赋值"
90 LET A = 10
100 B = 20
110 C$ = "Hello"
120 PRINT "A ="; A
130 PRINT "B ="; B
140 PRINT "C$ = "; C$
150 PRINT
160
170 REM ===== 算术运算 =====
180 PRINT "2. 算术运算"
190 PRINT "A + B ="; A + B
200 PRINT "A - B ="; A - B
210 PRINT "A * B ="; A * B
220 PRINT "A / B ="; A / B
230 PRINT "2 ^ 3 ="; 2 ^ 3
240 PRINT
250
260 REM ===== 比较运算 =====
270 PRINT "3. 比较运算"
280 D = 15
290 IF A = D THEN PRINT A; "等于"; D
300 IF A <> D THEN PRINT A; "不等于"; D
310 IF A < D THEN PRINT A; "小于"; D
320 IF B > D THEN PRINT B; "大于"; D
330 PRINT
340
350 REM ===== 逻辑运算 =====
360 PRINT "4. 逻辑运算"
370 IF A < B AND B > D THEN PRINT A; "<"; B; "AND"; B; ">"; D
380 IF A > 100 OR B < 50 THEN PRINT "条件满足：A > 100 OR B < 50"
390 PRINT
400
410 REM ===== 用户输入 =====
420 PRINT "5. 用户输入"
430 PRINT "请输入一个数字："
440 INPUT E
450 PRINT "您输入的数字是："; E
460 PRINT
470
480 REM ===== 条件语句 IF...THEN...ELSE =====
490 PRINT "6. 条件语句"
500 IF E > 50 THEN PRINT E; "大于 50" ELSE PRINT E; "小于等于 50"
510 PRINT
520
530 REM ===== 循环语句 FOR...NEXT =====
540 PRINT "7. 循环语句"
550 PRINT "倒计时："
560 FOR I = 5 TO 1 STEP -1
570 PRINT I
580 NEXT I
590 PRINT "发射！"
600 PRINT
610
620 REM ===== 嵌套循环 =====
630 PRINT "8. 嵌套循环（乘法表）"
640 FOR I = 1 TO 3
650 FOR J = 1 TO 3
660 PRINT I; "*"; J; "="; I * J
670 NEXT J
680 PRINT
690 NEXT I
700
710 REM ===== 子程序调用 GOSUB =====
720 PRINT "9. 子程序调用"
730 PRINT "调用子程序 1..."
740 GOSUB 1000
750 PRINT "返回主程序"
760 PRINT
770 PRINT "调用子程序 2..."
780 GOSUB 1100
790 PRINT "返回主程序"
800 PRINT
810
820 REM ===== GOTO 跳转 =====
830 PRINT "10. GOTO 跳转"
840 PRINT "跳转到结束部分..."
850 GOTO 2000
860
870 REM ===== 这段代码会被跳过 =====
880 PRINT "这行不会执行"
890
900 REM ===== 子程序定义 =====
1000 REM 子程序 1
1010 PRINT "  > 子程序 1 正在执行"
1020 PRINT "  > 执行一些计算..."
1030 LET F = A + B + E
1040 PRINT "  > A + B + E ="; F
1050 RETURN
1060
1100 REM 子程序 2
1110 PRINT "  > 子程序 2 正在执行"
1120 PRINT "  > 显示欢迎信息"
1130 PRINT "  > 感谢使用 BASIC 解释器！"
1140 RETURN
1150
1160 REM ===== 结束部分 =====
2000 PRINT
2010 PRINT "=== 程序演示完毕 ==="
2020 END
2030
