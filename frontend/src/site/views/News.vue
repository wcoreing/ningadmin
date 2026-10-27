<template>
  <div class="wrap page">
    <div class="crumb">
      <router-link to="/">首页</router-link>
      <span>资讯</span>
    </div>
    <section class="floor">
      <div class="floor-head">
        <h2>资讯</h2>
        <div class="filters">
          <button :class="{ on: !category }" @click="category = ''">全部</button>
          <button :class="{ on: category === 'company' }" @click="category = 'company'">公司动态</button>
          <button :class="{ on: category === 'industry' }" @click="category = 'industry'">行业资料</button>
        </div>
      </div>
      <p v-if="error" class="note">{{ error }}</p>
      <div class="news-list" v-else>
        <router-link v-for="item in list" :key="item.id" :to="`/news/${item.id}`">
          <img :src="item.cover" :alt="item.title">
          <div>
            <h3>{{ item.title }}</h3>
            <p>{{ item.summary }}</p>
            <span>{{ item.categoryName }} · {{ item.publishTime }}</span>
          </div>
        </router-link>
      </div>
      <p v-if="!error && !list.length" class="empty">还没有资讯。</p>
    </section>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { getArticles } from '@/site/api/site'

const list = ref([])
const category = ref('')
const error = ref('')

async function load() {
  error.value = ''
  try {
    const data = await getArticles({ category: category.value || undefined, page: 1, pageSize: 20 })
    list.value = data.list || []
  } catch (err) {
    error.value = err.message || '资讯读不出来'
  }
}

onMounted(load)
watch(category, load)
</script>
