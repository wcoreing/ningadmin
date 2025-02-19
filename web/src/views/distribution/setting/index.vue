<template>
  <div class="distribution-setting">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="分销模式" name="mode">
        <el-form :model="modeForm" label-width="120px">
          <!-- 分销启用 -->
          <el-form-item label="分销启用">
            <el-switch v-model="modeForm.enabled" />
            <div class="form-tip">商城分销功能开启关闭</div>
          </el-form-item>

          <!-- 分销层级 -->
          <el-form-item label="分销层级">
            <el-radio-group v-model="modeForm.level">
              <el-radio :label="1">一级分销</el-radio>
              <el-radio :label="2">二级分销</el-radio>
            </el-radio-group>
          </el-form-item>

          <!-- 分销模式 -->
          <el-form-item label="分销模式">
            <el-radio-group v-model="modeForm.mode">
              <el-radio label="specified">指定分销</el-radio>
              <el-radio label="all">人人分销</el-radio>
              <el-radio label="condition">满额分销</el-radio>
            </el-radio-group>
            <div class="form-tip">人人分销"默认每个人都可以分销，"指定分销"仅可后台手动设置推广人，"满额分销"指用户购买商品满足消费金额后自动开启分销权限</div>
          </el-form-item>

          <!-- 分销权限申请开关 -->
          <el-form-item label="分销权限申请">
            <el-switch v-model="modeForm.applyEnabled" />
            <div class="form-tip">在指定分销的模式下，此处开启后，移动端会显示申请分销的按钮，用户可以提交申请，后台审核通过后才会绑定推广关系</div>
          </el-form-item>

          <!-- 绑定用户 -->
          <el-form-item label="绑定用户">
            <el-radio-group v-model="modeForm.userType">
              <el-radio label="all">所有用户</el-radio>
              <el-radio label="new">新用户</el-radio>
            </el-radio-group>
            <div class="form-tip">所有用户"指所有注册用户都可以成为下级，"新用户"指首次进入系统的用户才可绑定推广关系</div>
          </el-form-item>

          <!-- 绑定模式 -->
          <el-form-item label="绑定模式">
            <el-radio-group v-model="modeForm.bindMode">
              <el-radio label="permanent">永久</el-radio>
              <el-radio label="period">有效期</el-radio>
              <el-radio label="temporary">临时</el-radio>
            </el-radio-group>
            <div class="form-tip">永久"一次绑定永久有效，"有效期"绑定后一段时间内有效，"临时"每次进入系统重新绑定推广关系</div>
          </el-form-item>

          <!-- 分销中心 -->
          <el-form-item label="分销中心">
            <el-radio-group v-model="modeForm.centerType">
              <el-radio label="normal">普通版</el-radio>
              <el-radio label="pro">专业版</el-radio>
            </el-radio-group>
            <div class="form-tip">普通版：展示详细分销数据，专业版：增加分销统计页面，帮助分销员快速了解分销情况</div>
          </el-form-item>

          <!-- 分销海报图 -->
          <el-form-item label="分销海报图">
            <el-upload
              action="/api/files/upload"
              list-type="picture-card"
              :limit="3"
              :before-upload="beforeUpload"
              :on-success="handleUploadSuccess"
            >
              <el-icon><Plus /></el-icon>
            </el-upload>
            <div class="form-tip">个人中心分销海报图片，建议尺寸：600 * 1000px</div>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="返佣设置" name="commission">
        <el-form :model="commissionForm" label-width="120px">
          <!-- 佣金计算方式 -->
          <el-form-item label="佣金计算方式">
            <el-radio-group v-model="commissionForm.calcType">
              <el-radio label="goods">商品售价</el-radio>
              <el-radio label="payment">支付金额</el-radio>
              <el-radio label="profit">商品利润</el-radio>
            </el-radio-group>
            <div class="form-tip">商品售价：订单商品售价*佣金比例，支付金额：订单实际支付金额*佣金比例，商品利润：(商品实付金额-商品成本价)*佣金比例</div>
          </el-form-item>

          <!-- 一级返佣比例 -->
          <el-form-item label="一级返佣比例">
            <el-input-number v-model="commissionForm.level1Rate" :min="0" :max="100" :precision="2" />
            <span class="unit">%</span>
            <div class="form-tip">订单交易成功后给上级返佣的比例0-100,例:5 = 返订单商品金额的5%</div>
          </el-form-item>

          <!-- 二级返佣比例 -->
          <el-form-item label="二级返佣比例">
            <el-input-number v-model="commissionForm.level2Rate" :min="0" :max="100" :precision="2" />
            <span class="unit">%</span>
            <div class="form-tip">订单交易成功后给上级返佣的比例0-100,例:5 = 返订单商品金额的5%</div>
          </el-form-item>

          <!-- 冻结时间 -->
          <el-form-item label="冻结时间">
            <el-input-number v-model="commissionForm.freezeDays" :min="0" />
            <span class="unit">天</span>
            <div class="form-tip">防止用户退款，佣金被提现，所以需要设置冻结时间(天)</div>
          </el-form-item>

          <!-- 自购返佣 -->
          <el-form-item label="自购返佣">
            <el-switch v-model="commissionForm.selfBuyEnabled" />
            <div class="form-tip">是否开启自购返佣（开启：分销员自己购买商品，享受一级返佣，上级享受二级返佣，上上级不享受返佣；关闭：分销员自己购买商品没有返佣，上级享受一级返佣，上上级正常享受二级返佣）</div>
          </el-form-item>

          <!-- 推广用户返佣 -->
          <el-form-item label="推广用户返佣">
            <el-switch v-model="commissionForm.promotionEnabled" />
            <div class="form-tip">分销推广用户获取佣金</div>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="提现设置" name="withdraw">
        <el-form :model="withdrawForm" label-width="120px">
          <!-- 提现最高金额 -->
          <el-form-item label="提现最高金额">
            <el-input-number v-model="withdrawForm.maxWithdraw" :min="0" />
            <span class="unit">元</span>
            <div class="form-tip">用户提现最高金额限制</div>
          </el-form-item>

          <!-- 提现最低金额 -->
          <el-form-item label="提现最低金额">
            <el-input-number v-model="withdrawForm.minWithdraw" :min="0" />
            <span class="unit">元</span>
            <div class="form-tip">用户提现最低金额限制</div>
          </el-form-item>

          <!-- 提现手续费 -->
          <el-form-item label="提现手续费">
            <el-input-number v-model="withdrawForm.withdrawFeeRate" :min="0" :max="100" :precision="2" />
            <span class="unit">%</span>
            <div class="form-tip">提现金额百分比(0~100)%，例如：提现100，手续费20%，实际到账80</div>
          </el-form-item>

          <!-- 银行卡提现 -->
          <el-form-item label="银行卡提现">
            <el-switch v-model="withdrawForm.bankWithdrawEnabled" />
            <div class="form-tip">银行卡提现</div>
          </el-form-item>

          <!-- 微信提现 -->
          <el-form-item label="微信提现">
            <el-switch v-model="withdrawForm.wechatWithdrawEnabled" />
            <div class="form-tip">微信提现开关</div>
          </el-form-item>

          <!-- 微信到账方式 -->
          <el-form-item label="微信到账方式">
            <el-radio-group v-model="withdrawForm.wechatPayMode">
              <el-radio label="manual">线下手动转账</el-radio>
              <el-radio label="auto">自动到账零钱</el-radio>
            </el-radio-group>
            <div class="form-tip">微信提现方式：手动线下转账，自动转账到微信零钱(需开通商家支付转账API及证书)</div>
          </el-form-item>

          <!-- 支付宝提现 -->
          <el-form-item label="支付宝提现">
            <el-switch v-model="withdrawForm.alipayWithdrawEnabled" />
            <div class="form-tip">支付宝提现开关</div>
          </el-form-item>

          <!-- 提现至余额 -->
          <el-form-item label="提现至余额">
            <el-switch v-model="withdrawForm.balanceWithdrawEnabled" />
            <div class="form-tip">佣金转入用户余额</div>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>

    <div class="form-footer">
      <el-button type="primary" @click="handleSubmit">保存设置</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getDistributionSetting, updateDistributionSetting } from '@/api/distribution/setting'

