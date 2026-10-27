<template>
  <div class="wrap home" v-if="home">
    <p v-if="error" class="note">{{ error }}</p>
    <section class="promo">
      <aside class="cate-panel">
        <router-link v-for="item in home.categories" :key="item.id" :to="`/products?categoryId=${item.id}`">
          <img v-if="item.image" :src="item.image" alt="" @error="hideImg">
          <span>{{ item.name }}</span>
        </router-link>
      </aside>
      <div class="promo-main">
        <a v-if="banner && banner.url && !banner.url.startsWith('/')" :href="banner.url">
          <img class="banner-img" :src="banner.image" :alt="banner.title">
        </a>
        <router-link v-else-if="banner" :to="banner.url || '/products'">
          <img class="banner-img" :src="banner.image" :alt="banner.title">
        </router-link>
        <div v-else class="promo-copy">
          <p class="tag">包装与箱包</p>
          <h1>{{ home.profile.shortName || '集亿玖泽' }}</h1>
          <p>{{ home.profile.summary }}</p>
          <router-link class="btn" to="/products">进入产品目录</router-link>
        </div>
        <div class="dots" v-if="home.banners && home.banners.length > 1">
          <button v-for="(item, index) in home.banners" :key="item.id" :class="{ on: index === bannerIndex }" @click="bannerIndex = index"></button>
        </div>
      </div>
      <aside class="service-card">
        <h2>采购咨询</h2>
        <p><b>联系人</b>{{ home.profile.contactName }}</p>
        <p><b>电话</b>{{ home.profile.phone }}</p>
        <p><b>地址</b>{{ home.profile.address }}</p>
        <router-link class="btn block" to="/inquiry">我要询盘</router-link>
      </aside>
    </section>

    <section class="floor">
      <div class="floor-head">
        <h2>产品分类</h2>
        <router-link to="/products">更多</router-link>
      </div>
      <div class="cat-grid">
        <router-link v-for="item in home.categories" :key="item.id" :to="`/products?categoryId=${item.id}`">
          <img v-if="item.image" :src="item.image" :alt="item.name" @error="hideImg">
          <span>{{ item.name }}</span>
        </router-link>
      </div>
    </section>

    <section class="floor">
      <div class="floor-head">
        <h2>推荐产品</h2>
        <router-link to="/products">更多</router-link>
      </div>
      <div class="goods-grid">
        <router-link class="goods" v-for="item in home.featured" :key="item.id" :to="`/products/${item.id}`">
          <div class="goods-pic"><img :src="item.cover" :alt="item.name"></div>
          <p class="goods-name">{{ item.name }}</p>
          <p class="goods-meta">{{ item.categoryName }}</p>
          <span class="goods-ask">询盘</span>
        </router-link>
      </div>
    </section>

    <section class="floor" v-if="home.pages.length">
      <div class="floor-head"><h2>公司资料</h2></div>
      <div class="info-grid">
        <router-link v-for="item in home.pages" :key="item.code" :to="`/page/${item.code}`">
          <img v-if="item.cover" :src="item.cover" :alt="item.title">
          <div>
            <h3>{{ item.title }}</h3>
            <p>{{ item.summary }}</p>
          </div>
        </router-link>
      </div>
    </section>

    <section class="floor" v-if="home.articles.length">
      <div class="floor-head">
        <h2>资讯</h2>
        <router-link to="/news">更多</router-link>
      </div>
      <div class="news-list">
        <router-link v-for="item in home.articles" :key="item.id" :to="`/news/${item.id}`">
          <img :src="item.cover" :alt="item.title">
          <div>
            <h3>{{ item.title }}</h3>
            <p>{{ item.summary }}</p>
            <span>{{ item.categoryName }} · {{ item.publishTime }}</span>
          </div>
        </router-link>
      </div>
    </section>
  </div>
  <p v-else-if="error" class="wrap note">{{ error }}</p>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { getHome } from '@/site/api/site'

const home = ref(null)
const error = ref('')
const bannerIndex = ref(0)
const banner = computed(() => (home.value?.banners || [])[bannerIndex.value] || null)

onMounted(async () => {
  try {
    home.value = await getHome()
    if (home.value?.profile?.seoTitle) document.title = home.value.profile.seoTitle
    const seo = home.value?.profile
    if (seo?.seoDescription) {
      ensureMeta('description', seo.seoDescription)
      ensureMeta('keywords', seo.seoKeywords)
    }
  } catch (err) {
    error.value = err.message || '首页暂时读不出来'
  }
})

function hideImg(event) {
  event.target.style.display = 'none'
}

function ensureMeta(name, content) {
  let tag = document.querySelector(`meta[name="${name}"]`)
  if (!tag) {
    tag = document.createElement('meta')
    tag.setAttribute('name', name)
    document.head.appendChild(tag)
  }
  tag.setAttribute('content', content || '')
}
</script>
