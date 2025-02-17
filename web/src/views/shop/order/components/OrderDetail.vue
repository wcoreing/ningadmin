<template>
  <div class="order-detail" v-if="order">
    <div class="detail-section">
      <h3>订单信息</h3>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="订单编号">{{ order.orderNo }}</el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ order.createTime }}</el-descriptions-item>
        <el-descriptions-item label="订单状态">
          <el-tag :type="getStatusType(order.status)">
            {{ getStatusText(order.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="支付方式">
          <el-tag size="small" :type="order.payType ? 'success' : ''">
            {{ getPayTypeText(order.payType) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item v-if="order.status === 'shipped'" label="发货时间">
          {{ order.shipTime }}
        </el-descriptions-item>
        <el-descriptions-item v-if="order.status === 'completed'" label="完成时间">
          {{ order.finishTime }}
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <div class="detail-section">
      <h3>收货信息</h3>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="收货人">{{ order.receiver.receiver }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ order.receiver.phone }}</el-descriptions-item>
        <el-descriptions-item label="收货地址">{{ order.receiver.address }}</el-descriptions-item>
        <el-descriptions-item v-if="order.status === 'shipped' || order.status === 'completed'" label="物流信息">
          <div class="logistics-detail">
            <div>物流公司：{{ getExpressCompanyName(order.receiver.logisticsCompany) }}</div>
            <div>物流单号：{{ order.receiver.logisticsNo }}</div>
          </div>
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <div class="detail-section">
      <h3>商品信息</h3>
      <el-table :data="order.items" border>
        <el-table-column label="商品图片" width="100">
          <template #default="{ row }">
            <el-image 
              :src="row.goodsImage" 
              fit="contain"
              style="width: 60px; height: 60px"
              :preview-src-list="[row.goodsImage]"
            />
          </template>
        </el-table-column>
        <el-table-column label="商品名称" prop="goodsName" min-width="200" />
        <el-table-column label="规格" width="150">
          <template #default="{ row }">
            <div v-if="row.specText">
              <div v-for="(spec, key) in JSON.parse(row.specText)" :key="key">
                {{ key }}：{{ spec.value }}
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="单价" width="100">
          <template #default="{ row }">
            <span class="price">¥{{ row.price }}</span>
          </template>
        </el-table-column>
        <el-table-column label="数量" prop="quantity" width="100" />
        <el-table-column label="小计" width="120">
          <template #default="{ row }">
            <span class="price">¥{{ row.price * row.quantity }}</span>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="detail-section">
      <h3>费用信息</h3>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="商品总额">
          <span class="price">¥{{ order.totalAmount }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="运费">
          <span class="price">¥{{ order.freightAmount || 0 }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="优惠金额">
          <span class="price">-¥{{ order.discountAmount || 0 }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="实付金额">
          <span class="price">¥{{ order.payAmount }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue'

const props = defineProps({
  order: {
    type: Object,
    required: true
  }
})

// 获取状态类型
const getStatusType = (status) => {
  const map = {
    'pending': 'warning',     // 待付款
    'unshipped': 'primary',   // 待发货
    'shipped': 'info',        // 待收货
    'completed': 'success',   // 已完成
    'cancelled': 'danger'     // 已取消
  }
  return map[status]
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    'pending': '待付款',
    'unshipped': '待发货',
    'shipped': '待收货',
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

// 获取物流公司名称
const getExpressCompanyName = (code) => {
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
  const company = expressCompanies.find(item => item.code === code)
  return company ? company.name : code
}
</script>

<style lang="scss" scoped>
.order-detail {
  .detail-section {
    margin-bottom: 24px;

    h3 {
      margin: 0 0 16px;
      font-size: 16px;
      font-weight: 500;
    }
  }
}

.price {
  color: #f56c6c;
  font-weight: bold;
}

.logistics-detail {
  div {
    line-height: 24px;
  }
}
</style> 