package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const systemPrompt = `
你是一个能使用工具的agent编程助手,名字叫kriscode。

请你必须按照以下规则输出：
在你需要使用工具时，严格按以下格式输出（不要输出别的）：
Thought: <你的思考>
Action: <工具名>
Action Input: <工具的参数>

当你不需要调用工具或者已经掌握足够信息可以回答时，按以下格式输出：
Thought: <你的思考>
Final Answer: <你的最终答案>
`

func buildToolList(tools map[string]Tool) string {
	list := ""
	for _, tool := range tools { // 遍历 map 的每个工具
		list += "- " + tool.Name() + ": " + tool.Description() + "\n"
	}
	return list
}

func main() {
	tools := map[string]Tool{ //建立工具注册表
		"read_file": ReadFileTool{},
		"list_dir":  ListDirTool{},
		"write_file": WriteFileTool{},
	}
	fullPrompt := systemPrompt + "\n\n你只能使用以下工具：\n" + buildToolList(tools)
	messages := []Message{
		{Role: "system", Content: fullPrompt},
	}

	scanner := bufio.NewScanner(os.Stdin) //创建一个扫描器
	for {
		fmt.Print("输入：")

		scanner.Scan() //读输入
		input := scanner.Text()
		if input == "exit" {
			break
		}

		messages = append(messages, Message{Role: "user", Content: input})
		//内层循环：针对这一个问题反复思考-行动
		for round := 0; round < 10; round++ { //最大循环思考十轮
			answer, err := chat(messages)
			if err != nil {
				fmt.Println("something wrong:", err)
				break //出错跳出循环，等待输入
			}

			messages = append(messages, Message{Role: "assistant", Content: answer})
			if strings.Contains(answer, "Final Answer:") {
				fmt.Println(answer)
				break //拿到最终答案就跳出循环
			} else {
				toolName, toolInput := parseAction(answer)
				tool, ok := tools[toolName] //查表找工具
				if !ok {
					messages = append(messages, Message{Role: "user", Content: "Observation:" + "未知工具：" + toolName})
					fmt.Println("[未知工具]" + toolName+"\n")
					continue //跳过本轮内循环，进入下一轮思考
				}
				result, err := tool.Execute(toolInput)
				if err != nil {
					result = "执行出错:" + err.Error()
				}
				messages = append(messages, Message{Role: "user", Content: "Observation:" + result})
				fmt.Printf("[调用%s]%s\n", toolName, toolInput)
			}
		}
	}

}

func chat(messages []Message) (string, error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY") //读apikey

	req := ChatRequest{ //构建请求
		Model:    "deepseek-v4-flash",
		Messages: messages,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	//HTTP发送
	reque, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewReader(data)) //构造post请求
	if err != nil {
		return "", err
	}

	reque.Header.Set("Content-Type", "application/json") //给请求添加header
	reque.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(reque) //发送请求
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // 把响应流读成字节切片
	if err != nil {                    // 老规矩检查
		return "", err
	}

	//解析返回
	var result ChatResponse             //准备容器
	err = json.Unmarshal(body, &result) //把body装进resp
	if err != nil {                     // 老规矩检查
		return "", err
	}
	if len(result.Choices) == 0 { //判断返回的choices长度是否为0
		return "", fmt.Errorf("API 未返回 choices，原始响应: %s", string(body))
	}
	return result.Choices[0].Message.Content, nil
}

func parseAction(text string) (string, string) {
	toolName := ""
	toolInput := ""

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Action:") {
			toolName = strings.TrimSpace(strings.TrimPrefix(line, "Action:"))
		}
		if strings.HasPrefix(line, "Action Input: ") {
			toolInput = strings.TrimSpace(strings.TrimPrefix(line, "Action Input: "))
		}
	}

	return toolName, toolInput
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type Tool interface { //定义tool接口
	Name() string
	Description() string
	Execute(input string) (string, error)
}

// 结构体定义
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
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

// write_file工具
type  WriteFileTool struct{}

func (r  WriteFileTool) Name() string {
	return "write_file"
}
func (r  WriteFileTool) Description() string {
	return "写入内容到一个文件（会覆盖原内容）。参数格式强制要求为：路径|||内容，例如 notes.txt|||hello world；如果内容有多行，必须用 \n 表示换行，全部写在一行内。" // 这段是给【模型】看的说明书
}
func (r  WriteFileTool) Execute(input string) (string, error) {
	parts :=strings.SplitN(input, "|||", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("参数格式错误，应为 路径|||内容")
	}
	content := strings.ReplaceAll(parts[1], "\\n", "\n")
	content = strings.ReplaceAll(content, "\\\"", "\"")   // 还原 \" → "
content = strings.ReplaceAll(content, "\\t", "\t")    // 还原 \t → Tab
	os.WriteFile(parts[0],[]byte(content),0644)
	return "已写入文件: " + parts[0], nil
}