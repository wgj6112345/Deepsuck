# Deepsuck - DeepSeek 网页版复刻

一个基于 Clean Architecture 的文本对话系统，复刻 DeepSeek 聊天界面的核心体验。

## 📋 项目简介

本项目旨在学习和理解网页 Agent 的实现原理，复刻 DeepSeek 聊天界面的核心功能，包括流式对话、思考模式、对话历史管理等。

### 核心特性

- ✨ **流式对话**：支持流式输出，实时展示 AI 回复
- 🛑 **智能停止**：支持中途停止 AI 生成，已接收内容自动保存
- 🧠 **思考模式**：展示 Agent 的思考过程，可折叠/展开，支持 Markdown 渲染
- 💬 **对话历史**：支持创建、切换、删除对话，支持置顶功能
- 📝 **Markdown 渲染**：支持代码高亮、列表、链接等格式
- 🎨 **一比一复刻**：精确复刻 DeepSeek 的 UI 设计
- 🔄 **智能标题**：基于助手回答内容自动生成标题，实时更新
- 🔧 **多 Provider 支持**：支持 Mimo、IFlow、OpenAI、Claude 等多个 AI 服务提供商
- 🎯 **深度思考开关**：支持前端控制是否显示思考过程

## 🛠 技术栈

### 前端
- **框架**：Vue 3.5.26 + TypeScript
- **状态管理**：Pinia 3.0.4
- **路由**：Vue Router 4.6.4
- **HTTP 客户端**：Axios 1.13.2
- **Markdown 渲染**：markdown-it 14.1.0
- **代码高亮**：Prism.js
- **构建工具**：Vite 7.3.0

### 后端
- **语言**：Golang 1.22.0
- **架构**：Clean Architecture（分层架构）
- **数据库**：SQLite 3
- **HTTP 框架**：标准库 net/http
- **ORM**：原生 SQL

## 📁 项目结构

```
deepsuck/
├── backend/                # 后端代码
│   ├── domain/            # 领域层（实体和接口）
│   ├── usecase/           # 业务逻辑层
│   ├── repository/        # 数据访问层
│   ├── handler/           # HTTP 处理层
│   ├── middleware/        # 中间件
│   ├── security/          # 安全相关
│   └── main.go            # 入口文件
├── frontend/             # 前端代码
│   ├── src/
│   │   ├── api/          # API 客户端
│   │   ├── assets/       # 静态资源
│   │   ├── components/   # Vue 组件
│   │   ├── router/       # 路由配置
│   │   ├── store/        # 状态管理
│   │   └── views/        # 页面视图
│   └── package.json
├── ProjectState.md       # 项目状态文档
├── OpenSpec.md           # 开放规格说明
├── ImplementationPlan.md # 实现计划
└── README.md             # 本文件
```

## 🚀 快速开始

### 环境要求

- **Node.js**: ^20.19.0 || >=22.12.0
- **Go**: 1.22.0
- **Git**: 任意版本

### 安装步骤

1. **克隆仓库**
```bash
git clone git@github.com:wgj6112345/Deepsuck.git
cd Deepsuck
```

2. **安装后端依赖**
```bash
cd backend
go mod download
```

3. **安装前端依赖**
```bash
cd frontend
npm install
```

### 运行项目

1. **启动后端服务**
```bash
cd backend
go run main.go
```
后端服务将在 `http://localhost:8080` 启动

2. **启动前端开发服务器**
```bash
cd frontend
npm run dev
```
前端服务将在 `http://localhost:5173` 启动

3. **访问应用**
打开浏览器访问 `http://localhost:5173`

## ⚙️ 配置说明

### API 配置

