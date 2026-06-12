package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	tools := map[string]Tool{ //建立工具注册表
		"read_file":   ReadFileTool{},
		"list_dir":    ListDirTool{},
		"write_file":  WriteFileTool{},
		"run_command": RunCommandTool{},
	}
	fullPrompt := systemPrompt + "\n\n你只能使用以下工具，并严格按照要求格式输出：\n" + buildToolList(tools)
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
				// 👇 新增：模型既没给 Final Answer，也没解析出有效 Action
				if toolName == "" {
					// 说明模型没按格式来。把它的原话直接当作回复打印
					messages = append(messages, Message{Role: "user", Content: "Observation:" + "你的输出为空或未按格式。请严格按照 Thought/Action/Action Input 或 Thought/Final Answer 格式重新输出。"})
					fmt.Printf("模型原文：%q\n", answer)
					continue
				}

				tool, ok := tools[toolName] //查表找工具
				if !ok {
					messages = append(messages, Message{Role: "user", Content: "Observation:" + "未知工具：" + toolName})
					fmt.Println("[未知工具]" + toolName + "\n")
					continue //跳过本轮内循环，进入下一轮思考
				}
				if tool.IsDangerous() {
					fmt.Printf("⚠️  即将执行危险操作 [%s]，参数：%s\n确认执行吗？(y/n): ", toolName, toolInput)
					scanner.Scan()
					confirm := scanner.Text()
					if confirm != "y" { //拒绝执行
						messages = append(messages, Message{Role: "user", Content: "Observation: 用户拒绝了这次操作"})
						fmt.Println("已取消")
						continue
					}
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
