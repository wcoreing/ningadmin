<template>
  <div class="wrap page">
    <div class="crumb">
      <router-link to="/">首页</router-link>
      <span>产品目录</span>
    </div>
    <div class="catalog">
      <aside class="side">
        <h2>全部品类</h2>
        <button :class="{ on: !categoryId }" @click="choose('')">全部</button>
        <button v-for="item in categories" :key="item.id" :class="{ on: categoryId === item.id }" @click="choose(item.id)">
          {{ item.name }}
        </button>
      </aside>
      <div>
        <div class="toolbar">
          <h1>{{ currentName }}</h1>
          <span v-if="!error">{{ total }} 款</span>
        </div>
        <p v-if="error" class="note">{{ error }}</p>
        <div class="goods-grid" v-else>
          <router-link class="goods" v-for="item in list" :key="item.id" :to="`/products/${item.id}`">
            <div class="goods-pic"><img :src="item.cover" :alt="item.name"></div>
            <p class="goods-name">{{ item.name }}</p>
            <p class="goods-meta">{{ item.categoryName }}</p>
            <p class="goods-desc">{{ item.description }}</p>
            <span class="goods-ask">询盘</span>
          </router-link>
        </div>
        <p v-if="!error && !list.length" class="empty">这个分类下还没有产品。</p>
        <div class="pager" v-if="pageCount > 1">
          <button type="button" :disabled="page <= 1" @click="go(page - 1)">上一页</button>
          <button type="button" v-for="n in pageCount" :key="n" :class="{ on: n === page }" @click="go(n)">{{ n }}</button>
          <button type="button" :disabled="page >= pageCount" @click="go(page + 1)">下一页</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCategories, getGoods } from '@/site/api/site'

const route = useRoute()
const router = useRouter()
const categories = ref([])
const list = ref([])
const total = ref(0)
const pageSize = 12
const categoryId = ref(route.query.categoryId ? Number(route.query.categoryId) : '')
const keyword = ref(route.query.keyword || '')
const page = ref(Math.max(1, Number(route.query.page) || 1))
const error = ref('')
const currentName = computed(() => categories.value.find(item => item.id === categoryId.value)?.name || (keyword.value ? `“${keyword.value}”` : '全部产品'))
const pageCount = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))

onMounted(async () => {
  document.title = '产品目录｜集亿玖泽'
  categories.value = await getCategories()
  await load()
})

watch(() => [route.query.categoryId, route.query.keyword, route.query.page], () => {
  categoryId.value = route.query.categoryId ? Number(route.query.categoryId) : ''
  keyword.value = route.query.keyword || ''
  page.value = Math.max(1, Number(route.query.page) || 1)
  load()
})

function choose(id) {
  router.replace({
    path: '/products',
    query: {
      ...(id ? { categoryId: id } : {}),
      ...(keyword.value ? { keyword: keyword.value } : {})
    }
  })
}

function go(n) {
  const next = Math.min(Math.max(1, n), pageCount.value)
  router.replace({
    path: '/products',
    query: {
      ...(categoryId.value ? { categoryId: categoryId.value } : {}),
      ...(keyword.value ? { keyword: keyword.value } : {}),
      ...(next > 1 ? { page: next } : {})
    }
  })
}

async function load() {
  error.value = ''
  try {
    const data = await getGoods({
      categoryId: categoryId.value || undefined,
      keyword: keyword.value || undefined,
      page: page.value,
      pageSize
    })
    list.value = data.list || []
    total.value = data.total || list.value.length
  } catch (err) {
    error.value = err.message || '产品读不出来'
  }
}
</script>
