<template>
  <div class="app-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :model="queryParams" ref="queryForm" :inline="true">
        <el-form-item label="标题" prop="title">
          <el-input
            v-model="queryParams.title"
            placeholder="请输入标题"
            clearable
            style="width: 240px"
            @keyup.enter="handleQuery"
          />
        </el-form-item>
        <el-form-item label="位置" prop="position">
          <el-select v-model="queryParams.position" placeholder="请选择位置" clearable style="width: 200px">
            <el-option
              v-for="dict in positionOptions"
              :key="dict.value"
              :label="dict.label"
              :value="dict.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable style="width: 200px">
            <el-option
              v-for="dict in statusOptions"
              :key="dict.value"
              :label="dict.label"
              :value="dict.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作按钮区域 -->
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <el-button type="primary" @click="handleAdd">新增</el-button>
        </div>
      </template>

      <!-- 表格区域 -->
      <el-table
        v-loading="loading"
        :data="bannerList"
        row-key="id"
        border
      >
        <el-table-column label="ID" prop="id" width="80" align="center" />
        <el-table-column label="标题" prop="title" min-width="120" align="left" show-overflow-tooltip />
        <el-table-column label="图片" width="100" align="center">
          <template #default="{ row }">
            <el-image
              style="width: 50px; height: 50px"
              :src="row.image"
              :preview-src-list="[row.image]"
              fit="cover"
            />
          </template>
        </el-table-column>
        <el-table-column label="位置" prop="position" width="100" align="center">
          <template #default="{ row }">
            {{ getPositionLabel(row.position) }}
          </template>
        </el-table-column>
        <el-table-column label="类型" prop="type" width="100" align="center">
          <template #default="{ row }">
            {{ getTypeLabel(row.type) }}
          </template>
        </el-table-column>
        <el-table-column label="目标商品" min-width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <template v-if="row.type === 1 && row.goodsInfo">
              {{ row.goodsInfo.name }}
            </template>
            <template v-else-if="row.type === 3">
              {{ row.url }}
            </template>
          </template>
        </el-table-column>
        <el-table-column label="排序" prop="sort" width="80" align="center" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="createTime" width="160" align="center" />
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleUpdate(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页区域 -->
      <el-pagination
        v-show="total > 0"
        v-model:current-page="queryParams.pageNum"
        v-model:page-size="queryParams.pageSize"
        :total="total"
        :page-sizes="[10, 20, 30, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>

    <!-- 添加或修改对话框 -->
    <el-dialog
      :title="title"
      v-model="open"
      width="700px"
      append-to-body
      destroy-on-close
      :modal-append-to-body="false"
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入标题" />
        </el-form-item>
        <el-form-item label="图片" prop="image">
          <div class="icon-select">
            <el-input v-model="form.image" placeholder="请选择图片">
              <template #append>
                <el-button @click="handleSelectImage">选择</el-button>
              </template>
            </el-input>
            <img v-if="form.image" :src="form.image" class="icon-preview" />
          </div>
        </el-form-item>
        <el-form-item label="位置" prop="position">
          <el-select v-model="form.position" placeholder="请选择位置">
            <el-option
              v-for="dict in positionOptions"
              :key="dict.value"
              :label="dict.label"
              :value="dict.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择类型">
            <el-option
              v-for="dict in typeOptions"
              :key="dict.value"
              :label="dict.label"
              :value="dict.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="目标ID" prop="targetId" v-if="form.type !== 3">
          <GoodsDisplay
            v-model="form.targetId"
            :goods-info="targetInfo"
            @change="handleSelectTarget"
            @clear="clearSelectedGoods"
          />
        </el-form-item>
        <el-form-item label="跳转链接" prop="url" v-if="form.type === 3">
          <el-input v-model="form.url" placeholder="请输入跳转链接" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
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
  </div>

  <!-- 将文件管理器移到外层 -->
  <FileManager
    v-model="fileManagerVisible"
    :multiple="false"
    @select="handleFileSelected"
  />

  <!-- 自定义图片预览组件 -->
  <div v-if="previewVisible" class="image-preview-container" @click="closePreview">
    <div class="image-preview-wrapper">
      <img :src="previewImage" class="preview-image" @click.stop />
      <div class="preview-close" @click="closePreview">×</div>
    </div>
  </div>

  <!-- 使用新的商品选择器组件 -->
  <GoodsSelector
    v-model="goodsSelectVisible"
    :multiple="false"
    @select="handleGoodsSelect"
  />
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getBannerList, createBanner, updateBanner, deleteBanner, changeBannerStatus, getBannerInfo } from '@/admin/api/shop/banner'
import { getGoodsList } from '@/admin/api/shop/goods'
import FileManager from '@/admin/components/FileManager/index.vue'
import GoodsSelector from '@/admin/components/GoodsSelector/index.vue'
import GoodsDisplay from '@/admin/components/GoodsDisplay/index.vue'

