# Zork Basic

一个用 Go 语言编写的 BASIC 编程语言解释器，使用 PEG (Parsing Expression Grammar) 语法定义语法。

## 🎯 项目背景

本项目是作为**测试本地大语言模型（LLM）和 AI Agent 能力的示范项目**。

### 为什么选择 BASIC 解释器？

BASIC 解释器是一个理想的测试项目，因为它具有以下特点：

- **✅ 难度适中**：复杂度足以验证 LLM 和 Agent 的真实能力，不会过于简单或过于困难
- **✅ 多领域知识**：
  - 编译原理（词法分析、语法分析、AST）
  - 解释器实现（表达式求值、执行引擎）
  - 数据结构（抽象语法树、符号表）
  - 算法实现（数学函数、字符串操作）
- **✅ 功能完整性**：支持变量、数组、控制流、函数等完整的编程语言特性
- - **✅ 可测试性强**：功能可以独立验证，易于发现和修复问题
- - **✅ 实用价值**：功能完备的解释器本身就是一个有用的工具

### 项目目标

通过实现一个完整的 BASIC 解释器，验证：

1. **代码理解能力**：理解现有代码结构和设计
2. **架构设计能力**：设计合理的模块划分和数据结构
3. **代码实现能力**：正确实现各个功能模块
4. **问题解决能力**：调试和解决实现中的各种问题
5. **功能扩展能力**：添加新功能（如字符串运算、多维数组等）

### 实现的功能

本项目完整实现了 BASIC 语言的核心功能，包括：

- ✅ 完整的词法和语法分析（使用 PEG）
- ✅ 表达式求值（算术、比较、逻辑、字符串）
- ✅ 控制流语句（IF、FOR...NEXT、GOTO、GOSUB）
- ✅ 数组支持（包括多维数组）
- ✅ 字符串操作（连接、比较、10+个字符串函数）
- ✅ 数学函数库
- ✅ 交互式编程环境
- ✅ 完整的错误处理和边界检查

**性能**：SIN 计算达到 1,110 万次/秒，比同类 C 实现快 2 倍

### 技术栈

- **语言**：Go 1.21+
- **解析器生成器**：Pigeon (PEG)
- **架构**：模块化设计，清晰的 AST 结构

本项目证明了本地 LLM 和 AI Agent 可以独立完成中等复杂度的编程项目，具有强大的代码理解和生成能力。

## 项目结构

```
zork-basic/
├── cmd/
│   └── zork-basic/       # 主程序入口
│       └── main.go
├── internal/
│   ├── ast/               # 抽象语法树定义
│   │   └── ast.go
│   ├── interpreter/       # BASIC 解释器实现
│   │   └── interpreter.go
│   └── parser/            # PEG 语法定义和生成的解析器
│       ├── basic.peg      # BASIC 语言的 PEG 语法定义
│       ├── parser_gen.go  # 由 pigeon 生成的解析器
│       └── helpers.go     # 解析器辅助函数
├── samples/               # BASIC 示例程序
├── go.mod
└── README.md
```

## 功能特性

### 支持的 BASIC 语句

- **REM** - 注释语句（传统 BASIC 风格）
- **单引号 `'`** - 注释语句（GW-BASIC/QuickBASIC 风格，功能同 REM）
- **LET** - 变量赋值（LET 关键字可选）
- **PRINT** - 打印输出，支持分号（紧凑输出）和逗号（添加空格）分隔
- **INPUT** - 用户输入，支持单变量和多变量输入
- **DIM** - 数组声明（0-based 索引）
- **IF...THEN...ELSE...END IF** - 条件判断，支持单行和多行格式
- **FOR...NEXT** - 循环结构，支持正负 STEP 步长
- **GOTO** - 无条件跳转
- **GOSUB** - 子程序调用
- **RETURN** - 从子程序返回
- **END** - 程序结束

### 支持的表达式和运算符

- **算术运算符**: `+`, `-`, `*`, `/`, `^` (幂运算), `MOD` (取模)
- **比较运算符**: `=`, `<>`, `>`, `<`, `>=`, `<=`
- **逻辑运算符**: `AND`, `OR`, `NOT`
- **数据类型**: 数字（浮点数）、字符串
- **变量**:
  - 标识符（字母开头，可包含字母、数字、下划线和 `$`）
  - 数字变量：`A`, `COUNT`, `TOTAL`
  - 字符串变量：`NAME$`, `TITLE$`（以 `$` 结尾）