在前端创建 `.env` 文件：

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_DEFAULT_AGENT_MODEL=mimo-v2-flash
VITE_DEFAULT_AGENT_BASE_URL=https://api.xiaomimimo.com/v1
```

### Agent 配置

应用支持多个 AI 服务提供商，在设置页面可以配置：

**支持的 Provider**：
- **Mimo**：小米自研 AI 助手，支持思考过程（`reasoning_content` 字段）
- **IFlow**：心流开放平台，支持思考过程（`reasoning_content` 字段）
- **OpenAI**：OpenAI API
- **Claude**：Anthropic Claude

**配置步骤**：
1. 进入设置页面
2. 选择要配置的 Provider 卡片
3. 点击"配置"按钮
4. 填写以下信息：
   - **API Key**: Provider 的 API 密钥
   - **模型名称**: 使用的模型名称（如 `mimo-v2-flash`、`iflow-model`）
   - **Base URL**: Provider API 的基础 URL
5. 点击"保存"
6. 点击"激活"按钮激活该 Provider

**特殊功能**：
- **深度思考模式**：Mimo 和 IFlow 支持展示思考过程，可通过输入框下方的"深度思考"按钮控制是否显示

## 💻 开发说明

### 后端开发

后端采用 Clean Architecture 分层架构：

```
handler（HTTP 处理）
  ↓
usecase（业务逻辑）
  ↓
repository（数据访问）
  ↓
