<template>
  <el-dialog
    :model-value="visible"
    @update:model-value="updateVisible"
    :title="type === 'add' ? '添加商品' : '编辑商品'"
    width="800px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
    >
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <!-- 基本信息 Tab -->
        <el-tab-pane label="基本信息" name="basic">
          <el-form-item label="商品名称" prop="goodsName">
            <el-input v-model="formData.goodsName" placeholder="请输入商品名称" />
          </el-form-item>
          
          <el-form-item label="商品描述" prop="description">
            <el-input v-model="formData.description" type="textarea" placeholder="请输入商品描述" />
          </el-form-item>

          <el-form-item label="封面图" prop="coverImg">
            <el-upload
              class="avatar-uploader"
              :action="uploadUrl"
              :show-file-list="false"
              :on-success="handleUploadSuccess"
            >
              <img v-if="formData.coverImg" :src="formData.coverImg" class="avatar">
              <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
            </el-upload>
          </el-form-item>

          <el-form-item label="商品分类" prop="categoryId">
            <el-tree-select
              v-model="formData.categoryId"
              :data="categoryTreeData"
              :props="{
                value: 'id',
                label: 'name',
                children: 'children'
              }"
              placeholder="请选择商品分类"
              clearable
              check-strictly
            />
          </el-form-item>

          <el-form-item label="商品品牌" prop="brandId">
            <el-select v-model="formData.brandId" placeholder="请选择商品品牌">
              <el-option
                v-for="item in brandOptions"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="供应商" prop="supplierId">
            <el-select v-model="formData.supplierId" placeholder="请选择供应商">
              <el-option
                v-for="item in supplierOptions"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="销售价" prop="price">
            <el-input-number
              v-model="formData.price"
              :precision="2"
              :step="0.1"
              :min="0"
            />
          </el-form-item>

          <el-form-item label="市场价" prop="originalPrice">
            <el-input-number
              v-model="formData.originalPrice"
              :precision="2"
              :step="0.1"
              :min="0"
            />
          </el-form-item>

          <el-form-item label="商品库存" prop="stock">
            <el-input-number
              v-model="formData.stock"
              :min="0"
              :precision="0"
            />
          </el-form-item>

          <el-form-item label="排序" prop="sort">
            <el-input-number
              v-model="formData.sort"
              :min="0"
              :precision="0"
            />
          </el-form-item>

          <el-form-item label="商品状态" prop="status">
            <el-radio-group v-model="formData.status">
              <el-radio label="0">上架</el-radio>
              <el-radio label="1">下架</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item label="商品标签" prop="tagIds">
            <el-select v-model="formData.tagIds" multiple placeholder="请选择商品标签">
              <el-option
                v-for="item in tagOptions"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
        </el-tab-pane>

        <!-- 商品图片 Tab -->
        <el-tab-pane label="商品图片" name="images">
          <el-form-item label="商品图片" prop="images">
            <el-upload
              class="goods-images-uploader"
              :action="uploadUrl"
              list-type="picture-card"
              :on-success="handleImagesUploadSuccess"
              :on-remove="handleImagesRemove"
              multiple
            >
              <el-icon><Plus /></el-icon>
            </el-upload>
          </el-form-item>
        </el-tab-pane>

        <!-- 规格设置 Tab -->
        <el-tab-pane label="规格设置" name="specs">
          <el-form-item label="规格设置">
            <el-radio-group v-model="formData.specType">
              <el-radio :label="1">单规格</el-radio>
              <el-radio :label="2">多规格</el-radio>
            </el-radio-group>
          </el-form-item>
          
          <!-- 多规格设置 -->
          <template v-if="formData.specType === 2">
            <div class="spec-container">
              <!-- 规格组列表 -->
              <div v-for="(group, groupIndex) in formData.specGroups" :key="groupIndex" class="spec-group-item">
                <div class="spec-group-header">
                  <div class="spec-group-title">规格组{{ groupIndex + 1 }}</div>
                  <el-input
                    v-model="group.name"
                    placeholder="请输入规格组名称"
                    class="spec-group-name"
                  />
                  <el-button type="danger" link @click="removeSpecGroup(groupIndex)">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
                
                <!-- 规格值列表 -->
                <div class="spec-values">
                  <div class="spec-values-list">
                    <el-tag
                      v-for="(value, valueIndex) in group.values"
                      :key="valueIndex"
                      closable
                      class="spec-value-tag"
                      @close="handleRemoveSpecValue(groupIndex, valueIndex)"
                    >
                      {{ value.name || value.value }}
                    </el-tag>
                  </div>
                  
                  <!-- 添加规格值 -->
                  <div class="spec-value-input">
                    <el-input
                      v-if="group.inputVisible"
                      ref="saveTagInput"
                      v-model="group.inputValue"
                      size="small"
                      class="spec-input"
                      placeholder="请输入规格值"
                      @keyup.enter="handleInputConfirm(groupIndex)"
                      @blur="handleInputConfirm(groupIndex)"
                    />
                    <el-button
                      v-else
                      size="small"
                      type="primary"
                      link
                      @click="showInput(groupIndex)"
                    >
                      <el-icon><Plus /></el-icon> 添加规格值
                    </el-button>
                  </div>
                </div>
              </div>

              <!-- 添加规格组按钮 -->
              <el-button type="primary" class="add-spec-group" @click="addSpecGroup">
                <el-icon><Plus /></el-icon> 添加规格组
              </el-button>
            </div>
          </template>
        </el-tab-pane>

        <!-- 商品详情 Tab -->
        <el-tab-pane label="商品详情" name="detail">
          <el-form-item label="商品详情" prop="detail">
            <div style="border: 1px solid #ccc">
              <Toolbar
                style="border-bottom: 1px solid #ccc"
                :editor="editorRef"
                :defaultConfig="toolbarConfig"
                mode="default"
              />
              <Editor
                style="height: 500px"
                v-model="formData.detail"
                :defaultConfig="editorConfig"
                mode="default"
                @onCreated="handleCreated"
              />
            </div>
          </el-form-item>
        </el-tab-pane>
      </el-tabs>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取 消</el-button>
      <el-button @click="saveDraft">保存草稿</el-button>
      <el-button type="primary" @click="handleSubmit">确定保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, defineProps, defineEmits, watch, computed, shallowRef, onBeforeUnmount, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { updateGoods, createGoods, getGoodsInfo } from '@/api/shop/goods'
