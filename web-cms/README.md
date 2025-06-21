# WnAdmin CMS 官网

这是 WnAdmin 系统的官方网站项目，基于 Vue 3 + Vite 构建的现代化官网。

## 项目特性

- 🎨 现代化设计，响应式布局
- ⚡ Vite 快速构建
- 🔥 Vue 3 组合式 API
- 📱 移动端适配
- 🎭 Element Plus UI 组件
- 🎪 AOS 动画效果
- 📊 优雅的数据展示

## 技术栈

- **框架**: Vue 3
- **构建工具**: Vite
- **UI 库**: Element Plus
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **样式**: CSS3 + 响应式设计
- **动画**: AOS (Animate On Scroll)

## 项目结构

```
web-cms/
├── public/                 # 静态资源
├── src/
│   ├── components/        # 公共组件
│   │   ├── Header.vue    # 头部导航
│   │   └── Footer.vue    # 底部信息
│   ├── views/            # 页面组件
│   │   ├── Home.vue      # 首页
│   │   ├── Features.vue  # 功能特性
│   │   ├── Technology.vue # 技术架构
│   │   ├── Demo.vue      # 在线演示
│   │   ├── Docs.vue      # 开发文档
│   │   └── Contact.vue   # 联系我们
│   ├── router/           # 路由配置
│   ├── styles/           # 样式文件
│   ├── App.vue          # 根组件
│   └── main.js          # 入口文件
├── package.json
├── vite.config.js
└── README.md
```

## 开发指南

### 环境要求

- Node.js >= 16
- npm >= 7

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

### 构建生产

```bash
npm run build
```

### 预览构建

```bash
npm run preview
```

## 页面说明

### 首页 (Home)
- 英雄区域展示
- 核心功能介绍
- 技术栈展示
- 统计数据
- CTA 行动区域

### 功能特性 (Features)
- 详细功能介绍
- 交互式演示
- 技术优势展示

### 技术架构 (Technology)
- 系统架构图
- 技术栈详情
- 核心特性说明

### 在线演示 (Demo)
- 演示系统入口
- 功能预览
- 使用指南

### 开发文档 (Docs)
- 文档分类导航
- 快速开始指南
- API 示例
- 技术支持

### 联系我们 (Contact)
- 联系方式展示
- 在线联系表单
- 团队信息介绍

## 自定义配置

### 修改主题色

在 `src/styles/main.css` 中修改 CSS 变量：

```css
:root {
  --primary-color: #667eea;
  --secondary-color: #764ba2;
}
```

### 添加新页面

1. 在 `src/views/` 下创建新的 Vue 组件
2. 在 `src/router/index.js` 中添加路由配置
3. 在导航组件中添加菜单项

## 部署说明

### 静态部署

构建后将 `dist` 目录部署到任何静态服务器即可。

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

## 许可证

MIT License

## 联系方式

- 邮箱: support@wnadmin.com
- 官网: https://wnadmin.com
- 文档: https://docs.wnadmin.com 