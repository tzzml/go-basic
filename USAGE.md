# zork-basic 使用指南

## 概述

`zork-basic` 是一个用 Go 编写的 BASIC 解释器，支持经典 BASIC 语言的大部分功能。

## 运行模式

`zork-basic` 支持两种运行模式：

### 1. 文件执行模式

直接执行 BASIC 程序文件：

```bash
./zork-basic program.bas
```

### 2. 交互模式（新增）

启动交互式 BASIC 环境：

```bash
# 方式1：使用 -i 参数
./zork-basic -i

# 方式2：无参数直接运行（默认进入交互模式）
./zork-basic
```

---

## 交互模式命令

在交互模式中，您可以使用以下命令：

### 程序编辑命令

| 命令 | 简写 | 说明 | 示例 |
|------|------|------|------|
| 输入 BASIC 代码 | - | 直接输入带行号的代码 | `10 PRINT "Hello"` |
| LIST | L | 列出所有程序行 | `LIST` |
| EDIT \<n\> | E \<n\> | 编辑第 n 行 | `EDIT 10` |
| DELETE \<n\> | D \<n\> | 删除第 n 行 | `DELETE 20` |
| FORMAT | F | 格式化程序（重新编号、大写关键字） | `FORMAT` |
| CLEAR | - | 清除所有程序 | `CLEAR` |
| NEW | - | 开始新程序 | `NEW` |

### 文件操作命令

| 命令 | 说明 | 示例 |
|------|------|------|
| LOAD \<file\> | 从文件加载程序 | `LOAD test.bas` |
| SAVE \<file\> | 保存程序到文件 | `SAVE test.bas` |

### 程序执行命令

| 命令 | 简写 | 说明 |
|------|------|------|
| RUN | R | 执行当前程序 |

### 帮助和退出

| 命令 | 简写 | 说明 |
|------|------|------|
| HELP | H, ? | 显示帮助信息 |
| EXIT | Q, QUIT | 退出解释器 |

---

## 交互模式使用示例

### 示例 1：编写并运行简单程序

```bash
$ ./zork-basic

=====================================
   zork-basic BASIC Interpreter
   Version 1.0.0
=====================================

Interactive mode. Type 'HELP' for commands.
Enter BASIC statements directly or use commands.

READY> 10 PRINT "Hello, World!"
Line 10 updated

READY> 20 X = 10
Line 20 updated

READY> 30 PRINT X * 2
Line 30 updated

READY> LIST
10 PRINT "Hello, World!"
20 X = 10
30 PRINT X * 2

READY> RUN
Hello, World!
20

Program complete.

READY> EXIT
Goodbye!
```

### 示例 2：加载、编辑和保存程序

```bash
$ ./zork-basic

READY> LOAD samples/11_arrays.bas
Loaded 24 lines from samples/11_arrays.bas

READY> LIST
10 REM DIM 数组声明和访问示例
20 PRINT
30 PRINT "=== BASIC 数组示例 ==="
...

READY> EDIT 20
Current line 20: PRINT
Enter new line (or press Enter to cancel): PRINT "Modified"
Line 20 updated

READY> DELETE 30
Line 30 deleted

READY> SAVE my_program.bas
Program saved to my_program.bas (23 lines)

READY> EXIT
```

### 示例 3：从零开始编写程序

```bash
$ ./zork-basic

READY> NEW
Ready for new program

READY> 10 REM 计算阶乘
READY> 20 INPUT "请输入数字:", N
READY> 30 RESULT = 1
READY> 40 FOR I = 1 TO N
READY> 50   RESULT = RESULT * I
READY> 60 NEXT I
READY> 70 PRINT N; "的阶乘是:"; RESULT
READY> 80 END

READY> LIST
10 REM 计算阶乘
20 INPUT "请输入数字:", N
30 RESULT = 1
40 FOR I = 1 TO N
50   RESULT = RESULT * I
60 NEXT I
70 PRINT N; "的阶乘是:"; RESULT
80 END

READY> RUN
```

---

### 示例 4：使用 FORMAT 命令格式化程序

```bash
$ ./zork-basic

READY> 5 print "hello"
READY> 15 a=10
READY> 25 if a>5 then print "large"
READY> LIST
5 PRINT "hello"
15 A = 10
25 IF A > 5 THEN PRINT "large"

READY> FORMAT
Program formatted: 3 lines renumbered

READY> LIST
10 PRINT "hello"
20 A = 10
30 IF A > 5 THEN PRINT "large"
```

**FORMAT 命令功能**:
1. **重新编号行号**: 将所有行号重新编号为 10, 20, 30, ...
2. **大写关键字**: 将所有 BASIC 关键字转换为大写（PRINT, IF, THEN 等）
3. **更新跳转目标**: 自动更新 GOTO 和 GOSUB 语句中的行号引用

**使用场景**:
- 清理手动输入的代码（行号不规则、关键字大小写混乱）
- 插入新行后重新编号（为新增代码留出空间）
- 统一代码风格

