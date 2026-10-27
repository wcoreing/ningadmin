<template>
  <div class="site-root">
    <div class="topbar">
      <div class="wrap topbar-inner">
        <p>欢迎来到{{ profile.shortName || '集亿玖泽' }}<template v-if="profile.phone"> · {{ profile.contactName }} {{ profile.phone }}</template></p>
        <nav>
          <router-link to="/">首页</router-link>
          <router-link to="/products">产品目录</router-link>
          <router-link to="/news">资讯</router-link>
          <router-link to="/page/about">关于我们</router-link>
          <router-link to="/inquiry">询盘</router-link>
        </nav>
      </div>
    </div>
    <header class="mast">
      <div class="wrap mast-inner">
        <router-link to="/" class="logo">
          <i>集</i>
          <span>
            <strong>{{ profile.shortName || '集亿玖泽' }}</strong>
            <em>包装 · 箱包</em>
          </span>
        </router-link>
        <form class="searchbox" @submit.prevent="goSearch">
          <input v-model="keyword" placeholder="搜索首饰盒、化妆包、背包" />
          <button type="submit">搜索</button>
        </form>
        <router-link class="ask" to="/inquiry">我要询盘</router-link>
      </div>
    </header>
    <nav class="catbar">
      <div class="wrap catbar-inner">
        <router-link class="all" to="/products">全部品类</router-link>
        <router-link v-for="item in categories" :key="item.id" :to="`/products?categoryId=${item.id}`">{{ item.name }}</router-link>
      </div>
    </nav>
    <main>
      <router-view :profile="profile" />
    </main>
    <footer>
      <div class="wrap foot-grid">
        <div>
          <h3>关于集亿</h3>
          <router-link to="/page/about">公司简介</router-link>
          <router-link to="/page/culture">做事方式</router-link>
          <router-link to="/page/honor">服务范围</router-link>
          <router-link to="/page/history">现在在做的事</router-link>
        </div>
        <div>
          <h3>采购</h3>
          <router-link to="/products">产品目录</router-link>
          <router-link to="/inquiry">提交询盘</router-link>
          <router-link to="/news">资讯</router-link>
          <router-link to="/privacy">隐私说明</router-link>
        </div>
        <div>
          <h3>联系</h3>
          <p v-if="profile.phone">{{ profile.contactName }} {{ profile.phone }}</p>
          <p v-if="profile.email">{{ profile.email }}</p>
          <p>{{ profile.address }}</p>
        </div>
      </div>
      <div class="foot-bar">
        <div class="wrap">
          <span>{{ profile.name || '深圳市集亿玖泽贸易有限公司' }}</span>
          <span>v0.1.7</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCategories, getProfile } from '@/site/api/site'
import '@/site/styles/site.css'

const route = useRoute()
const router = useRouter()
const keyword = ref(route.query.keyword || '')
const categories = ref([])
const profile = reactive({
  name: '',
  shortName: '',
  phone: '',
  contactName: '',
  email: '',
  address: '',
  summary: ''
})

watch(() => route.query.keyword, value => {
  keyword.value = value || ''
})

onMounted(async () => {
  try {
    Object.assign(profile, await getProfile())
  } catch (error) {
    console.error(error)
  }
  try {
    categories.value = (await getCategories()) || []
  } catch (error) {
    console.error(error)
  }
})

function goSearch() {
  router.push({ path: '/products', query: keyword.value ? { keyword: keyword.value } : {} })
}
</script>
