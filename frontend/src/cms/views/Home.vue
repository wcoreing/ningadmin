<template>
  <div class="home">
    <Header />
    
    <!-- 英雄区域 -->
    <section class="hero">
      <div class="container">
        <div class="hero-content" data-aos="fade-up">
          <h1 class="hero-title">
            现代化全栈<br>
            <span class="gradient-text">管理系统</span>
          </h1>
          <p class="hero-subtitle">
            基于Vue 3 + Go的企业级管理平台，提供完整的前后端解决方案，
            包含用户管理、权限控制、商品管理、订单处理等核心功能
          </p>
          <div class="hero-actions">
            <button class="btn-primary" @click="gotoDemo">
              <el-icon><VideoPlay /></el-icon>
              在线演示
            </button>
            <button class="btn-outline" @click="gotoFeatures">
              <el-icon><Document /></el-icon>
              了解更多
            </button>
          </div>
        </div>
        <div class="hero-image" data-aos="fade-left" data-aos-delay="200">
          <div class="dashboard-preview">
            <div class="preview-header">
              <div class="preview-dots">
                <span></span>
                <span></span>
                <span></span>
              </div>
            </div>
            <div class="preview-content">
              <div class="preview-sidebar"></div>
              <div class="preview-main">
                <div class="preview-cards">
                  <div class="preview-card"></div>
                  <div class="preview-card"></div>
                  <div class="preview-card"></div>
                  <div class="preview-card"></div>
                </div>
                <div class="preview-chart"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- 功能特性 -->
    <section class="features">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2 class="section-title">核心功能</h2>
          <p class="section-subtitle">为现代企业量身定制的管理解决方案</p>
        </div>
        
        <div class="features-grid">
          <div 
            v-for="(feature, index) in features" 
            :key="index"
            class="feature-card card"
            :data-aos="'fade-up'"
            :data-aos-delay="index * 100"
          >
            <div class="feature-icon gradient-bg">
              <el-icon><component :is="feature.icon" /></el-icon>
            </div>
            <h3 class="feature-title">{{ feature.title }}</h3>
            <p class="feature-desc">{{ feature.description }}</p>
          </div>
        </div>
      </div>
    </section>
    
    <!-- 技术栈 -->
    <section class="tech-stack">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2 class="section-title">技术架构</h2>
          <p class="section-subtitle">采用最新的技术栈，确保系统的稳定性和可扩展性</p>
        </div>
        
        <div class="tech-grid">
          <div 
            v-for="(tech, index) in techStack" 
            :key="index"
            class="tech-item"
            :data-aos="'zoom-in'"
            :data-aos-delay="index * 100"
          >
            <div class="tech-logo">
              <span class="tech-name">{{ tech.name }}</span>
            </div>
            <p class="tech-desc">{{ tech.description }}</p>
          </div>
        </div>
      </div>
    </section>
    
    <!-- 统计数据 -->
    <section class="stats gradient-bg">
      <div class="container">
        <div class="stats-grid">
          <div 
            v-for="(stat, index) in stats" 
            :key="index"
            class="stat-item"
            :data-aos="'fade-up'"
            :data-aos-delay="index * 100"
          >
            <div class="stat-number">{{ stat.number }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- CTA区域 -->
    <section class="cta">
      <div class="container">
        <div class="cta-content" data-aos="fade-up">
          <h2 class="cta-title">准备开始您的项目？</h2>
          <p class="cta-subtitle">
            立即体验WnAdmin的强大功能，快速构建您的管理系统
          </p>
          <div class="cta-actions">
            <button class="btn-primary large" @click="gotoDemo">
              立即试用
            </button>
            <button class="btn-outline large" @click="gotoContact">
              联系我们
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
import Header from '@/cms/components/Header.vue'
import Footer from '@/cms/components/Footer.vue'

const router = useRouter()

const features = ref([
  {
    icon: 'User',
    title: '用户管理',
    description: '完善的用户体系，支持多角色权限管理，灵活的组织架构配置'
  },
  {
    icon: 'Lock',
    title: '权限控制',
    description: 'RBAC权限模型，细粒度权限控制，保障系统安全'
  },
  {
    icon: 'ShoppingBag',
    title: '商品管理',
    description: '全面的商品管理功能，支持分类、库存、规格等多维度管理'
  },
  {
    icon: 'Document',
    title: '订单管理',
    description: '完整的订单流程管理，从下单到发货的全链路跟踪'
  },
  {
    icon: 'PieChart',
    title: '数据统计',
    description: '丰富的数据报表，实时统计分析，助力业务决策'
  },
  {
    icon: 'Iphone',
    title: '小程序端',
    description: '配套的微信小程序客户端，为用户提供便捷的移动体验'
  }
])

const techStack = ref([
  { name: 'Vue 3', description: '现代化前端框架' },
  { name: 'Go', description: '高性能后端语言' },
  { name: 'MySQL', description: '稳定的关系型数据库' },
  { name: 'Redis', description: '高效的缓存解决方案' },
  { name: 'Element Plus', description: '优雅的UI组件库' },
  { name: '微信小程序', description: '移动端应用平台' }
])

const stats = ref([
  { number: '50+', label: '核心功能模块' },
  { number: '99.9%', label: '系统稳定性' },
  { number: '1000+', label: '企业用户' },
  { number: '7x24', label: '技术支持' }
])

const gotoDemo = () => {
  router.push('/cms/demo')
}

const gotoFeatures = () => {
  router.push('/cms/features')
}

const gotoContact = () => {
  router.push('/cms/contact')
}
</script>

<style scoped>
.home {
  padding-top: 80px;
}

/* 英雄区域 */
.hero {
  min-height: 80vh;
  display: flex;
  align-items: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  position: relative;
  overflow: hidden;
}

.hero::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="50" cy="50" r="1" fill="%23667eea" fill-opacity="0.05"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>');
  pointer-events: none;
}

