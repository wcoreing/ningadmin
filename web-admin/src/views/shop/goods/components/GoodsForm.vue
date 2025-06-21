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
            :tag-options="tagOptions"
            @type-change="handleTypeChange"
          />
        </el-tab-pane>

        <!-- 商品图片 Tab -->
        <el-tab-pane label="商品图片" name="images">
          <goods-images :form-data="formData" />
        </el-tab-pane>

        <!-- 规格设置 Tab -->
        <el-tab-pane label="规格设置" name="specs">
          <template #label>规格设置</template>
          
          <!-- 规格类型选择器 - 始终显示 -->
          <div class="spec-type-card">
            <div class="spec-type-selector">
              <div 
                class="spec-type-option"
                :class="{ active: formData.specType === 1 }"
                @click="handleSpecTypeChange(1)"
              >
                <div class="option-icon">
                  <el-icon><Document /></el-icon>
                </div>
                <div class="option-content">
                  <h4>单规格</h4>
                  <p>商品只有一种规格，统一价格和库存</p>
                </div>
                <div class="option-check">
                  <el-icon v-if="formData.specType === 1"><Check /></el-icon>
                </div>
              </div>
              
              <div 
                class="spec-type-option"
                :class="{ active: formData.specType === 2 }"
                @click="handleSpecTypeChange(2)"
              >
                <div class="option-icon">
                  <el-icon><Grid /></el-icon>
                </div>
                <div class="option-content">
                  <h4>多规格</h4>
                  <p>商品有多种规格组合，如颜色、尺寸等</p>
                </div>
                <div class="option-check">
                  <el-icon v-if="formData.specType === 2"><Check /></el-icon>
                </div>
              </div>
              
              <div 
                class="spec-type-option"
                :class="{ active: formData.specType === 3 }"
                @click="handleSpecTypeChange(3)"
              >
                <div class="option-icon">
                  <el-icon><SetUp /></el-icon>
                </div>
                <div class="option-content">
                  <h4>套餐规格</h4>
                  <p>商品为套餐形式，包含多个组和选项</p>
                </div>
                <div class="option-check">
                  <el-icon v-if="formData.specType === 3"><Check /></el-icon>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 单规格和多规格设置 -->
          <SpecSettings 
            v-if="formData.specType === 1 || formData.specType === 2"
            :form-data="formData"
            :upload-url="uploadUrl"
            ref="specSettingsRef"
            @spec-type-change="handleSpecTypeChange"
            :show-type-selector="false"
          />
          
          <!-- 套餐规格设置 -->
          <combo-settings 
            v-if="formData.specType === 3"
            :form-data="formData"
            ref="comboSettingsRef"
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
import { Plus, Rank, Delete, Search, UploadFilled, Close, Document, Grid, SetUp, Check } from '@element-plus/icons-vue'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import FileManager from '@/components/FileManager/index.vue'
import BasicInfo from './BasicInfo.vue'
import GoodsImages from './GoodsImages.vue'
import SpecSettings from './SpecSettings.vue'
import ComboSettings from './ComboSettings.vue'
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
const comboSettingsRef = ref(null)

const formData = ref({
  goodsId: undefined,
  goodsName: '',
  categoryId: undefined,
  brand: '',
  supplier: '',
        type: 1, // 商品类型：1-普通商品，2-虚拟商品
  // 删除通用字段，每种规格使用自己对象中的字段
  sort: 0,
  status: 0,
  description: '',
  coverImg: '',
  tagNames: [],
  images: [],
  specType: 1,
  specGroups: [],
  skuList: [],
  comboGroups: [], // 套餐组配置
  detailDesc: '',
  
  // 多规格数据存储 - 支持同时保存多种规格类型的数据
  allSpecData: {
    singleSpecData: {
      price: 0,
      originalPrice: 0,
      costPrice: 0,
      stock: 0,
      weight: 0,
      volume: 0
    },
    multiSpecData: {
      specGroups: [],
      skuList: []
    },
    comboSpecData: {
      price: 0,
      comboGroups: []
    }
  }
})

