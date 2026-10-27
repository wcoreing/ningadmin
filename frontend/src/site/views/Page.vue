<template>
  <div class="wrap page">
    <div class="crumb">
      <router-link to="/">首页</router-link>
      <span>{{ page?.title || '资料' }}</span>
    </div>
    <p v-if="error" class="note">{{ error }}</p>
    <div class="story-grid" v-else-if="page">
      <aside class="story-nav">
        <h2>资料</h2>
        <router-link
          v-for="item in pages"
          :key="item.code"
          :to="`/page/${item.code}`"
          :class="{ on: item.code === page.code }"
        >{{ item.title }}</router-link>
        <div class="story-call">
          <b>直接问这一单</b>
          <p v-if="profile.contactName || profile.phone">{{ profile.contactName }} {{ profile.phone }}</p>
          <p v-else>留下数量和发往哪里，按电话回复。</p>
          <router-link class="btn" to="/inquiry">我要询盘</router-link>
        </div>
      </aside>
      <article class="story">
        <figure class="story-lid" v-if="coverOn && page.cover">
          <img :src="page.cover" :alt="page.title" @error="coverOn = false">
        </figure>
        <div class="story-main">
          <p class="story-kicker">公司资料</p>
          <h1>{{ page.title }}</h1>
          <p class="story-lead" v-if="page.summary">{{ page.summary }}</p>
          <div class="prose story-body" v-html="page.content"></div>
          <div class="story-actions">
            <router-link class="btn" to="/products">看产品目录</router-link>
            <router-link class="ask" to="/inquiry">我要询盘</router-link>
          </div>
        </div>
      </article>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getHome, getPage } from '@/site/api/site'

defineProps({
  profile: { type: Object, default: () => ({}) }
})

const route = useRoute()
const page = ref(null)
const pages = ref([])
const error = ref('')
const coverOn = ref(true)

async function load() {
  error.value = ''
  page.value = null
  coverOn.value = true
  try {
    page.value = await getPage(route.params.code)
    document.title = page.value.seoTitle || page.value.title
    setMeta('description', page.value.seoDescription)
    setMeta('keywords', page.value.seoKeywords)
  } catch (err) {
    error.value = err.message || '栏目不存在'
  }
}

function setMeta(name, content) {
  let tag = document.querySelector(`meta[name="${name}"]`)
  if (!tag) {
    tag = document.createElement('meta')
    tag.setAttribute('name', name)
    document.head.appendChild(tag)
  }
  tag.setAttribute('content', content || '')
}

onMounted(async () => {
  load()
  try {
    const home = await getHome()
    pages.value = home.pages || []
  } catch (err) {
    pages.value = []
  }
})
watch(() => route.params.code, load)
</script>
