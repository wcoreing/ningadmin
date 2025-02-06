<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="updateVisible"
    :title="type === 'add' ? '添加商品' : '编辑商品'"
    width="1200px"
    :show-close="false"
    @close="handleClose"
  >
    <template #header="{ close, titleId, titleClass }">
      <div class="my-header">
        <span :id="titleId" :class="titleClass">
          {{ type === 'add' ? '添加商品' : '编辑商品' }}
        </span>
        <div class="my-header-buttons">
          <el-button @click="handleClose">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
          <el-button link @click="close">
            <el-icon><Close /></el-icon>
          </el-button>
        </div>
      </div>
    </template>

    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
    >
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <!-- 基本信息 Tab -->
        <el-tab-pane label="基本信息" name="basic">
          <basic-info
            :form-data="formData"
            :category-options="categoryOptions"
            :brand-options="brandOptions"
            :supplier-options="supplierOptions"
            :tag-options="tagOptions"
          />
        </el-tab-pane>

        <!-- 商品图片 Tab -->
        <el-tab-pane label="商品图片" name="images">
          <goods-images :form-data="formData" />
        </el-tab-pane>

        <!-- 规格设置 Tab -->
        <el-tab-pane label="规格设置" name="specs">
          <template #label>规格设置</template>
          <SpecSettings 
            :form-data="formData"
            :upload-url="uploadUrl"
            ref="specSettingsRef"
          />
        </el-tab-pane>

        <!-- 商品详情 Tab -->
        <el-tab-pane label="商品详情" name="detail">
          <goods-detail :form-data="formData" :key="detailKey" />
        </el-tab-pane>
      </el-tabs>
    </el-form>
  </el-dialog>
</template>

<script setup>
import { ref, defineProps, defineEmits, watch, computed, shallowRef, onBeforeUnmount, nextTick, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { updateGoods, createGoods, getGoodsInfo } from '@/api/shop/goods'
import { Plus, Rank, Delete, Search, UploadFilled, Close } from '@element-plus/icons-vue'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import FileManager from '@/components/FileManager/index.vue'
import BasicInfo from './BasicInfo.vue'
import GoodsImages from './GoodsImages.vue'
import SpecSettings from './SpecSettings.vue'
import GoodsDetail from './GoodsDetail.vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  type: {
    type: String,
    default: 'add'
  },
  data: {
    type: Object,
    default: () => ({})
  },
  categoryOptions: {
    type: Array,
    default: () => []
  },
  brandOptions: {
    type: Array,
    default: () => []
  },
  supplierOptions: {
    type: Array,
    default: () => []
  },
  tagOptions: {
    type: Array,
    default: () => []
  },
  uploadUrl: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:visible', 'success', 'close'])

const formRef = ref(null)
const activeTab = ref('basic')

const formData = ref({
  goodsId: undefined,
  goodsName: '',
  categoryId: undefined,
  brandId: undefined,
  supplierId: undefined,
  price: 0,
  originalPrice: 0,
  stock: 0,
  sort: 0,
  status: '0',
  description: '',
  coverImg: '',
  tagIds: [],
  images: [],
  specType: 1,
  specGroups: [],
  skuList: [],
  detailDesc: ''
})

