<template>
  <div class="app-container">
    <div class="page-header"><h2>询盘记录</h2></div>
    <div class="filter">
      <el-select v-model="status" clearable placeholder="状态" style="width: 140px" @change="load">
        <el-option label="未联系" value="0" />
        <el-option label="已联系" value="1" />
      </el-select>
    </div>
    <el-table :data="list" border>
      <el-table-column prop="createTime" label="时间" width="150" />
      <el-table-column prop="name" label="联系人" width="100" />
      <el-table-column prop="phone" label="电话" width="130" />
      <el-table-column prop="company" label="公司" min-width="140" />
      <el-table-column prop="goodsName" label="产品" min-width="140" />
      <el-table-column prop="quantity" label="数量" width="100" />
      <el-table-column prop="message" label="需求" min-width="220" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">{{ row.status === 1 ? '已联系' : '未联系' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button link type="primary" @click="mark(row)">{{ row.status === 1 ? '改回未联系' : '标为已联系' }}</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-model:current-page="page"
      :page-size="20"
      :total="total"
      layout="total, prev, pager, next"
      style="margin-top: 16px"
      @current-change="load"
    />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getInquiries, updateInquiryStatus } from '@/admin/api/site'

const list = ref([])
const total = ref(0)
const page = ref(1)
const status = ref('')

async function load() {
  const data = await getInquiries({ page: page.value, pageSize: 20, status: status.value || undefined })
  list.value = data.list || []
  total.value = data.total || 0
}

async function mark(row) {
  const next = row.status === 1 ? 0 : 1
  await updateInquiryStatus(row.id, next)
  ElMessage.success('已更新')
  load()
}

onMounted(load)
</script>

<style scoped>
.filter { margin-bottom: 12px; }
</style>