// 遮罩层
const loading = ref(false)
// 总条数
const total = ref(0)
// 弹出层标题
const title = ref('')
// 是否显示弹出层
const open = ref(false)
// Banner列表
const bannerList = ref([])

// 查询参数
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  title: undefined,
  position: undefined,
  status: undefined
})

// 表单参数
const form = ref({
  id: undefined,
  title: undefined,
  image: undefined,
  url: undefined,
  position: undefined,
  type: 1, // 默认类型为商品
  targetId: undefined,
  sort: 0,
  status: 1
})

// 表单校验规则
const rules = ref({
  title: [{ required: true, message: '标题不能为空', trigger: 'blur' }],
  image: [{ required: true, message: '图片不能为空', trigger: 'blur' }],
  position: [{ required: true, message: '位置不能为空', trigger: 'change' }],
  type: [{ required: true, message: '类型不能为空', trigger: 'change' }]
})

// 位置选项
const positionOptions = [
  { label: '首页', value: 1 },
  { label: '分类页', value: 2 },
  { label: '活动页', value: 3 }
]

// 类型选项
const typeOptions = [
  { label: '商品', value: 1 },
  { label: '活动', value: 2 },
  { label: '外链', value: 3 }
]

// 状态选项
const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 0 }
]

// 上传地址
const uploadUrl = ref((import.meta.env.VITE_API_BASE_URL || '/api') + '/common/upload')

// 表单ref
const formRef = ref()

// 文件管理器相关
const fileManagerVisible = ref(false)

// 商品选择相关
const goodsSelectVisible = ref(false)
const targetInfo = ref(null)

// 选择图片
const handleSelectImage = () => {
  fileManagerVisible.value = true
}

// 文件选择回调
const handleFileSelected = (files) => {
  form.value.image = files.url
  fileManagerVisible.value = false
}

// 获取Banner列表
const getList = async () => {
  loading.value = true
  try {
    const res = await getBannerList(queryParams.value)
    console.log('Banner列表数据:', res)
    bannerList.value = res.list
    total.value = res.total
  } catch (error) {
    console.error('获取Banner列表失败:', error)
  }
  loading.value = false
}

// 表单重置
const reset = () => {
  form.value = {
    id: undefined,
    title: undefined,
    image: undefined,
    url: undefined,
    position: undefined,
    type: 1, // 设置默认类型为商品
    targetId: undefined,
    sort: 0,
    status: 1
  }
  targetInfo.value = null // 清除商品信息
  formRef.value?.resetFields()
}

// 取消按钮
const cancel = () => {
  open.value = false
  reset()
}

// 查询按钮操作
const handleQuery = () => {
  queryParams.value.pageNum = 1
  getList()
}

// 重置按钮操作
const resetQuery = () => {
  queryParams.value = {
    pageNum: 1,
    pageSize: 10,
    title: undefined,
    position: undefined,
    status: undefined
  }
  handleQuery()
}

// 新增按钮操作
const handleAdd = () => {
  reset()
  open.value = true
  title.value = '添加Banner'
}

// 修改按钮操作
const handleUpdate = async (row) => {
  reset()
  const id = row.id || row
  open.value = true
  title.value = '修改Banner'
  try {
    // 通过详情接口获取完整的 Banner 信息
    const res = await getBannerInfo(id)
    form.value = {
      id: res.id,
      title: res.title,
      image: res.image,
      url: res.url,
      position: res.position,
      type: res.type,
      targetId: res.targetId,
      sort: res.sort,
      status: res.status
    }
    // 如果有商品信息，设置到 targetInfo
    if (res.type === 1 && res.goodsInfo) {
      const goodsId = res.goodsInfo.goodsId || res.goodsInfo.id
      targetInfo.value = {
        goodsId: goodsId,
        name: res.goodsInfo.name,
        coverImg: res.goodsInfo.coverImg,
        price: res.goodsInfo.price
      }
    }
  } catch (error) {
    console.error('获取Banner详情失败:', error)
    ElMessage.error('获取Banner详情失败')
  }
}

// 监听表单类型变化，清除相关数据
watch(() => form.value.type, (newType) => {
  if (newType !== 1) {
    form.value.targetId = undefined
    targetInfo.value = null
  }
})