// 表单验证规则
const formRules = {
  goodsName: [
    { required: true, message: '请输入商品名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  categoryId: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ],
  brandId: [
    { required: true, message: '请选择商品品牌', trigger: 'change' }
  ],
  supplierId: [
    { required: true, message: '请选择供应商', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入商品价格', trigger: 'blur' }
  ],
  stock: [
    { required: true, message: '请输入商品库存', trigger: 'blur' }
  ]
}

// 添加 tab 验证规则映射
const tabRules = {
  basic: ['goodsName', 'categoryId', 'brandId', 'supplierId', 'price', 'stock'],
  specs: [], // 规格相关的验证
  images: [], // 图片相关的验证
  detail: [] // 详情相关的验证
}

// 编辑器实例
const editorRef = shallowRef(null)

// 编辑器配置
const toolbarConfig = {
  excludeKeys: []
}

const editorConfig = {
  placeholder: '请输入商品详情...',
  MENU_CONF: {
    uploadImage: {
      server: props.uploadUrl,
      fieldName: 'file',
      onSuccess(file, res) {
        return res.url
      }
    }
  }
}

// 修改编辑器创建回调
const handleCreated = (editor) => {
  editorRef.value = editor
}

// 修改 loadGoodsDetail 方法
const loadGoodsDetail = async (id) => {
  try {
    const data = await getGoodsInfo(id)
    
    // 先重置表单数据
    formData.value = {
      goodsId: undefined,
      goodsName: '',
      categoryId: undefined,
      brandId: undefined,
      supplierId: undefined,
      price: 0,
      originalPrice: 0,
      stock: 0,
      sort: 0,
      status: '0',
      description: '',
      coverImg: '',
      tagIds: [],
      images: [],
      specType: 1,
      specGroups: [],
      skuList: [],
      detailDesc: ''
    }

    // 等待下一个 tick，确保重置完成
    await nextTick()
    
    // 然后设置新数据
    formData.value = {
      ...formData.value,
      goodsId: data.goodsId,
      goodsName: data.goodsName,
      categoryId: data.categoryId,
      brandId: data.brandId,
      supplierId: data.supplierId,
      price: data.price,
      originalPrice: data.originalPrice,
      stock: data.stock,
      sort: data.sort,
      status: data.status,
      description: data.description,
      coverImg: data.coverImg,
      tagIds: data.tagIds || [],
      images: data.images || [],
      specType: data.specType || 1,
      specGroups: data.specGroups ? JSON.parse(JSON.stringify(data.specGroups)) : [],
      skuList: data.skuList ? JSON.parse(JSON.stringify(data.skuList)) : [],
      detailDesc: data.detailDesc || '<p><br></p>'
    }

   
    
    // 手动触发 SpecSettings 组件的初始化
    await nextTick()
    specSettingsRef.value?.initSkuList?.()
  } catch (error) {
    console.error('获取商品详情失败:', error)
    ElMessage.error('获取商品详情失败')
  }
}

// 添加 specSettingsRef
const specSettingsRef = ref(null)

// 修改 watch
watch(
  () => props.visible,
  async (val) => {
    if (val) {
      // 等待 Dialog 显示完成后再加载数据
      await nextTick()
      
      if (props.type === 'edit' && props.data?.goodsId) {
        console.log('开始加载商品详情:', props.data.goodsId)
        await loadGoodsDetail(props.data.goodsId)
      } else {
        // 新增时重置表单
        formData.value = {
          goodsId: undefined,
          goodsName: '',
          categoryId: undefined,
          brandId: undefined,
          supplierId: undefined,
          price: 0,
          originalPrice: 0,
          stock: 0,
          sort: 0,
          status: '0',
          description: '',
          coverImg: '',
          tagIds: [],
          images: [],
          specType: 1,
          specGroups: [],
          skuList: [],
          detailDesc: '<p><br></p>'
        }
        formRef.value?.resetFields()
      }
      activeTab.value = 'basic'
    } else {
      // Dialog 关闭时销毁编辑器
      if (editorRef.value) {
        editorRef.value.destroy()
        editorRef.value = null
      }
    }
  }
)

// 添加 updateVisible 方法处理可见性更新
const updateVisible = (value) => {
  emit('update:visible', value)
}

// 修改 activeTab 的处理
const handleTabChange = (tabName) => {
  activeTab.value = tabName
}

// 添加一个用于刷新商品详情组件的 key
const detailKey = ref(0)

// 修改提交表单方法
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    
    // 处理 SKU 数据
    if (formData.value.specType === 2) {
      formData.value.skuList = formData.value.skuList.map(sku => ({
        specs: sku.specs,
        price: Number(sku.price || 0),
        stock: Number(sku.stock || 0),
        image: sku.image || ''
      }))
    }

    const submitData = {
      ...formData.value
    }

    if (props.type === 'add') {
      await createGoods(submitData)
    } else {
      await updateGoods(submitData)
    }

    ElMessage.success(props.type === 'add' ? '添加成功' : '更新成功')
    emit('success')
    handleClose()
    // 刷新商品详情组件
    detailKey.value++
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('提交失败，请检查表单')
  }
}

// 处理文件选择
const handleImageSelect = (files) => {
  console.log('选择的文件:', files) // 调试日志
  
  // 确保 files 是一个对象或数组
  const fileData = Array.isArray(files) ? files : [files]
  
  if (currentSelectType.value === 'cover') {
    // 处理封面图片，单选
    if (fileData.length > 0) {
      const file = fileData[0]
      // 直接设置 url，不使用解构赋值
      formData.value.coverImg = file.url
      console.log('设置封面图片:', formData.value.coverImg)
    }
  } else {
    // 处理商品图片，多选
    if (!Array.isArray(formData.value.images)) {
      formData.value.images = []
    }
    
    // 直接使用 url
    const newImages = fileData.map(file => file.url)
    formData.value.images.push(...newImages)
    console.log('添加商品图片:', newImages)
  }
  
  // 关闭文件管理器
  fileManagerVisible.value = false
}

// 处理图片删除
const handleRemoveImage = (type, index) => {
  if (type === 'cover') {
    formData.value.coverImg = ''
  } else {
    formData.value.images.splice(index, 1)
  }
}

// 格式化文件大小
const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 修改初始化
onMounted(() => {
  if (props.type === 'edit' && props.data?.goodsId) {
    loadGoodsDetail(props.data.goodsId)
  }
})

// 生成笛卡尔积
const cartesian = (...arrays) => {
  return arrays.reduce((acc, curr) => {
    return acc.flatMap(a => curr.map(c => [...a, c]))
  }, [[]])
}

// 添加计算属性，过滤有效的规格组
const validSpecGroups = computed(() => {
  return formData.value.specGroups?.filter(group => 
    group?.values?.length > 0 && group.name
  ) || []
})

// 添加计算属性，判断是否有有效的规格组
const hasValidSpecGroups = computed(() => {
  return validSpecGroups.value.length > 0
})

// 修改 skuList 计算属性
const skuList = computed(() => {
  if (!hasValidSpecGroups.value) return []
  if (formData.value.specType !== 2) return []
  
  // 获取所有规格值数组
  const specArrays = validSpecGroups.value.map(group => group.values)
  
  // 生成规格组合
  const combinations = cartesian(...specArrays)
  
  // 生成 SKU 列表
  return combinations.map((combination) => {
    // 创建规格对象
    const specs = {}
    validSpecGroups.value.forEach((group, groupIndex) => {
      specs[group.name] = combination[groupIndex]
    })
    
    // 查找已存在的 SKU 数据
    const existingSku = formData.value.skuList?.find(
      sku => {
        // 只比较规格值，忽略 id 等其他属性
        const isSameSpecs = Object.entries(sku.specs).every(([groupName, spec]) => {
          const targetSpec = specs[groupName]
          return targetSpec && spec.value === targetSpec.value
        })
        return isSameSpecs
      }
    )
    return existingSku || {
      specs,
      price: formData.value.price || 0,
      stock: formData.value.stock || 0
    }
  })
})

// 监听规格组和规格类型变化，更新 SKU 列表
watch(
  [
    () => formData.value.specGroups,
    () => formData.value.specType
  ],
  () => {
    if (formData.value.specType === 2) {
      formData.value.skuList = skuList.value
    } else {
      formData.value.skuList = []
    }
  },
  { deep: true }
)

// 修改添加规格组方法
const addSpecGroup = () => {
  const newGroup = {
    id: Date.now(),
    name: '',
    values: [],
    inputVisible: false,
    inputValue: ''
  }
  formData.value.specGroups.push(newGroup)
}

// 添加 inputRef
const inputRefs = ref([])

// 修改显示输入框方法
const showInput = (groupIndex) => {
  formData.value.specGroups[groupIndex].inputVisible = true
  // 等待DOM更新后聚焦
  nextTick(() => {
    inputRefs.value[groupIndex].focus()
  })
}

// 修改处理输入确认方法
const handleInputConfirm = (groupIndex) => {
  const group = formData.value.specGroups[groupIndex]
  const inputValue = group.inputValue?.trim()
  if (inputValue) {
    // 检查是否重复
    const isDuplicate = group.values.some(item => item.name === inputValue)
    if (!isDuplicate) {
      group.values.push({
        id: Date.now(),
        name: inputValue,
        value: inputValue, // 默认使用名称作为值
        sort: group.values.length + 1
      })
    } else {
      ElMessage.warning('规格值不能重复')
    }
  }
  group.inputVisible = false
  group.inputValue = ''
}

// 修改移除规格值方法
const handleRemoveSpecValue = (groupIndex, valueIndex) => {
  if (formData.value.specGroups[groupIndex].values) {
    formData.value.specGroups[groupIndex].values.splice(valueIndex, 1)
  }
}

// 计算属性：转换分类数据为树形结构
const categoryTreeData = computed(() => {
  return props.categoryOptions || []
})

// 添加删除规格组方法
const removeSpecGroup = (groupIndex) => {
  ElMessageBox.confirm('确定要删除该规格组吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    formData.value.specGroups.splice(groupIndex, 1)
  })
}

