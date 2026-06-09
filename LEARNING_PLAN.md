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

| 项 | 选择 | 理由 |
|---|---|---|
| 语言 | Go | 学习目标之一；静态类型对初学者友好（错误暴露早）；标准库够用，几乎不需要第三方依赖 |
| LLM API | DeepSeek | 便宜、国内充值方便；兼容 OpenAI 格式（行业事实标准，学会即通用） |
| Agent 框架 | ReAct（自己实现） | Thought → Action → Observation 循环，是现代 agent 的思想源头 |
| 依赖策略 | 尽量只用标准库 | `net/http`、`encoding/json`、`os`、`os/exec` 全部够用 —— 亲手造轮子才能学到东西 |

---

## 四、阶段计划

### 阶段 0：起跑线（约 1 天）

**目标**：环境就绪，亲手写出第一个 Go 程序。

**Go 知识点**：安装与工具链、`go mod`（模块概念）、`package main` / `func main()` 程序骨架、变量与函数声明（类型后置）、`go run` 与 `go build` 的区别

**任务清单**：
- [x] 安装 Go，`go version` 验证通过（2026-06-08 完成）
- [x] `go mod init kriscode` 初始化项目（2026-06-08 完成）
- [x] 写出 `main.go` 打印 `Hello, kriscode!` 并用 `go run` 跑通（2026-06-08 完成）
- [x] 写一个 `add(a, b int) int` 函数并在 main 中调用打印（2026-06-08 完成，review 通过）

**验收标准**：不看任何参考，能独立写出"一个带函数调用的可运行 Go 程序"。

---

> **关于颗粒度**：每个阶段拆成若干「微步骤」，原则是 **一步只学一个新东西，且每一步都能 `go run` 看到结果**。不要跳步，按顺序走。卡住就停下问导师。每个微步骤后面括号里是「这步在学什么」。

---

### 阶段 1：会说话 —— 调通 LLM API

**大目标**：命令行里硬编码一个问题，运行后打印出 DeepSeek 的真实回答。

**先导准备**
- [x] 1.0a 注册 DeepSeek，创建 API key（2026-06-08）
- [x] 1.0b key 配进 `~/.zshenv` 环境变量，`echo $DEEPSEEK_API_KEY` 验证（2026-06-08，踩坑：.zshrc vs .zshenv）
- [x] 1.0c 用 curl 裸调一次 API 成功，看懂返回 JSON 结构（2026-06-08）

**第一组：把"请求体"用 Go 拼出来（只玩 JSON，先不发网络）**
- [x] 1.1 定义 `Message` struct 并打印一个实例（学：struct 定义、字段、实例化）
- [x] 1.2 用 `json.Marshal` 把 Message 变成 JSON 打印（学：Marshal=打包、`string(data)`、err 检查、tag 生效）
- [x] 1.3a 理解 tag 的作用（改成 `json:"AAA"` 做实验）（学：tag 是双向字段映射）
- [x] 1.3b 定义 `ChatRequest`（含 `Model string` 和 `Messages []Message`）（2026-06-08，踩坑：tag 漏引号，go vet 抓出）
- [x] 1.3c 构造含 system+user 两条消息的 ChatRequest，Marshal 输出与 curl 逐字一致（2026-06-08，踩坑：`:` vs `:=`、尾逗号）

**第二组：把请求真正发出去（HTTP）**
- [x] 1.4a 用 `os.Getenv` 读出 API key 并打印前几位确认（2026-06-08，踩坑：.env≠环境变量、改文件需重启shell）
- [x] 1.4b 用 `http.NewRequest` 构造一个 POST 请求（2026-06-08，踩坑：base_url vs 完整endpoint、SDK会自动补路径）
- [x] 1.4c 给请求加两个 Header：Content-Type 和 Authorization（2026-06-08，`"Bearer "+apiKey` 拼接一次过）
- [x] 1.4d 用 `http.DefaultClient.Do` 发送，读出响应 body 原样打印（2026-06-08，踩坑：忘import io、Do后漏检查err、defer位置）

**第三组：把回复从 JSON 里挖出来**
- [x] 1.5a 定义响应 struct：`ChatResponse` → `Choices []Choice` → `Choice{Message}`（2026-06-08，学会看 `[` vs `{` 判断切片/单体）
- [x] 1.5b 用 `json.Unmarshal` 解析响应，打印 `resp.Choices[0].Message.Content`（2026-06-08，Unmarshal+`&`+嵌套取值，输出干净答案）

**第四组：收尾成函数**
- [x] 1.6 把上面这套封装成 `chat(messages []Message) (string, error)` 函数，main 里调用它（2026-06-08，函数抽取/参数化/错误向上抛，编译+vet通过）