// 动态表单验证规则
const formRules = computed(() => {
  const baseRules = {
  goodsName: [
    { required: true, message: '请输入商品名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  categoryId: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ],
    brand: [
    { required: true, message: '请选择商品品牌', trigger: 'change' }
  ],
    supplier: [
      { required: true, message: '请输入供应商名称', trigger: 'blur' }
    ]
  }

  return baseRules
})

// 更新 tab 验证规则映射
const tabRules = computed(() => {
  return {
    basic: ['goodsName', 'categoryId', 'brand', 'supplier'],
    specs: [], // 规格相关的验证在提交时单独处理
  images: [], // 图片相关的验证
  detail: [] // 详情相关的验证
}
})

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
      brand: '',
      supplier: '',
      type: 1,
      // 删除通用price字段，每种规格使用自己的价格
      originalPrice: 0,
      costPrice: 0,
      stock: 0,
      weight: 0,
      volume: 0,
      sort: 0,
      status: 0,
      description: '',
      coverImg: '',
      tagNames: [],
      images: [],
      specType: 1,
      specGroups: [],
      skuList: [],
      comboGroups: [],
      detailDesc: '',
      allSpecData: {
        singleSpecData: {
          price: 0,
          originalPrice: 0,
          costPrice: 0,
          stock: 0,
          weight: 0,
          volume: 0
        },
        multiSpecData: {
          specGroups: [],
          skuList: []
        },
        comboSpecData: {
          price: 0,
          comboGroups: []
        }
      }
    }

    // 等待下一个 tick，确保重置完成
    await nextTick()
    
    // 删除通用字段的提取，数据直接存储在各自规格对象中
    
    // 然后设置新数据
    formData.value = {
      ...formData.value,
      goodsId: data.goodsId,
      goodsName: data.goodsName,
      categoryId: data.categoryId,
      brand: data.brand || '',
      supplier: data.supplier || '',
      type: data.type || 1, // 商品类型
      // 删除通用字段，每种规格使用自己对象中的字段
      sort: data.sort,
      status: Number(data.status),
      description: data.description,
      coverImg: data.coverImg,
      tagNames: data.tagNames || [],
      images: data.images || [],
      specType: data.specType || 1,
      specGroups: data.multiSpecData?.specGroups ? JSON.parse(JSON.stringify(data.multiSpecData.specGroups)) : [],
      skuList: data.multiSpecData?.skuList ? JSON.parse(JSON.stringify(data.multiSpecData.skuList)) : [],
      comboGroups: data.comboSpecData?.comboGroups ? JSON.parse(JSON.stringify(data.comboSpecData.comboGroups)) : [], // 套餐组数据
      detailDesc: data.detailDesc || '<p><br></p>',
      
      // 加载多规格数据存储
      allSpecData: {
        singleSpecData: data.singleSpecData || null,
        multiSpecData: data.multiSpecData || null,
        comboSpecData: data.comboSpecData || null
      }
    }
    
    console.log('GoodsForm - 设置套餐组数据:', formData.value.comboGroups)
    console.log('GoodsForm - 原始数据:', data.comboGroups)

   
    
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
          brand: '',
          supplier: '',
          type: 1,
          // 删除通用字段，每种规格使用自己对象中的字段
          sort: 0,
          status: 0,
          description: '',
          coverImg: '',
          tagIds: [],
          images: [],
          specType: 1,
          specGroups: [],
          skuList: [],
          comboGroups: [],
          detailDesc: '<p><br></p>',
          allSpecData: {
            singleSpecData: {
              price: 0,
              originalPrice: 0,
              costPrice: 0,
              stock: 0,
              weight: 0,
              volume: 0
            },
            multiSpecData: {
              specGroups: [],
              skuList: []
            },
            comboSpecData: {
              price: 0,
              comboGroups: []
            }
          }
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
  
  // 如果切换到规格设置Tab且是套餐规格，确保套餐组已初始化
  if (tabName === 'specs' && formData.value.specType === 3) {
    nextTick(() => {
      if (!formData.value.comboGroups || formData.value.comboGroups.length === 0) {
        comboSettingsRef.value?.initComboGroups?.()
      }
    })
  }
}

// 添加一个用于刷新商品详情组件的 key
const detailKey = ref(0)

