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
              :class="{ 'is-selected': isSelected(item) }"
              @click="handleSelect(item)"
            >
              <el-image 
                :src="item.coverImg" 
                fit="contain"
                :preview-teleported="true"
                :preview-src-list="[item.coverImg]"
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
                  <span class="price">¥{{ item.price }}</span>
                  <span class="stock">库存: {{ item.stock }}</span>
                </div>
              </div>
              <div class="select-mask" v-if="isSelected(item)">
                <el-icon><Select /></el-icon>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 分页 -->
      <div class="pagination-container">
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
import { ref, computed, watch } from 'vue'
import { getGoodsList } from '@/api/shop/goods'
import { getCategoryList } from '@/api/shop/category'
import { ElMessage } from 'element-plus'
import { Picture, Select, Search } from '@element-plus/icons-vue'

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
    goodsList.value = res.list
    total.value = res.total
    // 如果有选中商品，保持选中状态
    if (selectedGoods.value.length > 0) {
      const selected = selectedGoods.value[0]
    }
  } catch (error) {
    console.error('获取商品列表失败:', error)
  }
  loading.value = false
}

// 获取分类列表
const getCategories = async () => {
  try {
    const res = await getCategoryList({
      parentId: 0
    })
    // 将多级分类转换为带层级的一维数组
    const flattenCategories = (categories, level = 0) => {
      let result = []
      categories.forEach(item => {
        result.push({
          categoryId: item.id,
          categoryName: item.name,
          level: level  // 添加层级信息
        })
        if (item.children && item.children.length) {
          result = result.concat(flattenCategories(item.children, level + 1))
        }
      })
      return result
    }
    categoryOptions.value = flattenCategories(res)
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 查询操作
const handleQuery = () => {
  // 如果输入的是数字，则认为是在搜索商品ID
  if (queryParams.value.name && /^\d+$/.test(queryParams.value.name)) {
    queryParams.value.goodsId = parseInt(queryParams.value.name)
    queryParams.value.name = ''
  } else {
    queryParams.value.goodsId = undefined
  }
  queryParams.value.pageNum = 1
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
  if (props.multiple) {
    const index = selectedGoods.value.findIndex(i => i.goodsId === item.goodsId)
    if (index > -1) {
      selectedGoods.value.splice(index, 1)
    } else {
      if (props.max && selectedGoods.value.length >= props.max) {
        ElMessage.warning(`最多只能选择 ${props.max} 个商品`)
        return
      }
      selectedGoods.value.push(item)
    }
  } else {
    selectedGoods.value = [item]
    ElMessage.success('已选择商品：' + item.name)
    confirmSelection()
  }
}

// 判断是否选中
const isSelected = (item) => {
  return selectedGoods.value.some(i => i.goodsId === item.goodsId)
}

// 确认选择
const confirmSelection = () => {
  if (selectedGoods.value.length === 0) {
    ElMessage.warning('请选择商品')
    return
  }
  emit('select', props.multiple ? selectedGoods.value : selectedGoods.value[0])
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

// 监听选中ID变化
watch(() => props.selectedIds, (newVal) => {
  if (newVal && newVal.length > 0) {
    // 获取所有选中商品的详情
    const getSelectedGoods = async () => {
      try {
        const res = await getGoodsList({ ids: newVal })
        if (res.list && res.list.length > 0) {
          selectedGoods.value = res.list
        }
      } catch (error) {
        console.error('获取商品详情失败:', error)
      }
    }
    getSelectedGoods()
  } else {
    selectedGoods.value = []
  }
}, { immediate: true })

// 初始化
getCategories()
getList()
</script>

<style lang="scss" scoped>
.goods-selector {
  min-height: 400px;
  padding: 0 4px;
  
  .filter-container {
    margin-bottom: 12px;
  }

  .goods-list {
    .goods-item {
      position: relative;
      margin-bottom: 12px;
      border-radius: 4px;
      overflow: hidden;
      cursor: pointer;
      transition: all 0.3s;
      border: 1px solid #ebeef5;
      background-color: #fff;

      &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        transform: translateY(-2px);
      }

      &.is-selected {
        border-color: var(--el-color-primary);
        box-shadow: 0 0 0 2px var(--el-color-primary-light-8);
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

      .select-mask {
        position: absolute;
        top: 8px;
        right: 8px;
        width: 24px;
        height: 24px;
        background: var(--el-color-primary);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
      }
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