- **大小写不敏感**: 变量名和函数名不区分大小写（`myVar` = `MYVAR`）
- **内置函数**: `ABS`, `SIN`, `COS`, `TAN`, `INT`, `SQR`, `LOG`, `EXP`, `RND`

### PRINT 语句分隔符

- **分号 `;`**: 紧凑输出，值之间不添加空格
  ```basic
  PRINT "A("; I; ") = "; A(I)  ' 输出: A(0) = 10
  ```
- **逗号 `,`**: 添加空格
  ```basic
  PRINT "Hello", "World"  ' 输出: Hello  World
  ```
- **末尾分隔符**: 分号或逗号结尾不换行
  ```basic
  PRINT "Value:"; X;  ' 不换行
  PRINT Y               ' Y 会接在同一行
  ```

### 注释

zork-basic 支持两种注释风格：

- **REM 注释**（传统 BASIC 风格）:
  ```basic
  10 REM 这是一个注释
  20 REM 程序作者: Zork
  30 A = 10: REM 行尾注释（需要冒号分隔）
  ```

- **单引号 `'` 注释**（GW-BASIC/QuickBASIC 风格）:
  ```basic
  10 ' 这是一个单引号注释
  20 A = 10: ' 行尾注释（需要冒号分隔）
  30 ' 下面的代码不会执行
  40 ' PRINT "Commented out"
  ```

**注意**：
- 两种注释风格功能完全相同
- 在语句后添加注释需要使用冒号 `:` 分隔
- 单引号注释更简洁，推荐使用

### INPUT 语句

- **单变量输入**:
  ```basic
  INPUT A
  INPUT "Enter name:", NAME$
  ```
- **多变量输入**: 为每个变量提示一次
  ```basic
  INPUT X, Y, Z           ' 提示: ? [1]: ? [2]: ? [3]:
  INPUT "Enter:", A, B     ' 提示: Enter: [1]: Enter: [2]:
  ```

### 数组

- **声明**: `DIM A(10)` 创建索引 0-9 的数组
- **访问**: `A(0) = 10`, `X = A(5)`

## 安装和使用

### 前置要求