// 修改提交表单方法
const handleSubmit = async () => {
  try {
    // 基础表单验证
    await formRef.value.validate()
    
    // 套餐规格验证
    if (formData.value.specType === 3) {
      // 验证套餐价格
      const comboPrice = formData.value.allSpecData?.comboSpecData?.price
      if (comboPrice === undefined || comboPrice < 0) {
        ElMessage.error('请设置套餐价格，价格必须大于等于0')
        activeTab.value = 'specs'
        return
      }
      
      // 验证套餐组
      if (!formData.value.comboGroups || formData.value.comboGroups.length === 0) {
        ElMessage.error('套餐规格必须设置套餐组')
        activeTab.value = 'specs'
        return
      }
      
      // 验证套餐组是否完整
      const hasValidGroups = formData.value.comboGroups.some(group => 
        group.name && group.options && group.options.length > 0
      )
      if (!hasValidGroups) {
        ElMessage.error('请完善套餐组设置，至少需要一个完整的套餐组')
        activeTab.value = 'specs'
        return
      }
      
      // 验证每个套餐组的选项
      for (const group of formData.value.comboGroups) {
        if (!group.name) {
          ElMessage.error('请为所有套餐组设置名称')
          activeTab.value = 'specs'
          return
        }
        
        if (!group.options || group.options.length === 0) {
          ElMessage.error(`套餐组"${group.name}"必须至少有一个选项`)
          activeTab.value = 'specs'
          return
        }
        
        // 验证必选组是否有默认选项
        if (group.isRequired) {
          const hasDefault = group.options.some(option => option.isDefault)
          if (!hasDefault) {
            ElMessage.error(`必选套餐组"${group.name}"必须至少有一个默认选项`)
            activeTab.value = 'specs'
            return
          }
        }
        
        // 验证选项的完整性
        for (const option of group.options) {
          if (!option.name || !option.value) {
            ElMessage.error(`套餐组"${group.name}"中有选项名称或值为空`)
            activeTab.value = 'specs'
            return
          }
          
          if (option.addPrice < 0 || !Number.isFinite(Number(option.addPrice))) {
            ElMessage.error(`套餐组"${group.name}"中的选项"${option.name}"加价金额设置错误`)
            activeTab.value = 'specs'
            return
          }
          
          if (option.quantity <= 0 || !Number.isInteger(Number(option.quantity))) {
            ElMessage.error(`套餐组"${group.name}"中的选项"${option.name}"数量必须为正整数`)
            activeTab.value = 'specs'
            return
          }
        }
      }
    }
    
    // 单规格商品验证
    else if (formData.value.specType === 1) {
      // 验证价格 - 从allSpecData.singleSpecData中获取
      const singleSpecPrice = formData.value.allSpecData?.singleSpecData?.price
      if (!singleSpecPrice || singleSpecPrice <= 0) {
        ElMessage.error('请在规格设置页面设置商品价格，价格必须大于0')
        activeTab.value = 'specs'
        return
      }
      // 验证库存 - 从allSpecData.singleSpecData中获取
      const singleSpecStock = formData.value.allSpecData?.singleSpecData?.stock
      if (singleSpecStock < 0 || !Number.isInteger(Number(singleSpecStock))) {
        ElMessage.error('请在规格设置页面设置商品库存，库存必须为非负整数')
        activeTab.value = 'specs'
        return
      }
    }
    
    // 多规格商品的额外验证
    if (formData.value.specType === 2) {
      // 验证是否有规格组
      if (!formData.value.specGroups || formData.value.specGroups.length === 0) {
        ElMessage.error('多规格商品必须设置规格组')
        activeTab.value = 'specs'
        return
      }

      // 验证规格组是否完整
      const hasValidGroups = formData.value.specGroups.some(group => 
        group.name && group.values && group.values.length > 0
      )
      if (!hasValidGroups) {
        ElMessage.error('请完善规格组设置，至少需要一个完整的规格组')
        activeTab.value = 'specs'
        return
      }

      // 验证SKU列表
      if (!formData.value.skuList || formData.value.skuList.length === 0) {
        ElMessage.error('多规格商品必须有SKU数据')
        activeTab.value = 'specs'
        return
      }

      // 验证每个SKU的价格和库存
      const invalidSku = formData.value.skuList.find(sku => 
        !sku.price || sku.price < 0 || !Number.isFinite(Number(sku.price)) ||
        sku.stock < 0 || !Number.isInteger(Number(sku.stock)) ||
        (sku.originalPrice && (sku.originalPrice < 0 || !Number.isFinite(Number(sku.originalPrice)))) ||
        (sku.costPrice && (sku.costPrice < 0 || !Number.isFinite(Number(sku.costPrice)))) ||
        (sku.weight && (sku.weight < 0 || !Number.isFinite(Number(sku.weight)))) ||
        (sku.volume && (sku.volume < 0 || !Number.isFinite(Number(sku.volume))))
      )
      if (invalidSku) {
        ElMessage.error('请检查SKU的价格、库存等设置，价格必须大于0，库存必须为非负整数，其他数值字段必须为非负数')
        activeTab.value = 'specs'
        return
      }

      // 处理 SKU 数据格式
      formData.value.skuList = formData.value.skuList.map(sku => ({
        specs: sku.specs,
        price: Number(sku.price || 0),
        originalPrice: Number(sku.originalPrice || 0),
        costPrice: Number(sku.costPrice || 0),
        stock: Number(sku.stock || 0),
        weight: Number(sku.weight || 0),
        volume: Number(sku.volume || 0),
        image: sku.image || ''
      }))
    }

    // 在提交前，保存当前规格类型的数据到allSpecData
    saveCurrentSpecData()

    const submitData = {
      ...formData.value,
      // 确保包含所有规格类型的数据
      singleSpecData: formData.value.allSpecData.singleSpecData,
      multiSpecData: formData.value.allSpecData.multiSpecData,
      comboSpecData: formData.value.allSpecData.comboSpecData
    }
    
    console.log('GoodsForm - 提交数据:', submitData)

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
    if (error.message) {
      ElMessage.error(error.message)
    } else {
    ElMessage.error('提交失败，请检查表单')
    }
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

// 商品类型变化处理
const handleTypeChange = (type) => {
  formData.value.type = type
}

// 规格类型变化处理 - 支持多规格数据同时存储
const handleSpecTypeChange = (specType) => {
  console.log('GoodsForm - 规格类型变化:', formData.value.specType, '->', specType)
  
  // 在切换前，保存当前规格类型的数据
  saveCurrentSpecData()
  
  // 切换规格类型
  const oldSpecType = formData.value.specType
  formData.value.specType = specType
  
  // 加载新规格类型的数据
  loadSpecDataByType(specType)
  
  // 如果切换到套餐规格且没有套餐组数据，初始化套餐组
  if (specType === 3 && (!formData.value.comboGroups || formData.value.comboGroups.length === 0)) {
      nextTick(() => {
        comboSettingsRef.value?.initComboGroups?.()
      })
    }
  
  console.log('GoodsForm - 规格类型切换完成:', oldSpecType, '->', specType)
  console.log('GoodsForm - 当前所有规格数据:', formData.value.allSpecData)
}

// 保存当前规格类型的数据
const saveCurrentSpecData = () => {
  const currentSpecType = formData.value.specType
  
  switch (currentSpecType) {
    case 1: // 单规格
      // 单规格数据已经直接存储在allSpecData.singleSpecData中，无需额外保存
      break
      
    case 2: // 多规格
      formData.value.allSpecData.multiSpecData = {
        specGroups: JSON.parse(JSON.stringify(formData.value.specGroups || [])),
        skuList: JSON.parse(JSON.stringify(formData.value.skuList || []))
      }
      break
      
    case 3: // 套餐规格
      // 套餐规格数据已经直接存储在allSpecData.comboSpecData中，无需额外保存
      formData.value.allSpecData.comboSpecData = {
        ...formData.value.allSpecData.comboSpecData,
        comboGroups: JSON.parse(JSON.stringify(formData.value.comboGroups || []))
      }
      break
  }
  
  console.log('GoodsForm - 保存规格数据:', currentSpecType, formData.value.allSpecData)
}

// 加载指定规格类型的数据（三种规格类型数据完全独立，不相互干扰）
const loadSpecDataByType = (specType) => {
  switch (specType) {
    case 1: // 单规格
      // 单规格数据直接存储在allSpecData.singleSpecData中，无需加载到formData
      // 清空当前规格类型不需要的显示数据
      formData.value.specGroups = []
      formData.value.skuList = []
      formData.value.comboGroups = []
      break
      
    case 2: // 多规格
      if (formData.value.allSpecData.multiSpecData) {
        const data = formData.value.allSpecData.multiSpecData
        formData.value.specGroups = JSON.parse(JSON.stringify(data.specGroups || []))
        formData.value.skuList = JSON.parse(JSON.stringify(data.skuList || []))
      } else {
        // 如果没有保存的数据，使用默认值
        formData.value.specGroups = []
        formData.value.skuList = []
      }
      // 多规格数据直接存储在allSpecData.multiSpecData中
      formData.value.comboGroups = []
      break
      
    case 3: // 套餐规格
      console.log('GoodsForm - 套餐规格数据检查:', formData.value.allSpecData.comboSpecData)
      if (formData.value.allSpecData.comboSpecData) {
        const data = formData.value.allSpecData.comboSpecData
        console.log('GoodsForm - 套餐规格数据存在:', data)
        console.log('GoodsForm - 套餐价格:', data.price)
        formData.value.comboGroups = JSON.parse(JSON.stringify(data.comboGroups || []))
      } else {
        // 如果没有保存的数据，使用默认值
        console.log('GoodsForm - 套餐规格数据不存在，使用默认值')
        formData.value.comboGroups = []
        // 初始化comboSpecData结构
        if (!formData.value.allSpecData.comboSpecData) {
          formData.value.allSpecData.comboSpecData = { price: 0, comboGroups: [] }
        }
      }
      // 套餐规格数据直接存储在allSpecData.comboSpecData中，清空其他规格数据
      formData.value.specGroups = []
      formData.value.skuList = []
      break
  }
  
  console.log('GoodsForm - 加载规格数据:', specType, {
    specGroups: formData.value.specGroups,
    comboGroups: formData.value.comboGroups,
    allSpecData: formData.value.allSpecData
  })
}

// 添加 handleClose 方法
const handleClose = () => {
  // 重置表单数据
  formData.value = {
    goodsId: undefined,
    goodsName: '',
    categoryId: undefined,
    brand: '',
    supplier: '',
    type: 1,
    // 删除通用字段，每种规格使用自己对象中的字段
    sort: 0,
    status: 0,
    description: '',
    coverImg: '',
    tagNames: [],
    images: [],
    specType: 1,
    specGroups: [],
    skuList: [],
    comboGroups: [],
    detailDesc: '',
    // 初始化三种规格类型的数据结构
    allSpecData: {
      singleSpecData: {
        price: 0,
        originalPrice: 0,
        costPrice: 0,
        stock: 0,
        weight: 0,
        volume: 0
      },
      multiSpecData: {
        specGroups: [],
        skuList: []
      },
      comboSpecData: {
        price: 0,
        comboGroups: []
      }
    }
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

.spec-type-card {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #e4e7ed;
  margin-bottom: 20px;
  overflow: hidden;
  
  .spec-type-selector {
    padding: 24px;
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 16px;
    
    .spec-type-option {
      display: flex;
      align-items: center;
      padding: 20px;
      border: 2px solid #e4e7ed;
      border-radius: 12px;
      cursor: pointer;
      transition: all 0.3s ease;
      background: #fafbfc;
      
      &:hover {
        border-color: #409eff;
        background: #f0f9ff;
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
      }
      
      &.active {
        border-color: #409eff;
        background: #f0f9ff;
        box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
      }
      
      .option-icon {
        margin-right: 16px;
        font-size: 24px;
        color: #409eff;
      }
      
      .option-content {
        flex: 1;
        
        h4 {
          margin: 0 0 4px 0;
          font-size: 16px;
          font-weight: 600;
          color: #303133;
        }
        
        p {
          margin: 0;
          font-size: 14px;
          color: #606266;
          line-height: 1.4;
        }
      }
      
      .option-check {
        font-size: 20px;
        color: #67c23a;
      }
    }
  }
}
</style> 