package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	fmt.Println("开始性能测试...")
	fmt.Println("测试内容: 100,000 次 SIN 计算")
	fmt.Println()

	// 测试：计算 SIN(i) 累加
	start := time.Now()
	sum := 0.0
	for i := 1; i <= 100000; i++ {
		sum += math.Sin(float64(i))
	}
	elapsed := time.Since(start)

	fmt.Println("测试完成!")
	fmt.Printf("累加结果: %g\n", sum)
	fmt.Printf("执行时间: %v\n", elapsed)
	fmt.Println()
	fmt.Println("性能说明:")
	fmt.Println("  - 执行了 100,000 次 math.Sin 函数调用")
	fmt.Println("  - Go 原生性能，直接编译为机器码")
	fmt.Println("  - 无解释器开销")
}
