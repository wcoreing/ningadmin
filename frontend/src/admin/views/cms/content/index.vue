<template>
  <div class="app-container">
    <div class="page-header">
      <h2>内容</h2>
    </div>
    <div class="filter-container">
      <el-select v-model="query.channelId" clearable placeholder="栏目" class="filter-item" style="width: 180px">
        <el-option v-for="item in channels" :key="item.id" :label="item.name" :value="item.id" />
      </el-select>
      <el-input v-model="query.keyword" placeholder="标题" class="filter-item" style="width: 200px" @keyup.enter="search" />
      <el-select v-model="query.status" clearable placeholder="状态" class="filter-item" style="width: 120px">
        <el-option label="已发布" value="1" />
        <el-option label="草稿" value="0" />
      </el-select>
      <el-button class="filter-item" type="primary" @click="search">搜索</el-button>
      <el-button class="filter-item" type="primary" :icon="Plus" @click="openForm()">新增内容</el-button>
    </div>
    <el-table v-loading="loading" :data="list" border>
      <el-table-column prop="title" label="标题" min-width="220" />
      <el-table-column prop="channelName" label="栏目" width="140" />
      <el-table-column prop="publishTime" label="发布时间" width="170" />
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '已发布' : '草稿' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="160" align="center">
        <template #default="{ row }">
          <el-button link type="primary" @click="openForm(row)">编辑</el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="pager">
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="total"
        :page-size="query.pageSize"
        :current-page="query.page"
        @current-change="changePage"
      />
    </div>

    <el-dialog v-model="visible" class="content-editor-dialog" :title="form.id ? '编辑内容' : '新增内容'" width="920px" append-to-body destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="栏目" prop="channelId">
          <el-select v-model="form.channelId" placeholder="选择栏目" style="width: 100%">
            <el-option v-for="item in channels" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="摘要">
          <el-input v-model="form.summary" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="封面">
          <el-input v-model="form.cover" placeholder="图片地址" />
        </el-form-item>
        <el-form-item label="正文">
          <RichText v-if="visible" v-model="form.body" placeholder="写下正文，可插入图片" :height="360" />
        </el-form-item>
        <el-form-item label="发布时间">
          <el-date-picker v-model="form.publishTime" type="datetime" value-format="YYYY-MM-DD HH:mm:ss" placeholder="默认现在" />
        </el-form-item>
        <el-form-item label="SEO标题">
          <el-input v-model="form.seoTitle" />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="form.seoKeywords" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.seoDescription" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">发布</el-radio>
            <el-radio :value="0">草稿</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button type="primary" @click="submit">保存</el-button>
        <el-button @click="visible = false">取消</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { deleteContent, getContentInfo, getContentList, saveContent } from '@/admin/api/cms'
import { getCategoryList } from '@/admin/api/shop/category'
import RichText from '@/admin/components/RichText/index.vue'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const channels = ref([])
const visible = ref(false)
const formRef = ref()
const query = reactive({ channelId: undefined, keyword: '', status: '', page: 1, pageSize: 10 })
const form = reactive(emptyForm())
const rules = {
  channelId: [{ required: true, message: '请选择栏目', trigger: 'change' }],
  title: [{ required: true, message: '请填写标题', trigger: 'blur' }]
}

function emptyForm() {
  return {
    id: 0,
    channelId: undefined,
    title: '',
    summary: '',
    cover: '',
    body: '',
    seoTitle: '',
    seoKeywords: '',
    seoDescription: '',
    sort: 0,
    status: 1,
    publishTime: ''
  }
}

async function load() {
  loading.value = true
  try {
    const data = await getContentList({
      ...query,
      channelId: query.channelId || undefined
    })
    list.value = data?.list || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

function search() {
  query.page = 1
  load()
}

function changePage(page) {
  query.page = page
  load()
}

async function openForm(row) {
  Object.assign(form, emptyForm())
  if (row) {
    const detail = await getContentInfo(row.id)
    Object.assign(form, detail)
  }
  visible.value = true
}

async function submit() {
  await formRef.value.validate()
  await saveContent({ ...form })
  ElMessage.success('已保存')
  visible.value = false
  load()
}

async function remove(row) {
  await ElMessageBox.confirm(`删除「${row.title}」？`, '删除内容', { type: 'warning' })
  await deleteContent(row.id)
  ElMessage.success('已删除')
  load()
}

function flattenCategories(nodes, acc = []) {
  for (const node of nodes || []) {
    acc.push(node)
    if (node.children && node.children.length) flattenCategories(node.children, acc)
  }
  return acc
}

onMounted(async () => {
  channels.value = flattenCategories(await getCategoryList({ module: 'cms' }))
  load()
})
</script>

<style scoped>
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>

<style>
.content-editor-dialog {
  margin-top: 4vh;
  max-height: 92vh;
  display: flex;
  flex-direction: column;
}
.content-editor-dialog .el-dialog__body {
  overflow: auto;
}
</style>