// 附件管理相关
const fileManagerVisible = ref(false)
const currentSelectType = ref('')  // 添加回来，用于区分封面和商品图片

// 打开文件管理器
const openFileManager = (type) => {
  currentSelectType.value = type
  fileManagerVisible.value = true
}

// 添加 handleClose 方法
const handleClose = () => {
  // 重置表单数据
  formData.value = {
    goodsId: undefined,
    goodsName: '',
    categoryId: undefined,
    brandId: undefined,
    supplierId: undefined,
    price: 0,
    originalPrice: 0,
    stock: 0,
    sort: 0,
    status: '0',
    description: '',
    coverImg: '',
    tagIds: [],
    images: [],
    specType: 1,
    specGroups: [],
    skuList: [],
    detailDesc: ''
  }
  
  // 重置表单验证
  formRef.value?.resetFields()
  
  // 重置当前标签页
  activeTab.value = 'basic'
  
  // 关闭对话框
  emit('update:visible', false)
  
  // 销毁编辑器实例
  if (editorRef.value) {
    editorRef.value.destroy()
    editorRef.value = null
  }
}
</script>

<style lang="scss" scoped>
:deep(.el-tree-select) {
  width: 100%;
}

.spec-container {
  padding: 20px;
  background: #f5f7fa;
  border-radius: 4px;

  .spec-group-item {
    margin-bottom: 20px;
    padding: 20px;
    background: #fff;
    border: 1px solid #dcdfe6;
    border-radius: 4px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);

    .spec-group-header {
      display: flex;
      align-items: center;
      margin-bottom: 20px;
      gap: 16px;

      .spec-group-title {
        min-width: 80px;
        font-weight: 500;
        color: #606266;
      }

      .spec-group-name {
        width: 240px;
      }
    }

    .spec-values {
      padding-left: 96px;  // 与规格组标题对齐

      .spec-values-list {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        margin-bottom: 12px;

        .spec-value-tag {
          margin: 0;
        }
      }

      .spec-value-input {
        display: flex;
        align-items: center;
        gap: 8px;

        .spec-input {
          width: 240px;
        }
      }
    }
  }

  .add-spec-group {
    width: 100%;
    height: 40px;
    margin-top: 16px;
    border: 1px dashed #dcdfe6;
    
    &:hover {
      border-color: var(--el-color-primary);
      color: var(--el-color-primary);
    }
  }
}

