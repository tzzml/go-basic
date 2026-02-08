# FreeBASIC 安装指南

## macOS 安装方案

由于 FreeBASIC 在 macOS (特别是 Apple Silicon) 上没有官方预编译版本，这里提供几种替代方案：

---

## 🎯 推荐方案 1: 使用在线编译器

### **FBIDE Online**
- 网址: https://www.onlinegdb.com/online_basic_compiler
- 优点: 无需安装，直接在浏览器中使用
- 缺点: 需要网络连接

### **JDoodle**
- 网址: https://www.jdoodle.com/execute-freebasic-online
- 优点: 简单易用
- 缺点: 免费版有使用限制

---

## 🔧 方案 2: 虚拟机/容器

### 使用 Docker (推荐给熟悉 Linux 的用户)

```bash
# 拉取 Ubuntu 镜像
docker pull ubuntu:22.04

# 运行容器
docker run -it ubuntu:22.04 bash

# 在容器中安装 FreeBASIC
apt update
apt install -y build-essential libx11-dev libxext-dev libxpm-dev
apt install -y git

# 克隆并编译 FreeBASIC
git clone https://github.com/freebasic/fbc.git
cd fbc
make
make install

# 运行 FreeBASIC 程序
fbc hello.bas
./hello
```

---

## 🍺 方案 3: 在 Linux/Windows 上体验

如果你有：
- **Windows 电脑**: 直接下载安装包 https://www.freebasic.net/get
- **Linux 电脑**: `sudo apt install freebasic` (Debian/Ubuntu)

---

## 💻 方案 4: 在 macOS 上编译 (高级)

### 前置要求
```bash
# 安装 Xcode Command Line Tools
xcode-select --install

# 安装依赖
brew install xquartz
```

### 从源码编译
```bash
# 克隆仓库
git clone https://github.com/freebasic/fbc.git
cd fbc

# 编译 (可能需要调整)
bootstrap.osx
make
```

**注意**: 这可能比较复杂，可能需要修改一些编译配置。

---

## 🎮 快速体验 FreeBASIC

### 示例 1: Hello World
```freebasic
' hello.bas
PRINT "Hello from FreeBASIC!"
SLEEP
```

### 示例 2: 循环和计算
```freebasic
' loop.bas
DIM sum AS DOUBLE = 0
DIM i AS INTEGER

FOR i = 1 TO 100000
    sum += SIN(i)
NEXT i

PRINT "Sum = "; sum
SLEEP
```

### 示例 3: 图形界面
```freebasic
' graphics.bas
SCREEN 12
CIRCLE (320, 240), 100, 15
PAINT (320, 240), 4, 15
SLEEP
```

---

## 📊 FreeBASIC vs zork-basic 语法对比

| 特性 | FreeBASIC | zork-basic |
|------|-----------|------------|
| **类型声明** | `DIM x AS INTEGER` | 不需要 |
| **变量作用域** | 支持局部/全局 | 全局 |
| **指针** | 支持 | 不支持 |
| **结构体** | 支持 | 不支持 |
| **面向对象** | 支持 | 不支持 |
| **内联汇编** | 支持 | 不支持 |
| **跨平台** | Windows/Linux/DOS | macOS/Linux/Windows |

### 代码对比

**zork-basic**:
```basic
10 SUM = 0
20 FOR I = 1 TO 100
30   SUM = SUM + I
40 NEXT I
50 PRINT SUM
```

**FreeBASIC**:
```freebasic
DIM sum AS INTEGER = 0
DIM i AS INTEGER

FOR i = 1 TO 100
    sum += i
NEXT i

PRINT sum
SLEEP
```

---

## 🎯 实用建议

### **学习目的**
- ✅ 使用 **zork-basic** - 已经有完整的 BASIC 特性
- ✅ 无需复杂安装
- ✅ 性能足够学习使用

### **高性能需求**
- ✅ 使用 **FreeBASIC** (在 Linux/Windows 上)
- ✅ 或使用 **Go/C++** (更适合 macOS)

### **图形界面开发**
- ✅ **FreeBASIC** - 有图形库支持
- ✅ **zork-basic** - 专注于文本交互

---

## 🔗 相关资源

### **FreeBASIC 官方**
- 官网: https://www.freebasic.net
- 文档: https://www.freebasic.net/wiki
- 论坛: https://www.freebasic.net/forum
- 下载: https://www.freebasic.net/get

### **学习资源**
- FreeBASIC 官方教程: https://www.freebasic.net/wiki/doc/tutorials
- 示例代码: https://www.freebasic.net/wiki/code

---

## 💡 总结

对于 macOS 用户，我推荐：

1. **学习和教学**: 使用 **zork-basic** ✅ (已安装，功能完整)
2. **性能测试**: 使用在线 FreeBASIC 编译器
3. **深度开发**: 使用 Linux 虚拟机或容器

**zork-basic 已经提供了完整的 BASIC 体验，对于大多数学习场景完全够用！** 🎉
