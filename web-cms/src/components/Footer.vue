<template>
  <footer class="footer gradient-bg">
    <div class="container">
      <div class="footer-content">
        <div class="footer-section">
          <h3 class="footer-title">WnAdmin</h3>
          <p class="footer-desc">
            现代化全栈管理系统，提供完整的前后端解决方案，
            助力企业快速构建高效的管理平台。
          </p>
          <div class="social-links">
            <a href="#" class="social-link">
              <el-icon><Platform /></el-icon>
            </a>
            <a href="#" class="social-link">
              <el-icon><Connection /></el-icon>
            </a>
            <a href="#" class="social-link">
              <el-icon><Message /></el-icon>
            </a>
          </div>
        </div>
        
        <div class="footer-section">
          <h4 class="section-title">产品功能</h4>
          <ul class="footer-links">
            <li><router-link to="/features">用户管理</router-link></li>
            <li><router-link to="/features">权限控制</router-link></li>
            <li><router-link to="/features">商品管理</router-link></li>
            <li><router-link to="/features">订单管理</router-link></li>
            <li><router-link to="/features">数据统计</router-link></li>
          </ul>
        </div>
        
        <div class="footer-section">
          <h4 class="section-title">技术栈</h4>
          <ul class="footer-links">
            <li><a href="#">Vue 3</a></li>
            <li><a href="#">Go语言</a></li>
            <li><a href="#">MySQL</a></li>
            <li><a href="#">Redis</a></li>
            <li><a href="#">微信小程序</a></li>
          </ul>
        </div>
        
        <div class="footer-section">
          <h4 class="section-title">帮助与支持</h4>
          <ul class="footer-links">
            <li><router-link to="/docs">开发文档</router-link></li>
            <li><router-link to="/demo">在线演示</router-link></li>
            <li><router-link to="/contact">技术支持</router-link></li>
            <li><router-link to="/contact">联系我们</router-link></li>
          </ul>
        </div>
      </div>
      
      <div class="footer-bottom">
        <div class="copyright">
          <p>&copy; 2024 WnAdmin. 保留所有权利</p>
        </div>
        <div class="footer-nav">
          <router-link to="/contact">隐私政策</router-link>
          <router-link to="/contact">服务条款</router-link>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup>
// Footer组件
</script>

<style scoped>
.footer {
  color: white;
  padding: 3rem 0 1rem;
  margin-top: 4rem;
}

.footer-content {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 2rem;
  margin-bottom: 2rem;
}

.footer-section {
  display: flex;
  flex-direction: column;
}

.footer-title {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 1rem;
  letter-spacing: -1px;
}

.footer-desc {
  color: rgba(255, 255, 255, 0.9);
  line-height: 1.6;
  margin-bottom: 1.5rem;
}

.social-links {
  display: flex;
  gap: 1rem;
}

.social-link {
  width: 40px;
  height: 40px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  text-decoration: none;
  transition: all 0.3s ease;
}

.social-link:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.section-title {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: white;
}

.footer-links {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.footer-links a {
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  transition: color 0.3s ease;
}

.footer-links a:hover {
  color: white;
}

.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 2rem;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
}

.copyright p {
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.9rem;
}

.footer-nav {
  display: flex;
  gap: 2rem;
}

.footer-nav a {
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  font-size: 0.9rem;
  transition: color 0.3s ease;
}

.footer-nav a:hover {
  color: white;
}

@media (max-width: 968px) {
  .footer-content {
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
  }
}

@media (max-width: 768px) {
  .footer-content {
    grid-template-columns: 1fr;
    text-align: center;
  }
  
  .footer-bottom {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .footer-nav {
    justify-content: center;
  }
}
</style> 