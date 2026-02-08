10 REM 性能对比测试 - 大量 SIN 计算
20 REM 用于对比 BASIC 解释器 vs Go 原生性能
30
40 PRINT "开始性能测试..."
50 PRINT "测试内容: 1,000,000 次 SIN 计算"
60 PRINT
70
80 REM 测试：计算 SIN(i) 累加
90 SUM = 0
100 FOR I = 1 TO 1000000
110   SUM = SUM + SIN(I)
120 NEXT I
130
140 PRINT "测试完成!"
150 PRINT "累加结果: "; SUM
160 PRINT
170 PRINT "性能说明:"
180 PRINT "  - 执行了 1,000,000 次 SIN 函数调用"
190 PRINT "  - 每次调用包含: 数组查找 + 类型转换 + math.Sin 调用"
200 END
