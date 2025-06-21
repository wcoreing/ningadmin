<template>
  <el-dialog
    :title="title"
    v-model="dialogVisible"
    width="1600px"
    center
    :close-on-click-modal="false"
    append-to-body
    destroy-on-close
  >
    <div class="goods-selector">
      <!-- 搜索和筛选区 -->
      <div class="filter-container">
        <el-form :inline="true" :model="queryParams" @keyup.enter="handleQuery">
          <el-form-item>
            <el-input
              v-model="queryParams.name"
              placeholder="输入商品名称或ID搜索"
              clearable
              prefix-icon="Search"
              style="width: 220px"
            />
          </el-form-item>
          <el-form-item>
            <el-select
              v-model="queryParams.categoryId" 
              placeholder="商品分类"
              clearable
              style="width: 180px"
            >
              <el-option
                v-for="item in categoryOptions"
                :key="item.categoryId"
                :label="item.categoryName"
                :value="item.categoryId"
              >
                <span :style="{ paddingLeft: item.level * 20 + 'px' }">{{ item.categoryName }}</span>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleQuery">查询</el-button>
            <el-button @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 商品列表 -->
      <div class="goods-list" v-loading="loading">
        <el-row :gutter="12">
          <el-col 
            :span="4"
            v-for="item in goodsList" 
            :key="item.id"
          >
            <div 
              class="goods-item"
              :class="{ 'selected': isSelected(item) }"
              @click="handleSelect(item)"
            >
              <div class="use-button" v-if="isSelected(item)">
                <el-icon><Check /></el-icon>
              </div>
              <el-image 
                :src="item.coverImg" 
                fit="contain"
                :preview-disabled="true"
                class="goods-image"
              >
                <template #error>
                  <div class="image-placeholder">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <div class="goods-info">
                <div class="goods-name" :title="item.name">{{ item.name }}</div>
                <div class="goods-price">
                  <span class="price">¥{{ item.minPrice || item.price || 0 }}</span>
                  <span class="stock">库存: {{ item.totalStock || item.stock || 0 }}</span>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
        
        <!-- 空状态 -->
        <div v-if="!loading && goodsList.length === 0" class="empty-state">
          <el-empty description="暂无商品数据" />
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-container" v-if="total > 0">
        <el-pagination
          v-model:current-page="queryParams.pageNum"
          v-model:page-size="queryParams.pageSize"
          :total="total"
          :page-sizes="[18, 36, 54, 72]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="confirmSelection">确 定</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { getGoodsList } from '@/api/shop/goods'
import { getCategoryList } from '@/api/shop/category'
import { ElMessage } from 'element-plus'
import { Picture, Select, Search, Check } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: Boolean,
  title: {
    type: String,
    default: '选择商品'
  },
  multiple: {
    type: Boolean,
    default: false
  },
  max: {
    type: Number,
    default: 0
  },
  selectedIds: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'select'])

// 控制对话框显示
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// 数据列表
const loading = ref(false)
const goodsList = ref([])
const total = ref(0)
const categoryOptions = ref([])

// 查询参数
const queryParams = ref({
  pageNum: 1,
  pageSize: 18,
  name: '',  // 用于搜索商品名称
  goodsId: undefined,  // 用于搜索商品ID
  categoryId: undefined
})

// 选中数据
const selectedGoods = ref([])

// 获取商品列表
const getList = async () => {
  loading.value = true
  try {
    const res = await getGoodsList(queryParams.value)
    console.log('商品列表数据:', res)
    goodsList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取商品列表失败:', error)
    ElMessage.error('获取商品列表失败')
    goodsList.value = []
    total.value = 0
  }
  loading.value = false
}

// 获取分类列表
const getCategories = async () => {
  try {
    const res = await getCategoryList({ module: 'goods' })
    console.log('分类列表数据:', res)
    // 将多级分类转换为带层级的一维数组
    const flattenCategories = (categories, level = 0) => {
      let result = []
      if (Array.isArray(categories)) {
        categories.forEach(item => {
          result.push({
            categoryId: item.id || item.categoryId,
            categoryName: item.name || item.categoryName,
            level: level  // 添加层级信息
          })
          if (item.children && item.children.length) {
            result = result.concat(flattenCategories(item.children, level + 1))
          }
        })
      }
      return result
    }
    const categories = res.list || res || []
    categoryOptions.value = flattenCategories(categories)
  } catch (error) {
    console.error('获取分类列表失败:', error)
    categoryOptions.value = []
  }
}

// 查询操作
const handleQuery = () => {
  // 构建查询参数
  const params = {
    pageNum: 1,
    pageSize: queryParams.value.pageSize,
    categoryId: queryParams.value.categoryId
  }
  
  // 如果输入的是数字，则认为是在搜索商品ID
  if (queryParams.value.name && /^\d+$/.test(queryParams.value.name)) {
    params.goodsId = parseInt(queryParams.value.name)
  } else if (queryParams.value.name) {
    params.name = queryParams.value.name
  }
  
  queryParams.value = params
  getList()
}

// 重置操作
const resetQuery = () => {
  queryParams.value = {
    pageNum: 1,
    pageSize: 18,
    name: '',
    goodsId: undefined,
    categoryId: undefined
  }
  handleQuery()
}

