<template>
  <header class="header" :class="{ 'scrolled': isScrolled }">
    <div class="container">
      <nav class="navbar">
        <div class="logo">
          <router-link to="/" class="logo-link">
            <span class="logo-text gradient-text">WnAdmin</span>
          </router-link>
        </div>
        
        <div class="nav-menu" :class="{ 'active': isMenuOpen }">
          <router-link 
            v-for="item in menuItems" 
            :key="item.path" 
            :to="item.path" 
            class="nav-link"
            @click="closeMenu"
          >
            {{ item.name }}
          </router-link>
        </div>
        
        <div class="nav-actions">
          <button class="btn-outline demo-btn" @click="gotoDemo">
            在线演示
          </button>
          <button class="menu-toggle" @click="toggleMenu">
            <span></span>
            <span></span>
            <span></span>
          </button>
        </div>
      </nav>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isScrolled = ref(false)
const isMenuOpen = ref(false)

const menuItems = [
  { name: '首页', path: '/' },
  { name: '功能特性', path: '/features' },
  { name: '技术架构', path: '/technology' },
  { name: '开发文档', path: '/docs' },
  { name: '联系我们', path: '/contact' }
]

const handleScroll = () => {
  isScrolled.value = window.scrollY > 50
}

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value
}

const closeMenu = () => {
  isMenuOpen.value = false
}

const gotoDemo = () => {
  router.push('/demo')
  closeMenu()
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
}

.header.scrolled {
  background: rgba(255, 255, 255, 0.98);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
}

.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 0;
}

.logo-link {
  text-decoration: none;
}

.logo-text {
  font-size: 2rem;
  font-weight: 700;
  letter-spacing: -1px;
}

.nav-menu {
  display: flex;
  align-items: center;
  gap: 2rem;
}

.nav-link {
  text-decoration: none;
  color: #333;
  font-weight: 500;
  position: relative;
  padding: 0.5rem 0;
  transition: color 0.3s ease;
}

.nav-link:hover {
  color: #667eea;
}

.nav-link.router-link-active {
  color: #667eea;
}

.nav-link::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 0;
  height: 2px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s ease;
}

.nav-link:hover::after,
.nav-link.router-link-active::after {
  width: 100%;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.demo-btn {
  font-size: 14px;
}

.menu-toggle {
  display: none;
  flex-direction: column;
  justify-content: space-around;
  width: 24px;
  height: 24px;
  background: transparent;
  border: none;
  cursor: pointer;
}

.menu-toggle span {
  width: 100%;
  height: 2px;
  background: #333;
  transition: all 0.3s ease;
}

@media (max-width: 768px) {
  .nav-menu {
    position: fixed;
    top: 70px;
    left: 0;
    right: 0;
    background: white;
    flex-direction: column;
    padding: 2rem;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
    transform: translateY(-100%);
    opacity: 0;
    visibility: hidden;
    transition: all 0.3s ease;
  }
  
  .nav-menu.active {
    transform: translateY(0);
    opacity: 1;
    visibility: visible;
  }
  
  .menu-toggle {
    display: flex;
  }
  
  .demo-btn {
    display: none;
  }
}
</style> 