<template>
  <div class="features-page">
    <Header />
    
    <!-- 页面头部 -->
    <section class="page-hero gradient-bg">
      <div class="container">
        <div class="hero-content" data-aos="fade-up">
          <h1 class="page-title">功能特性</h1>
          <p class="page-subtitle">
            WnAdmin提供全面的管理功能，满足企业各种业务需求
          </p>
        </div>
      </div>
    </section>
    
    <!-- 功能详情 -->
    <section class="feature-details">
      <div class="container">
        <div 
          v-for="(feature, index) in detailedFeatures" 
          :key="index"
          class="feature-section"
          :class="{ 'reverse': index % 2 === 1 }"
        >
          <div class="feature-content" data-aos="fade-right">
            <div class="feature-badge">
              <el-icon><component :is="feature.icon" /></el-icon>
            </div>
            <h2 class="feature-name">{{ feature.title }}</h2>
            <p class="feature-description">{{ feature.description }}</p>
            <ul class="feature-list">
              <li v-for="item in feature.features" :key="item">
                <el-icon><Check /></el-icon>
                {{ item }}
              </li>
            </ul>
          </div>
          <div class="feature-visual" data-aos="fade-left">
            <div class="feature-mockup">
              <div class="mockup-content">
                <div class="mockup-header">
                  <h3>{{ feature.mockupTitle }}</h3>
                </div>
                <div class="mockup-body">
                  <div v-for="n in 4" :key="n" class="mockup-item"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- 技术优势 -->
    <section class="tech-advantages">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2 class="section-title">技术优势</h2>
          <p class="section-subtitle">基于现代化技术栈，确保系统的高性能和稳定性</p>
        </div>
        
        <div class="advantages-grid">
          <div 
            v-for="(advantage, index) in advantages" 
            :key="index"
            class="advantage-card card"
            :data-aos="'zoom-in'"
            :data-aos-delay="index * 100"
          >
            <div class="advantage-icon">
              <el-icon><component :is="advantage.icon" /></el-icon>
            </div>
            <h3 class="advantage-title">{{ advantage.title }}</h3>
            <p class="advantage-desc">{{ advantage.description }}</p>
          </div>
        </div>
      </div>
    </section>
    
    <Footer />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Header from '@/cms/components/Header.vue'
import Footer from '@/cms/components/Footer.vue'

const detailedFeatures = ref([
  {
    icon: 'User',
    title: '用户管理系统',
    description: '完善的用户管理体系，支持多层级组织架构，灵活的角色权限配置，让用户管理变得简单高效。',
    features: [
      '多级组织架构管理',
      '灵活的角色权限分配',
      '用户行为日志追踪',
      '批量用户操作',
      '用户状态管理'
    ],
    mockupTitle: '用户管理界面'
  },
  {
    icon: 'ShoppingBag',
    title: '商品管理系统',
    description: '全面的商品管理功能，支持多规格商品、库存管理、分类管理等，满足电商业务的各种需求。',
    features: [
      '多规格商品管理',
      '库存实时监控',
      '商品分类体系',
      '批量导入导出',
      '商品审核流程'
    ],
    mockupTitle: '商品管理界面'
  },
  {
    icon: 'Document',
    title: '订单管理系统',
    description: '完整的订单处理流程，从订单创建到发货完成的全流程管理，支持多种支付方式和物流配送。',
    features: [
      '订单全流程管理',
      '多支付方式支持',
      '物流信息追踪',
      '退款退货处理',
      '订单数据分析'
    ],
    mockupTitle: '订单管理界面'
  },
  {
    icon: 'PieChart',
    title: '数据统计分析',
    description: '强大的数据统计分析功能，提供丰富的报表和图表，帮助企业做出更好的业务决策。',
    features: [
      '实时数据统计',
      '多维度数据分析',
      '自定义报表生成',
      '数据可视化图表',
      '数据导出功能'
    ],
    mockupTitle: '数据统计界面'
  }
])

const advantages = ref([
  {
    icon: 'Lightning',
    title: '高性能',
    description: 'Go语言后端提供卓越性能，Vue 3前端框架保证流畅体验'
  },
  {
    icon: 'Shield',
    title: '安全可靠',
    description: '完善的权限控制机制，数据加密传输，保障系统安全'
  },
  {
    icon: 'Setting',
    title: '易于扩展',
    description: '模块化设计，支持二次开发，轻松扩展业务功能'
  },
  {
    icon: 'Monitor',
    title: '响应式设计',
    description: '适配各种设备屏幕，提供一致的用户体验'
  }
])
</script>

<style scoped>
.features-page {
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

/* 功能详情 */
.feature-details {
  padding: 5rem 0;
}

.feature-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 4rem;
  align-items: center;
  margin-bottom: 6rem;
}

.feature-section.reverse {
  direction: rtl;
}

.feature-section.reverse > * {
  direction: ltr;
}

.feature-badge {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.5rem;
  margin-bottom: 1.5rem;
}

.feature-name {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 1rem;
  color: #2d3748;
}

.feature-description {
  font-size: 1.1rem;
  color: #718096;
  line-height: 1.6;
  margin-bottom: 2rem;
}

.feature-list {
  list-style: none;
}

.feature-list li {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.8rem;
  color: #4a5568;
}

.feature-list li .el-icon {
  color: #48bb78;
  font-size: 1.2rem;
}

.feature-visual {
  display: flex;
  justify-content: center;
}

.feature-mockup {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 15px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.mockup-content {
  padding: 2rem;
}

.mockup-header {
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 1rem;
  margin-bottom: 1.5rem;
}

.mockup-header h3 {
  font-size: 1.2rem;
  font-weight: 600;
  color: #2d3748;
}

.mockup-body {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.mockup-item {
  height: 40px;
  background: linear-gradient(90deg, #f7fafc 0%, #edf2f7 50%, #f7fafc 100%);
  border-radius: 8px;
  position: relative;
  overflow: hidden;
}

.mockup-item::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.6), transparent);
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% { left: -100%; }
  100% { left: 100%; }
}

/* 技术优势 */
.tech-advantages {
  padding: 5rem 0;
  background: #f8fafc;
}

.advantages-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
}

.advantage-card {
  padding: 2.5rem 2rem;
  text-align: center;
}

.advantage-icon {
  width: 70px;
  height: 70px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1.5rem;
  color: white;
  font-size: 1.8rem;
}

.advantage-title {
  font-size: 1.3rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #2d3748;
}

.advantage-desc {
  color: #718096;
  line-height: 1.6;
}

/* 响应式 */
@media (max-width: 768px) {
  .page-title {
    font-size: 2rem;
  }
  
  .feature-section {
    grid-template-columns: 1fr;
    gap: 2rem;
    text-align: center;
  }
  
  .feature-section.reverse {
    direction: ltr;
  }
  
  .advantages-grid {
    grid-template-columns: 1fr;
  }
}
</style> 