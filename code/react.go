package main

import (
	"strings"
)

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
