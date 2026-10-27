<template>
  <div class="app-container">
    <div class="page-header">
      <h2>内容分类</h2>
    </div>
    <div class="filter-container">
      <el-input v-model="query.name" placeholder="名称或编码" class="filter-item" style="width: 200px" @keyup.enter="load" />
      <el-select v-model="query.status" clearable placeholder="状态" class="filter-item" style="width: 120px">
        <el-option label="启用" :value="1" />
        <el-option label="禁用" :value="0" />
      </el-select>
      <el-button class="filter-item" type="primary" @click="load">搜索</el-button>
      <el-button class="filter-item" type="primary" :icon="Plus" @click="openForm()">新增分类</el-button>
    </div>
    <el-table v-loading="loading" :data="list" row-key="id" border default-expand-all :tree-props="{ children: 'children' }">
      <el-table-column prop="name" label="名称" min-width="180">
        <template #default="{ row }">
          <img
            v-if="coverOf(row)"
            v-show="coverVisible[row.id]"
            :src="coverOf(row)"
            class="cover"
            alt=""
            @load="coverVisible[row.id] = true"
            @error="coverVisible[row.id] = false"
          />
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="编码" width="120">
        <template #default="{ row }">{{ row.extra?.code }}</template>
      </el-table-column>
      <el-table-column label="类型" width="80" align="center">
        <template #default="{ row }">{{ row.extra?.type === 'list' ? '列表' : '单页' }}</template>
      </el-table-column>
      <el-table-column prop="sort" label="排序" width="80" align="center" />
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" align="center">
        <template #default="{ row }">
          <el-button link type="primary" @click="openForm(null, row.id)">添加子分类</el-button>
          <el-button link type="primary" @click="openForm(row)">编辑</el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="visible" :title="form.id ? '编辑内容分类' : '新增内容分类'" width="640px" append-to-body>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="上级分类">
          <el-tree-select
            v-model="form.parent_id"
            :data="parentOptions"
            :props="{ label: 'name', value: 'id' }"
            value-key="id"
            check-strictly
            clearable
            placeholder="无"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="编码" prop="extra.code">
          <el-input v-model="form.extra.code" placeholder="英文编码，如 about" />
        </el-form-item>
        <el-form-item label="类型">
          <el-radio-group v-model="form.extra.type">
            <el-radio value="page">单页</el-radio>
            <el-radio value="list">列表</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="封面">
          <el-input v-model="form.image" placeholder="图片地址" />
        </el-form-item>
        <el-form-item label="摘要">
          <el-input v-model="form.extra.summary" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="SEO标题">
          <el-input v-model="form.extra.seoTitle" />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="form.extra.seoKeywords" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.extra.seoDescription" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" :max="9999" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
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
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { createCategory, deleteCategory, getCategoryList, updateCategory } from '@/admin/api/shop/category'

const loading = ref(false)
const coverVisible = ref({})
const source = ref([])
const list = ref([])
const visible = ref(false)
const formRef = ref()
const query = reactive({ name: '', status: undefined })
const form = reactive(emptyForm())
const rules = {
  name: [{ required: true, message: '请填写名称', trigger: 'blur' }],
  'extra.code': [{ required: true, message: '请填写编码', trigger: 'blur' }]
}

const parentOptions = computed(() => [{
  id: 0,
  name: '无',
  children: stripSelf(source.value, form.id)
}])

function emptyExtra() {
  return { code: '', type: 'page', summary: '', seoTitle: '', seoKeywords: '', seoDescription: '' }
}

function emptyForm() {
  return {
    id: undefined,
    parent_id: 0,
    name: '',
    image: '',
    icon: '',
    extra: emptyExtra(),
    sort: 0,
    status: 1,
    module: 'cms'
  }
}

function stripSelf(nodes, id) {
  return (nodes || []).filter(node => node.id !== id).map(node => ({
    ...node,
    children: stripSelf(node.children, id)
  }))
}

function filterTree(nodes) {
  const keyword = query.name.trim()
  const result = []
  for (const node of nodes || []) {
    const children = filterTree(node.children)
    const text = `${node.name} ${node.extra?.code || ''}`
    const statusUnset = query.status === undefined || query.status === '' || query.status === null
    const matched = (!keyword || text.includes(keyword)) && (statusUnset || node.status === query.status)
    if (matched || children.length) {
      result.push({ ...node, children })
    }
  }
  return result
}

function coverOf(row) {
  return (row.image || row.icon || '').trim()
}

async function load() {
  loading.value = true
  coverVisible.value = {}
  try {
    source.value = (await getCategoryList({ module: 'cms' })) || []
    list.value = filterTree(source.value)
  } finally {
    loading.value = false
  }
}

function openForm(row, parentId) {
  Object.assign(form, emptyForm())
  if (row) {
    Object.assign(form, row, { extra: { ...emptyExtra(), ...(row.extra || {}) }, module: 'cms' })
  } else if (parentId) {
    form.parent_id = parentId
  }
  visible.value = true
}

async function submit() {
  await formRef.value.validate()
  const payload = { ...form, parent_id: form.parent_id || 0, module: 'cms' }
  if (form.id) {
    await updateCategory(form.id, payload)
  } else {
    await createCategory(payload)
  }
  ElMessage.success('已保存')
  visible.value = false
  load()
}

async function remove(row) {
  await ElMessageBox.confirm(`删除内容分类「${row.name}」？有子分类或内容时不能删除。`, '删除分类', { type: 'warning' })
  await deleteCategory(row.id, 'cms')
  ElMessage.success('已删除')
  load()
}

onMounted(load)
</script>

<style scoped>
.cover {
  width: 20px;
  height: 20px;
  margin-right: 8px;
  vertical-align: middle;
}
</style>