**注意事项**:
- FORMAT 会修改程序，操作无法撤销
- 建议在 FORMAT 前使用 SAVE 保存原始代码
- 如果程序有语法错误，FORMAT 会显示错误信息

---

## 命令行参数

```bash
用法: zork-basic [选项] <程序文件>

选项:
  -i, --interactive    交互模式
  -v, --version        显示版本信息
  -h, --help           显示帮助信息

示例:
  zork-basic program.bas      执行 BASIC 程序
  zork-basic -i               启动交互模式
  zork-basic                 启动交互模式（默认）
```

---

## 支持的 BASIC 语言特性

### 基础功能（已实现）

- **变量赋值**: `LET X = 10` 或 `X = 10`
- **输入输出**: `INPUT`, `PRINT`（支持分号/逗号分隔符）
- **条件语句**: `IF...THEN...ELSE...END IF`
- **循环语句**: `FOR...NEXT`（支持正负 STEP）
- **跳转语句**: `GOTO`, `GOSUB`, `RETURN`
- **数组**: `DIM A(10)`, `A(0) = 10`, `X = A(0)`
- **注释**: `REM 注释内容` 或 `' 注释内容`（两种风格功能相同）

### 运算符

- **算术运算**: `+`, `-`, `*`, `/`, `^`（幂）, `MOD`（取模）
- **比较运算**: `=`, `<>`, `>`, `<`, `>=`, `<=`
- **逻辑运算**: `AND`, `OR`, `NOT`

### 内置函数

- **数学函数**: `ABS`, `SIN`, `COS`, `TAN`, `INT`, `SQR`, `LOG`, `EXP`
- **随机数**: `RND()`
- **大小写不敏感**: 所有函数名和变量名自动转换为大写

### PRINT 语句分隔符

- **分号 `;`**: 紧凑输出，值之间不添加空格
  ```basic
  PRINT "A("; I; ") = "; A(I)  ' 输出: A(0) = 10
  ```
- **逗号 `,`**: 添加空格分隔
  ```basic
  PRINT "Hello", "World"  ' 输出: Hello  World
  ```
- **末尾分隔符**: 分号或逗号结尾抑制换行
  ```basic
  PRINT "Value:"; X;  ' 不换行
  PRINT Y               ' Y 会接在同一行
  ```

### INPUT 语句

- **单变量输入**: `INPUT A`, `INPUT "Name:", N$`
- **多变量输入**: `INPUT X, Y, Z`（为每个变量提示一次）
- **字符串变量**: 支持 `$` 后缀（如 `NAME$`）
- **自动类型识别**: 尝试解析为数字，失败则作为字符串

### 扩展功能

- **科学计数法**: `1.5E3`, `2.5E-2`
- **INPUT 提示**: `INPUT "提示:", X`
- **多变量 INPUT**: `INPUT X, Y, Z`
- **字符串变量**: `NAME$`, `TITLE$`

---

## 示例程序

### 基础功能示例

- `01_hello.bas` - Hello World
- `02_variables.bas` - 变量赋值
- `03_input.bas` - INPUT 基础输入
- `04_arithmetic.bas` - 算术运算
- `05_comparison.bas` - 比较运算
- `06_logical.bas` - 逻辑运算（AND, OR, NOT）
- `07_ifstmt.bas` - IF 条件语句
- `08_forloop.bas` - FOR 循环
- `09_goto.bas` - GOTO 跳转
- `10_gosub.bas` - GOSUB 子程序
- `11_arrays.bas` - 数组
- `99_comprehensive.bas` - 基础功能综合

### 扩展功能示例

- `12_operators.bas` - MOD 和 NOT 运算符
- `13_functions.bas` - 内置数学函数
- `14_scientific_notation.bas` - 科学计数法
- `15_input_advanced.bas` - INPUT 高级功能
- `16_extensions.bas` - 扩展功能综合

---

## 运行示例程序

```bash
# 执行基础功能综合示例
./zork-basic samples/99_comprehensive.bas

# 执行扩展功能综合示例
./zork-basic samples/16_extensions.bas

# 执行数组示例
./zork-basic samples/11_arrays.bas
```

---

## 注意事项

1. **数组索引从 0 开始**：`DIM A(5)` 创建索引 0-4 的元素
2. **交互模式中输入代码**：必须包含行号（如 `10 PRINT X`）
3. **命令不区分大小写**：`LIST`、`list`、`List` 都可以
4. **使用 Ctrl+D 退出**：在交互模式中按 Ctrl+D 或输入 `EXIT` 退出

---

## 技术实现

- **语言**: Go 1.x
- **解析器**: PEG (pigeon)
- **架构**:
  - PEG 语法定义
  - AST 节点定义
  - 解释器实现

---

## 更多信息

详细的实现文档请参阅 [FEATURES.md](FEATURES.md)

项目源码：https://github.com/yourusername/zork-basic