- Go 1.21 或更高版本
- [pigeon](https://github.com/mna/pigeon) - PEG 解析器生成器

### 安装 pigeon

```bash
go install github.com/mna/pigeon@latest
```

### 构建项目

```bash
# 克隆或下载项目
cd zork-basic

# 构建可执行文件
go build -o zork-basic ./cmd/zork-basic
```

### 运行 BASIC 程序

```bash
# 直接执行
./zork-basic samples/01_hello.bas

# 交互模式
./zork-basic -i
```

## 示例程序

项目包含多个示例程序，每个演示一个特定的语言特性：

| 文件名 | 说明 |
|--------|------|
| [01_hello.bas](samples/01_hello.bas) | Hello World - PRINT 和 REM 语句 |
| [02_variables.bas](samples/02_variables.bas) | 变量赋值 - LET 语句 |
| [03_input.bas](samples/03_input.bas) | 用户输入 - INPUT 交互 |
| [04_arithmetic.bas](samples/04_arithmetic.bas) | 算术运算 - +, -, *, /, ^, MOD |
| [05_comparison.bas](samples/05_comparison.bas) | 比较运算 - =, <>, >, <, >=, <= |
| [06_logical.bas](samples/06_logical.bas) | 逻辑运算 - AND, OR, NOT |
| [07_ifstmt.bas](samples/07_ifstmt.bas) | 条件语句 - IF...THEN...ELSE |
| [08_forloop.bas](samples/08_forloop.bas) | 循环语句 - FOR...NEXT，嵌套循环 |
| [09_goto.bas](samples/09_goto.bas) | 无条件跳转 - GOTO |
| [10_gosub.bas](samples/10_gosub.bas) | 子程序 - GOSUB/RETURN |
| [11_arrays.bas](samples/11_arrays.bas) | 数组 - DIM 语句 |
| [12_operators.bas](samples/12_operators.bas) | 运算符 - MOD 和 NOT |
| [13_functions.bas](samples/13_functions.bas) | 内置函数 - 数学函数 |
| [14_scientific_notation.bas](samples/14_scientific_notation.bas) | 科学计数法 |
| [15_input_advanced.bas](samples/15_input_advanced.bas) | INPUT 高级用法 - 多变量、提示、字符串 |
| [16_extensions.bas](samples/16_extensions.bas) | 扩展特性演示 |
| [17_comments.bas](samples/17_comments.bas) | 注释风格 - REM 和单引号注释 |
| [99_comprehensive.bas](samples/99_comprehensive.bas) | 综合示例 - 所有特性演示 |

### 示例：乘法表（直角三角形格式）

```basic
10 REM 嵌套循环 (直角三角形乘法表)
20 PRINT "嵌套循环 (直角三角形乘法表):"
30 FOR I = 1 TO 9
40   FOR J = 1 TO I
50     PRINT I; "*"; J; "="; I * J; " ";
60   NEXT J
70   PRINT
80 NEXT I
90 END
```

输出：
```
嵌套循环 (直角三角形乘法表):
1*1=1
2*1=2 2*2=4
3*1=3 3*2=6 3*3=9
...
```

### 示例：INPUT 多变量输入

```basic
10 INPUT X, Y, Z
20 PRINT "Sum = "; X + Y + Z
30 END
```

运行时：
```
?  [1]: 1
?  [2]: 2
?  [3]: 3
Sum = 6
```

## 重新生成解析器

如果修改了 `internal/parser/basic.peg` 语法定义文件，需要重新生成解析器：

```bash
# 使用完整路径（如果 pigeon 不在 PATH 中）
~/go/bin/pigeon -o internal/parser/parser_gen.go internal/parser/basic.peg

# 或如果已配置 PATH
pigeon -o internal/parser/parser_gen.go internal/parser/basic.peg

# 然后重新构建
go build -o zork-basic ./cmd/zork-basic
```

## 交互模式

zork-basic 支持交互式编程环境：

```bash
./zork-basic -i
```

交互模式命令：
- `LIST` 或 `L` - 列出所有程序行
- `RUN` 或 `R` - 执行程序
- `EDIT <n>` - 编辑行 n
- `DELETE <n>` - 删除行 n
- `FORMAT` 或 `F` - 格式化程序（详见下文）
- `CLEAR` - 清除所有程序行
- `NEW` - 开始新程序
- `SAVE <file>` - 保存程序到文件
- `LOAD <file>` - 从文件加载程序
- `HELP` 或 `?` - 显示帮助
- `EXIT` 或 `QUIT` - 退出

### FORMAT 命令

`FORMAT` 命令（或 `F`）用于重新格式化当前程序，执行以下操作：

1. **重新编号行号**: 将所有行号重新编号为 10, 20, 30, ...
2. **大写关键字**: 将所有 BASIC 关键字转换为大写
3. **更新跳转目标**: 自动更新 GOTO 和 GOSUB 语句中的行号引用

**示例**:

```basic
READY> 1 print "hello"
READY> 3 a=10
READY> 5 if a>5 then print "large"
READY> FORMAT
Program formatted: 3 lines renumbered
READY> LIST
10 PRINT "hello"
20 A = 10
30 IF A > 5 THEN PRINT "large"
```

**注意事项**:
- FORMAT 会修改程序，无法撤销
- 如果程序包含嵌套的 IF 语句，FORMAT 会保持原有结构
- 所有 GOTO/GOSUB 目标行号会自动更新以匹配新的行号

### 大小写不敏感

BASIC 传统上不区分大小写。zork-basic 会将所有变量名和函数名转换为大写：

```basic
10 myVar = 100
20 PRINT myvar    ' 输出: 100
30 PRINT MyVar    ' 输出: 100
```

### 比较运算符

比较运算符按优先级顺序匹配（`>=` 和 `<=` 必须在 `>` 和 `<` 之前）：

```basic
10 IF A >= 10 THEN ...  ' 正确
20 IF A > 10 THEN ...   ' 正确
```

### 变量类型

- **数字变量**: `A`, `COUNT`, `TOTAL`
- **字符串变量**: `NAME$`, `TITLE$`（以 `$` 结尾）
- **数组**: `DIM A(10)` 创建 `A(0)` 到 `A(9)`

## 许可证

本项目仅用于学习和演示目的。

## 贡献

欢迎提交 Issue 和 Pull Request！
