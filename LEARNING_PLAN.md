# kriscode 学习计划

> 用 Go 从零实现一个基于 ReAct 框架的 coding agent
>
> 开始日期：2026-06-07 ｜ 预计周期：2~4 周

---

## 一、项目目标

**双重学习目标，完成学习比完成产品更重要：**

1. **理解 Agent 工作原理**：通过亲手实现 ReAct 循环、工具系统、上下文管理，真正搞懂 coding agent（如 Claude Code）是如何运作的
2. **学会 Go 基础开发**：在真实项目中掌握 Go 的核心语法和工程实践，摆脱"看得懂写不出"的状态

**最终成果（Demo 场景）**：在终端运行 `kriscode`，输入一句需求，比如：

```
> 帮我看看这个项目里有哪些 .go 文件，然后在 README 里补充一份文件说明
```

agent 会自己思考 → 调用工具列目录、读文件 → 观察结果 → 继续思考 → 写入文件 → 汇报完成。
本质上是一个**迷你版 Claude Code**。

---

## 二、教学约定（导师模式）

1. **导师讲概念 + 布置小任务 → 学员自己写代码**
2. **写完导师 review**：指出问题但不直接改，由学员自己修正（模拟真实工作中的 code review）
3. 导师只给最小示例片段或语法提示，**不给完整实现**；卡住超过 20 分钟可以要提示
4. 每一步先讲"为什么"，再讲"怎么做"
5. 每完成一个阶段，在本文档的任务清单里打勾 ✅，并做一次小复盘

---

## 三、技术选型


| 项        | 选择          | 理由                                                            |
| -------- | ----------- | ------------------------------------------------------------- |
| 语言       | Go          | 学习目标之一；静态类型对初学者友好（错误暴露早）；标准库够用，几乎不需要第三方依赖                     |
| LLM API  | DeepSeek    | 便宜、国内充值方便；兼容 OpenAI 格式（行业事实标准，学会即通用）                          |
| Agent 框架 | ReAct（自己实现） | Thought → Action → Observation 循环，是现代 agent 的思想源头             |
| 依赖策略     | 尽量只用标准库     | `net/http`、`encoding/json`、`os`、`os/exec` 全部够用 —— 亲手造轮子才能学到东西 |


---

## 四、阶段计划

### 阶段 0：起跑线（约 1 天）

**目标**：环境就绪，亲手写出第一个 Go 程序。

**Go 知识点**：安装与工具链、`go mod`（模块概念）、`package main` / `func main()` 程序骨架、变量与函数声明（类型后置）、`go run` 与 `go build` 的区别

**任务清单**：

- [x] 安装 Go，`go version` 验证通过（2026-06-08 完成）
- [x] `go mod init kriscode` 初始化项目
- [x] 写出 `main.go` 打印 `Hello, kriscode!` 并用 `go run` 跑通
- [x] 写一个 `add(a, b int) int` 函数并在 main 中调用打印

**验收标准**：不看任何参考，能独立写出"一个带函数调用的可运行 Go 程序"。

---

### 阶段 1：会说话 —— 调通 LLM API（约 2~3 天）

**目标**：用 Go 发起 HTTP 请求调用 DeepSeek API，命令行发一句话，模型回一句。

**Go 知识点**：

- struct 定义与 **struct tag**（`json:"..."`，Go 处理 JSON 的核心机制）
- `encoding/json` 的 `Marshal` / `Unmarshal`
- `net/http` 发起 POST 请求、设置 Header
- **error 处理**：`if err != nil` 模式（Go 没有 try/except，错误是返回值）
- 环境变量读取（`os.Getenv`，API key 不能写死在代码里）

**Agent 原理**：

- LLM API 的本质：一次无状态的"文本进、文本出"
- OpenAI 兼容格式：`messages` 数组、`role`（system/user/assistant）、`model`、`temperature` 等参数的含义

**任务清单**：

- [x] 注册 DeepSeek 开放平台，创建 API key，充值少量金额
- [x] 用 `curl` 先手动调通一次 API（理解裸协议，再写代码）
- [x] 定义请求/响应的 struct（ChatRequest、Message、ChatResponse）
- [x] 实现 `chat(messages) (string, error)` 函数：组装请求 → 发送 → 解析响应
- [x] main 中硬编码一个问题，打印模型回答

