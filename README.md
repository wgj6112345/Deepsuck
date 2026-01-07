# Deepsuck - DeepSeek 网页版复刻

一个基于 Clean Architecture 的文本对话系统，复刻 DeepSeek 聊天界面的核心体验。

## 📋 项目简介

本项目旨在学习和理解网页 Agent 的实现原理，复刻 DeepSeek 聊天界面的核心功能，包括流式对话、思考模式、对话历史管理等。

### 核心特性

- ✨ **流式对话**：支持流式输出，实时展示 AI 回复
- 🧠 **思考模式**：展示 Agent 的思考过程，可折叠/展开
- 💬 **对话历史**：支持创建、切换、删除对话
- 📝 **Markdown 渲染**：支持代码高亮、列表、链接等格式
- 🎨 **一比一复刻**：精确复刻 DeepSeek 的 UI 设计
- 🔄 **智能标题**：基于对话内容自动生成标题

## 🛠 技术栈

### 前端
- **框架**：Vue 3.5.26 + TypeScript
- **状态管理**：Pinia 3.0.4
- **路由**：Vue Router 4.6.4
- **HTTP 客户端**：Axios 1.13.2
- **Markdown 渲染**：markdown-it 14.1.0
- **代码高亮**：highlight.js 11.11.1
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

在应用设置页面配置以下信息：
- **API Key**: 你的 Agent API 密钥
- **模型名称**: 使用的模型名称（如 `mimo-v2-flash`）
- **Base URL**: Agent API 的基础 URL

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

- [x] 基础对话功能
- [x] 流式输出
- [x] 思考模式
- [x] 对话历史管理
- [x] Markdown 渲染
- [x] 代码高亮
- [x] 对话持久化
- [x] Agent 配置管理
- [x] UI 一比一复刻
- [x] 自动标题生成
- [x] "正在思考中"提示
- [x] 悬浮定位按钮
- [x] 字体大小优化
- [x] 输入框自动清空（切换对话时）
- [ ] 单元测试
- [ ] 集成测试
- [ ] 性能优化
- [ ] 安全加固

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

### 待完成的功能 🚧

#### 高优先级（核心功能）

1. **消息操作功能**（UI 已完成，逻辑待实现）
   - ⏳ 重试功能
   - ⏳ 点赞功能
   - ⏳ 点踩功能
   - ⏳ 分享功能

2. **输入框功能**
   - ⏳ 附件上传功能
   - ⏳ 深度思考功能（UI 已完成）
   - ⏳ 联网搜索功能（UI 已完成）

3. **用户认证系统**
   - ⏳ 登录/注册
   - ⏳ 退出登录（UI 已完成）
   - ⏳ 用户信息管理

4. **代码渲染逻辑**（重要！）
   - ⏳ **Agent 渲染代码的逻辑还未实现**（目前 Markdown 渲染器已集成，但 agent 输出代码块时的特殊处理逻辑还未添加）

#### 中优先级（增强功能）

5. **对话管理**
   - ⏳ 对话导出（导出为 Markdown/文本）
   - ⏳ 对话搜索
   - ⏳ 对话标签/分类
   - ⏳ 批量删除对话

6. **设置页面**
   - ⏳ 下载手机应用功能
   - ⏳ 联系我们功能
   - ⏳ 更多设置选项

7. **界面优化**
   - ⏳ 深色模式
   - ⏳ 移动端适配优化
   - ⏳ 键盘快捷键
   - ⏳ 主题切换

#### 低优先级（锦上添花）

8. **高级功能**
   - ⏳ 对话云端同步
   - ⏳ 协作功能
   - ⏳ 多语言支持
   - ⏳ 语音输入

9. **性能优化**
   - ⏳ 虚拟滚动（大量消息时）
   - ⏳ 图片懒加载
   - ⏳ 缓存优化

### 建议的迭代顺序

**第一阶段（核心功能补全）：**
1. 实现 agent 渲染代码的逻辑
2. 实现消息操作功能（重试、点赞、点踩、分享）
3. 实现用户认证系统
4. 完善设置页面

**第二阶段（功能增强）：**
1. 附件上传功能
2. 对话导出功能
3. 深色模式

**第三阶段（体验优化）：**
1. 移动端适配
2. 性能优化
3. 键盘快捷键

**第四阶段（高级功能）：**
1. 对话云端同步
2. 多语言支持
3. 语音输入

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