<template>
  <div class="wrap page">
    <div class="crumb">
      <router-link to="/">首页</router-link>
      <router-link to="/news">资讯</router-link>
      <span>{{ article?.title || '资讯' }}</span>
    </div>
    <p v-if="error" class="note">{{ error }}</p>
    <template v-else-if="article">
      <article class="story">
        <figure class="story-lid" v-if="coverOn && article.cover">
          <img :src="article.cover" :alt="article.title" @error="coverOn = false">
        </figure>
        <div class="story-main">
          <p class="story-kicker">{{ article.categoryName || '资讯' }}</p>
          <h1>{{ article.title }}</h1>
          <p class="story-meta">{{ article.publishTime }}</p>
          <p class="story-lead" v-if="article.summary">{{ article.summary }}</p>
          <div class="prose story-body" v-html="article.content"></div>
          <div class="story-actions">
            <router-link class="btn" to="/news">返回资讯</router-link>
            <router-link class="ask" to="/inquiry">按这篇询盘</router-link>
          </div>
        </div>
      </article>
      <section class="story-more" v-if="more.length">
        <h2>其他资讯</h2>
        <router-link v-for="item in more" :key="item.id" :to="`/news/${item.id}`">
          <img v-if="item.cover" :src="item.cover" :alt="item.title">
          <div>
            <h3>{{ item.title }}</h3>
            <p>{{ item.categoryName }} · {{ item.publishTime }}</p>
          </div>
        </router-link>
      </section>
    </template>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getArticle, getArticles } from '@/site/api/site'

const route = useRoute()
const article = ref(null)
const more = ref([])
const error = ref('')
const coverOn = ref(true)

async function load() {
  error.value = ''
  article.value = null
  coverOn.value = true
  try {
    article.value = await getArticle(route.params.id)
    document.title = article.value.seoTitle || article.value.title
    const data = await getArticles({ page: 1, pageSize: 6 })
    more.value = (data.list || []).filter(item => String(item.id) !== String(route.params.id)).slice(0, 3)
  } catch (err) {
    error.value = err.message || '资讯不存在'
  }
}

onMounted(load)
watch(() => route.params.id, load)
</script>