**验收标准**：`go run .` 能打印出 DeepSeek 的真实回复；能解释 struct tag 的作用和 `if err != nil` 的逻辑。

---

### 阶段 2：有记忆 —— 多轮对话（约 1~2 天）

**目标**：做成一个能在终端连续聊天的 REPL（读取-求值-打印循环）。

**Go 知识点**：

- slice 的追加（`append`）与遍历
- `for` 循环（Go 唯一的循环关键字）
- `bufio.Scanner` 读取标准输入
- 简单的代码组织：拆分多个函数/文件

**Agent 原理**：

- **LLM 没有记忆**：所谓"上下文"，就是每次把完整历史 messages 重新发一遍
- system prompt 的作用与写法
- 上下文窗口的概念（为什么历史不能无限长）

**任务清单**：

- [x] 实现输入循环：读取用户输入 → 追加到 messages → 调 API → 打印回复 → 追加回复到 messages
- [x] 加一条 system prompt（比如"你是一个编程助手"）
- [x] 支持输入 `exit` 退出
- [x] 实验：故意把历史截断，观察模型"失忆"现象，体会上下文的本质

**验收标准**：能连续多轮对话且模型记得前文；能向别人解释"LLM 的记忆是怎么回事"。

---

### 阶段 3：会动手 —— ReAct 循环 + 第一个工具 ⭐ 核心阶段（约 3~5 天）

**目标**：实现 Thought → Action → Observation 循环，agent 能自主调用 `read_file` 工具读取文件后回答问题。

**Go 知识点**：

- 字符串处理（`strings` 包：`Contains`、`TrimPrefix`、`Split`）
- `os.ReadFile` 文件读取
- 多返回值的灵活运用
- `switch` 语句

**Agent 原理**（本项目的灵魂）：

- **ReAct prompt 设计**：如何用 system prompt 约定 `Thought:` / `Action:` / `Final Answer:` 输出格式
- **解析模型输出**：从自由文本里提取出"要调用的工具+参数"
- **Observation 回灌**：把工具结果作为新消息塞回历史，驱动下一轮思考
- **循环终止条件**：什么时候停（出现 Final Answer / 达到最大轮数）
- stop 参数的妙用：防止模型自己编造 Observation

**任务清单**：

- [x] 设计 ReAct system prompt（先在 DeepSeek 网页版手工实验格式是否被遵守）
- [x] 实现输出解析：判断是 Action 还是 Final Answer，提取工具名和参数
- [x] 实现 `read_file` 工具
- [x] 实现主循环：思考 → 解析 → 执行工具 → 回灌结果 → 再思考，直到 Final Answer
- [x] 加最大轮数保护（比如 10 轮强制停止）
- [x] 端到端测试：问"main.go 里有几个函数？"，agent 自己读文件后回答

**验收标准**：完整跑通一次"模型自主决定读文件并基于内容回答"的过程；能画出 ReAct 循环的流程图并解释每一步。

---

### 阶段 4：工具箱 —— 多工具与抽象（约 2~3 天）

**目标**：扩展到 4 个工具，用 interface 重构出可扩展的工具系统 —— agent 从"能读"进化到"能写、能执行"。

**Go 知识点**：

- **interface**（Go 的精髓）：定义 `Tool` 接口，所有工具实现它
- map 的使用（工具注册表：名字 → 工具实例）
- `os.WriteFile`、`os.ReadDir`
- `os/exec` 执行外部命令、捕获输出

**Agent 原理**：

- 工具的三要素：名称、描述（给模型看的"说明书"）、执行逻辑
- 工具描述的质量直接决定模型会不会用、用得对不对（这就是 prompt engineering）
- **安全边界**：`write_file` 和 `run_command` 是危险操作 —— 执行前向用户确认（对照 Claude Code 的权限确认机制）

**任务清单**：

