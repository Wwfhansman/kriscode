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
你是一个能使用工具的agent编程助手,名字叫kriscode。你可以使用以下工具：
- read_file: 读取一个文件的内容。参数是文件路径。
请你必须按照以下规则输出：
在你需要使用工具时，严格按以下格式输出（不要输出别的）：
Thought: <你的思考>
Action: <工具名>
Action Input: <工具的参数>

当你不需要调用工具或者已经掌握足够信息可以回答时，按以下格式输出：
Thought: <你的思考>
Final Answer: <你的最终答案>
`

func main() {

	messages := []Message{
		{Role: "system", Content: systemPrompt},
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
				if toolName == "read_file" {
					result, err := readFile(toolInput)
					if err != nil {
						result = "读取文件出错：" + err.Error() //出错也要告诉模型
					}
					//读取成功
					// 👇 把结果包成 Observation，append 回 messages
					observation := "Observation:" + result
					messages = append(messages, Message{Role: "user", Content: observation})
					fmt.Println("【调用 read_file】" + toolInput)
				}
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