// 选择商品
const handleSelect = (item) => {
  // 统一使用 goodsId 作为商品ID，同时兼容 id
  const goodsId = item.goodsId || item.id
  // 统一价格字段，优先使用 minPrice，然后是 price
  const price = item.minPrice || item.price || 0
  // 统一库存字段，优先使用 totalStock，然后是 stock
  const stock = item.totalStock || item.stock || 0
  // 统一商品名称字段
  const goodsName = item.goodsName || item.name
  // 统一图片字段
  const goodsImage = item.goodsImage || item.coverImg || item.image
  
  const goodsData = {
    ...item,
    goodsId: goodsId,
    id: goodsId,
    name: goodsName,
    goodsName: goodsName,
    price: price,
    stock: stock,
    coverImg: goodsImage,
    goodsImage: goodsImage,
    image: goodsImage
  }
  
  console.log('选择商品数据:', goodsData)
  
  if (props.multiple) {
    const index = selectedGoods.value.findIndex(i => (i.goodsId || i.id) === goodsId)
    if (index > -1) {
      selectedGoods.value.splice(index, 1)
    } else {
      if (props.max && selectedGoods.value.length >= props.max) {
        ElMessage.warning(`最多只能选择 ${props.max} 个商品`)
        return
      }
      selectedGoods.value.push(goodsData)
    }
  } else {
    selectedGoods.value = [goodsData]
    ElMessage.success('已选择商品：' + item.name)
    confirmSelection()
  }
}

// 判断是否选中
const isSelected = (item) => {
  const goodsId = item.goodsId || item.id
  return selectedGoods.value.some(i => (i.goodsId || i.id) === goodsId)
}

// 确认选择
const confirmSelection = () => {
  console.log('确认选择 - 当前选中商品:', selectedGoods.value)
  console.log('确认选择 - 是否多选模式:', props.multiple)
  
  if (selectedGoods.value.length === 0) {
    ElMessage.warning('请选择商品')
    return
  }
  
  const emitData = props.multiple ? selectedGoods.value : selectedGoods.value[0]
  console.log('确认选择 - 发送的数据:', emitData)
  
  emit('select', emitData)
  dialogVisible.value = false
}

// 分页方法
const handleSizeChange = (val) => {
  queryParams.value.pageSize = val
  getList()
}

const handleCurrentChange = (val) => {
  queryParams.value.pageNum = val
  getList()
}

// 监听对话框打开，初始化数据
watch(dialogVisible, (newVal) => {
  if (newVal) {
    // 重置选中状态
    selectedGoods.value = []
    // 重置查询参数
    queryParams.value = {
      pageNum: 1,
      pageSize: 18,
      name: '',
      goodsId: undefined,
      categoryId: undefined
    }
    // 加载数据
    getCategories()
    getList()
  }
})

// 添加方法获取商品详情
const getGoodsDetail = async (id) => {
  try {
    const res = await getGoodsList({ ids: [id] })
    if (res.list && res.list.length > 0) {
      selectedGoods.value = [res.list[0]]
    }
  } catch (error) {
    console.error('获取商品详情失败:', error)
  }
}
</script>

<style lang="scss" scoped>
.goods-selector {
  min-height: 400px;
  padding: 0 4px;
  
  .filter-container {
    margin-bottom: 12px;
  }

  .goods-list {
    min-height: 300px;
    
    .goods-item {
      position: relative;
      margin-bottom: 12px;
      border-radius: 4px;
      overflow: hidden;
      cursor: pointer;
      transition: all 0.3s;
      border: 1px solid #ebeef5;
      background-color: #fff;

      &.selected {
        border-color: var(--el-color-primary);
        box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
        background-color: #f0f9ff;
      }

      .use-button {
        position: absolute;
        top: 8px;
        right: 8px;
        background: var(--el-color-success);
        color: #fff;
        width: 24px;
        height: 24px;
        border-radius: 50%;
        font-size: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.3s;
        z-index: 2;
        
        .el-icon {
          font-size: 12px;
        }
      }

      &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        transform: translateY(-2px);
      }

      .goods-image {
        width: 100%;
        height: 110px;
        padding: 8px;
        box-sizing: border-box;
        background-color: #fff;
        display: block;

        :deep(.el-image__inner) {
          width: 100%;
          height: 100%;
          object-fit: contain;
        }
      }

      .image-placeholder {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #f5f7fa;
        color: #909399;
        font-size: 24px;
      }

      .goods-info {
        padding: 6px 8px;

        .goods-name {
          font-size: 14px;
          line-height: 1.4;
          height: 28px;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
        }

        .goods-price {
          margin-top: 4px;
          display: flex;
          justify-content: space-between;
          align-items: center;

          .price {
            color: #f56c6c;
            font-size: 14px;
            font-weight: bold;
          }

          .stock {
            color: #909399;
            font-size: 12px;
          }
        }
      }
    }
    
    .empty-state {
      text-align: center;
      padding: 40px 0;
    }
  }

  .pagination-container {
    margin-top: 12px;
    text-align: right;
  }

  .el-loading-mask {
    background-color: rgba(255, 255, 255, 0.8);
  }
}

:deep(.el-dialog__body) {
  padding: 10px;
  max-height: calc(100vh - 200px);
  overflow-y: auto;
}

:deep(.el-dialog__header) {
  padding: 12px;
  margin-right: 0;
}

:deep(.el-dialog__footer) {
  padding: 10px;
}
</style> 