.sku-table {
  margin-top: 20px;
  
  .sku-specs {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  :deep(.el-input-number) {
    width: 160px;
  }
}

.goods-images-uploader {
  :deep(.el-upload--picture-card) {
    width: 100px;
    height: 100px;
    line-height: 100px;
  }
}

:deep(.w-e-text-container) {
  background-color: #fff;
}

.image-list {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;

  .image-item {
    position: relative;
    width: 120px;
    height: 120px;
    border: 1px solid var(--el-border-color);
    border-radius: 4px;
    overflow: hidden;

    .el-image {
      width: 100%;
      height: 100%;
    }

    .image-actions {
      position: absolute;
      top: 8px;
      right: 8px;
      opacity: 0;
      transition: opacity 0.2s;
    }

    &:hover .image-actions {
      opacity: 1;
    }
  }

  .image-upload {
    width: 120px;
    height: 120px;
    border: 1px dashed var(--el-border-color);
    border-radius: 4px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: border-color 0.2s;
    color: var(--el-text-color-secondary);

    &:hover {
      border-color: var(--el-color-primary);
      color: var(--el-color-primary);
    }

    .el-icon {
      font-size: 24px;
      margin-bottom: 8px;
    }
  }
}

:deep(.el-pagination) {
  margin-top: 16px;
  justify-content: center;
}

.my-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.my-header-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style> 