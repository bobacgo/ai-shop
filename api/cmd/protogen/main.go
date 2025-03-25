package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	protoDir   = "proto"   // proto文件目录
	pbDir      = "gen/go"  // 生成pb.go代码的目录
	swaggerDir = "openapi" // swagger文档目录
)

// go run main.go <module_name>
func main() {
	// 检查命令行参数
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <module_name>")
		os.Exit(1)
	}

	// 获取模块名
	module := os.Args[1]

	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	// 构建proto文件路径
	protoPath := filepath.Join(wd, protoDir, module)

	// 检查proto文件目录是否存在
	if _, err := os.Stat(protoPath); os.IsNotExist(err) {
		fmt.Printf("Proto directory not found: %s\n", protoPath)
		os.Exit(1)
	}

	// 获取版本目录
	versions, err := os.ReadDir(protoPath)
	if err != nil {
		fmt.Printf("Error reading proto directory: %v\n", err)
		os.Exit(1)
	}

	// 检查是否有版本目录
	if len(versions) == 0 {
		fmt.Printf("No version directories found in %s\n", protoPath)
		os.Exit(1)
	}

	// 递归查找所有.proto文件
	var protoFiles []string
	err = filepath.Walk(protoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".proto" {
			protoFiles = append(protoFiles, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error finding proto files: %v\n", err)
		os.Exit(1)
	}

	if len(protoFiles) == 0 {
		fmt.Printf("No proto files found in %s\n", protoPath)
		os.Exit(1)
	}

	// 确保swagger目录存在
	if err := os.MkdirAll(filepath.Join(wd, swaggerDir), 0755); err != nil {
		fmt.Printf("Error creating swagger directory: %v\n", err)
		os.Exit(1)
	}

	// 构建protoc命令
	args := []string{
		"--go_out=" + pbDir,
		"--go-grpc_out=" + pbDir,
		"--grpc-gateway_out=" + pbDir,
		"--openapiv2_out=" + swaggerDir,
		"--openapiv2_opt=allow_merge=true,merge_file_name=" + module,
		"--proto_path=" + wd,                                        // 使用工作目录作为proto_path的根目录
		"--proto_path=" + filepath.Join(wd, "proto", "third_party"), // 添加third_party目录
	}

	// 将proto文件路径转换为相对于proto_path的路径
	relativeProtoFiles := make([]string, len(protoFiles))
	for i, file := range protoFiles {
		relPath, err := filepath.Rel(wd, file)
		if err != nil {
			fmt.Printf("Error converting to relative path: %v\n", err)
			os.Exit(1)
		}
		relativeProtoFiles[i] = relPath
	}

	// 添加所有proto文件
	args = append(args, relativeProtoFiles...)

	// 执行protoc命令
	cmd := exec.Command("protoc", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Generating %s gRPC code...\n", module)
	fmt.Printf("Command: protoc %s\n", strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error generating gRPC code: %v\n", err)
		fmt.Printf("Please ensure all import paths in your proto files are relative to the project root\n")
		os.Exit(1)
	}

	fmt.Printf("%s gRPC code generated successfully!\n", module)

	// 执行protoc-go-inject-tag命令
	// 遍历版本目录，为每个版本执行标签注入
	for _, version := range versions {
		if !version.IsDir() {
			continue
		}
		pbPattern := filepath.Join(pbDir, module, version.Name(), "*.pb.go")
		injectCmd := exec.Command("protoc-go-inject-tag", "-remove_tag_comment", "-input", pbPattern)
		injectCmd.Stdout = os.Stdout
		injectCmd.Stderr = os.Stderr

		fmt.Printf("Injecting tags for %s...\n", module)
		fmt.Printf("Command: protoc-go-inject-tag -remove_tag_comment -input=%s\n", pbPattern)
		if err := injectCmd.Run(); err != nil {
			fmt.Printf("Error injecting tags: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s tags injected successfully!\n", module)
	}
}
