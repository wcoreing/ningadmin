<template>
  <div class="app-container">
    <!-- 搜索工具栏 -->
    <div class="filter-container">
      <el-input
        v-model="queryParams.name"
        placeholder="分类名称"
        style="width: 200px"
        class="filter-item"
        @keyup.enter="handleQuery"
      />
      <el-select
        v-model="queryParams.status"
        placeholder="状态"
        clearable
        class="filter-item"
        style="width: 130px"
      >
        <el-option label="启用" :value="1" />
        <el-option label="禁用" :value="0" />
      </el-select>
      <el-button
        class="filter-item"
        type="primary"
        :icon="Search"
        @click="handleQuery"
      >
        搜索
      </el-button>
      <el-button
        class="filter-item"
        type="primary"
        :icon="Plus"
        @click="handleAdd"
      >
        新增
      </el-button>
    </div>

    <!-- 分类表格 -->
    <el-table
      v-loading="loading"
      :data="categoryList"
      row-key="id"
      border
      default-expand-all
      :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
    >
      <template #empty>
        <span>暂无数据</span>
      </template>
      <el-table-column prop="name" label="分类名称" min-width="200">
        <template #default="{ row }">
          <img v-if="row.icon" :src="row.icon" class="category-icon" />
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="sort" label="排序" width="80" align="center" />
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        label="操作"
        align="center"
        width="280"
        class-name="small-padding fixed-width"
      >
        <template #default="{ row }">
          <el-button
            link
            type="primary"
            :icon="Plus"
            @click="handleAdd(row)"
          >
            添加子分类
          </el-button>
          <el-button
            link
            type="primary"
            :icon="Edit"
            @click="handleUpdate(row)"
          >
            编辑
          </el-button>
          <el-button
            link
            type="danger"
            :icon="Delete"
            @click="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加或修改分类对话框 -->
    <el-dialog
      :title="dialog.title"
      v-model="dialog.visible"
      width="500px"
      append-to-body
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="上级分类">
          <el-tree-select
            v-model="form.parent_id"
            :data="categoryOptions"
            :props="{ label: 'name', value: 'id' }"
            value-key="id"
            placeholder="选择上级分类"
            check-strictly
            clearable
          />
        </el-form-item>
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <div class="icon-select">
            <el-input v-model="form.icon" placeholder="请选择图标">
              <template #append>
                <el-button @click="handleSelectIcon">选择</el-button>
              </template>
            </el-input>
            <img v-if="form.icon" :src="form.icon" class="icon-preview" />
          </div>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitForm">确 定</el-button>
          <el-button @click="cancel">取 消</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 文件管理组件 -->
    <FileManager
      v-model="fileManagerVisible"
      :multiple="false"
      @select="handleFileSelected"
      @confirm="(files) => console.log('confirm事件触发', files)"
      @change="(files) => console.log('change事件触发', files)"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, Edit, Delete } from '@element-plus/icons-vue'
import { getCategoryList, createCategory, updateCategory, deleteCategory } from '@/api/shop/category'
import FileManager from '@/components/FileManager/index.vue'

const loading = ref(false)
const categoryList = ref([])
const categoryOptions = ref([])
const queryParams = ref({
  name: '',
  status: undefined
})

const dialog = ref({
  title: '',
  visible: false
})

const form = ref({
  id: undefined,
  parent_id: 0,
  name: '',
  icon: '',
  sort: 0,
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const formRef = ref()
const fileManagerVisible = ref(false)

// 查询分类列表
const getList = async () => {
  try {
    loading.value = true
    const res = await getCategoryList(queryParams.value)
    console.log('分类列表响应:', res)
    console.log('分类数据结构:', JSON.stringify(res.data, null, 2))
    // 确保数据是数组
    if (Array.isArray(res)) {
      categoryList.value = res
      categoryOptions.value = [{ 
        id: 0, 
        name: '顶级分类', 
        children: res
      }]
    } else {
      categoryList.value = []
      categoryOptions.value = [{ 
        id: 0, 
        name: '顶级分类', 
        children: []
      }]
    }
  } catch (error) {
    console.error('获取分类列表失败:', error)
    categoryList.value = []
    categoryOptions.value = [{ 
      id: 0, 
      name: '顶级分类', 
      children: []
    }]
  } finally {
    loading.value = false
  }
}

// 搜索按钮操作
const handleQuery = () => {
  getList()
}

// 重置表单
const resetForm = () => {
  form.value = {
    id: undefined,
    parent_id: 0,
    name: '',
    icon: '',
    sort: 0,
    status: 1
  }
  formRef.value?.resetFields()
}

// 取消按钮
const cancel = () => {
  resetForm()
  dialog.value.visible = false
}

// 新增按钮操作
const handleAdd = (row?: any) => {
  resetForm()
  if (row) {
    form.value.parent_id = row.id
  }
  dialog.value.title = '添加分类'
  dialog.value.visible = true
}

// 修改按钮操作
const handleUpdate = (row: any) => {
  resetForm()
  form.value = { ...row }
  dialog.value.title = '修改分类'
  dialog.value.visible = true
}

// 提交表单
const submitForm = async () => {
  formRef.value?.validate(async (valid: boolean) => {
    if (!valid) return
    try {
      if (form.value.id) {
        await updateCategory(form.value.id, form.value)
        ElMessage.success('修改成功')
      } else {
        await createCategory(form.value)
        ElMessage.success('新增成功')
      }
      dialog.value.visible = false
      getList()
    } catch (error) {
      console.error('提交失败:', error)
    }
  })
}

// 删除按钮操作
const handleDelete = (row: any) => {
  ElMessageBox.confirm('是否确认删除该分类?', '警告', {
    type: 'warning'
  }).then(async () => {
    try {
      await deleteCategory(row.id)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error('删除失败:', error)
    }
  })
}

// 选择图标
const handleSelectIcon = () => {
  console.log('点击选择图标按钮')
  fileManagerVisible.value = true
}

// 文件选择回调
const handleFileSelected = (files: any) => {
  console.log('文件选择回调被触发', files)
  form.value.icon = files.url
 
  fileManagerVisible.value = false
}

onMounted(() => {
  getList()
})
</script>

<style scoped>
.category-icon {
  width: 20px;
  height: 20px;
  margin-right: 8px;
  vertical-align: middle;
}
.filter-container {
  padding-bottom: 10px;
}
.filter-item {
  margin-right: 10px;
}
.icon-select {
  display: flex;
  align-items: center;
  gap: 8px;
}

.icon-preview {
  width: 32px;
  height: 32px;
  object-fit: contain;
}
</style> 