import { Plus, Rank, Delete } from '@element-plus/icons-vue'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import Draggable from 'vuedraggable'

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
  detail: ''
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

// 获取商品详情数据
const loadGoodsDetail = async (id) => {
  try {
    const data = await getGoodsInfo(id)
    console.log('获取商品详情原始数据:', data)
    
    // 重置表单数据
    formData.value = {
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
      detail: data.detail || '',
      specType: data.specType || 1,
      specGroups: [],
      skuList: [],
      images: []
    }

    // 处理规格组数据
    if (data.specGroups && data.specGroups.length > 0) {
      console.log('规格组原始数据:', data.specGroups)
      formData.value.specGroups = data.specGroups.map(group => {
        console.log('处理规格组:', group)
        return {
          id: group.id,
          name: group.name,
          sort: group.sort || 0,
          values: (group.values || []).map(value => ({
            id: value.id,
            name: value.name,
            value: value.value,
            sort: value.sort || 0
          })),
          inputVisible: false,
          inputValue: ''
        }
      })
    }

    console.log('处理后的规格组数据:', formData.value.specGroups)

    // 处理 SKU 列表数据
    if (data.skuList && data.skuList.length > 0) {
      formData.value.skuList = data.skuList.map(sku => ({
        id: sku.id,
        specs: sku.specs || {},
        price: sku.price || 0,
        stock: sku.stock || 0
      }))
    }

    console.log('处理后的表单数据:', formData.value)

    // 设置富文本内容
    nextTick(() => {
      if (editorRef.value && data.detail) {
        editorRef.value.setHtml(data.detail)
      }
    })
  } catch (error) {
    console.error('获取商品详情失败:', error)
    ElMessage.error('获取商品详情失败')
  }
}

// 修改 watch
watch(
  () => props.visible,
  async (val) => {
    if (val) {
      if (props.type === 'edit' && props.data?.goodsId) {
        console.log('开始加载商品详情:', props.data.goodsId)
        await loadGoodsDetail(props.data.goodsId)
      } else {
        // 新增时重置表单
        handleClose()
      }
      // 重置当前标签页
      activeTab.value = 'basic'
    }
  },
  { immediate: true }
)

// 添加 updateVisible 方法处理可见性更新
const updateVisible = (value) => {
  emit('update:visible', value)
}