- [x] 定义 `Tool` interface（Name / Description / Execute）
- [x] 把 `read_file` 改造为接口实现，新增 `list_dir`、`write_file`、`run_command`
- [x] 工具注册表：system prompt 中的工具列表由注册表自动生成
- [x] 危险工具执行前的用户确认（y/n）
- [x] 端到端测试：让 agent 完成一个真实小任务（如"统计本项目代码行数并写入 stats.txt"）（2026-06-11 验收：跑通多工具链路；过程中发现并修复"空回复静默退出"问题）
- [x] 拆文件重构（4.10）：main.go 拆成 main/llm/react/tools 四文件，删除了阶段 3 遗留的旧 `readFile` 函数（2026-06-12 完成，阶段 4 收官 🎉）

**验收标准**：新增一个工具只需写一个 struct + 注册一行；能解释 interface 解决了什么问题。

---

### 阶段 5：进化 —— 现代化改造（选做，约 3~5 天）

**目标**：把"靠 prompt 约定格式"升级为工业级做法，体验现代 agent 的真实形态。

**Go 知识点**：

- 更复杂的 JSON 结构（嵌套、`json.RawMessage`）
- goroutine 与 channel 初探（流式输出）
- 项目结构整理（多包组织）

**Agent 原理**：

- **原生 Function Calling**：DeepSeek 的 `tools` 参数 + `tool_calls` 响应，对比手工解析的可靠性差异 —— 理解为什么业界都迁移到了这条路
- SSE 流式输出（打字机效果是怎么实现的）
- 上下文管理策略：历史过长时的截断/摘要

**任务清单**：

- [ ] 把工具定义改为 JSON Schema，走原生 function calling
- [ ] 对比实验：同一个任务，prompt 解析版 vs 原生版的稳定性
- [ ] （选）实现流式输出
- [ ] （选）简单的上下文截断策略
- [ ] 写一份项目 README，总结架构和学到的东西

**验收标准**：能向面试官清晰讲出"ReAct 的原理、prompt 解析与原生 function calling 的差异、为什么需要工具抽象" —— 这是 AI 产品经理的硬通货。

---

## 五、进度记录


| 日期         | 内容           | 备注             |
| ---------- | ------------ | -------------- |
| 2026-06-07 | 项目启动，制定计划    |                |
| 2026-06-08 | 阶段 0：Go 安装完成 | brew 卡住，改用手动安装 |
| 2026-06-10 | 阶段 4 端到端首测：多工具链路跑通，但发现 3 个问题 | 空回复静默退出 / stats.txt 数据自相矛盾 / 模型绕过 write_file 用 shell 写文件 |
| 2026-06-11 | 修复空回复：回灌 Observation 重试 + %q 调试打印；4.9 验收通过 | 学到：Println 不解析 %q 要用 Printf；stdin 缓冲会让提前输入穿透确认提示 |
| 2026-06-12 | 4.10 拆文件重构完成（main/llm/react/tools），阶段 4 收官 | 学到：一个目录一个包、import 按文件算、go run ./code；"一起变的放一起"设计原则；相对路径相对于工作目录 |


> 每次学习结束后在此追加一行，并更新对应阶段的任务勾选。

---

## 六、参考资源

- **Go 官方教程**：[A Tour of Go](https://go.dev/tour/)（中文版 [tour.go-zh.org](https://tour.go-zh.org/)）—— 配合各阶段按需查阅，不必通读
- **Go 标准库文档**：[pkg.go.dev](https://pkg.go.dev/)
- **ReAct 论文**：[ReAct: Synergizing Reasoning and Acting in Language Models](https://arxiv.org/abs/2210.03629)（Yao et al., 2022）
- **DeepSeek API 文档**：[api-docs.deepseek.com](https://api-docs.deepseek.com/)
- **OpenAI Chat 格式参考**：DeepSeek 完全兼容，看 DeepSeek 文档即可

---

## 七、常见约定

- 每次开始学习，直接说"继续"即可，导师会根据本文档和进度记录接续
- 遇到报错先自己读错误信息、尝试 10~20 分钟，再求助 —— 读懂编译器报错本身就是 Go 学习的一部分
- 所有代码自己敲，不复制粘贴（肌肉记忆很重要）