const activeTab = ref('mode')

// 分销模式表单数据
const modeForm = reactive({
  enabled: false,
  level: 1,
  mode: 'specified',
  applyEnabled: false,
  userType: 'all',
  bindMode: 'permanent',
  centerType: 'normal',
  posterImages: ''
})

// 返佣设置表单数据
const commissionForm = reactive({
  calcType: 'goods',
  level1Rate: 20.65,
  level2Rate: 5,
  freezeDays: 0,
  selfBuyEnabled: false,
  promotionEnabled: false
})

// 提现设置表单数据
const withdrawForm = reactive({
  maxWithdraw: 500,
  minWithdraw: 1,
  withdrawFeeRate: 20,
  bankWithdrawEnabled: false,
  wechatWithdrawEnabled: true,
  wechatPayMode: 'manual',
  alipayWithdrawEnabled: false,
  balanceWithdrawEnabled: false
})

// 获取设置
const getSetting = async () => {
  try {
    const  data = await getDistributionSetting()
    // 更新表单数据
    Object.assign(modeForm, {
      enabled: data.enabled,
      level: data.level,
      mode: data.mode,
      applyEnabled: data.applyEnabled,
      userType: data.userType,
      bindMode: data.bindMode,
      centerType: data.centerType,
      posterImages: data.posterImages
    })

    Object.assign(commissionForm, {
      calcType: data.calcType,
      level1Rate: data.level1Rate,
      level2Rate: data.level2Rate,
      freezeDays: data.freezeDays,
      selfBuyEnabled: data.selfBuyEnabled,
      promotionEnabled: data.promotionEnabled
    })

    Object.assign(withdrawForm, {
      maxWithdraw: data.maxWithdraw,
      minWithdraw: data.minWithdraw,
      withdrawFeeRate: data.withdrawFeeRate,
      bankWithdrawEnabled: data.bankWithdrawEnabled,
      wechatWithdrawEnabled: data.wechatWithdrawEnabled,
      wechatPayMode: data.wechatPayMode,
      alipayWithdrawEnabled: data.alipayWithdrawEnabled,
      balanceWithdrawEnabled: data.balanceWithdrawEnabled
    })
  } catch (error) {
    ElMessage.error('获取设置失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    const data = {
      ...modeForm,
      ...commissionForm,
      ...withdrawForm
    }
    await updateDistributionSetting(data)
    ElMessage.success('保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 图片上传相关
const handleUploadSuccess = (response, file) => {
  if (response.code === 0) {
    const url = response.data.url
    modeForm.posterImages = url
  } else {
    ElMessage.error('上传失败')
  }
}

const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
    return false
  }
  return true
}

// 页面加载时获取设置
onMounted(() => {
  getSetting()
})
</script>

<style scoped>
.distribution-setting {
  padding: 20px;
  background-color: #fff;
  border-radius: 4px;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.unit {
  margin-left: 8px;
  color: #666;
}

.form-footer {
  margin-top: 20px;
  text-align: center;
}

:deep(.el-form-item) {
  margin-bottom: 24px;
}

:deep(.el-upload--picture-card) {
  width: 120px;
  height: 120px;
  line-height: 120px;
}
</style> 