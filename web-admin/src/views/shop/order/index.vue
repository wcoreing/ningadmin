<template>
  <div class="app-container">
    <!-- 搜索区域 -->
    <el-card class="filter-container">
      <el-form :inline="true" :model="queryParams" class="demo-form-inline">
        <el-form-item label="订单编号">
          <el-input v-model="queryParams.orderNo" placeholder="请输入订单编号" clearable />
        </el-form-item>
        <el-form-item label="订单状态" style="width: 200px;">
          <el-select v-model="queryParams.status" placeholder="请选择订单状态" clearable>
            <el-option label="待付款" value="unpaid" />
            <el-option label="待发货" value="unshipped" />
            <el-option label="待收货" value="delivered" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>
        <el-form-item label="用餐类型" style="width: 200px;">
          <el-select v-model="queryParams.diningType" placeholder="请选择用餐类型" clearable>
            <el-option label="堂食" value="dine_in" />
            <el-option label="外卖" value="takeout" />
          </el-select>
        </el-form-item>
        <el-form-item label="下单时间">
          <el-date-picker
            v-model="queryParams.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 订单列表 -->
    <el-card>
      <div class="table-operations">
        <el-button type="primary" icon="Plus">新增订单</el-button>
        <el-button type="success" icon="Download">导出订单</el-button>
      </div>

      <el-table v-loading="loading" :data="orderList" border>
        <el-table-column label="订单信息" min-width="180">
          <template #default="{ row }">
            <div class="order-info">
              <div class="order-no">{{ row.orderNo }}</div>
              <div class="order-time">{{ row.createTime }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="商品信息" min-width="380">
          <template #default="{ row }">
            <div v-for="(item, index) in row.items" :key="index" class="goods-item">
              <el-image 
                :src="item.goodsImage"
                fit="contain"
                style="width: 60px; height: 60px"
                :preview-src-list="[item.goodsImage]"
                class="goods-image"
              />
              <div class="goods-info">
                <div class="goods-main">
                  <div class="goods-name">{{ item.goodsName }}</div>
                  <div class="goods-spec">
                    <template v-if="item.specText">
                      <span v-for="(spec, key) in JSON.parse(item.specText)" :key="key">
                        {{ key }}：{{ spec.value }}
                        <el-divider direction="vertical" v-if="key < JSON.parse(item.specText).length - 1" />
                      </span>
                    </template>
                  </div>
                </div>
                <div class="goods-price">
                  <span class="price">¥{{ item.price }}</span>
                  <span class="qty">x{{ item.quantity }}</span>
                </div>
              </div>
            </div>
            <div v-if="row.items.length > 1" class="goods-total">
              共 {{ row.items.length }} 件商品
            </div>
          </template>
        </el-table-column>
        <el-table-column label="实付金额" width="120" align="center">
          <template #default="{ row }">
            <div class="amount-info">
              <div class="pay-amount">¥{{ row.payAmount }}</div>
              <div class="total-amount">总额: ¥{{ row.totalAmount }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="订单状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
            <div v-if="row.status === 'delivered'" class="logistics-info">
              <el-tooltip effect="dark" placement="top">
                <template #content>
                  <div>物流公司：{{ row.receiver.logisticsCompany }}</div>
                  <div>物流单号：{{ row.receiver.logisticsNo }}</div>
                </template>
                <el-tag size="small" type="info" style="margin-top: 4px">
                  已发货
                </el-tag>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="用餐类型" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="getDiningTypeColor(row.diningType)">
              {{ getDiningTypeText(row.diningType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="支付方式" width="100" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="row.payType ? 'success' : ''">
              {{ getPayTypeText(row.payType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="收货人信息" min-width="180">
          <template #default="{ row }">
            <div class="receiver-info">
              <div class="receiver-main">
                <span class="name">{{ row.receiver.receiver }}</span>
                <span class="phone">{{ row.receiver.phone }}</span>
              </div>
              <div class="address">{{ row.receiver.address }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleDetail(row)">详情</el-button>
            <el-button 
              v-if="['delivered', 'completed'].includes(row.status)"
              link 
              type="primary" 
              @click="handleViewLogistics(row)"
            >
              物流
            </el-button>
            <el-button 
              v-if="row.status === 'unshipped'" 
              link 
              type="primary" 
              @click="handleShip(row)"
            >
              发货
            </el-button>
            <el-button 
              v-if="row.status === 'unpaid'" 
              link 
              type="danger" 
              @click="handleCancel(row)"
            >
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.pageNum"
          v-model:page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 30, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 订单详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      title="订单详情"
      width="800px"
      destroy-on-close
    >
      <OrderDetail :order="currentOrder" />
    </el-dialog>

    <!-- 发货弹窗 -->
    <el-dialog
      v-model="shipVisible"
      title="订单发货"
      width="500px"
      destroy-on-close
    >
      <el-form 
        ref="shipFormRef"
        :model="shipForm" 
        :rules="shipRules"
        label-width="100px"
      >
        <el-form-item label="物流公司">
          <el-select v-model="shipForm.expressCompany" placeholder="请选择物流公司">
            <el-option 
              v-for="item in expressCompanies" 
              :key="item.code" 
              :label="item.name" 
              :value="item.code" 
            />
          </el-select>
        </el-form-item>
        <el-form-item label="物流单号">
          <el-input v-model="shipForm.expressNo" placeholder="请输入物流单号" />
        </el-form-item>
        <el-form-item label="发货备注">
          <el-input 
            v-model="shipForm.remark" 
            type="textarea" 
            :rows="3"
            placeholder="请输入发货备注信息（选填）" 
          />
        </el-form-item>
      </el-form>
      <div class="ship-info">
        <div class="info-item">
          <span class="label">收货人：</span>
          <span>{{ currentOrder?.receiver.receiver }}</span>
        </div>
        <div class="info-item">
          <span class="label">联系电话：</span>
          <span>{{ currentOrder?.receiver.phone }}</span>
        </div>
        <div class="info-item">
          <span class="label">收货地址：</span>
          <span>{{ currentOrder?.receiver.address }}</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="shipVisible = false">取 消</el-button>
        <el-button type="primary" @click="confirmShip">确 定</el-button>
      </template>
    </el-dialog>

    <!-- 物流信息弹窗 -->
    <el-dialog
      v-model="logisticsVisible"
      title="物流信息"
      width="600px"
    >
      <div v-loading="logisticsLoading" class="logistics-track">
        <div class="track-header">
          <div class="company">物流公司：{{ logisticsInfo.logisticsCompany }}</div>
          <div class="number">物流单号：{{ logisticsInfo.logisticsNo }}</div>
          <div class="status">物流状态：{{ logisticsInfo.status }}</div>
        </div>
        <el-timeline>
          <el-timeline-item
            v-for="(trace, index) in logisticsInfo.traces"
            :key="index"
            :timestamp="trace.acceptTime"
          >
            {{ trace.acceptStation }}
            <div v-if="trace.remark" class="remark">{{ trace.remark }}</div>
          </el-timeline-item>
        </el-timeline>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getOrderList, shipOrder, cancelOrder, getLogisticsTrack } from '@/api/shop/order'
import OrderDetail from './components/OrderDetail.vue'

// 查询参数
const queryParams = reactive({
  orderNo: '',
  status: '',
  diningType: '',
  dateRange: [],
  pageNum: 1,
  pageSize: 10
})

const orderList = ref([])
const loading = ref(false)
const total = ref(0)
const detailVisible = ref(false)
const shipVisible = ref(false)
const currentOrder = ref(null)

// 物流公司列表
const expressCompanies = [
  { code: 'SF', name: '顺丰速运' },
  { code: 'ZTO', name: '中通快递' },
  { code: 'YTO', name: '圆通速递' },
  { code: 'YD', name: '韵达快递' },
  { code: 'STO', name: '申通快递' },
  { code: 'EMS', name: 'EMS' },
  { code: 'HTKY', name: '百世快递' },
  { code: 'JD', name: '京东物流' }
]

const shipFormRef = ref(null)
const shipForm = reactive({
  orderId: null,
  expressCompany: '',
  expressNo: '',
  remark: ''
})

const shipRules = {
  expressCompany: [
    { required: true, message: '请选择物流公司', trigger: 'change' }
  ],
  expressNo: [
    { required: true, message: '请输入物流单号', trigger: 'blur' },
    { min: 5, message: '物流单号长度不能小于5位', trigger: 'blur' }
  ]
}

// 物流信息相关
const logisticsVisible = ref(false)
const logisticsLoading = ref(false)
const logisticsInfo = ref({})

// 获取状态类型
const getStatusType = (status) => {
  const map = {
    'unpaid': 'warning',      // 待付款
    'unshipped': 'primary',   // 待发货
    'delivered': 'info',      // 待收货
    'completed': 'success',   // 已完成
    'cancelled': 'danger'     // 已取消
  }
  return map[status]
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    'unpaid': '待付款',
    'unshipped': '待发货',
    'delivered': '待收货',
    'completed': '已完成',
    'cancelled': '已取消'
  }
  return map[status]
}

// 获取支付方式文本
const getPayTypeText = (payType) => {
  const map = {
    1: '微信支付',
    2: '支付宝',
    3: '余额支付'
  }
  return map[payType] || '未支付'
}

// 获取用餐类型文本
const getDiningTypeText = (diningType) => {
  const map = {
    'dine_in': '堂食',
    'takeout': '外卖'
  }
  return map[diningType] || '外卖'
}

// 获取用餐类型颜色
const getDiningTypeColor = (diningType) => {
  const map = {
    'dine_in': 'warning',  // 堂食用橙色
    'takeout': 'success'   // 外卖用绿色
  }
  return map[diningType] || 'success'
}

// 获取物流公司名称
const getExpressCompanyName = (code) => {
  const company = expressCompanies.find(item => item.code === code)
  return company ? company.name : code
}

// 获取订单列表
const getList = async () => {
  try {
    loading.value = true
    const params = {
      ...queryParams,
      startTime: queryParams.dateRange?.[0],
      endTime: queryParams.dateRange?.[1]
    }
    const data = await getOrderList(params)
    orderList.value = data.list
    total.value = data.total
  } catch (error) {
    console.error('获取订单列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 查询
const handleQuery = () => {
  queryParams.pageNum = 1
  getList()
}

// 重置
const resetQuery = () => {
  queryParams.orderNo = ''
  queryParams.status = ''
  queryParams.diningType = ''
  queryParams.dateRange = []
  handleQuery()
}

// 查看详情
const handleDetail = (row) => {
  currentOrder.value = row
  detailVisible.value = true
}

// 发货
const handleShip = (row) => {
  currentOrder.value = row
  shipForm.orderId = row.id
  shipVisible.value = true
}

// 确认发货
const confirmShip = () => {
  if (!shipFormRef.value) return
  
  shipFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await shipOrder({
          orderId: shipForm.orderId,
          logisticsCompany: shipForm.expressCompany,
          logisticsNo: shipForm.expressNo,
          remark: shipForm.remark
        })
        ElMessage.success('发货成功')
        shipVisible.value = false
        getList()
      } catch (error) {
        console.error('发货失败:', error)
      }
    }
  })
}

// 取消订单
const handleCancel = (row) => {
  ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
    type: 'warning'
  }).then(() => {
    cancelOrder(row.id).then(() => {
      ElMessage.success('订单已取消')
      getList()
    })
  })
}

// 查看物流
const handleViewLogistics = async (row) => {
  logisticsVisible.value = true
  logisticsLoading.value = true
  try {
    const { data } = await getLogisticsTrack({ orderId: row.id })
    logisticsInfo.value = data
  } catch (error) {
    console.error('获取物流信息失败:', error)
  } finally {
    logisticsLoading.value = false
  }
}

// 页面加载时获取列表
onMounted(() => getList())

// 分页方法
const handleSizeChange = (val) => {
  queryParams.pageSize = val
  getList()
}

const handleCurrentChange = (val) => {
  queryParams.pageNum = val
  getList()
}
</script>

<style lang="scss" scoped>
// 文本省略混入
@mixin ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.filter-container {
  margin-bottom: 20px;
}

.price {
  color: #f56c6c;
  font-weight: bold;
}

.address-text {
  color: #909399;
  font-size: 13px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.order-detail {
  .detail-section {
    margin-bottom: 24px;

    h3 {
      margin: 0 0 16px;
      font-size: 16px;
      font-weight: 500;
    }

    .logistics-detail {
      div {
        line-height: 24px;
      }
    }
  }
}

.table-operations {
  margin-bottom: 16px;
  display: flex;
  gap: 8px;
}

.order-info {
  .order-no {
    font-weight: 500;
    margin-bottom: 4px;
  }
  .order-time {
    font-size: 12px;
    color: #909399;
  }
}

.goods-item {
  display: flex;
  padding: 4px 0;
  border-bottom: 1px solid #f0f0f0;
  align-items: center;
  
  &:last-child {
    border-bottom: none;
  }

  .goods-image {
    margin-right: 12px;
    border: 1px solid #eee;
    border-radius: 4px;
  }

  .goods-info {
    flex: 1;
    min-width: 0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
    
    .goods-main {
      flex: 1;
      min-width: 0;
      
      .goods-name {
        font-size: 14px;
        font-weight: 500;
        margin-bottom: 4px;
        @include ellipsis;
      }
      
      .goods-spec {
        font-size: 12px;
        color: #909399;
        @include ellipsis;
      }
    }
    
    .goods-price {
      flex-shrink: 0;
      text-align: right;
      
      .price {
        display: block;
        font-size: 14px;
        color: #f56c6c;
        font-weight: 500;
        margin-bottom: 4px;
      }
      
      .qty {
        display: block;
        font-size: 12px;
        color: #909399;
      }
    }
  }
}

.goods-total {
  font-size: 12px;
  color: #999;
  text-align: right;
  margin-top: 8px;
}

.amount-info {
  .pay-amount {
    font-size: 15px;
    color: #f56c6c;
    font-weight: 500;
    margin-bottom: 4px;
  }
  .total-amount {
    font-size: 12px;
    color: #909399;
  }
}

.receiver-info {
  .receiver-main {
    margin-bottom: 4px;
    
    .name {
      font-weight: 500;
      margin-right: 8px;
    }
    
    .phone {
      color: #606266;
    }
  }
  
  .address {
    font-size: 12px;
    color: #909399;
    @include ellipsis;
  }
}

.ship-info {
  margin-top: 20px;
  padding: 16px;
  background-color: #f5f7fa;
  border-radius: 4px;

  .info-item {
    line-height: 24px;
    
    .label {
      color: #606266;
      margin-right: 8px;
    }
  }
}

.logistics-detail {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.logistics-track {
  .track-header {
    margin-bottom: 20px;
    padding: 16px;
    background-color: #f5f7fa;
    border-radius: 4px;
    
    > div {
      line-height: 24px;
    }
    
    .status {
      color: #409eff;
      font-weight: 500;
    }
  }
  
  .remark {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
  }
}
</style> 