**阶段验收**：✅ `go run .` 打印出 DeepSeek 真实回复；能解释 Marshal/Unmarshal 和 `if err != nil`。
**Agent 原理收获**：LLM API 是无状态的「文本进文本出」；messages 数组 + role 的含义。
**🏁 阶段 1 完成于 2026-06-08。**

---

### 阶段 2：有记忆 —— 多轮对话

**大目标**：终端里能连续聊天，模型记得前文，输入 exit 退出。

- [x] 2.1 用 `bufio.Scanner` 读一行用户输入（2026-06-08，三行套路 NewScanner/Scan/Text）
- [x] 2.2 套上 `for` 无限循环，输入 `exit` 时 `break` 退出（2026-06-08，自己判断 scanner 该放循环外，理由正确）
- [x] 2.3+2.4 维护 `messages` 切片：每轮 append 用户输入和模型回复（2026-06-08，🎉记忆成功，模型答出"根据之前的对话"）
- [x] 2.5 system prompt（2026-06-08，自己改成给助手起名 kriscode）
- [x] 2.6 实验：故意只发最后一条消息，观察"失忆"，对比体会上下文（2026-06-08，切片 `[len-1:]`，理解记忆在应用层不在模型层）

**阶段验收**：✅ 连续多轮对话模型记得前文；能解释"记忆在应用层、模型无状态"。
**Agent 原理收获**：所谓上下文 = 每次把完整历史重发一遍；记忆是应用层工程实现；token 随轮数累积。
**🏁 阶段 2 完成于 2026-06-08。**

---

### 阶段 3：会动手 —— ReAct 循环 + 第一个工具 ⭐ 核心阶段

**大目标**：agent 自己决定调用 `read_file` 工具读文件，再基于内容回答。

**第一组：让模型按格式说话**
- [x] 3.1 网页版手工调教 ReAct prompt，验证格式遵守（2026-06-09，模型会自主判断要不要用工具、且不脑补 Observation）
- [x] 3.2 把 prompt 写成 Go const（反引号多行字符串），验证模型在复杂任务输出 Action（2026-06-09，发现"简单任务模型会脱离格式"=软约束不可靠）

**第二组：解析模型想干嘛**
- [x] 3.3 用 `strings.Contains` 判断 Action 还是 Final Answer（2026-06-09，if-else 分叉，无论模型规不规范都能分类）
- [x] 3.4 parseAction 函数：Split 按行 + HasPrefix + TrimPrefix + TrimSpace 切出工具名和参数（2026-06-09，学 range 遍历；埋点：HasPrefix 前缀包含陷阱）

**第三组：第一个工具 + 接上循环**
- [x] 3.5 写 `readFile(path) (string, error)` 工具，用 `os.ReadFile`（2026-06-09，独立写出，同 chat 模式）
- [x] 3.6 把结果包成 `Observation: ...` 追加进 messages 回灌（2026-06-09，回灌成功；踩坑：Choices[0] 越界 panic→加 len 防御+fmt.Errorf 带原始 body；发现 DeepSeek 对文本ReAct 的 tool_call_id 校验=阶段5动机；手动测试要勤重启避免历史污染）
- [x] 3.7 用内层 `for` 把「思考→解析→执行→回灌」串成自动循环，Final Answer 就 break（2026-06-09，🎉 ReAct 闭环跑通，无需人工"继续"）
- [x] 3.8 最大轮数保护（写法B：`for round:=0; round<10; round++`）（2026-06-09，给 agent 装刹车防死循环烧钱）

**阶段验收**：✅ 完整跑通"agent 自主读文件并回答"，自主判断要不要用工具。
**Agent 原理收获**：ReAct 三件套、Observation 必须由程序提供（不能让模型脑补）、模型决策vs程序执行、双层循环、最大轮数保护、文本约定=软约束（不可靠→阶段5动机）。
**🏁 阶段 3 完成于 2026-06-09。整个项目的灵魂阶段。**

---

### 阶段 4：工具箱 —— 多工具与抽象

**大目标**：用 interface 重构成可扩展工具系统，agent 会读、会写、会执行命令。