// 修改 activeTab 的处理
const handleTabChange = async (tabName) => {
  // 验证当前 tab 的数据
  const currentTabFields = tabRules[activeTab.value] || []
  if (currentTabFields.length > 0) {
    try {
      await formRef.value.validateFields(currentTabFields)
      activeTab.value = tabName
    } catch (error) {
      ElMessageBox.confirm(
        '当前标签页有未完成的必填项，是否继续切换？',
        '提示',
        {
          confirmButtonText: '继续',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(() => {
        activeTab.value = tabName
      }).catch(() => {})
    }
  } else {
    activeTab.value = tabName
  }
}

// 修改提交表单方法
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    // 添加规格验证
    if (formData.value.specType === 1 && validSpecGroups.value.length === 0) {
      ElMessage.error('请至少添加一个规格组')
      activeTab.value = 'specs'
      return
    }
    
    // 添加图片验证
    if (!formData.value.coverImg) {
      ElMessage.error('请上传商品封面图')
      activeTab.value = 'basic'
      return
    }
    
    // 添加草稿标记
    const submitData = {
      ...formData.value,
      isDraft: false // 正式保存
    }
    
    try {
      if (props.type === 'add') {
        await createGoods(submitData)
        ElMessage.success('添加成功')
      } else {
        await updateGoods(submitData)
        ElMessage.success('更新成功')
      }
      emit('success')
      handleClose()
    } catch (error) {
      console.error(props.type === 'add' ? '添加失败:' : '更新失败:', error)
      ElMessage.error(props.type === 'add' ? '添加失败' : '更新失败')
    }
  } catch (error) {
    // 获取第一个验证失败的字段
    const firstErrorField = Object.keys(error.fields)[0]
    // 找到对应的 tab
    for (const [tabName, fields] of Object.entries(tabRules)) {
      if (fields.includes(firstErrorField)) {
        activeTab.value = tabName
        break
      }
    }
  }
}

// 添加保存草稿方法
const saveDraft = async () => {
  try {
    const draftData = {
      ...formData.value,
      isDraft: true // 标记为草稿
    }
    
    if (props.type === 'add') {
      await createGoods(draftData)
      ElMessage.success('草稿保存成功')
    } else {
      await updateGoods(draftData)
      ElMessage.success('草稿更新成功')
    }
  } catch (error) {
    console.error('保存草稿失败:', error)
    ElMessage.error('保存草稿失败')
  }
}

// 编辑器实例，必须用 shallowRef
const editorRef = shallowRef()

// 编辑器配置
const toolbarConfig = {
  excludeKeys: [] // 可以配置需要排除的菜单
}

const editorConfig = {
  placeholder: '请输入商品详情...',
  MENU_CONF: {
    uploadImage: {
      server: props.uploadUrl,
      fieldName: 'file',
      // 单个文件上传成功后
      onSuccess(file, res) {
        return res.url // 返回图片 URL
      },
    }
  }
}

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor) {
    editor.destroy()
  }
})

// 编辑器回调函数
const handleCreated = (editor) => {
  editorRef.value = editor
}

// 修改 handleClose 方法
const handleClose = () => {
  formRef.value?.resetFields()
  // 先销毁编辑器实例
  const editor = editorRef.value
  if (editor) {
    editor.destroy()
  }
  
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
    detail: ''
  }
  emit('update:visible', false)
  emit('close')
}

// 处理图片上传成功
const handleUploadSuccess = (res) => {
  formData.value.coverImg = res.url
}

// 处理多图片上传成功
const handleImagesUploadSuccess = (res) => {
  formData.value.images.push(res.url)
}

// 处理图片移除
const handleImagesRemove = (file) => {
  const index = formData.value.images.indexOf(file.url)
  if (index !== -1) {
    formData.value.images.splice(index, 1)
  }
}

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
  return combinations.map((combination, index) => {
    // 创建规格对象
    const specs = {}
    validSpecGroups.value.forEach((group, groupIndex) => {
      specs[group.name] = combination[groupIndex]
    })
    
    // 查找已存在的 SKU 数据
    const existingSku = formData.value.skuList?.find(
      sku => JSON.stringify(sku.specs) === JSON.stringify(specs)
    ) || {
      specs,
      price: formData.value.price || 0,
      stock: formData.value.stock || 0
    }
    
    // 返回新的对象，避免直接修改原对象
    return {
      ...existingSku,
      specs: { ...specs }
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
    // 只在多规格模式下更新 SKU 列表
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
</style> 