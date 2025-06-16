<template>
  <div class="app-container">
 

    <!-- 搜索区域 -->
    <el-card class="filter-container">
      <el-form :inline="true" :model="queryParams" class="demo-form-inline">
        <el-form-item label="商品分类">
          <el-tree-select
            v-model="queryParams.categoryId" 
            placeholder="请选择商品分类" 
            clearable
            class="filter-select"
            :data="categoryOptions"
            :props="{
              label: 'name',
              value: 'id',
              children: 'children'
            }"
            default-expand-all
            check-strictly
            node-key="id"
          />
        </el-form-item>
        <el-form-item label="商品标签">
          <el-select 
            v-model="queryParams.tagId" 
            placeholder="选择商品标签" 
            clearable
            class="filter-select"
          >
            <el-option
              v-for="item in tagOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="商品搜索">
          <el-input 
            v-model="queryParams.name" 
            placeholder="请输入商品名称,关键字,ID"
            class="filter-input"
          />
        </el-form-item>
        <el-form-item label="商品类型">
          <el-select v-model="queryParams.type" placeholder="请选择商品类型" clearable>
            <el-option label="全部" value="" />
            <el-option label="普通商品" :value="1" />
            <el-option label="虚拟商品" :value="2" />
            
          </el-select>
        </el-form-item>
        <el-form-item label="商品状态">
          <el-select v-model="queryParams.status" placeholder="商品状态" clearable>
            <el-option label="全部" value="" />
            <el-option label="上架" value="1" />
            <el-option label="下架" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
        </el-form-item>
      </el-form>
    </el-card>

 

    <!-- 操作按钮 -->
    <div class="table-operations">
      <el-button type="primary" @click="handleAdd">添加商品</el-button>
      <el-button type="primary" @click="handleAddCombo">新增套餐</el-button>
      <el-button 
        type="primary" 
        icon="Position" 
        :disabled="selectedRows.length === 0"
        @click="handleBatchMove"
      >
        批量移动
      </el-button>
      <el-button 
        type="danger" 
        icon="Delete"
        :disabled="selectedRows.length === 0"
        @click="handleBatchDelete"
      >
        批量删除
      </el-button>
    </div>

    <!-- 表格 -->
    <el-table
      v-loading="loading"
      :data="goodsList"
      style="width: 100%"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column label="商品ID" prop="goodsId" width="80" />
      <el-table-column label="商品图片" width="100">
        <template #default="{ row }">
          <el-image
            style="width: 50px; height: 50px"
            :src="row.coverImg"
            fit="cover"
          >
            <template #error>
              <div class="image-slot">
                <el-icon><Picture /></el-icon>
              </div>
            </template>
          </el-image>
        </template>
      </el-table-column>
      <el-table-column label="商品名称" prop="name" min-width="200" show-overflow-tooltip />
      <el-table-column label="商品分类" prop="categoryName" width="120" show-overflow-tooltip />
      <el-table-column label="商品类型" width="100">
        <template #default="{ row }">
          <el-tag 
            :type="row.type === 1 ? 'primary' : row.type === 2 ? 'warning' : 'success'"
            size="small"
          >
            {{ row.type === 1 ? '普通商品' : '虚拟商品' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="品牌" prop="brand" width="120" show-overflow-tooltip />
      <el-table-column label="标签" width="120" show-overflow-tooltip>
        <template #default="{ row }">
          <el-tag 
            v-for="(tag, index) in row.tagNames" 
            :key="index"
            size="small"
            style="margin-right: 4px"
          >
            {{ tag }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="价格范围" width="120">
        <template #default="{ row }">
          <span class="price">
            <template v-if="row.minPrice === row.maxPrice">
              ¥{{ row.minPrice }}
            </template>
            <template v-else>
              ¥{{ row.minPrice }} - ¥{{ row.maxPrice }}
            </template>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="销量" prop="sales" width="80" />
      <el-table-column label="总库存" prop="totalStock" width="80" />
      <el-table-column label="排序" prop="sort" width="80" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-switch
            v-model="row.status"
            :active-value="1"
            :inactive-value="0"
            @change="(val) => handleStatusToggle(row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
          <el-button link type="primary" @click="handleCopy(row)">复制</el-button>
          <el-button 
            v-if="row.type === 3" 
            link 
            type="success" 
            @click="handleViewCombo(row)"
          >
            查看套餐
          </el-button>
          <el-button 
            link 
            type="danger" 
            @click="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
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

    <!-- 使用商品表单组件 -->
    <goods-form
      v-model:visible="dialogVisible"
      :type="dialogType"
      :data="formData"
      :category-options="categoryOptions"
      :tag-options="tagOptions"
      :upload-url="uploadUrl"
      @success="handleFormSuccess"
    />

    <!-- 批量移动对话框 -->
    <el-dialog
      v-model="moveDialogVisible"
      title="批量移动商品"
      width="500px"
    >
      <el-form :model="moveForm" label-width="100px">
        <el-form-item label="目标分类">
          <el-cascader
            v-model="moveForm.categoryId"
            :options="categoryOptions"
            :props="{ 
              checkStrictly: true,
              label: 'name',
              value: 'id',
              emitPath: false
            }"
            placeholder="请选择分类"
            clearable
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="moveDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmMove">确定</el-button>
      </template>
    </el-dialog>

    <!-- 套餐配置对话框 -->
    <el-dialog title="套餐配置" v-model="comboDialogVisible" width="80%" :before-close="handleComboDialogClose">
      <div v-if="comboConfig">
        <div class="combo-header">
          <h3>{{ comboConfig.name }}</h3>
          <p>{{ comboConfig.description }}</p>
          <p>基础价格: ¥{{ comboConfig.price }}</p>
        </div>
        
        <div v-for="group in comboConfig.groups" :key="group.id" class="combo-group">
          <div class="group-header">
            <h4>{{ group.name }}</h4>
            <span class="group-desc">{{ group.description }}</span>
            <el-tag v-if="group.isRequired" type="danger" size="mini">必选</el-tag>
            <el-tag v-else type="info" size="mini">可选</el-tag>
            <span class="select-rule">
              ({{ group.minSelect }}-{{ group.maxSelect }}选)
            </span>
          </div>
          
          <div class="group-options">
            <div v-for="option in group.options" :key="option.id" class="option-item">
              <div class="option-info">
                <span class="option-name">{{ option.name }}</span>
                <span class="option-price">
                  +¥{{ option.addPrice }}
                  <span v-if="option.originalPrice">(原价: ¥{{ option.originalPrice }})</span>
                </span>
                <span class="option-quantity">数量: {{ option.quantity }}</span>
                <span class="option-stock">库存: {{ option.stock }}</span>
                <el-tag v-if="option.isDefault" type="success" size="mini">默认</el-tag>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { 
  getGoodsList, 
  deleteGoods, 
  updateGoodsStatus, 
  getGoodsInfo,
  updateGoods,
  copyGoods,
  batchMoveGoods,
  batchDeleteGoods,
  getComboConfig
} from '@/api/shop/goods'
import { getCategoryList } from '@/api/shop/category'
import { getTagList } from '@/api/shop/tag'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown, Picture } from '@element-plus/icons-vue'
import GoodsForm from './components/GoodsForm.vue'

// 查询参数
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  name: '',
  categoryId: '',
  tagId: '',
  type: '',
  status: ''
})

