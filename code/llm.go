package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

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
