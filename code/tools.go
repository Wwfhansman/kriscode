package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Tool interface { //定义tool接口
	Name() string
	Description() string
	Execute(input string) (string, error)
	IsDangerous() bool
}

func buildToolList(tools map[string]Tool) string {
	list := ""
	for _, tool := range tools { // 遍历 map 的每个工具
		list += "- " + tool.Name() + ": " + tool.Description() + "\n"
	}
	return list
}

// readfile工具
type ReadFileTool struct{}

func (r ReadFileTool) Name() string {
	return "read_file"
}
func (r ReadFileTool) Description() string {
	return "读取一个文件的内容。参数是文件路径。" // 这段是给【模型】看的说明书
}
func (r ReadFileTool) Execute(input string) (string, error) {
	data, err := os.ReadFile(input)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func (r ReadFileTool) IsDangerous() bool { return false }

// listdirtool工具
type ListDirTool struct{}

func (r ListDirTool) Name() string {
	return "list_dir"
}
func (r ListDirTool) Description() string {
	return "列出一个目录下的所有文件和子目录。参数是目录路径。" // 这段是给【模型】看的说明书
}
func (r ListDirTool) Execute(input string) (string, error) {
	entries, err := os.ReadDir(input)
	if err != nil {
		return "", err
	}
	result := ""
	for _, entry := range entries {
		result += entry.Name() + "\n"
	}
	return result, nil
}
func (r ListDirTool) IsDangerous() bool { return false }

// write_file工具
type WriteFileTool struct{}

func (r WriteFileTool) Name() string {
	return "write_file"
}
func (r WriteFileTool) Description() string {
	return "写入内容到一个文件（会覆盖原内容）。参数格式强制要求为：路径|||内容，例如 notes.txt|||hello world；如果内容有多行，必须用 \n 表示换行，全部写在一行内。" // 这段是给【模型】看的说明书
}
func (r WriteFileTool) Execute(input string) (string, error) {
	parts := strings.SplitN(input, "|||", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("参数格式错误，应为 路径|||内容")
	}
	content := strings.ReplaceAll(parts[1], "\\n", "\n")
	content = strings.ReplaceAll(content, "\\\"", "\"") // 还原 \" → "
	content = strings.ReplaceAll(content, "\\t", "\t")  // 还原 \t → Tab
	os.WriteFile(parts[0], []byte(content), 0644)
	return "已写入文件: " + parts[0], nil
}
func (r WriteFileTool) IsDangerous() bool { return true }

// RunCommand工具
type RunCommandTool struct{}

func (r RunCommandTool) Name() string {
	return "run_command"
}
func (r RunCommandTool) Description() string {
	return "执行一条 shell 命令并返回输出。参数必须是完整的命令字符串，例如 go run hello.go" // 这段是给【模型】看的说明书
}
func (r RunCommandTool) Execute(input string) (string, error) {
	cmd := exec.Command("bash", "-c", input)
	result, err := cmd.CombinedOutput()
	if err != nil {
		// 注意：命令执行失败（比如编译错误）也要把 output 返回给模型，
		// 因为报错信息就在 output 里，模型需要看到它
		return string(result) + "\n命令出错: " + err.Error(), nil
	}

	return string(result), nil
}
func (r RunCommandTool) IsDangerous() bool { return true }