const loading = ref(false)
const total = ref(0)
const goodsList = ref([])
const activeStatus = ref('')
const selectedRows = ref([])

// 添加下拉选项数据
const categoryOptions = ref([])
const tagOptions = ref([])

// 添加上传地址配置
const uploadUrl = ref((import.meta.env.VITE_API_BASE_URL || '/api') + '/common/upload') // 根据实际的上传接口配置

// 获取所有下拉选项数据
const getOptions = async () => {
  try {
    console.log('开始获取分类和标签数据...')
    
    const [categoryRes, tagRes] = await Promise.all([
      getCategoryList({ module: 'goods' }),
      getTagList()
    ])
    
    console.log('获取到的分类数据:', categoryRes)
    console.log('获取到的标签数据:', tagRes)
    
    // 分类数据已经是树形结构，直接使用
    categoryOptions.value = categoryRes || []
    tagOptions.value = tagRes || []
    
    console.log('设置后的分类选项:', categoryOptions.value)
    console.log('设置后的标签选项:', tagOptions.value)
  } catch (error) {
    console.error('获取选项数据失败:', error)
    categoryOptions.value = []
    tagOptions.value = []
  }
}

// 商品列表数据结构调整
const formatGoodsData = (list) => {
  if (!list || !Array.isArray(list)) {
    console.warn('商品数据格式不正确:', list)
    return []
  }
  
  return list.map(item => {
    return {
      ...item,
      goodsId: item.goodsId,
      goodsName: item.name,
      categoryId: item.categoryId,
      coverImg: item.coverImg || '', // 从主图中获取
      price: item.price,
      originalPrice: item.originalPrice,
      sales: item.sales,
      stock: item.stock,
      sort: item.sort,
      status: Number(item.status), // 直接转换为数字
      type: item.type, // 保持数字类型
      source: '自营',
      description: item.description || '', // 确保有值
      tagNames: item.tagNames || [],
      categoryName: item.categoryName,
      tagNames: item.tagNames || [],
      minPrice: item.minPrice,
      maxPrice: item.maxPrice,
      brand: item.brand,
      totalStock: item.totalStock
    }
  })
}

