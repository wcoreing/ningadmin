<template>
  <div class="wrap page">
    <div class="crumb">
      <router-link to="/">首页</router-link>
      <span>询盘</span>
    </div>
    <section class="floor form-panel">
      <div class="floor-head">
        <h2>提交询盘</h2>
      </div>
      <p class="lede" v-if="goodsName">正在询盘：{{ goodsName }}</p>
      <p class="lede" v-else>写下要装的东西、数量和发往哪里。我们按电话回复。</p>
      <form class="form" @submit.prevent="submit">
        <div class="form-row">
          <label>联系人<input v-model="form.name" required></label>
          <label>电话<input v-model="form.phone" required></label>
        </div>
        <div class="form-row">
          <label>公司<input v-model="form.company"></label>
          <label>邮箱<input v-model="form.email" type="email"></label>
        </div>
        <div class="form-row">
          <label>目标市场<input v-model="form.country" placeholder="例如：国内 / 东南亚"></label>
          <label>数量<input v-model="form.quantity" placeholder="例如：500 个"></label>
        </div>
        <label>需求说明<textarea v-model="form.message" required placeholder="尺寸、材质、颜色，或者你想改的地方"></textarea></label>
        <p v-if="error" class="note">{{ error }}</p>
        <p v-if="done" class="ok">已经记下。请保持电话畅通。</p>
        <button class="btn" type="submit" :disabled="sending">{{ sending ? '提交中' : '提交询盘' }}</button>
      </form>
    </section>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import { createInquiry, getGoodsDetail } from '@/site/api/site'

const route = useRoute()
const goodsName = ref('')
const error = ref('')
const done = ref(false)
const sending = ref(false)
const form = reactive({
  goodsId: route.query.goodsId ? Number(route.query.goodsId) : 0,
  name: '',
  phone: '',
  company: '',
  email: '',
  country: '',
  quantity: '',
  message: ''
})

onMounted(async () => {
  document.title = '询盘｜集亿玖泽'
  if (!form.goodsId) return
  try {
    const goods = await getGoodsDetail(form.goodsId)
    goodsName.value = goods.name
    if (!form.message) form.message = `想了解「${goods.name}」的规格和起订量。`
  } catch (err) {
    error.value = err.message || '这款产品没有找到'
  }
})

async function submit() {
  error.value = ''
  done.value = false
  sending.value = true
  try {
    await createInquiry(form)
    done.value = true
  } catch (err) {
    error.value = err.message || '提交失败'
  } finally {
    sending.value = false
  }
}
</script>
