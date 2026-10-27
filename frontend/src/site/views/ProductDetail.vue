<template>
  <div class="wrap page">
    <p v-if="error" class="note">{{ error }}</p>
    <template v-else-if="goods">
      <div class="crumb">
        <router-link to="/">首页</router-link>
        <router-link to="/products">产品目录</router-link>
        <router-link :to="`/products?categoryId=${goods.categoryId}`">{{ goods.categoryName }}</router-link>
        <span>{{ goods.name }}</span>
      </div>
      <article class="detail">
        <div class="gallery">
          <img :src="current" :alt="goods.name">
          <div class="thumbs" v-if="pictures.length > 1">
            <button v-for="img in pictures" :key="img" :class="{ on: current === img }" @click="current = img">
              <img :src="img" :alt="goods.name">
            </button>
          </div>
        </div>
        <div class="buy">
          <p class="goods-meta">{{ goods.categoryName }}</p>
          <h1>{{ goods.name }}</h1>
          <p class="lede">{{ goods.description }}</p>
          <dl>
            <div v-if="goods.unit"><dt>单位</dt><dd>{{ goods.unit }}</dd></div>
            <div v-if="goods.specDesc"><dt>规格</dt><dd class="pre">{{ goods.specDesc }}</dd></div>
            <div v-if="goods.serviceDesc"><dt>服务</dt><dd>{{ goods.serviceDesc }}</dd></div>
          </dl>
          <router-link class="btn" :to="`/inquiry?goodsId=${goods.id}`">询盘这款</router-link>
        </div>
      </article>
      <section class="floor" v-if="goods.detailDesc">
        <div class="floor-head"><h2>产品说明</h2></div>
        <div class="prose">{{ goods.detailDesc }}</div>
      </section>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getGoodsDetail } from '@/site/api/site'

const route = useRoute()
const goods = ref(null)
const error = ref('')
const current = ref('')
const pictures = computed(() => {
  if (!goods.value) return []
  return goods.value.images?.length ? goods.value.images : [goods.value.cover].filter(Boolean)
})

async function load() {
  error.value = ''
  goods.value = null
  try {
    goods.value = await getGoodsDetail(route.params.id)
    current.value = pictures.value[0] || ''
    document.title = `${goods.value.name}｜集亿玖泽`
  } catch (err) {
    error.value = err.message || '没有这款产品'
  }
}

onMounted(load)
watch(() => route.params.id, load)
</script>
