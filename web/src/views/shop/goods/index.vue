<template>
  <div class="app-container">
    <!-- 温馨提示 -->
    <el-alert
      title="温馨提示: 1.新增商品时可选统一规格或者多规格，满足商品不同销售层次场景；2.商品销售状态分为销售中且库存足够时才可下单购买"
      type="info"
      :closable="true"
      show-icon
    />

    <!-- 搜索区域 -->
    <el-card class="filter-container">
      <el-form :inline="true" :model="queryParams" class="demo-form-inline">
        <el-form-item label="商品分类">
          <el-select 
            v-model="queryParams.categoryId" 
            placeholder="请选择商品分类" 
            clearable
            class="filter-select"
          >
            <el-option
              v-for="item in categoryOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="供应商">
          <el-select 
            v-model="queryParams.supplierId" 
            placeholder="请选择供应商" 
            clearable
            class="filter-select"
          >
            <el-option
              v-for="item in supplierOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="商品品牌">
          <el-select 
            v-model="queryParams.brandId" 
            placeholder="请选择商品品牌" 
            clearable
            class="filter-select"
          >
            <el-option
              v-for="item in brandOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
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
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 状态标签页 -->
    <div class="status-tabs">
      <el-radio-group v-model="activeStatus" @change="handleStatusChange">
        <el-radio-button label="">全部(0)</el-radio-button>
        <el-radio-button label="1">销售中(133)</el-radio-button>
        <el-radio-button label="2">仓库中(9)</el-radio-button>
        <el-radio-button label="3">已售罄(38)</el-radio-button>
        <el-radio-button label="4">库存预警(13)</el-radio-button>
        <el-radio-button label="5">回收站(176)</el-radio-button>
        <el-radio-button label="6">待审核(48)</el-radio-button>
        <el-radio-button label="7">审核未通过(1)</el-radio-button>
        <el-radio-button label="8">强制下架(0)</el-radio-button>
      </el-radio-group>
    </div>

    <!-- 操作按钮 -->
    <div class="table-operations">
      <el-button type="primary" @click="handleAdd">添加商品</el-button>
      <el-button type="success">商品采集</el-button>
      <el-button>批量下架</el-button>
      <el-button>导出</el-button>
      <el-button>批量设置</el-button>
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
      <el-table-column label="商品图" width="80">
        <template #default="{ row }">
          <el-image
            :src="row.coverImg"
            :preview-src-list="[row.coverImg]"
            fit="cover"
            style="width: 50px; height: 50px"
          />
        </template>
      </el-table-column>
      <el-table-column label="商品名称" prop="goodsName" min-width="200" show-overflow-tooltip />
      <el-table-column label="商品来源" prop="source" width="100" />
      <el-table-column label="商品类型" prop="type" width="100" />
      <el-table-column label="商品售价" prop="price" width="100">
        <template #default="{ row }">
          <span class="price">{{ row.price }}</span>
        </template>
      </el-table-column>
      <el-table-column label="销量" prop="sales" width="80" />
      <el-table-column label="库存" prop="stock" width="80" />
      <el-table-column label="排序" prop="sort" width="80" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-switch
            v-model="row.status"
            :active-value="'0'"
            :inactive-value="'1'"
            @change="handleStatusToggle(row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
          <el-button type="primary" link @click="handlePreview(row)">预览</el-button>
          <el-dropdown>
            <el-button type="primary" link>
              更多<el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleDelete(row)">删除</el-dropdown-item>
                <el-dropdown-item>复制</el-dropdown-item>
                <el-dropdown-item>分享</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <pagination
      v-show="total > 0"
      :total="total"
      v-model:page="queryParams.pageNum"
      v-model:limit="queryParams.pageSize"
      @pagination="getList"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getGoodsList, deleteGoods, updateGoodsStatus } from '@/api/shop/goods'
import { getCategoryList } from '@/api/shop/category'
import { getBrandList } from '@/api/shop/brand'
import { getTagList } from '@/api/shop/tag'
import { getSupplierList } from '@/api/shop/supplier'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

// 查询参数
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  name: '',
  categoryId: '',
  brandId: '',
  tagId: '',
  supplierId: '',
  status: ''
})

const loading = ref(false)
const total = ref(0)
const goodsList = ref([])
const activeStatus = ref('')
const selectedRows = ref([])

// 添加下拉选项数据
const categoryOptions = ref([])
const supplierOptions = ref([])
const brandOptions = ref([])
const tagOptions = ref([])

// 获取所有下拉选项数据
const getOptions = async () => {
  try {
    const [categoryRes, supplierRes, brandRes, tagRes] = await Promise.all([
      getCategoryList(),
      getSupplierList(),
      getBrandList(),
      getTagList()
    ])
    
    // 添加调试日志
    console.log('分类数据:', categoryRes)
    console.log('供应商数据:', supplierRes)
    console.log('品牌数据:', brandRes)
    console.log('标签数据:', tagRes)
    
    // 直接使用返回的数据，不需要 .data
    categoryOptions.value = categoryRes || []
    supplierOptions.value = supplierRes || []
    brandOptions.value = brandRes || []
    tagOptions.value = tagRes || []
  } catch (error) {
    console.error('获取选项数据失败:', error)
    // 出错时设置为空数组
    categoryOptions.value = []
    supplierOptions.value = []
    brandOptions.value = []
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
    console.log('格式化单个商品数据:', item) // 添加调试日志
    return {
      goodsId: item.goodsId,
      goodsName: item.name,
      coverImg: '', // 暂时不显示图片
      price: item.price,
      originalPrice: item.originalPrice,
      sales: item.sales,
      stock: item.stock,
      sort: item.sort,
      status: item.status.toString(),
      type: item.type === 1 ? '普通商品' : '虚拟商品',
      source: '自营'
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
      brandId: queryParams.value.brandId || undefined,
      tagId: queryParams.value.tagId || undefined,
      supplierId: queryParams.value.supplierId || undefined,
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
      goodsList.value = formatGoodsData(res.list)
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
  try {
    await updateGoodsStatus({
      id: row.goodsId,
      status: row.status
    })
    ElMessage.success('状态更新成功')
  } catch (error) {
    row.status = row.status === '0' ? '1' : '0'
    ElMessage.error('状态更新失败')
    console.error('状态更新失败:', error)
  }
}

// 选择行
const handleSelectionChange = (rows) => {
  selectedRows.value = rows
}

// 新增
const handleAdd = () => {
  // TODO: 跳转到新增页面
}

// 编辑
const handleEdit = (row) => {
  // TODO: 跳转到编辑页面
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
}
</style> 