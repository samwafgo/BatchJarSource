package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func checkJava() error {
	// 执行 java -version 命令
	cmd := exec.Command("java", "-version")
	output, err := cmd.CombinedOutput()

	if err != nil {
		// 输出错误信息
		fmt.Println("Java 环境检查失败：", err)
		return err
	}

	// 检查输出是否包含 "java version" 或 "openjdk version" 字符串
	if !strings.Contains(string(output), "java version") && !strings.Contains(string(output), "openjdk version") {
		fmt.Println("未找到有效的 Java 环境。")
		return fmt.Errorf("未找到有效的 Java 环境")
	}

	fmt.Println("Java 环境检查通过。")
	return nil
}

// 反编译 JAR 文件
func decompileJar(jarFile, outputDir string) error {
	// 检查输出目录是否存在，如果不存在则创建
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, os.ModePerm)
	}
	// 检查目标文件是否已经存在
	outputFile := filepath.Join(outputDir, filepath.Base(jarFile))
	if _, err := os.Stat(outputFile); err == nil {
		fmt.Println("目标文件已存在，跳过反编译：", outputFile)
		return nil
	}

	// 创建一个 Buffer 用来存储命令执行的输出
	var out bytes.Buffer
	var stderr bytes.Buffer

	// 执行命令调用 Fernflower 进行反编译
	cmd := exec.Command("java", "-jar", "intellij-fernflower-1.2.1.16.jar", jarFile, outputDir)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		// 打印错误信息
		fmt.Println("执行命令时出错：", err)
		fmt.Println("错误信息：", stderr.String())
	}
	return err
}

// 批量反编译 JAR 文件
func batchDecompileJars(inputDir, outputDir string) {
	// 遍历输入目录下的所有文件
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		fmt.Println("读取输入目录失败：", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".jar") {
			jarFile := filepath.Join(inputDir, file.Name())
			fmt.Println("正在反编译：", jarFile)
			err := decompileJar(jarFile, outputDir)
			if err != nil {
				fmt.Println("反编译失败：", err)
			}
		}
	}
}

func main() {
	fmt.Println("批量反编译jar包辅助工具，用于代码安全审计。请勿用于非法用途。By SamWaf")
	err := checkJava()
	if err != nil {
		fmt.Println("请确保 Java 已正确安装并已添加到系统的 PATH 环境变量中。")
		return
	}
	var inputDir, outputDir string
	fmt.Print("请输入包含 JAR 文件的目录路径：")
	fmt.Scanln(&inputDir)
	fmt.Print("请输入保存反编译结果的目录路径：")
	fmt.Scanln(&outputDir)

	batchDecompileJars(inputDir, outputDir)
}