.hero .container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 4rem;
  align-items: center;
  position: relative;
  z-index: 1;
}

.hero-title {
  font-size: 3.5rem;
  font-weight: 800;
  line-height: 1.2;
  margin-bottom: 1.5rem;
  color: #2d3748;
}

.hero-subtitle {
  font-size: 1.2rem;
  color: #4a5568;
  margin-bottom: 2.5rem;
  line-height: 1.6;
}

.hero-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.hero-actions .btn-primary,
.hero-actions .btn-outline {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1rem;
}

.dashboard-preview {
  background: white;
  border-radius: 15px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transform: perspective(1000px) rotateY(-15deg) rotateX(10deg);
}

.preview-header {
  background: #f7fafc;
  padding: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.preview-dots {
  display: flex;
  gap: 0.5rem;
}

.preview-dots span {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #cbd5e0;
}

.preview-dots span:first-child {
  background: #f56565;
}

.preview-dots span:nth-child(2) {
  background: #ed8936;
}

.preview-dots span:last-child {
  background: #48bb78;
}

.preview-content {
  display: flex;
  height: 300px;
}

.preview-sidebar {
  width: 80px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.preview-main {
  flex: 1;
  padding: 1.5rem;
  background: #f8fafc;
}

.preview-cards {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.preview-card {
  height: 60px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.preview-chart {
  height: 120px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  position: relative;
  overflow: hidden;
}

.preview-chart::after {
  content: '';
  position: absolute;
  bottom: 20px;
  left: 20px;
  right: 20px;
  height: 60px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px;
  opacity: 0.8;
}

/* 功能特性 */
.features {
  padding: 5rem 0;
  background: white;
}

.section-header {
  text-align: center;
  margin-bottom: 4rem;
}

.section-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
  color: #2d3748;
}

.section-subtitle {
  font-size: 1.1rem;
  color: #718096;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 2rem;
}

.feature-card {
  padding: 2.5rem;
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
  line-height: 1.6;
}

/* 技术栈 */
.tech-stack {
  padding: 5rem 0;
  background: #f8fafc;
}

.tech-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 2rem;
}

.tech-item {
  text-align: center;
  padding: 2rem 1rem;
}

.tech-logo {
  width: 80px;
  height: 80px;
  background: white;
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1rem;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.tech-name {
  font-weight: 600;
  color: #667eea;
}

.tech-desc {
  color: #718096;
  font-size: 0.9rem;
}

/* 统计数据 */
.stats {
  padding: 4rem 0;
  color: white;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 2rem;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.stat-label {
  font-size: 1.1rem;
  opacity: 0.9;
}

/* CTA区域 */
.cta {
  padding: 5rem 0;
  background: white;
  text-align: center;
}

.cta-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
  color: #2d3748;
}

.cta-subtitle {
  font-size: 1.1rem;
  color: #718096;
  margin-bottom: 2.5rem;
}

.cta-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.btn-primary.large,
.btn-outline.large {
  padding: 15px 40px;
  font-size: 1.1rem;
}

/* 响应式 */
@media (max-width: 768px) {
  .hero .container {
    grid-template-columns: 1fr;
    text-align: center;
  }
  
  .hero-title {
    font-size: 2.5rem;
  }
  
  .hero-actions {
    justify-content: center;
    flex-wrap: wrap;
  }
  
  .dashboard-preview {
    transform: none;
    max-width: 400px;
    margin: 0 auto;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .tech-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .cta-actions {
    flex-direction: column;
    align-items: center;
  }
}
</style> 