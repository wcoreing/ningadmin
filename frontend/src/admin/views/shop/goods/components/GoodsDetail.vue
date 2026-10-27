<template>
  <div class="goods-detail">
    <el-form-item label="商品详情">
      <RichText v-model="formData.detailDesc" placeholder="请输入商品详情..." :height="800" />
    </el-form-item>
  </div>
</template>

<script setup>
import RichText from '@/admin/components/RichText/index.vue'

defineProps({
  formData: {
    type: Object,
    required: true
  }
})
</script>