// 获取商品列表
const getList = async () => {
  loading.value = true
  try {
    // 构建查询参数
    const params = {
      pageNum: queryParams.value.pageNum,
      pageSize: queryParams.value.pageSize,
      name: queryParams.value.name,
      categoryId: queryParams.value.categoryId || undefined,
      tagId: queryParams.value.tagId || undefined,
      type: queryParams.value.type || undefined,
      status: queryParams.value.status || undefined
    }
    
    // 移除值为undefined的参数
    Object.keys(params).forEach(key => {
      if (params[key] === undefined || params[key] === '') {
        delete params[key]
      }
    })

    console.log('发送查询参数:', params) // 调试日志
    const res = await getGoodsList(params)
    console.log('商品列表接口返回数据:', res)
    
    // 确保清空旧数据
    goodsList.value = []
    
    // 处理返回的数据
    if (res && res.list) {
      goodsList.value = formatGoodsData(res.list).map(item => ({
        ...item,
        _isInitializing: true // 标记为初始化状态
      }))
      total.value = res.total
    } else {
      console.warn('返回数据格式不正确:', res)
    }
  } catch (error) {
    console.error('获取商品列表失败:', error)
    goodsList.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 查询
const handleQuery = () => {
  queryParams.value.pageNum = 1
  getList()
}

// 状态切换
const handleStatusChange = (value) => {
  queryParams.value.status = value
  handleQuery()
}

// 状态开关切换
const handleStatusToggle = async (row) => {
  const newStatus = Number(row.status)
  console.log('切换状态:', {
    goodsId: row.goodsId,
    status: newStatus,
    rawStatus: row.status,
    typeof: typeof newStatus
  })
  try {
    await updateGoodsStatus({
      goodsId: row.goodsId,
      status: newStatus || 0  // 确保不会发送 undefined 或 null
    })
    ElMessage.success('修改成功')
  } catch (error) {
    console.log('错误详情:', error.response?.data)
    console.error('修改状态失败:', error)
    // 切换失败时恢复原状态
    row.status = Number(!newStatus)  // 使用取反操作
    ElMessage.error('修改失败')
  }
}

// 选择行
const handleSelectionChange = (rows) => {
  selectedRows.value = rows
}

// 新增
const handleAdd = () => {
  dialogType.value = 'add'
  formData.value = {
    goodsId: undefined,
    goodsName: '',
    categoryId: undefined,
    price: 0,
    originalPrice: 0,
    stock: 0,
    sort: 0,
    status: '0',
    description: '',
    coverImg: '',
    tagIds: []
  }
  dialogVisible.value = true
}

// 编辑
const handleEdit = async (row) => {
  dialogType.value = 'edit'
  try {
    // 通过详情接口获取完整的商品信息
    const res = await getGoodsInfo(row.goodsId)
    console.log('商品详情数据:', res)
    
    formData.value = {
      goodsId: res.goodsId,
      goodsName: res.goodsName,
      categoryId: res.categoryId,
      price: res.price,
      originalPrice: res.originalPrice,
      stock: res.stock,
      sort: res.sort,
      status: res.status,
      description: res.description,
      coverImg: res.coverImg,
      tagIds: res.tagIds || []
    }
    dialogVisible.value = true
  } catch (error) {
    console.error('获取商品详情失败:', error)
    ElMessage.error('获取商品详情失败')
  }
}

// 预览
const handlePreview = (row) => {
  // TODO: 预览商品
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm('确认要删除该商品吗?', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await deleteGoods(row.goodsId)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error('删除失败:', error)
    }
  })
}

// 复制商品
const handleCopy = (row) => {
  ElMessageBox.prompt('请输入复制数量', '复制商品', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /^[1-9]\d*$/,
    inputErrorMessage: '请输入正整数',
    inputValue: '1'
  }).then(async ({ value }) => {
    try {
      await copyGoods(row.goodsId, parseInt(value))
      ElMessage.success('复制成功')
      getList()
    } catch (error) {
      console.error('复制商品失败:', error)
      ElMessage.error('复制失败')
    }
  })
}

