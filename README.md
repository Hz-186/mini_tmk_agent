# 🎙️ Mini-TMK-Agent (同声传译小工具)

**Mini-TMK-Agent** 是一个基于 Go 语言开发的轻量级同声传译命令行工具。它集成了实时语音识别（ASR）、大语言模型翻译（LLM）以及语音合成（TTS）功能，致力于提供开箱即用的跨语言沟通体验。

本项目全面接入了 **SiliconFlow（硅基流动）** 的 AI 生态，底层模型包括 `SenseVoice`（识别）、`Qwen2.5-7B-Instruct`（翻译）和 `CosyVoice2`（播报）。

##  核心功能
* **终端流式同传 (Stream):** 边说边译，终端实时打印双语字幕，支持直接键盘敲击模拟输入。
* **Web 沉浸式看板 (Serve):** 启动本地 Web 服务，通过浏览器获得更直观的动态翻译字幕面板。
* **WebRTC P2P 直连 (RTC):** 支持两人异地组网，一人建房一人加入，实现跨越网络的同声传译对讲。
* **音频离线转写 (Transcript):** 导入本地音频文件，一键转写为纯文本。

---

## 🛠 环境与准备工作 

### 1. 安装 Go 环境
请确保你的电脑上已经安装了 **Go 1.20 或更高版本**。
可以在终端输入 `go version` 检查。如果没有，请前往 [Go 官网](https://go.dev/) 下载安装。

### 2. 申请 AI 密钥
本项目依赖 SiliconFlow 的云端 API 来完成 AI 运算：
1. 访问 [SiliconFlow 官网](https://siliconflow.cn/) 注册账号。
2. 在控制台生成你的 API Key（以 `sk-` 开头）。

### 3. 配置环境变量
在项目的**根目录**（与 `go.mod` 同级）下，新建一个名为 `.env` 的文件（注意不要有 `.txt` 后缀），并在里面填入你的密钥：
```env
TMK_AI_KEY=sk-你的专属密钥请粘贴在这里
```
*(注意：等号两边不要有空格，也不要加引号)*

---

##  编译与安装

在终端中进入项目根目录，依次执行以下命令：

1. **下载并同步依赖：**
   ```bash
   go mod tidy
   ```
2. **编译生成可执行文件：**
   ```bash
   go build -o mini-tmk-agent.exe ./cmd/mini-tmk-agent
   ```
   *(Mac/Linux 用户请去掉 `.exe` 后缀，执行 `go build -o mini-tmk-agent ./cmd/mini-tmk-agent`)*

---

##  使用指南 & 参数详解

你可以通过 `mini-tmk-agent -h` 查看所有的全局帮助信息。
**全局可用参数：**
* `-d, --debug`: 开启 Debug 模式，打印底层的网络握手和请求日志。

### 1. 启动 Web 看板服务 (`serve`)
在本地开启一个 Web 界面，方便通过浏览器查看实时的双语字幕流。

**基础命令：**
```bash
./mini-tmk-agent.exe serve -p 18080 --source-lang zh --target-lang en --tts
```
**参数详解：**
* `-p, --port`: 指定 Web 服务的本地端口。**强烈建议使用高位端口（如 18080 或 8888）**，避免被 Windows 系统拦截。
* `--source-lang`: 你说话的源语言代码（例如 `zh` 中文, `en` 英文, `ja` 日文）。
* `--target-lang`: 想要翻译成的目标语言代码。
* `--tts`: （可选）带上此参数则开启语音播报，不带则只显示文字字幕。

**如何使用：**
启动后，在浏览器中访问 `http://localhost:18080`，然后对着麦克风说话即可。

---

### 2. WebRTC P2P 远程对讲 (`rtc`)
支持两人在不同电脑上建立点对点（Peer-to-Peer）连接进行带翻译的语音通话。

#### A. 作为房主建房 (Host)
```bash
./mini-tmk-agent.exe rtc host --source-lang zh --target-lang en --tts
```
* **用途**：创建房间，你可以设置自己的输入语言和对方传过来时你希望听到的语言。
* **操作**：运行后，终端会打印出一长串 `Offer` 代码（JSON 格式）。请**复制这串代码**并发送给你的小伙伴。然后等待输入对方的回信。

#### B. 作为访客加入 (Join)
```bash
./mini-tmk-agent.exe rtc join "房主发来的Offer代码"
```
* **用途**：加入房主的房间。
* **操作**：将 `房主发来的Offer代码` 替换为你收到的长串 JSON，并执行命令。运行后，你的终端会生成一串 `Answer` 代码。请**复制这串 Answer 代码**发回给房主。
* **房主操作**：房主将这串 Answer 粘贴到自己的终端并回车，连接即可建立。

---

### 3. 纯终端流式同传 (`stream`)
如果你不需要 Web 界面，只希望在黑框框里看到实时翻译，使用此模式。

```bash
./mini-tmk-agent.exe stream --source-lang zh --target-lang en --tts
```
支持麦克风直接收音，**同时也支持直接在终端敲击键盘回车来模拟语音输入进行翻译。**

---

### 4. 离线音频转写 (`transcript`)
识别本地音频文件并提取文字。

```bash
./mini-tmk-agent.exe transcript -f audio.mp3 -o output.txt
```
**参数详解：**
* `-f, --file`: （可选）需要提取文字的音频文件路径。如果不传此参数，程序会弹出一个**图形化文件选择框**让你手动挑选文件！
* `-o, --output`: （可选）输出文本的文件名，默认保存为 `transcript.txt`。

