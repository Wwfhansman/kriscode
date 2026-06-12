# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## ⚠️ 最重要：这是一个教学项目，不是普通开发项目

本仓库是用户（AI 产品经理校招生，编程基础弱）的**项目制学习项目**：用 Go 从零实现一个基于 ReAct 的 coding agent。**完成学习比完成代码更重要。**

**教学约定（必须遵守）：**

- **不要直接替用户写大段实现代码。** 正确流程：讲概念 → 给最小示例/填空模板 → 用户自己写 → review 用户代码指出问题（让用户自己改）→ 答疑
- 一次只教一个新知识点，每步都要能 `go run` 看到结果
- review 时先肯定做对的部分，再指出问题并解释"为什么"
- 全程使用中文交流
- 详细的阶段计划、任务清单和进度记录在 `LEARNING_PLAN.md` —— **每次会话从它接续教学，完成任务后更新勾选和进度表**
- 每个阶段结束时，把学习内容总结写入用户的 Obsidian 笔记库：`/Users/test/Library/Mobile Documents/iCloud~md~obsidian/Documents/kris/agent实战/`

## 常用命令

```bash
go run ./code         # 运行 agent（交互式 REPL，输入 exit 退出）；多文件项目，不能只 run 单个文件
gofmt -w ./code       # 格式化（用户常忘，提醒他自己跑）
go vet ./...          # 静态检查（能抓坏 struct tag、不可达代码等）
```

- Go 安装在 `/usr/local/go/bin/go`（Claude 的 shell 里 PATH 可能没有 go，用绝对路径）
- 运行依赖环境变量 `DEEPSEEK_API_KEY`（配置在用户的 `~/.zshenv`；`.env` 文件仅作备份，程序不读它）
- 无测试、无第三方依赖 —— 全部用标准库手写（这是刻意的教学选择，不要建议引入 SDK/框架）

## 架构（四个文件，按"一起变化的代码放一起"划分）

| 文件 | 职责 |
| --- | --- |
| `code/main.go` | main() 双层循环（REPL + ReAct 循环） |
| `code/llm.go` | 传输层：chat() + API 相关 struct（换模型供应商只动这里） |
| `code/react.go` | 协议层：systemPrompt + parseAction()（改 ReAct 格式只动这里） |
| `code/tools.go` | Tool 接口 + 4 个工具实现 + buildToolList() |

这是一个手写的 ReAct agent，靠**纯文本约定**（非原生 function calling）驱动：

```
main() 双层循环：
  外层 for：读用户输入（bufio.Scanner）
  内层 for（最多10轮）：
    chat(messages) 调 DeepSeek API
      → 含 "Final Answer:" → 打印，break
      → parseAction() 解析出工具名/参数
        → toolName 为空（输出为空或没按格式）→ %q 打印原文，回灌纠正提示，continue 重试
        → 查 tools map 注册表 → 未知工具 → 回灌 Observation，continue
        → tool.IsDangerous() → 终端 y/n 确认，拒绝则回灌"用户拒绝"
        → tool.Execute() → 结果包成 "Observation: ..." 以 user 角色 append 回 messages
```

关键设计：

- **Tool interface**（Name/Description/Execute/IsDangerous）+ `map[string]Tool` 注册表。加新工具 = 写一个实现接口的 struct + 注册表加一行；system prompt 的工具清单由 `buildToolList()` 自动生成
- **Observation 必须由程序提供**，绝不让模型自己编（agent 核心铁律）
- write_file 的参数约定 `路径|||内容`，多行内容用字面 `\n` 压成一行，Execute 里 ReplaceAll 还原（已知的纯文本传参局限，阶段 5 会用原生 function calling 替代）
- 模型用 `deepseek-v4-flash`（OpenAI 兼容格式，endpoint 为 api.deepseek.com/chat/completions）

## 当前进度（详见 LEARNING_PLAN.md）

阶段 0-4 全部完成（2026-06-12 收官：4.9 端到端验收 + 4.10 拆文件重构，阶段 4 总结已写入 Obsidian `工具系统与重构.md`）。下一步：阶段 5（原生 function calling、流式输出），尚未开始。

## 已知注意事项

- 用户终端输入中文偶尔乱码/粘连，导致对话历史污染 —— 建议用户 exit 重开，不是代码 bug
- stdin 缓冲问题：API 调用期间提前输入的内容会被后续 Scan() 消费，可能穿透危险确认提示 —— 已向用户解释，暂不修，提示符出现前不要输入
- 模型偶发返回空串或丢 `Final Answer:` 前缀 —— 主循环已有回灌重试机制兜底，属正常现象