// 编辑相关的响应式变量
const dialogVisible = ref(false)
const dialogType = ref('add')
const formData = ref({})

// 表单提交成功
const handleFormSuccess = () => {
  getList() // 刷新列表
}

// 添加处理方法
const handleSizeChange = (val) => {
  queryParams.pageSize = val
  getList()
}

const handleCurrentChange = (val) => {
  queryParams.pageNum = val
  getList()
}

// 批量移动
const moveDialogVisible = ref(false)
const moveForm = ref({
  categoryId: null
})

const handleBatchMove = () => {
  moveForm.value.categoryId = null
  moveDialogVisible.value = true
}

const confirmMove = async () => {
  if (!moveForm.value.categoryId) {
    ElMessage.warning('请选择目标分类')
    return
  }
  
  try {
    const goodsIds = selectedRows.value.map(item => item.goodsId)
    await batchMoveGoods({
      goodsIds,
      categoryId: moveForm.value.categoryId
    })
    ElMessage.success('移动成功')
    moveDialogVisible.value = false
    getList()
  } catch (error) {
    console.error('批量移动失败:', error)
    ElMessage.error('移动失败')
  }
}

// 批量删除
const handleBatchDelete = () => {
  ElMessageBox.confirm(`确认要删除选中的 ${selectedRows.value.length} 个商品吗?`, '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      const goodsIds = selectedRows.value.map(item => item.goodsId)
      await batchDeleteGoods(goodsIds)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error('批量删除失败:', error)
      ElMessage.error('删除失败')
    }
  })
}

// 新增套餐
const handleAddCombo = () => {
  // 实现新增套餐的逻辑
}

// 查看套餐配置
const comboDialogVisible = ref(false)
const comboConfig = ref(null)

const handleComboDialogClose = () => {
  comboDialogVisible.value = false
  comboConfig.value = null
}

// 查看套餐配置
const handleViewCombo = async (row) => {
  try {
    const res = await getComboConfig(row.goodsId)
    comboConfig.value = res
    comboDialogVisible.value = true
  } catch (error) {
    console.error('获取套餐配置失败:', error)
    ElMessage.error('获取套餐配置失败')
  }
}

onMounted(() => {
  getList()
  getOptions() // 获取所有下拉选项数据
})
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;

  .filter-container {
    margin: 20px 0;

    // 设置下拉框和输入框的宽度
    :deep(.filter-select) {
      width: 200px;  // 调整下拉框宽度
    }

    :deep(.filter-input) {
      width: 250px;  // 调整输入框宽度
    }

    // 调整表单项间距
    :deep(.el-form-item) {
      margin-right: 20px;
      margin-bottom: 18px;
    }
  }

  .status-tabs {
    margin-bottom: 20px;
  }

  .table-operations {
    margin-bottom: 20px;
    
    .el-button {
      margin-right: 10px;
    }
  }

  .price {
    color: #f56c6c;
    font-weight: bold;
    &::before {
      content: '¥';
    }
  }

  .image-slot {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
    background: #f5f7fa;
    color: #909399;
    font-size: 20px;
  }

  .combo-header {
    margin-bottom: 20px;
    padding: 15px;
    background-color: #f5f7fa;
    border-radius: 4px;
  }

  .combo-header h3 {
    margin: 0 0 10px 0;
    color: #303133;
  }

  .combo-header p {
    margin: 5px 0;
    color: #606266;
  }

  .combo-group {
    margin-bottom: 25px;
    border: 1px solid #e4e7ed;
    border-radius: 4px;
    overflow: hidden;
  }

  .group-header {
    padding: 12px 15px;
    background-color: #fafafa;
    border-bottom: 1px solid #e4e7ed;
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .group-header h4 {
    margin: 0;
    color: #303133;
  }

  .group-desc {
    color: #909399;
    font-size: 12px;
  }

  .select-rule {
    color: #909399;
    font-size: 12px;
    margin-left: auto;
  }

  .group-options {
    padding: 15px;
  }

  .option-item {
    padding: 10px;
    margin-bottom: 10px;
    border: 1px solid #e4e7ed;
    border-radius: 4px;
    background-color: #fff;
  }

  .option-item:last-child {
    margin-bottom: 0;
  }

  .option-info {
    display: flex;
    align-items: center;
    gap: 15px;
  }

  .option-name {
    font-weight: 500;
    color: #303133;
  }

  .option-price {
    color: #e6a23c;
    font-weight: 500;
  }

  .option-quantity,
  .option-stock {
    color: #909399;
    font-size: 12px;
  }
}
</style> 