<template>
  <div class="app-container">
    <div class="page-header"><h2>企业资料</h2></div>
    <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" style="max-width: 720px">
      <el-form-item label="公司名称" prop="name"><el-input v-model="form.name" /></el-form-item>
      <el-form-item label="简称"><el-input v-model="form.shortName" /></el-form-item>
      <el-form-item label="联系人"><el-input v-model="form.contactName" /></el-form-item>
      <el-form-item label="电话" prop="phone"><el-input v-model="form.phone" /></el-form-item>
      <el-form-item label="地址"><el-input v-model="form.address" /></el-form-item>
      <el-form-item label="询盘邮箱">
        <el-input v-model="form.email" placeholder="新询盘发到这个邮箱" />
        <p class="hint">服务器填好发信账号后才会真正发出。没填时，询盘仍会记在「询盘记录」里。</p>
      </el-form-item>
      <el-form-item label="简介"><el-input v-model="form.summary" type="textarea" :rows="3" /></el-form-item>
      <el-form-item label="SEO标题"><el-input v-model="form.seoTitle" /></el-form-item>
      <el-form-item label="关键词"><el-input v-model="form.seoKeywords" /></el-form-item>
      <el-form-item label="描述"><el-input v-model="form.seoDescription" type="textarea" :rows="2" /></el-form-item>
      <el-form-item><el-button type="primary" @click="submit">保存</el-button></el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getSiteProfile, saveSiteProfile } from '@/admin/api/site'

const formRef = ref()
const form = reactive({
  name: '', shortName: '', contactName: '', phone: '', address: '', email: '',
  summary: '', seoTitle: '', seoKeywords: '', seoDescription: ''
})
const rules = {
  name: [{ required: true, message: '请填写公司名称', trigger: 'blur' }],
  phone: [{ required: true, message: '请填写电话', trigger: 'blur' }]
}

onMounted(async () => {
  Object.assign(form, await getSiteProfile())
})

async function submit() {
  await formRef.value.validate()
  await saveSiteProfile(form)
  ElMessage.success('已保存')
}
</script>

<style scoped>
.hint { margin: 6px 0 0; color: #909399; font-size: 12px; line-height: 1.5; }
</style>