- [x] 4.1 定义 `Tool` interface（Name/Description/Execute）（2026-06-09，理解隐式实现：方法凑齐自动满足接口=鸭子类型）
- [x] 4.2 把 read_file 改写成实现 Tool 接口的 struct（2026-06-09，学方法/接收者 `(r ReadFileTool)`，用 `var _ Tool = ...` 验证接口实现）
- [x] 4.3 新增 `list_dir` 工具，用 `os.ReadDir`（2026-06-09，照模板填空，体会 interface 让加工具变机械）
- [x] 4.4 用 map 做工具注册表，主循环查表统一 `tool.Execute()`（2026-06-09，学 map、`tool,ok` 双返回判断未知工具、continue）
- [x] 4.5 工具清单从注册表 `buildToolList` 自动生成拼进 prompt（2026-06-09，踩坑：const 不能放函数调用/作用域顺序；🎉 多工具复合任务跑通：list_dir→read_file→回答）
- [x] 4.6 新增 `write_file` 工具，`os.WriteFile`（2026-06-09，学 SplitN 拆多参数、Description=与模型的契约；踩坑：多行参数被parseAction按行截断→约定\n压一行；转义打地鼠 \n\"\t→体会纯文本传参缺陷=阶段5动机；🎉 agent 写出真能 go run 的代码）
- [ ] 4.7 新增 `run_command` 工具，`os/exec`（学：执行外部命令、捕获输出）
- [ ] 4.8 危险工具（write/run）执行前要用户 y/n 确认（学：安全边界，对照 Claude Code 权限机制）
- [ ] 4.9 端到端：让 agent 完成"统计本项目代码行数并写入 stats.txt"
- [ ] 4.10 重构练习：把单文件 main.go 拆成多个 .go 文件（同 package 内函数可直接互调，无需 import）。功能不变只改结构，改完验证行为一致——体验"安全重构"

**阶段验收**：新增工具只需一个 struct + 注册一行；能解释 interface 解决了什么问题。
**Agent 原理收获**：工具三要素、工具描述就是 prompt engineering、危险操作的人类确认。

---

### 阶段 5：进化 —— 现代化改造（选做）

**大目标**：升级到工业级做法，理解现代 agent 的真实形态。

- [ ] 5.1 把工具定义改成 JSON Schema，走 DeepSeek 原生 function calling（学：`tools` 参数、`tool_calls` 响应）
- [ ] 5.2 对比实验：同一任务，prompt 解析版 vs 原生版的稳定性（学：理解业界为何迁移）
- [ ] 5.3（选）SSE 流式输出，打字机效果（学：流式、goroutine/channel 初探）
- [ ] 5.4（选）历史过长时的截断/摘要策略（学：上下文管理）
- [ ] 5.5 写项目 README，总结架构和收获

**阶段验收**：能向面试官清晰讲出"ReAct 原理、prompt 解析 vs 原生 function calling 的差异、为什么需要工具抽象" —— AI 产品经理的硬通货。

---

## 五、进度记录

| 日期 | 内容 | 备注 |
|---|---|---|
| 2026-06-07 | 项目启动，制定计划 | |
| 2026-06-08 | 阶段 0：Go 安装完成 | brew 卡住，改用手动安装 |
| 2026-06-08 | 阶段 0 完成 🎓：hello world + add 函数 review 通过 | 亮点：自发用了参数类型简写；待养成 gofmt 习惯 |
| 2026-06-08 | 阶段 1 第一组完成：struct/tag/Marshal/切片，能用 Go 拼出完整请求体 JSON | 踩坑：.zshenv、tag漏引号(go vet)、`:`vs`:=`、尾逗号 |
| 2026-06-08 | 🎉 阶段 1 第二组完成：用自己写的 Go 程序首次调通 DeepSeek！net/http 全链路 | 踩坑：.env≠环境变量、base_url、忘import io、漏检查err、defer位置 |
| 2026-06-08 | 🏁 阶段 1 完成：封装 chat() 函数，已写 Obsidian 笔记 api调用.md | |
| 2026-06-08 | 🏁 阶段 2 完成：for循环+append实现多轮记忆，REPL 能连续聊天 | 高光：模型答出"根据之前的对话"；理解记忆在应用层 |
| 2026-06-08 | 推送 GitHub：Wwfhansman/kriscode，建了正确 .gitignore 保护 .env | 踩坑：原 .gitignore 文件名带空格形同虚设；待重置泄露的 key |
| 2026-06-09 | 🏁 阶段 3 完成 ⭐：实现完整 ReAct 循环，agent 能自主读文件回答 | 高光：agent review 了自己的源码、揪出重复return；踩坑：Choices[0]越界、tool_call_id校验、历史污染 |
| | | |

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
- **每个阶段结束时**，导师把该阶段学到的内容（概念、语法、踩坑记录）总结为笔记，写入 Obsidian 库 `kris/agent实战/`（阶段 1 → `api调用.md`）
- 遇到报错先自己读错误信息、尝试 10~20 分钟，再求助 —— 读懂编译器报错本身就是 Go 学习的一部分
- 所有代码自己敲，不复制粘贴（肌肉记忆很重要）
