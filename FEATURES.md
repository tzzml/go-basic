# Zork Basic 功能特性文档

本文档详细描述 zork-basic BASIC 解释器支持的所有功能特性。

## 目录

- [语句参考](#语句参考)
- [运算符](#运算符)
- [内置函数](#内置函数)
- [数据类型](#数据类型)
- [变量](#变量)
- [控制流](#控制流)
- [输入输出](#输入输出)
- [高级特性](#高级特性)

---

## 语句参考

### REM - 注释

**语法**: `REM <注释文本>`

注释语句，用于在程序中添加说明。REM 后面的所有内容都会被忽略。

```basic
10 REM 这是一个注释
20 REM 程序作者: Zork
```

### 单引号注释 `'`

**语法**: `' <注释文本>`

单引号注释是 GW-BASIC 和 QuickBASIC 中常用的注释风格，功能与 REM 完全相同。

```basic
10 ' 这是一个单引号注释
20 ' 程序作者: Zork
30 A = 10: ' 行尾注释（需要冒号分隔）
40 ' 下面的代码被注释掉了
50 ' PRINT "This will not execute"
```

**两种注释风格对比**:

| 风格 | 优点 | 缺点 |
|------|------|------|
| `REM` | 传统 BASIC 兼容 | 输入较长 |
| `'` | 简洁，现代 BASIC 风格 | 不兼容极早期 BASIC |

**注释多语句行**:
```basic
10 A = 10: B = 20: ' 计算输入值
20 C = A + B: REM 计算总和
30 PRINT C: ' 输出结果
```

### LET - 变量赋值

**语法**:
```
[LET] <变量名> = <表达式>
```

LET 关键字是可选的。用于给变量赋值。

```basic
10 LET A = 10
20 B = 20 + A      ' LET 可以省略
30 NAME$ = "John"  ' 字符串变量
```

### PRINT - 输出

**语法**:
```
PRINT [<表达式>[;|,] [<表达式>[;|,] ... [;|,]]
```

输出一个或多个表达式的值。

**分隔符**:
- **分号 `;`**: 紧凑输出，不添加空格
- **逗号 `,`**: 添加空格
- **末尾分隔符**: 抑制换行

```basic
10 PRINT "Hello"        ' 输出: Hello（换行）
20 PRINT "Hello";       ' 输出: Hello（不换行）
30 PRINT "World"        ' 输出: World（接在上一行）
40 PRINT "A"; 1; 2; 3   ' 输出: A123（紧凑）
50 PRINT "A", 1, 2, 3   ' 输出: A 1 2 3（有空格）
60 PRINT "A("; I; ")    ' 输出: A(0) （分号紧凑）
```

### INPUT - 用户输入

**语法**:
```
INPUT ["<提示字符串>",] <变量1>[, <变量2>, ...]
```

从用户读取输入并存储到变量中。

```basic
10 INPUT A                  ' 提示: ?
20 INPUT "Enter name:", N$  ' 提示: Enter name:
30 INPUT X, Y, Z            ' 多变量: ? [1]: ? [2]: ? [3]:
```

**多变量行为**: 为每个变量单独提示一次，显示序号。

### DIM - 数组声明

**语法**:
```
DIM <数组名>(<大小>)
```

声明一个数组。数组索引从 0 开始。

```basic
10 DIM A(10)    ' 创建 A(0) 到 A(9)，共 10 个元素
20 A(0) = 100
30 A(5) = 200
40 PRINT A(0)   ' 输出: 100
```

### IF...THEN...ELSE - 条件判断

**语法**:
```
# 单行格式
IF <条件> THEN <语句> [ELSE <语句>]

# 多行格式
IF <条件> THEN
    <语句块>
[ELSE
    <语句块>]
END IF
```

条件判断语句。

```basic
10 REM 单行 IF
20 IF A > 10 THEN PRINT "Large" ELSE PRINT "Small"

30 REM 多行 IF
40 IF SCORE >= 90 THEN
50   PRINT "优秀"
60 ELSE
70   PRINT "需要努力"
80 END IF
```

### FOR...NEXT - 循环

**语法**:
```
FOR <变量> = <起始值> TO <结束值> [STEP <步长>]
    <循环体>
NEXT [<变量>]
```

循环执行语句块。

```basic
10 REM 正步长循环
20 FOR I = 1 TO 5
30   PRINT I
40 NEXT I

50 REM 负步长循环
60 FOR I = 10 TO 1 STEP -1
70   PRINT I
80 NEXT I

90 REM 嵌套循环（直角三角形乘法表）
100 FOR I = 1 TO 9
110   FOR J = 1 TO I
120     PRINT I; "*"; J; "="; I * J; " ";
130   NEXT J
140   PRINT
150 NEXT I
```

### GOTO - 无条件跳转

**语法**: `GOTO <行号>`

跳转到指定行号继续执行。

```basic
10 INPUT N
20 IF N < 0 THEN GOTO 50  ' 跳转到行 50
30 PRINT "Positive"
40 END
50 PRINT "Negative"
60 END
```

### GOSUB...RETURN - 子程序

**语法**:
```
GOSUB <行号>   ' 跳转到子程序
...            ' 子程序代码
RETURN         ' 从子程序返回
```

调用子程序。

```basic
10 REM 主程序
20 GOSUB 100  ' 调用子程序
30 PRINT "Back from subroutine"
40 END

100 REM 子程序
110 PRINT "In subroutine"
120 RETURN
```

### END - 程序结束

**语法**: `END`

结束程序执行。

```basic
10 PRINT "Start"
20 END          ' 程序在这里结束
30 PRINT "Never executed"  ' 这行不会执行
```

---

## 运算符

### 算术运算符

| 运算符 | 说明 | 示例 |
|--------|------|------|
| `+` | 加法 | `A + B` |
| `-` | 减法 | `A - B` |
| `*` | 乘法 | `A * B` |
| `/` | 除法 | `A / B` |
| `^` | 幂运算 | `2 ^ 3` 结果为 8 |
| `MOD` | 取模 | `10 MOD 3` 结果为 1 |

**优先级** (从高到低):
1. `^`
2. `*`, `/`, `MOD`
3. `+`, `-`

### 比较运算符

| 运算符 | 说明 | 示例 |
|--------|------|------|
| `=` | 等于 | `A = B` |
| `<>` | 不等于 | `A <> B` |
| `>` | 大于 | `A > B` |
| `<` | 小于 | `A < B` |
| `>=` | 大于等于 | `A >= B` |
| `<=` | 小于等于 | `A <= B` |

**注意**: `>=` 和 `<=` 必须在 `>` 和 `<` 之前匹配。

### 逻辑运算符

| 运算符 | 说明 | 示例 |
|--------|------|------|
| `AND` | 逻辑与 | `A > 10 AND A < 20` |
| `OR` | 逻辑或 | `A < 5 OR A > 15` |
| `NOT` | 逻辑非 | `NOT (A > 10)` |

**优先级**: `NOT` > `AND` > `OR`

---

## 内置函数

### 数学函数

| 函数 | 说明 | 示例 |
|------|------|------|
| `ABS(x)` | 绝对值 | `ABS(-5)` 结果为 5 |
| `SIN(x)` | 正弦（弧度） | `SIN(1.5708)` 结果约 1 |
| `COS(x)` | 余弦（弧度） | `COS(0)` 结果为 1 |
| `TAN(x)` | 正切（弧度） | `TAN(0.7854)` 结果约 1 |
| `INT(x)` | 取整（向零截断） | `INT(3.7)` 结果为 3 |
| `SQR(x)` | 平方根 | `SQR(16)` 结果为 4 |
| `LOG(x)` | 自然对数 | `LOG(2.718)` 结果约 1 |
| `EXP(x)` | e 的 x 次方 | `EXP(1)` 结果约 2.718 |
| `RND()` | 随机数（0-1） | `RND()` |

**示例**:
```basic
10 PRINT ABS(-5)      ' 输出: 5
20 PRINT SIN(1.57)    ' 输出: 1（约）
30 PRINT SQR(25)      ' 输出: 5
40 PRINT INT(3.9)     ' 输出: 3
```

---

## 数据类型

### 数字类型

所有数字都是浮点数（float64）。

```basic
10 A = 10        ' 整数
20 B = 3.14      ' 浮点数
30 C = 1.5E10    ' 科学计数法
```

### 字符串类型

字符串用双引号括起来。

```basic
10 NAME$ = "John Doe"
20 GREETING$ = "Hello, World!"
30 PRINT NAME$      ' 输出: John Doe
```

---

## 变量

### 命名规则

- 必须以字母或下划线开头
- 可包含字母、数字、下划线和 `$`
- 大小写不敏感（`myVar` = `MYVAR`）

### 变量类型

**数字变量**:
```basic
10 A = 10
20 COUNT = 100
30 TOTAL = A + COUNT
```

**字符串变量** (以 `$` 结尾):
```basic
10 NAME$ = "Alice"
20 TITLE$ = "Manager"
30 PRINT NAME$     ' 输出: Alice
```

**数组变量**:
```basic
10 DIM A(10)
20 A(0) = 100
30 PRINT A(0)    ' 输出: 100
```

### 大小写不敏感

所有变量名和函数名都转换为大写：

```basic
10 myvar = 100
20 PRINT myvar     ' 输出: 100
30 PRINT MYVAR     ' 输出: 100（同一变量）
40 PRINT MyVar     ' 输出: 100（同一变量）
```

---

## 控制流

### 条件判断

```basic
10 REM 单行 IF
20 IF SCORE >= 60 THEN PRINT "Pass"

30 REM 单行 IF...ELSE
40 IF SCORE >= 60 THEN PRINT "Pass" ELSE PRINT "Fail"

50 REM 多行 IF
60 IF SCORE >= 90 THEN
70   PRINT "优秀"
80 ELSE
90   PRINT "良好"
100 END IF
```

### 循环

```basic
10 REM 基本循环
20 FOR I = 1 TO 10
30   PRINT I
40 NEXT I

50 REM 指定步长
60 FOR I = 0 TO 100 STEP 10
70   PRINT I
80 NEXT I

90 REM 倒序循环
100 FOR I = 10 TO 1 STEP -1
110   PRINT I
120 NEXT I

130 REM 嵌套循环
140 FOR I = 1 TO 3
150   FOR J = 1 TO 3
160     PRINT I; "*"; J; "="; I * J
170   NEXT J
180   PRINT
190 NEXT I
```

---

## 输入输出

### 输入

**单变量输入**:
```basic
10 INPUT A
20 PRINT "A = "; A
```

**带提示的输入**:
```basic
10 INPUT "Enter your name:", NAME$
20 PRINT "Hello, "; NAME$
```

**多变量输入**:
```basic
10 INPUT X, Y, Z
20 PRINT "Sum = "; X + Y + Z
```

### 输出

**基本输出**:
```basic
10 PRINT "Hello, World!"
```

**格式化输出**:
```basic
10 A = 10
20 PRINT "Value:"; A     ' 紧凑输出: Value:10
30 PRINT "Value:", A     ' 带空格: Value: 10
40 PRINT "A("; A; ") = "; A * 2  ' A(10) = 20
```

**不换行输出**:
```basic
10 PRINT "Value:";
20 PRINT 10        ' 输出: Value:10（在同一行）
30 PRINT           ' 空行
```

---

## 高级特性

### 交互模式

启动交互模式：
```bash
./zork-basic -i
```

交互模式命令：
- `LIST` 或 `L` - 列出程序
- `RUN` 或 `R` - 执行程序
- `EDIT <n>` - 编辑行 n
- `DELETE <n>` - 删除行 n
- `FORMAT` 或 `F` - 格式化程序
- `NEW` - 新建程序
- `SAVE <file>` - 保存程序
- `LOAD <file>` - 加载程序
- `HELP` - 帮助
- `EXIT` - 退出

### 表达式复杂度

支持嵌套和复杂表达式：

```basic
10 A = 10
20 B = 20
30 RESULT = (A + B) * 2 - (A / B) ^ 2
40 PRINT RESULT
```

### 逻辑组合

```basic
10 AGE = 25
20 SCORE = 85
30 IF AGE >= 18 AND SCORE >= 60 THEN PRINT "Adult Pass"
40 IF AGE < 18 OR SCORE < 60 THEN PRINT "Failed"
50 IF NOT (SCORE >= 90) THEN PRINT "Not excellent"
```

### 数组操作

```basic
10 DIM A(5)
20 REM 初始化数组
30 FOR I = 0 TO 4
40   A(I) = I * 10
50 NEXT I
60 REM 读取数组
70 FOR I = 0 TO 4
80   PRINT "A("; I; ") = "; A(I)
90 NEXT I
```

---

## 最佳实践

### 1. 使用行号间隔

建议使用 10 的倍数作为行号，便于插入新行：

```basic
10 REM 主程序
20 ...
30 ...
40 ...
```

### 2. 注释

使用 REM 添加注释说明代码逻辑：

```basic
10 REM 计算阶乘
20 INPUT "Enter n:", N
30 RESULT = 1
40 FOR I = 1 TO N
50   RESULT = RESULT * I
60 NEXT I
70 PRINT N; "! = "; RESULT
```

### 3. 变量命名

使用有意义的变量名：

```basic
10 REM 好的变量名
20 TotalScore = 100
30 StudentName$ = "Alice"

40 REM 避免使用单字母变量（除循环变量外）
50 X = 10  ' 不推荐
61 COUNT = 10  ' 推荐
```

### 4. PRINT 格式化

使用分号实现紧凑输出，逗号添加空格：

```basic
10 REM 数组输出（紧凑格式）
20 PRINT "A("; I; ") = "; A(I)

30 REM 表格输出（带空格）
40 PRINT "Name", "Score"
```

---

## 限制和注意事项

1. **数组索引从 0 开始**: `DIM A(10)` 创建 `A(0)` 到 `A(9)`
2. **所有数字都是浮点数**: 即使整数也存储为浮点数
3. **字符串必须用双引号**: 不支持单引号字符串
4. **行号必须是整数**: 行号不支持表达式
5. **每行可以有多个语句**: 使用冒号 `:` 分隔
   ```basic
   10 A = 10: PRINT A: B = 20
   ```

---

## 示例程序

完整示例请参考 `samples/` 目录：
- [01_hello.bas](samples/01_hello.bas) - Hello World
- [08_forloop.bas](samples/08_forloop.bas) - 循环和乘法表
- [15_input_advanced.bas](samples/15_input_advanced.bas) - INPUT 高级用法
- [99_comprehensive.bas](samples/99_comprehensive.bas) - 综合示例
