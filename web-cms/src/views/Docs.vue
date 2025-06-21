<template>
  <div class="docs-page">
    <Header />
    
    <!-- 页面头部 -->
    <section class="page-hero gradient-bg">
      <div class="container">
        <div class="hero-content" data-aos="fade-up">
          <h1 class="page-title">开发文档</h1>
          <p class="page-subtitle">
            详细的技术文档和API说明，助您快速上手和二次开发
          </p>
        </div>
      </div>
    </section>
    
    <!-- 文档导航 -->
    <section class="docs-navigation">
      <div class="container">
        <div class="docs-grid">
          <div 
            v-for="(category, index) in docCategories" 
            :key="index"
            class="doc-category card"
            :data-aos="'fade-up'"
            :data-aos-delay="index * 100"
          >
            <div class="category-header">
              <div class="category-icon">
                <el-icon><component :is="category.icon" /></el-icon>
              </div>
              <h3 class="category-title">{{ category.title }}</h3>
              <p class="category-desc">{{ category.description }}</p>
            </div>
            
            <div class="doc-list">
              <a 
                v-for="doc in category.docs" 
                :key="doc.title"
                href="#"
                class="doc-link"
                @click.prevent="openDoc(doc)"
              >
                <span class="doc-title">{{ doc.title }}</span>
                <el-icon><ArrowRight /></el-icon>
              </a>
            </div>
            
            <div class="category-actions">
              <button class="btn-outline" @click="viewCategory(category)">
                查看全部
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- 快速开始 -->
    <section class="quick-start">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2 class="section-title">快速开始</h2>
          <p class="section-subtitle">5分钟快速搭建开发环境</p>
        </div>
        
        <div class="start-steps">
          <div class="step-timeline">
            <div 
              v-for="(step, index) in quickStartSteps" 
              :key="index"
              class="timeline-item"
              :data-aos="'fade-right'"
              :data-aos-delay="index * 100"
            >
              <div class="timeline-marker">
                <span class="step-number">{{ index + 1 }}</span>
              </div>
              <div class="timeline-content">
                <h3 class="step-title">{{ step.title }}</h3>
                <p class="step-desc">{{ step.description }}</p>
                <div v-if="step.code" class="code-block">
                  <pre><code>{{ step.code }}</code></pre>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- API预览 -->
    <section class="api-preview">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2 class="section-title">API接口预览</h2>
          <p class="section-subtitle">RESTful API设计，标准化接口规范</p>
        </div>
        
        <div class="api-demo">
          <div class="api-tabs">
            <button 
              v-for="(api, index) in apiExamples" 
              :key="index"
              class="api-tab"
              :class="{ 'active': activeApiTab === index }"
              @click="activeApiTab = index"
            >
              <span class="method" :class="api.method.toLowerCase()">{{ api.method }}</span>
              {{ api.title }}
            </button>
          </div>
          
          <div class="api-content" data-aos="fade-up">
            <div class="api-info">
              <h3 class="api-title">{{ apiExamples[activeApiTab].title }}</h3>
              <p class="api-desc">{{ apiExamples[activeApiTab].description }}</p>
              <div class="api-url">
                <span class="method" :class="apiExamples[activeApiTab].method.toLowerCase()">
                  {{ apiExamples[activeApiTab].method }}
                </span>
                <span class="url">{{ apiExamples[activeApiTab].url }}</span>
              </div>
            </div>
            
            <div class="api-example">
              <div class="example-section">
                <h4>请求示例</h4>
                <div class="code-block">
                  <pre><code>{{ apiExamples[activeApiTab].request }}</code></pre>
                </div>
              </div>
              
              <div class="example-section">
                <h4>响应示例</h4>
                <div class="code-block">
                  <pre><code>{{ apiExamples[activeApiTab].response }}</code></pre>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- 技术支持 -->
    <section class="tech-support">
      <div class="container">
        <div class="support-content" data-aos="fade-up">
          <h2 class="support-title">需要技术支持？</h2>
          <p class="support-desc">
            我们的技术团队随时为您提供专业的技术支持和咨询服务
          </p>
          <div class="support-actions">
            <button class="btn-primary" @click="contactSupport">
              <el-icon><Service /></el-icon>
              联系技术支持
            </button>
            <button class="btn-outline" @click="joinCommunity">
              <el-icon><ChatDotRound /></el-icon>
              加入社区
            </button>
          </div>
        </div>
      </div>
    </section>
    
    <Footer />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'

const router = useRouter()
const activeApiTab = ref(0)

