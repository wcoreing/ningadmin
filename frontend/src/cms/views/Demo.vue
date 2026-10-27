<template>
  <div class="demo-page">
    <Header />
    
    <!-- 页面头部 -->
    <section class="page-hero gradient-bg">
      <div class="container">
        <div class="hero-content" data-aos="fade-up">
          <h1 class="page-title">在线演示</h1>
          <p class="page-subtitle">
            体验WnAdmin的完整功能，感受现代化管理系统的强大能力
          </p>
        </div>
      </div>
    </section>
    
    <!-- 演示入口 -->
    <section class="demo-entries">
      <div class="container">
        <div class="entries-grid">
          <div 
            v-for="(demo, index) in demoList" 
            :key="index"
            class="demo-card card"
            :data-aos="'fade-up'"
            :data-aos-delay="index * 100"
          >
            <div class="demo-header">
              <div class="demo-icon">
                <el-icon><component :is="demo.icon" /></el-icon>
              </div>
              <h3 class="demo-title">{{ demo.title }}</h3>
              <p class="demo-desc">{{ demo.description }}</p>
            </div>
            
            <div class="demo-info">
              <div class="info-item">
                <span class="info-label">账号:</span>
                <span class="info-value">{{ demo.account }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">密码:</span>
                <span class="info-value">{{ demo.password }}</span>
              </div>
            </div>
            
            <div class="demo-actions">
              <button class="btn-primary full-width" @click="openDemo(demo.url)">
                <el-icon><View /></el-icon>
                立即体验
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- 功能预览 -->
    <section class="feature-preview">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2 class="section-title">功能预览</h2>
          <p class="section-subtitle">了解系统的核心功能模块</p>
        </div>
        
        <div class="preview-tabs">
          <div class="tab-headers">
            <button 
              v-for="(tab, index) in previewTabs" 
              :key="index"
              class="tab-header"
              :class="{ 'active': activeTab === index }"
              @click="activeTab = index"
            >
              <el-icon><component :is="tab.icon" /></el-icon>
              {{ tab.title }}
            </button>
          </div>
          
          <div class="tab-content" data-aos="fade-up">
            <div class="content-left">
              <h3 class="content-title">{{ previewTabs[activeTab].title }}</h3>
              <p class="content-desc">{{ previewTabs[activeTab].description }}</p>
              <ul class="feature-list">
                <li v-for="feature in previewTabs[activeTab].features" :key="feature">
                  <el-icon><Check /></el-icon>
                  {{ feature }}
                </li>
              </ul>
            </div>
            <div class="content-right">
              <div class="screenshot">
                <div class="screenshot-header">
                  <div class="dots">
                    <span></span>
                    <span></span>
                    <span></span>
                  </div>
                  <div class="url-bar">{{ previewTabs[activeTab].url }}</div>
                </div>
                <div class="screenshot-body">
                  <div class="mock-interface">
                    <div class="interface-sidebar"></div>
                    <div class="interface-main">
                      <div class="interface-header"></div>
                      <div class="interface-content">
                        <div v-for="n in 6" :key="n" class="content-block"></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- 使用说明 -->
    <section class="usage-guide">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2 class="section-title">使用说明</h2>
          <p class="section-subtitle">快速上手指南</p>
        </div>
        
        <div class="guide-steps">
          <div 
            v-for="(step, index) in guideSteps" 
            :key="index"
            class="step-item"
            :data-aos="'fade-right'"
            :data-aos-delay="index * 100"
          >
            <div class="step-number">{{ index + 1 }}</div>
            <div class="step-content">
              <h3 class="step-title">{{ step.title }}</h3>
              <p class="step-desc">{{ step.description }}</p>
            </div>
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

const activeTab = ref(0)

const demoList = ref([
  {
    icon: 'Monitor',
    title: '管理后台',
    description: '完整的后台管理系统，包含所有核心功能',
    account: 'admin',
    password: '123456',
    url: 'http://admin.wnadmin.com'
  },
  {
    icon: 'Iphone',
    title: '小程序端',
    description: '微信小程序客户端，扫码体验移动端功能',
    account: '微信扫码',
    password: '无需密码',
    url: '#'
  },
  {
    icon: 'Globe',
    title: 'API文档',
    description: '完整的API接口文档，方便二次开发',
    account: '无需登录',
    password: '直接访问',
    url: 'http://api.wnadmin.com/docs'
  }
])

const previewTabs = ref([
  {
    icon: 'User',
    title: '用户管理',
    description: '完善的用户管理体系，支持多层级组织架构和角色权限管理',
    url: '/admin/user/list',
    features: [
      '用户信息管理',
      '角色权限分配',
      '组织架构管理',
      '用户行为日志',
      '批量操作功能'
    ]
  },
  {
    icon: 'ShoppingBag',
    title: '商品管理',
    description: '全面的商品管理功能，支持多规格商品和库存管理',
    url: '/admin/goods/list',
    features: [
      '商品信息维护',
      '多规格商品',
      '库存管理',
      '商品分类',
      '批量导入导出'
    ]
  },
  {
    icon: 'Document',
    title: '订单管理',
    description: '完整的订单处理流程，支持订单状态跟踪和处理',
    url: '/admin/order/list',
    features: [
      '订单列表管理',
      '订单状态更新',
      '物流信息',
      '退款处理',
      '订单统计'
    ]
  },
  {
    icon: 'PieChart',
    title: '数据统计',
    description: '丰富的数据统计和分析功能，助力业务决策',
    url: '/admin/statistics',
    features: [
      '销售数据统计',
      '用户行为分析',
      '商品销量排行',
      '收入趋势图',
      '自定义报表'
    ]
  }
])

const guideSteps = ref([
  {
    title: '选择演示系统',
    description: '根据您的需求选择管理后台、小程序或API文档进行体验'
  },
  {
    title: '使用演示账号',
    description: '使用提供的演示账号和密码登录系统，开始体验功能'
  },
  {
    title: '探索核心功能',
    description: '浏览用户管理、商品管理、订单处理等核心功能模块'
  },
  {
    title: '联系我们',
    description: '如有疑问或需要定制开发，随时联系我们的技术团队'
  }
])

const openDemo = (url) => {
  if (url === '#') {
    alert('小程序演示请扫描二维码')
  } else {
    window.open(url, '_blank')
  }
}
</script>

<style scoped>
.demo-page {
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

/* 演示入口 */
.demo-entries {
  padding: 5rem 0;
  background: white;
}

.entries-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.demo-card {
  padding: 2.5rem 2rem;
  text-align: center;
}

.demo-header {
  margin-bottom: 2rem;
}

.demo-icon {
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

.demo-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #2d3748;
}

.demo-desc {
  color: #718096;
  line-height: 1.6;
}

.demo-info {
  margin-bottom: 2rem;
  text-align: left;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 0.8rem 0;
  border-bottom: 1px solid #e2e8f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-weight: 500;
  color: #4a5568;
}

.info-value {
  color: #667eea;
  font-weight: 500;
}

.btn-primary.full-width {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

/* 功能预览 */
.feature-preview {
  padding: 5rem 0;
  background: #f8fafc;
}

.preview-tabs {
  max-width: 1000px;
  margin: 0 auto;
}

.tab-headers {
  display: flex;
  background: white;
  border-radius: 15px;
  padding: 0.5rem;
  margin-bottom: 3rem;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  overflow-x: auto;
}

.tab-header {
  flex: 1;
  padding: 1rem 1.5rem;
  border: none;
  background: transparent;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  white-space: nowrap;
  min-width: 120px;
}

.tab-header.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.tab-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3rem;
  align-items: center;
}

.content-title {
  font-size: 1.8rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #2d3748;
}

.content-desc {
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

.screenshot {
  background: white;
  border-radius: 15px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.screenshot-header {
  background: #f7fafc;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.dots {
  display: flex;
  gap: 0.3rem;
}

.dots span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #cbd5e0;
}

.dots span:first-child {
  background: #f56565;
}

.dots span:nth-child(2) {
  background: #ed8936;
}

.dots span:last-child {
  background: #48bb78;
}

.url-bar {
  background: white;
  padding: 0.3rem 0.8rem;
  border-radius: 15px;
  font-size: 0.8rem;
  color: #718096;
  border: 1px solid #e2e8f0;
}

.screenshot-body {
  height: 300px;
  background: #f8fafc;
}

.mock-interface {
  display: flex;
  height: 100%;
}

.interface-sidebar {
  width: 60px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.interface-main {
  flex: 1;
  padding: 1rem;
}

.interface-header {
  height: 40px;
  background: white;
  border-radius: 6px;
  margin-bottom: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.interface-content {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.8rem;
}

.content-block {
  height: 60px;
  background: white;
  border-radius: 6px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

/* 使用说明 */
.usage-guide {
  padding: 5rem 0;
  background: white;
}

.guide-steps {
  max-width: 800px;
  margin: 0 auto;
}

.step-item {
  display: flex;
  align-items: flex-start;
  gap: 2rem;
  margin-bottom: 3rem;
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
  flex-shrink: 0;
}

.step-title {
  font-size: 1.3rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: #2d3748;
}

.step-desc {
  color: #718096;
  line-height: 1.6;
}

/* 响应式 */
@media (max-width: 768px) {
  .page-title {
    font-size: 2rem;
  }
  
  .entries-grid {
    grid-template-columns: 1fr;
  }
  
  .tab-content {
    grid-template-columns: 1fr;
    gap: 2rem;
  }
  
  .step-item {
    gap: 1rem;
  }
  
  .step-number {
    width: 50px;
    height: 50px;
    font-size: 1.2rem;
  }
}
</style> 