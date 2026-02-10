# Zork Basic

一个用 Go 语言编写的高性能 BASIC 编程语言解释器。采用双引擎架构，支持直接 AST 解释执行和字节码虚拟机（Bytecode VM）执行。

## 🎯 项目背景

本项目是作为**测试本地大语言模型（LLM）和 AI Agent 能力的示范项目**。

### 为什么选择 BASIC 解释器？

BASIC 解释器是一个理想的测试项目，因为它具有以下特点：

- **✅ 难度适中**：复杂度足以验证 LLM 和 Agent 的真实能力，不会过于简单或过于困难。
- **✅ 多领域知识**：
  - 编译原理（词法分析、语法分析、AST、字节码生成）。
  - 解释器与虚拟机（栈式机设计、指令集、符号表、热路径优化）。
  - 数据结构（抽象语法树、常量池、循环栈帧）。
  - 算法实现（数学函数、字符串操作）。
- **✅ 功能完整性**：支持变量、数组、控制流、函数等完整的编程语言特性。
- **✅ 性能挑战**：通过虚拟化和底层优化，验证 LLM 在性能调优方面的极限。

### 项目成果

本项目已证明 AI Agent 不仅能实现一个“能跑”的解释器，还能独立完成**从 AST 到 Bytecode VM 的工程升级**，并实施多项深度性能优化。

- **🚀 极致性能**：循环计算吞吐量达到 **4,000 万次/秒**（VM 模式），性能逼近原生 Go 的 1/4。
- **🛠️ 现代架构**：支持符号表、常量池、以及专用的 FOR/NEXT 循环指令。

## 🚀 核心架构：双引擎驱动

zork-basic 采用灵活的双引擎设计：

1. **Bytecode VM (推荐)**: 
   - 流程：`Parse` -> `Compile` -> `Execute`
   - 特点：使用符号表和栈式虚拟机，消除运行时哈希查找，性能卓越。
2. **AST Interpreter (经典)**:
   - 流程：`Parse` -> `Visit AST`
   - 特点：直接遍历语法树，零延迟启动，适合简单的交互式任务。

## 项目结构

```
zork-basic/
├── cmd/
│   └── zork-basic/       # 主程序入口
├── internal/
│   ├── ast/               # 抽象语法树定义
│   ├── parser/            # PEG 语法定义及解析器
│   ├── compiler/          # 字节码编译器 (AST -> Bytecode) 🆕
│   ├── bytecode/          # 指令集、Chunk、常量池定义 🆕
│   ├── vm/                # 高性能虚拟执行引擎 🆕
│   ├── interpreter/       # 经典 AST 解释执行引擎
│   ├── repl/              # 交互式编程环境
│   └── formatter/         # 代码格式化与重编号
├── samples/               # BASIC 示例程序
└── PERFORMANCE.md         # 详细的性能优化报告记录
```

## 功能特性

- **完整的 BASIC 语句**: `LET`, `PRINT`, `INPUT`, `IF...THEN...ELSE`, `FOR...NEXT`, `GOTO`, `GOSUB/RETURN`, `DIM`, `END`, `REM` 等。
- **数据结构**: 支持多维数组、字符串（$ 结尾）、全局变量。
- **表达式引擎**: 支持算术 (+, -, *, /, ^, MOD)、逻辑 (AND, OR, NOT) 和比较运算。
- **内置函数**: 完备的数学函数库（ABS, SIN, COS, TAN, SQR...）和字符串函数库（LEN, LEFT$, MID$, INSTR...）。
- **专业环境**: 具有代码重编号 (FORMAT)、自动大小写规范化、SAVE/LOAD 功能。

## 安装和使用

### 构建

```bash
# 获取 pigeon 语法生成器 (仅在修改 basic.peg 时需要)
go install github.com/mna/pigeon@latest

# 构建
go build -o zork-basic ./cmd/zork-basic
```

### 运行

```bash
# 使用高性能 VM 模式执行 (默认)
./zork-basic -mode vm samples/08_forloop.bas

# 使用 AST 解释模式执行
./zork-basic -mode ast samples/08_forloop.bas

# 启动交互式 REPL
./zork-basic -i
```

## 性能表现

在 Apple Silicon 芯片上，zork-basic 的表现如下：

- **VM 模式**: ~40.3M operations/sec (SIN 计算循环)
- **AST 模式**: ~10.9M operations/sec
- **对比**: 性能是传统 AST 解释器的 **4 倍**，是同类动态语言解释器的佼佼者。

详细优化细节见：[PERFORMANCE.md](PERFORMANCE.md)

## 示例程序

| 文件 | 特性 |
| :--- | :--- |
| [08_forloop.bas](samples/08_forloop.bas) | 演示高性能 FOR/NEXT 循环 |
| [19_multidim_arrays.bas](samples/19_multidim_arrays.bas) | 演示多维数组和矩阵计算 |
| [for_test.bas](samples/for_test.bas) | 包含循环、嵌套循环、STEP 步长的基准测试 |

---

**声明**：本项目由 Antigravity (AI Agent) 深度参与并完成核心底层架构的编写与优化。