domain（领域模型）
```

**添加新功能流程**：
1. 在 `domain/models.go` 定义数据模型
2. 在 `repository/` 实现数据访问接口
3. 在 `usecase/` 实现业务逻辑
4. 在 `handler/` 添加 HTTP 处理器
5. 在 `main.go` 注册路由

### 前端开发

前端采用 Vue 3 Composition API + TypeScript + Pinia：

**组件结构**：
- `components/`: 可复用组件
- `views/`: 页面级组件
- `store/`: Pinia 状态管理
- `api/`: API 客户端封装

**添加新功能流程**：
1. 在 `api/` 添加 API 客户端方法
2. 在 `store/` 添加状态管理逻辑
3. 在 `components/` 或 `views/` 实现 UI 组件
4. 在 `router/` 添加路由（如需要）

### 代码规范

- **后端**：遵循 Go 官方代码规范和 Effective Go
- **前端**：遵循 Vue 3 官方风格指南和 TypeScript 最佳实践
- **提交信息**：使用语义化提交信息（如 `feat:`, `fix:`, `docs:` 等）

## 🧪 测试

### 手动测试场景

1. **基础对话**
   - 输入文字问题
   - 验证 AI 回复正常显示
   - 验证流式输出效果

2. **思考模式**
   - 开启思考模式
   - 输入问题
   - 验证思考过程展示
   - 验证思考过程可折叠/展开

3. **对话历史**
   - 创建新对话
   - 切换对话
   - 删除对话
   - 验证对话持久化

4. **边界场景**
   - API Key 未配置
   - 网络连接失败
   - Token 欠费

## 📝 功能清单

### 已实现功能

- [x] 基础对话功能
- [x] 流式输出
- [x] 智能停止功能（中途停止 AI 生成）
- [x] 停止时自动保存已接收内容
- [x] 思考模式（支持 Mimo 和 IFlow）
- [x] 思考过程 Markdown 渲染
- [x] 思考过程折叠/展开
- [x] 深度思考开关（前端控制显示）
- [x] 对话历史管理
- [x] 对话置顶功能
- [x] Markdown 渲染
- [x] 代码高亮（Prism.js）
- [x] 代码复制功能
- [x] 对话持久化
- [x] 多 Provider 支持（Mimo、IFlow、OpenAI、Claude）
- [x] Provider 激活切换
- [x] UI 一比一复刻
- [x] 自动标题生成（基于助手回答）
- [x] 标题实时更新
- [x] "正在思考中"提示
- [x] 悬浮定位按钮
- [x] 字体大小优化
- [x] 输入框自动清空（切换对话时）

### 待完成功能

- [ ] 单元测试
- [ ] 集成测试
- [ ] 性能优化
- [ ] 安全加固
- [ ] 消息操作功能（重试、点赞、点踩、分享）
- [ ] 附件上传功能
- [ ] 联网搜索功能
- [ ] 用户认证系统
- [ ] 对话导出功能
- [ ] 深色模式
- [ ] 移动端适配优化

## 🗺️ 迭代计划

### 已完成的核心功能 ✅

1. **基础对话功能**
   - ✅ 消息发送和接收
   - ✅ SSE 流式输出
   - ✅ 输入框自动调整高度
   - ✅ 允许在 AI 输出时继续输入

2. **界面布局**
   - ✅ 侧边栏（对话列表、新建、删除、置顶）
   - ✅ 消息列表（用户消息、agent 消息）
   - ✅ 输入框区域
   - ✅ 响应式布局（侧边栏折叠）

3. **用户体验优化**
   - ✅ "正在思考中"提示
   - ✅ 悬浮定位按钮
   - ✅ 字体大小优化
   - ✅ 消息宽度自适应
   - ✅ 思考过程展示

4. **多 Provider 支持**（2025年1月）
   - ✅ 支持 Mimo、IFlow、OpenAI、Claude
   - ✅ Provider 卡片式配置界面
   - ✅ Provider 激活/切换功能
   - ✅ Mimo 思考过程支持（`reasoning_content` 字段）
   - ✅ IFlow 思考过程支持（`reasoning_content` 字段）
   - ✅ 深度思考开关（前端控制显示）

5. **思考过程功能**（2025年1月）
   - ✅ ThinkingPanel 组件（参考 ChatGPT/Claude 设计）
   - ✅ 思考过程 Markdown 渲染
   - ✅ 思考过程折叠/展开
   - ✅ 修复刷新页面后数据丢失
   - ✅ 修复标题生成逻辑（基于助手回答）
   - ✅ 修复标题更新实时性

6. **代码渲染优化**（2025年1月）
   - ✅ 从 highlight.js 切换到 Prism.js
   - ✅ 添加代码块头部（语言标签 + 复制按钮）
   - ✅ 实现代码复制功能
   - ✅ 支持多种编程语言

### 待完成的功能 🚧

#### 高优先级（核心功能）

1. **消息操作功能**（UI 已完成，逻辑待实现）
   - ⏳ 重试功能
   - ⏳ 点赞功能
   - ⏳ 点踩功能
   - ⏳ 分享功能

2. **输入框功能**
   - ⏳ 附件上传功能
   - ⏳ 联网搜索功能（UI 已完成）

3. **用户认证系统**
   - ⏳ 登录/注册
   - ⏳ 退出登录（UI 已完成）
   - ⏳ 用户信息管理

4. **对话管理增强**
   - ⏳ 对话导出（导出为 Markdown/文本）
   - ⏳ 对话搜索
   - ⏳ 对话标签/分类
   - ⏳ 批量删除对话

#### 中优先级（增强功能）

5. **设置页面**
   - ⏳ 下载手机应用功能
   - ⏳ 联系我们功能
   - ⏳ 更多设置选项

6. **界面优化**
   - ⏳ 深色模式
   - ⏳ 移动端适配优化
   - ⏳ 键盘快捷键
   - ⏳ 主题切换

#### 低优先级（锦上添花）

7. **高级功能**
   - ⏳ 对话云端同步
   - ⏳ 协作功能
   - ⏳ 多语言支持
   - ⏳ 语音输入

8. **性能优化**
   - ⏳ 虚拟滚动（大量消息时）
   - ⏳ 图片懒加载
   - ⏳ 缓存优化

### 建议的迭代顺序

**第一阶段（核心功能补全）：**
1. 实现消息操作功能（重试、点赞、点踩、分享）
2. 实现用户认证系统
3. 完善设置页面
4. 实现对话导出功能

**第二阶段（功能增强）：**
1. 附件上传功能
2. 对话搜索功能
3. 深色模式

**第三阶段（体验优化）：**
1. 移动端适配
2. 性能优化
3. 键盘快捷键

**第四阶段（高级功能）：**
1. 对话云端同步
2. 多语言支持
3. 语音输入

## 📈 更新日志

### 2025年1月

**v1.2.0 - 智能停止功能完善**

**新增功能**：
- 🛑 添加显式停止接口 POST /api/chat/stop
- 🛑 使用 Context 传播实现 Agent 请求的即时取消
- 🛑 停止时自动保存已接收的助手消息内容
- 🛑 停止时根据已接收内容生成标题
- 🛑 前端停止按钮实时反馈（发送图标 ↔ 停止图标）

**问题修复**：
- 🐛 修复停止时助手消息未保存到数据库的问题
- 🐛 修复停止时 Context 取消导致错误发送的问题
- 🐛 修复停止时发送 done 事件阻塞的问题
- 🐛 修复标题更新超时导致前端无法实时更新的问题
- 🐛 修复标题更新阻塞问题

**技术改进**：
- 🔄 在 AgentRequest 中添加 Context 字段
- 🔄 在 ChatUseCase 中添加 activeRequests 管理活跃请求
- 🔄 Agent Provider 使用 NewRequestWithContext 创建请求
- 🔄 流式读取时检查 Context 取消状态
- 🔄 扫描器错误时区分 Context 取消和其他错误
- 🔄 移除标题更新超时限制，实现真正的实时推送

**v1.1.0 - 多 Provider 支持与思考过程完善**

**新增功能**：
- ✨ 支持多个 AI 服务提供商（Mimo、IFlow、OpenAI、Claude）
- ✨ Provider 卡片式配置界面
- ✨ Provider 激活/切换功能
- ✨ 深度思考开关（前端控制是否显示思考过程）
- ✨ ThinkingPanel 组件（参考 ChatGPT/Claude 设计）
- ✨ 思考过程 Markdown 渲染
- ✨ 代码复制功能
- ✨ 代码块头部（语言标签 + 复制按钮）

**问题修复**：
- 🐛 修复刷新页面后 thinking 内容丢失的问题
- 🐛 修复标题生成逻辑（改为基于助手回答）
- 🐛 修复标题更新不是实时的问题
- 🐛 修复不开启思考模式时"正在思考中"提示不显示的问题
- 🐛 修复深度思考按钮颜色与品牌色不一致
- 🐛 修复 Provider 卡片激活状态显示问题
- 🐛 修复 Mimo 思考过程显示问题

**技术改进**：
- 🔄 从 highlight.js 切换到 Prism.js
- 🔄 优化代码高亮样式
- 🔄 优化消息字号和输入框字号
- 🔄 删除重复的样式定义

**Provider 支持**：
- Mimo：✅ 支持思考过程（`reasoning_content` 字段，需要 `thinking.type: "enabled"` 参数）
- IFlow：✅ 支持思考过程（`reasoning_content` 字段，无需参数）
- OpenAI：⏳ 待实现
- Claude：⏳ 待实现

## 🤝 贡献指南

欢迎贡献代码、报告问题或提出改进建议！

### 提交 Issue
1. 在 GitHub 上创建 Issue
2. 详细描述问题或建议
3. 提供复现步骤（如适用）

### 提交 Pull Request
1. Fork 本仓库
2. 创建特性分支（`git checkout -b feature/AmazingFeature`）
3. 提交更改（`git commit -m 'feat: Add some AmazingFeature'`）
4. 推送到分支（`git push origin feature/AmazingFeature`）
5. 创建 Pull Request

### 代码审查要求
- 代码通过 ESLint 检查
- 遵循项目代码规范
- 添加必要的注释
- 更新相关文档

## 📄 许可证

本项目仅用于学习和研究目的。

## 🙏 致谢

- [DeepSeek](https://chat.deepseek.com/) - 灵感来源
- [Vue.js](https://vuejs.org/) - 前端框架
- [Golang](https://golang.org/) - 后端语言
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - 架构参考

## 📧 联系方式

- GitHub: [@wgj6112345](https://github.com/wgj6112345)

---

**注意**：本项目仅供学习使用，请勿用于商业用途。