const docCategories = ref([
  {
    icon: 'Document',
    title: '基础文档',
    description: '项目介绍、安装部署、基本配置等基础文档',
    docs: [
      { title: '项目介绍', url: '/docs/intro' },
      { title: '快速安装', url: '/docs/install' },
      { title: '配置说明', url: '/docs/config' },
      { title: '目录结构', url: '/docs/structure' }
    ]
  },
  {
    icon: 'Code',
    title: '前端开发',
    description: 'Vue 3前端开发指南、组件使用、状态管理等',
    docs: [
      { title: '开发环境搭建', url: '/docs/frontend/setup' },
      { title: '组件开发规范', url: '/docs/frontend/components' },
      { title: '路由配置', url: '/docs/frontend/router' },
      { title: '状态管理', url: '/docs/frontend/store' }
    ]
  },
  {
    icon: 'Server',
    title: '后端开发',
    description: 'Go语言后端开发指南、API开发、数据库操作等',
    docs: [
      { title: '项目架构', url: '/docs/backend/architecture' },
      { title: 'API开发规范', url: '/docs/backend/api' },
      { title: '数据库设计', url: '/docs/backend/database' },
      { title: '中间件使用', url: '/docs/backend/middleware' }
    ]
  },
  {
    icon: 'Tools',
    title: '部署运维',
    description: '生产环境部署、性能优化、监控告警等运维文档',
    docs: [
      { title: '生产部署', url: '/docs/deploy/production' },
      { title: '性能优化', url: '/docs/deploy/optimize' },
      { title: '监控配置', url: '/docs/deploy/monitor' },
      { title: '故障排查', url: '/docs/deploy/troubleshoot' }
    ]
  }
])

const quickStartSteps = ref([
  {
    title: '克隆项目',
    description: '从Git仓库克隆项目到本地',
    code: 'git clone https://github.com/your-repo/wnadmin.git\ncd wnadmin'
  },
  {
    title: '安装依赖',
    description: '分别安装前端和后端依赖',
    code: '# 安装前端依赖\ncd web-admin\nnpm install\n\n# 安装后端依赖\ncd ../server\ngo mod tidy'
  },
  {
    title: '配置数据库',
    description: '创建数据库并导入初始数据',
    code: '# 创建数据库\nmysql -u root -p\nCREATE DATABASE wnsys;\n\n# 导入数据\nmysql -u root -p wnsys < docs/sql/wnsys.sql'
  },
  {
    title: '启动项目',
    description: '启动前端和后端服务',
    code: '# 启动后端服务\ncd server\ngo run main.go\n\n# 启动前端服务\ncd web-admin\nnpm run dev'
  }
])

const apiExamples = ref([
  {
    method: 'GET',
    title: '获取用户列表',
    description: '分页获取系统用户列表',
    url: '/api/users',
    request: `curl -X GET "http://localhost:8080/api/users?page=1&size=10" \\
  -H "Authorization: Bearer your-token" \\
  -H "Content-Type: application/json"`,
    response: `{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin",
        "nickname": "管理员",
        "email": "admin@example.com",
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 10
  }
}`
  },
  {
    method: 'POST',
    title: '创建用户',
    description: '创建新的系统用户',
    url: '/api/users',
    request: `curl -X POST "http://localhost:8080/api/users" \\
  -H "Authorization: Bearer your-token" \\
  -H "Content-Type: application/json" \\
  -d '{
    "username": "newuser",
    "nickname": "新用户",
    "email": "newuser@example.com",
    "password": "123456",
    "role_ids": [2]
  }'`,
    response: `{
  "code": 200,
  "message": "用户创建成功",
  "data": {
    "id": 2,
    "username": "newuser",
    "nickname": "新用户",
    "email": "newuser@example.com",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}`
  },
  {
    method: 'PUT',
    title: '更新用户',
    description: '更新指定用户的信息',
    url: '/api/users/:id',
    request: `curl -X PUT "http://localhost:8080/api/users/2" \\
  -H "Authorization: Bearer your-token" \\
  -H "Content-Type: application/json" \\
  -d '{
    "nickname": "更新的用户",
    "email": "updated@example.com",
    "status": 1
  }'`,
    response: `{
  "code": 200,
  "message": "用户更新成功",
  "data": {
    "id": 2,
    "username": "newuser",
    "nickname": "更新的用户",
    "email": "updated@example.com",
    "status": 1,
    "updated_at": "2024-01-01T12:00:00Z"
  }
}`
  },
  {
    method: 'DELETE',
    title: '删除用户',
    description: '删除指定的用户',
    url: '/api/users/:id',
    request: `curl -X DELETE "http://localhost:8080/api/users/2" \\
  -H "Authorization: Bearer your-token"`,
    response: `{
  "code": 200,
  "message": "用户删除成功",
  "data": null
}`
  }
])

const openDoc = (doc) => {
  // 在实际项目中，这里应该跳转到具体的文档页面
  console.log('打开文档:', doc.title)
}

const viewCategory = (category) => {
  console.log('查看分类:', category.title)
}

const contactSupport = () => {
  router.push('/contact')
}