// 提交按钮
const submitForm = async () => {
  formRef.value?.validate(async (valid) => {
    if (valid) {
      try {
        // 准备提交数据
        const submitData = { ...form.value }
        
        // 确保必填字段有值
        if (!submitData.type) {
          ElMessage.error('请选择Banner类型')
          return
        }
        
        // 转换 targetId 为数字
        if (submitData.targetId) {
          submitData.targetId = parseInt(submitData.targetId)
        } else if (submitData.type === 1) {
          // 如果是商品类型但没有选择商品，提示用户
          ElMessage.error('请选择目标商品')
          return
        }
        
        // 如果是外链类型，确保有URL
        if (submitData.type === 3 && !submitData.url) {
          ElMessage.error('请输入跳转链接')
          return
        }
        
        console.log('提交数据:', submitData)
        console.log('表单原始数据:', form.value)
        
        if (submitData.id) {
          console.log('执行更新操作')
          await updateBanner(submitData)
          ElMessage.success('修改成功')
        } else {
          console.log('执行创建操作')
          await createBanner(submitData)
          ElMessage.success('新增成功')
        }
        open.value = false
        getList()
      } catch (error) {
        console.error('提交失败:', error.response?.data || error)
        ElMessage.error(error.response?.data?.message || '提交失败')
      }
    }
  })
}

// 删除按钮操作
const handleDelete = (row) => {
  const ids = row.id || row
  ElMessageBox.confirm('是否确认删除?', '警告', {
    type: 'warning'
  }).then(async () => {
    try {
      await deleteBanner(ids)
      getList()
      ElMessage.success('删除成功')
    } catch (error) {
      console.error('删除失败:', error)
    }
  })
}

// 状态修改
const handleStatusChange = async (row) => {
  try {
    await changeBannerStatus(row.id, row.status)
    ElMessage.success('修改成功')
  } catch (error) {
    console.error('修改状态失败:', error)
    row.status = row.status === 0 ? 1 : 0
  }
}

// 获取位置标签
const getPositionLabel = (value) => {
  const option = positionOptions.find(opt => opt.value === value)
  return option ? option.label : ''
}

// 获取类型标签
const getTypeLabel = (value) => {
  const option = typeOptions.find(opt => opt.value === value)
  return option ? option.label : ''
}

// 图片预览相关
const previewVisible = ref(false)
const previewImage = ref('')

// 预览图片
const handlePreviewImage = (url) => {
  previewImage.value = url
  previewVisible.value = true
}

// 关闭预览
const closePreview = () => {
  previewVisible.value = false
  previewImage.value = ''
}

// 打开商品选择
const handleSelectTarget = () => {
  goodsSelectVisible.value = true
}

// 修改商品选择处理方法
const handleGoodsSelect = (goods) => {
  console.log('选择的商品数据:', goods)
  const goodsId = goods.goodsId || goods.id
  form.value.targetId = goodsId
  targetInfo.value = {
    goodsId: goodsId,
    name: goods.name,
    coverImg: goods.coverImg,
    price: goods.price
  }
  goodsSelectVisible.value = false
}

// 清除选中的商品
const clearSelectedGoods = () => {
  form.value.targetId = undefined
  targetInfo.value = null
}

// 分页处理函数
const handleSizeChange = (val) => {
  queryParams.value.pageSize = val
  queryParams.value.pageNum = 1
  getList()
}

const handleCurrentChange = (val) => {
  queryParams.value.pageNum = val
  getList()
}

// 页面加载时获取列表数据
onMounted(() => {
  getList()
})
</script>

<style lang="scss" scoped>
/* 文本省略 mixin */
@mixin text-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.banner-uploader {
  :deep(.el-upload) {
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: var(--el-transition-duration-fast);

    &:hover {
      border-color: var(--el-color-primary);
    }
  }
}

.banner-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
  line-height: 178px;
}

.banner-image {
  width: 80px;
  height: 80px;
  border-radius: 4px;
  border: 1px solid #ebeef5;
  cursor: pointer;
}

.search-card {
  margin-bottom: 15px;

  :deep(.el-form-item) {
    margin-bottom: 12px;
    margin-right: 12px;
  }

  :deep(.el-form-item__label) {
    width: 60px;
  }
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

.el-dialog {
  :deep(.el-dialog__body) {
    max-height: 70vh;
    overflow-y: auto;
  }
}

/* 调整弹窗层级 */
:deep(.el-dialog__wrapper) {
  z-index: 2000 !important;
}

:deep(.el-overlay) {
  z-index: 1999 !important;
}

/* 文件管理器层级 */
:deep(.file-manager-dialog) {
  z-index: 2100 !important;
}

:deep(.file-manager-overlay) {
  z-index: 2099 !important;
}

/* 图片预览样式 */
.image-preview-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  z-index: 9999;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
}

.image-preview-wrapper {
  position: relative;
  max-width: 90%;
  max-height: 90%;
}

.preview-image {
  max-width: 100%;
  max-height: 90vh;
  object-fit: contain;
  cursor: default;
}

.preview-close {
  position: absolute;
  top: -40px;
  right: -40px;
  width: 40px;
  height: 40px;
  line-height: 40px;
  text-align: center;
  color: #fff;
  font-size: 24px;
  cursor: pointer;
  transition: all 0.3s;
  
  &:hover {
    transform: rotate(90deg);
  }
}
</style> 