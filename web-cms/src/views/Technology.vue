<template>
  <div class="technology-page">
    <Header />
    
    <!-- 页面头部 -->
    <section class="page-hero gradient-bg">
      <div class="container">
        <div class="hero-content" data-aos="fade-up">
          <h1 class="page-title">技术架构</h1>
          <p class="page-subtitle">
            采用最新的技术栈，构建高性能、可扩展的现代化管理系统
          </p>
        </div>
      </div>
    </section>
    
    <!-- 架构图 -->
    <section class="architecture">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2 class="section-title">系统架构</h2>
          <p class="section-subtitle">分层架构设计，确保系统的稳定性和可维护性</p>
        </div>
        
        <div class="arch-diagram" data-aos="fade-up" data-aos-delay="200">
          <div class="arch-layer">
            <h3>前端层</h3>
            <div class="layer-items">
              <div class="layer-item">Vue 3</div>
              <div class="layer-item">Element Plus</div>
              <div class="layer-item">微信小程序</div>
            </div>
          </div>
          <div class="arch-arrow">
            <el-icon><ArrowDown /></el-icon>
          </div>
          <div class="arch-layer">
            <h3>应用层</h3>
            <div class="layer-items">
              <div class="layer-item">Go Gin</div>
              <div class="layer-item">JWT认证</div>
              <div class="layer-item">RESTful API</div>
            </div>
          </div>
          <div class="arch-arrow">
            <el-icon><ArrowDown /></el-icon>
          </div>
          <div class="arch-layer">
            <h3>数据层</h3>
            <div class="layer-items">
              <div class="layer-item">MySQL</div>
              <div class="layer-item">Redis</div>
              <div class="layer-item">文件存储</div>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- 技术栈详情 -->
    <section class="tech-stack-details">
      <div class="container">
        <div class="tech-categories">
          <div 
            v-for="(category, index) in techCategories" 
            :key="index"
            class="tech-category"
            :data-aos="'fade-up'"
            :data-aos-delay="index * 100"
          >
            <div class="category-header">
              <div class="category-icon">
                <el-icon><component :is="category.icon" /></el-icon>
              </div>
              <h3 class="category-title">{{ category.title }}</h3>
            </div>
            <div class="tech-items">
              <div 
                v-for="tech in category.technologies" 
                :key="tech.name"
                class="tech-item-card"
              >
                <div class="tech-name">{{ tech.name }}</div>
                <div class="tech-version">{{ tech.version }}</div>
                <div class="tech-desc">{{ tech.description }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- 核心特性 -->
    <section class="core-features">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2 class="section-title">核心特性</h2>
          <p class="section-subtitle">现代化技术带来的核心优势</p>
        </div>
        
        <div class="features-grid">
          <div 
            v-for="(feature, index) in coreFeatures" 
            :key="index"
            class="feature-card card"
            :data-aos="'zoom-in'"
            :data-aos-delay="index * 100"
          >
            <div class="feature-icon gradient-bg">
              <el-icon><component :is="feature.icon" /></el-icon>
            </div>
            <h3 class="feature-title">{{ feature.title }}</h3>
            <p class="feature-desc">{{ feature.description }}</p>
            <ul class="feature-benefits">
              <li v-for="benefit in feature.benefits" :key="benefit">
                {{ benefit }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>
    
    <Footer />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'

const techCategories = ref([
  {
    icon: 'Monitor',
    title: '前端技术',
    technologies: [
      { name: 'Vue 3', version: '3.3.4', description: '现代化前端框架，组合式API' },
      { name: 'Element Plus', version: '2.3.8', description: '企业级UI组件库' },
      { name: 'Vite', version: '4.4.5', description: '快速构建工具' },
      { name: 'Pinia', version: '2.1.6', description: '状态管理库' }
    ]
  },
  {
    icon: 'Server',
    title: '后端技术',
    technologies: [
      { name: 'Go', version: '1.21', description: '高性能编程语言' },
      { name: 'Gin', version: '1.9.1', description: 'Web框架' },
      { name: 'GORM', version: '1.25.2', description: 'ORM框架' },
      { name: 'JWT', version: '5.0.0', description: '身份认证' }
    ]
  },
  {
    icon: 'Database',
    title: '数据存储',
    technologies: [
      { name: 'MySQL', version: '8.0', description: '关系型数据库' },
      { name: 'Redis', version: '7.0', description: '内存数据库' },
      { name: '阿里云OSS', version: 'Latest', description: '对象存储服务' },
      { name: 'MinIO', version: '2023', description: '私有云存储' }
    ]
  },
  {
    icon: 'Iphone',
    title: '移动端',
    technologies: [
      { name: '微信小程序', version: 'Latest', description: '原生小程序开发' },
      { name: 'uni-app', version: '3.8.7', description: '跨平台框架' },
      { name: 'WeUI', version: '2.5.11', description: '微信UI组件库' }
    ]
  }
])

const coreFeatures = ref([
  {
    icon: 'Lightning',
    title: '高性能',
    description: 'Go语言高并发特性结合前端优化策略',
    benefits: ['毫秒级响应时间', '支持万级并发', '内存使用优化', '数据库连接池']
  },
  {
    icon: 'Shield',
    title: '安全性',
    description: '多层次安全保障，全方位数据保护',
    benefits: ['JWT身份验证', 'RBAC权限控制', 'SQL注入防护', 'XSS攻击防护']
  },
  {
    icon: 'Setting',
    title: '可扩展',
    description: '模块化架构设计，支持业务快速扩展',
    benefits: ['插件化架构', '中间件支持', 'API标准化', '微服务就绪']
  },
  {
    icon: 'Monitor',
    title: '易维护',
    description: '清晰的代码结构和完善的文档体系',
    benefits: ['代码规范统一', '自动化测试', '详细文档', '错误日志追踪']
  }
])
</script>

<style scoped>
.technology-page {
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

/* 架构图 */
.architecture {
  padding: 5rem 0;
  background: white;
}

.arch-diagram {
  max-width: 800px;
  margin: 0 auto;
  text-align: center;
}

.arch-layer {
  background: white;
  border: 2px solid #e2e8f0;
  border-radius: 15px;
  padding: 2rem;
  margin-bottom: 2rem;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.arch-layer h3 {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
  color: #2d3748;
}

.layer-items {
  display: flex;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.layer-item {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 0.8rem 1.5rem;
  border-radius: 25px;
  font-weight: 500;
  font-size: 0.9rem;
}

.arch-arrow {
  font-size: 2rem;
  color: #667eea;
  margin: 1rem 0;
}

/* 技术栈详情 */
.tech-stack-details {
  padding: 5rem 0;
  background: #f8fafc;
}

.tech-categories {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.tech-category {
  background: white;
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.category-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 2rem;
  text-align: center;
}

.category-icon {
  font-size: 2.5rem;
  margin-bottom: 1rem;
}

.category-title {
  font-size: 1.5rem;
  font-weight: 600;
}

.tech-items {
  padding: 1.5rem;
}

.tech-item-card {
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.tech-item-card:last-child {
  border-bottom: none;
}

.tech-name {
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 0.3rem;
}

.tech-version {
  font-size: 0.8rem;
  color: #667eea;
  margin-bottom: 0.5rem;
}

.tech-desc {
  font-size: 0.9rem;
  color: #718096;
}

/* 核心特性 */
.core-features {
  padding: 5rem 0;
  background: white;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 2rem;
}

.feature-card {
  padding: 2.5rem 2rem;
  text-align: center;
}

.feature-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1.5rem;
  color: white;
  font-size: 2rem;
}

.feature-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #2d3748;
}

.feature-desc {
  color: #718096;
  margin-bottom: 1.5rem;
  line-height: 1.6;
}

.feature-benefits {
  list-style: none;
  text-align: left;
}

.feature-benefits li {
  padding: 0.3rem 0;
  color: #4a5568;
  position: relative;
  padding-left: 1.5rem;
}

.feature-benefits li::before {
  content: '•';
  color: #667eea;
  position: absolute;
  left: 0;
  font-size: 1.2rem;
}

/* 响应式 */
@media (max-width: 768px) {
  .page-title {
    font-size: 2rem;
  }
  
  .arch-layer {
    padding: 1.5rem;
  }
  
  .layer-items {
    flex-direction: column;
    align-items: center;
  }
  
  .tech-categories {
    grid-template-columns: 1fr;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
}
</style> 