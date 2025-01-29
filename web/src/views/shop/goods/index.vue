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
          <el-select v-model="queryParams.categoryId" placeholder="请选择商品分类" clearable>
            <el-option label="请选择商品分类" value="" />
          </el-select>
        </el-form-item>
        <el-form-item label="供应商">
          <el-select v-model="queryParams.supplierId" placeholder="请选择" clearable>
            <el-option label="请选择" value="" />
          </el-select>
        </el-form-item>
        <el-form-item label="商品品牌">
          <el-select v-model="queryParams.brandId" placeholder="请选择商品品牌" clearable>
            <el-option label="请选择商品品牌" value="" />
          </el-select>
        </el-form-item>
        <el-form-item label="商品标签">
          <el-select v-model="queryParams.tagId" placeholder="选择商品标签" clearable>
            <el-option label="选择商品标签" value="" />
          </el-select>
        </el-form-item>
        <el-form-item label="商品搜索">
          <el-input v-model="queryParams.goodsName" placeholder="请输入商品名称,关键字,ID" />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

// 查询参数
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  goodsName: '',
  categoryId: '',
  supplierId: '',
  brandId: '',
  tagId: '',
  status: ''
})

const loading = ref(false)
const total = ref(0)
const goodsList = ref([])
const activeStatus = ref('')
const selectedRows = ref([])

// 获取商品列表
const getList = async () => {
  loading.value = true
  try {
    const res = await getGoodsList(queryParams.value)
    goodsList.value = res.data.list
    total.value = res.data.total
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
})
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;

  .filter-container {
    margin: 20px 0;
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