const joinCommunity = () => {
  // 跳转到社区页面或显示二维码
  console.log('加入社区')
}
</script>

<style scoped>
.docs-page {
  padding-top: 80px;
}

/* 页面头部 */
.page-hero {
  padding: 6rem 0 4rem;
  color: white;
  text-align: center;
}

.page-title {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 1rem;
}

.page-subtitle {
  font-size: 1.2rem;
  opacity: 0.9;
}

/* 文档导航 */
.docs-navigation {
  padding: 5rem 0;
  background: white;
}

.docs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.doc-category {
  padding: 2.5rem 2rem;
}

.category-header {
  text-align: center;
  margin-bottom: 2rem;
}

.category-icon {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1.5rem;
  color: white;
  font-size: 2rem;
}

.category-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #2d3748;
}

.category-desc {
  color: #718096;
  line-height: 1.6;
}

.doc-list {
  margin-bottom: 2rem;
}

.doc-link {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.8rem 0;
  text-decoration: none;
  color: #4a5568;
  border-bottom: 1px solid #e2e8f0;
  transition: color 0.3s ease;
}

.doc-link:hover {
  color: #667eea;
}

.doc-link:last-child {
  border-bottom: none;
}

.doc-title {
  font-weight: 500;
}

.category-actions {
  text-align: center;
}

/* 快速开始 */
.quick-start {
  padding: 5rem 0;
  background: #f8fafc;
}

.step-timeline {
  max-width: 800px;
  margin: 0 auto;
  position: relative;
}

.step-timeline::before {
  content: '';
  position: absolute;
  left: 30px;
  top: 30px;
  bottom: 30px;
  width: 2px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.timeline-item {
  display: flex;
  align-items: flex-start;
  gap: 2rem;
  margin-bottom: 3rem;
  position: relative;
}

.timeline-marker {
  flex-shrink: 0;
  position: relative;
  z-index: 1;
}

.step-number {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 600;
}

.timeline-content {
  flex: 1;
  background: white;
  padding: 2rem;
  border-radius: 15px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.step-title {
  font-size: 1.3rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: #2d3748;
}

.step-desc {
  color: #718096;
  margin-bottom: 1rem;
  line-height: 1.6;
}

.code-block {
  background: #1a202c;
  color: #e2e8f0;
  padding: 1.5rem;
  border-radius: 8px;
  overflow-x: auto;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace;
  font-size: 0.9rem;
  line-height: 1.5;
}

.code-block pre {
  margin: 0;
}

/* API预览 */
.api-preview {
  padding: 5rem 0;
  background: white;
}

.api-demo {
  max-width: 1000px;
  margin: 0 auto;
}

.api-tabs {
  display: flex;
  background: #f8fafc;
  border-radius: 15px;
  padding: 0.5rem;
  margin-bottom: 2rem;
  overflow-x: auto;
}

.api-tab {
  flex: 1;
  padding: 1rem 1.5rem;
  border: none;
  background: transparent;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  white-space: nowrap;
  min-width: 150px;
}

.api-tab.active {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.method {
  padding: 0.2rem 0.6rem;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 600;
  color: white;
}

.method.get {
  background: #48bb78;
}

.method.post {
  background: #4299e1;
}

.method.put {
  background: #ed8936;
}

.method.delete {
  background: #f56565;
}

.api-content {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 15px;
  overflow: hidden;
}

.api-info {
  padding: 2rem;
  border-bottom: 1px solid #e2e8f0;
}

.api-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: #2d3748;
}

.api-desc {
  color: #718096;
  margin-bottom: 1rem;
}

.api-url {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace;
}

.url {
  color: #2d3748;
  font-weight: 500;
}

.api-example {
  padding: 2rem;
}

.example-section {
  margin-bottom: 2rem;
}

.example-section:last-child {
  margin-bottom: 0;
}

.example-section h4 {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #2d3748;
}

/* 技术支持 */
.tech-support {
  padding: 5rem 0;
  background: #f8fafc;
  text-align: center;
}

.support-title {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 1rem;
  color: #2d3748;
}

.support-desc {
  color: #718096;
  margin-bottom: 2rem;
  font-size: 1.1rem;
}

.support-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.support-actions .btn-primary,
.support-actions .btn-outline {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* 响应式 */
@media (max-width: 768px) {
  .page-title {
    font-size: 2rem;
  }
  
  .docs-grid {
    grid-template-columns: 1fr;
  }
  
  .timeline-item {
    gap: 1rem;
  }
  
  .step-number {
    width: 50px;
    height: 50px;
    font-size: 1.2rem;
  }
  
  .step-timeline::before {
    left: 25px;
  }
  
  .api-tabs {
    flex-direction: column;
  }
  
  .api-url {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .support-actions {
    flex-direction: column;
    align-items: center;
  }
}
</style> 