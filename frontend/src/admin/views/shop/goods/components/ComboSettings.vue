<template>
  <div class="combo-settings-container">
    <!-- 套餐基础配置 -->
    <div class="combo-basic-config">
      <el-form :model="formData" label-width="100px" size="default">
        <el-form-item label="套餐价格" required>
          <el-input-number
            v-model="comboPrice"
            :min="0"
            :precision="2"
            placeholder="请输入套餐价格"
            style="width: 200px"
          />
          <span class="unit">元</span>
          <div class="form-tip">套餐的价格，最终价格 = 套餐价格 + 选项加价</div>
        </el-form-item>
      </el-form>
    </div>

    <!-- 套餐组设置 -->
    <div class="groups-content-card">
      <div class="card-header">
        <h3 class="card-title">
          <el-icon><Collection /></el-icon>
          套餐组配置
        </h3>
        <el-button type="primary" size="small" @click="addGroup">
          <el-icon><Plus /></el-icon>
          添加套餐组
        </el-button>
      </div>
      
      <div v-if="groups.length > 0" class="groups-container">
     
        <div 
          v-for="(group, groupIndex) in groups" 
          :key="group.id || groupIndex"
          class="group-card"
        >
          <div class="group-header">
            <div class="group-title">
              <el-icon><Folder /></el-icon>
              <span>{{ group.name || `套餐组 ${groupIndex + 1}` }}</span>
              <el-tag v-if="group.options && group.options.length > 0" size="small" type="info">
                {{ group.options.length }} 个选项
              </el-tag>
            </div>
            <el-button
              type="danger"
              size="small"
              plain
              @click="removeGroup(groupIndex)"
              :disabled="groups.length === 1"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
          
          <div class="group-content">
            <!-- 套餐组基本信息 -->
            <div class="group-basic-info">
              <el-form :model="group" label-width="80px" label-position="left" size="small">
                <el-row :gutter="16">
                  <el-col :span="8">
                    <el-form-item label="组名称" required>
                      <el-input
                        v-model="group.name"
                        placeholder="如：主食、饮料、小食"
                        size="small"
                      />
                    </el-form-item>
                  </el-col>
                  
                  <el-col :span="8">
                    <el-form-item label="组描述">
                      <el-input
                        v-model="group.description"
                        placeholder="选择说明"
                        size="small"
                      />
                    </el-form-item>
                  </el-col>
                  
                  <el-col :span="8">
                    <el-form-item label="是否必选">
                      <el-switch
                        v-model="group.isRequired"
                        active-text="必选"
                        inactive-text="可选"
                        size="small"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>
                
                <el-row :gutter="16">
                  <el-col :span="8">
                    <el-form-item label="最少选择">
                      <el-input-number
                        v-model="group.minSelect"
                        :min="0"
                        :max="group.maxSelect || 99"
                        :precision="0"
                        size="small"
                        style="width: 100%"
                      />
                    </el-form-item>
                  </el-col>
                  
                  <el-col :span="8">
                    <el-form-item label="最多选择">
                      <el-input-number
                        v-model="group.maxSelect"
                        :min="group.minSelect || 1"
                        :max="99"
                        :precision="0"
                        size="small"
                        style="width: 100%"
                      />
                    </el-form-item>
                  </el-col>
                  
                  <el-col :span="8">
                    <el-form-item label="排序">
                      <el-input-number
                        v-model="group.sort"
                        :min="0"
                        :precision="0"
                        size="small"
                        style="width: 100%"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>
              </el-form>
            </div>

            <!-- 套餐组选项配置 -->
            <div class="group-options-section">
              <div class="section-header">
                <label class="form-label">组内选项</label>
                <el-button
                  type="primary"
                  size="small"
                  @click="openGoodsSelector(groupIndex)"
                >
                  <el-icon><Plus /></el-icon>
                  添加选项
                </el-button>
              </div>
              
              <div v-if="group.options && group.options.length > 0" class="options-list">
                <div 
                  v-for="(option, optionIndex) in group.options"
                  :key="option.id || optionIndex"
                  class="option-item"
                >
                  <div class="option-info">
                    <el-image
                      :src="option.goodsImage"
                      fit="cover"
                      class="option-image"
                    >
                      <template #error>
                        <div class="image-slot">
                          <el-icon><Picture /></el-icon>
                        </div>
                      </template>
                    </el-image>
                    <div class="option-details">
                      <div class="option-name">{{ option.name }}</div>
                      <div class="option-goods" v-if="option.goodsName">
                        商品：{{ option.goodsName }}
                      </div>
                    </div>
                  </div>
                  
                  <div class="option-config">
                    <div class="config-item">
                      <label>选项值</label>
                      <el-input
                        v-model="option.value"
                        placeholder="选项值"
                        size="small"
                        style="width: 120px"
                      />
                    </div>
                    
                    <div class="config-item">
                      <label>加价</label>
                      <el-input-number
                        v-model="option.addPrice"
                        :precision="2"
                        :step="0.01"
                        :min="0"
                        size="small"
                        style="width: 100px"
                      />
                      <span class="unit">元</span>
                    </div>
                    
                    <div class="config-item">
                      <label>数量</label>
                      <el-input-number
                        v-model="option.quantity"
                        :min="1"
                        :precision="0"
                        size="small"
                        style="width: 80px"
                      />
                    </div>
                    
                    <div class="config-item">
                      <label>默认</label>
                      <el-switch
                        v-model="option.isDefault"
                        @change="handleDefaultChange(groupIndex, optionIndex)"
                      />
                    </div>
                    
                    <div class="config-item">
                      <el-button
                        type="danger"
                        size="small"
                        plain
                        @click="removeOption(groupIndex, optionIndex)"
                      >
                        <el-icon><Delete /></el-icon>
                      </el-button>
                    </div>
                  </div>
                </div>
              </div>
              
              <el-empty 
                v-else 
                description="暂无选项，请点击添加选项按钮"
                :image-size="80"
              />
            </div>
          </div>
        </div>
      </div>
      
      <el-empty 
        v-else 
        description="暂无套餐组，请点击右上角按钮添加"
        :image-size="100"
      />
    </div>

    <!-- 商品选择器 -->
    <goods-selector
      v-model="goodsSelectorVisible"
      title="选择商品"
      :multiple="true"
      @select="handleGoodsSelect"
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Plus, Delete, Picture, Collection, Folder } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import GoodsSelector from '@/admin/components/GoodsSelector/index.vue'

const props = defineProps({
  formData: {
    type: Object,
    required: true
  }
})

// 添加调试信息
console.log('=== ComboSettings 组件初始化 ===')
console.log('ComboSettings - 接收到的formData:', props.formData)
console.log('ComboSettings - formData.allSpecData:', props.formData.allSpecData)
console.log('ComboSettings - formData.allSpecData.comboSpecData:', props.formData.allSpecData?.comboSpecData)

// 响应式数据
const goodsSelectorVisible = ref(false)
const currentGroupIndex = ref(-1)

// 套餐价格计算属性 - 直接操作allSpecData.comboSpecData.price
const comboPrice = computed({
  get() {
    const price = props.formData.allSpecData?.comboSpecData?.price || 0
    console.log('ComboSettings - 读取套餐价格:', price)
    return price
  },
  set(newPrice) {
    console.log('ComboSettings - 设置套餐价格:', newPrice)
    // 确保comboSpecData对象存在
    if (!props.formData.allSpecData) {
      props.formData.allSpecData = {}
    }
    if (!props.formData.allSpecData.comboSpecData) {
      props.formData.allSpecData.comboSpecData = { price: 0, comboGroups: [] }
    }
    // 只更新comboSpecData中的价格
    props.formData.allSpecData.comboSpecData.price = newPrice
  }
})

// ID生成器
let idCounter = 0
const generateId = () => {
  return Date.now() + (++idCounter)
}

// 直接使用formData.comboGroups，不使用本地groups状态
const groups = computed({
  get() {
    if (!props.formData.comboGroups) {
      props.formData.comboGroups = []
    }
    console.log('ComboSettings - 获取套餐组数据:', props.formData.comboGroups)
    return props.formData.comboGroups
  },
  set(newGroups) {
    if (!props.formData.comboGroups) {
      props.formData.comboGroups = []
    }
    console.log('ComboSettings - 设置套餐组数据:', newGroups)
    props.formData.comboGroups.splice(0, props.formData.comboGroups.length, ...newGroups)
  }
})

// 添加套餐组
const addGroup = () => {
  const newGroup = {
    id: generateId(),
    name: '',
    description: '',
    minSelect: 1,
    maxSelect: 1,
    isRequired: true,
    sort: groups.value.length + 1,
    options: []
  }
  groups.value.push(newGroup)
}

// 删除套餐组
const removeGroup = (index) => {
  groups.value.splice(index, 1)
}

// 打开商品选择器
const openGoodsSelector = (groupIndex) => {
  currentGroupIndex.value = groupIndex
  goodsSelectorVisible.value = true
}

// 处理商品选择
const handleGoodsSelect = (selectedGoods) => {
  console.log('接收到的商品数据:', selectedGoods)
  
  if (currentGroupIndex.value === -1) {
    console.error('当前组索引无效:', currentGroupIndex.value)
    return
  }
  
  const group = groups.value[currentGroupIndex.value]
  if (!group.options) {
    group.options = []
  }
  
  // 处理选中的商品
  const goodsArray = Array.isArray(selectedGoods) ? selectedGoods : [selectedGoods]
  console.log('处理的商品数组:', goodsArray)
  
  goodsArray.forEach(goods => {
    console.log('处理单个商品:', goods)
    
    // 统一字段名称，兼容不同的数据结构
    const goodsId = goods.goodsId || goods.id
    const goodsName = goods.goodsName || goods.name
    const goodsImage = goods.goodsImage || goods.coverImg || goods.image
    
    // 检查是否已存在
    const exists = group.options.some(item => 
      item.goodsId === goodsId && 
      (!goods.skuId || item.skuId === goods.skuId)
    )
    
    if (!exists) {
      const newOption = {
        id: generateId(),
        name: goodsName,
        value: `goods_${goodsId}`,
        goodsId: goodsId,
        goodsName: goodsName,
        goodsImage: goodsImage,
        skuId: goods.skuId || null,
        skuInfo: goods.skuInfo || null,
        addPrice: 0,
        quantity: 1,
        isDefault: group.options.length === 0, // 第一个选项默认为默认选项
        sort: group.options.length + 1
      }
      console.log('添加新选项:', newOption)
      group.options.push(newOption)
    } else {
      console.log('选项已存在，跳过添加')
    }
  })
  
  console.log('更新后的组选项:', group.options)
  goodsSelectorVisible.value = false
}

// 删除选项
const removeOption = (groupIndex, optionIndex) => {
  const group = groups.value[groupIndex]
  if (group.options && group.options[optionIndex]) {
    group.options.splice(optionIndex, 1)
  }
}

// 处理默认选项变化
const handleDefaultChange = (groupIndex, optionIndex) => {
  const group = groups.value[groupIndex]
  const option = group.options[optionIndex]
  
  if (option.isDefault) {
    // 如果设置为默认，取消其他选项的默认状态
    group.options.forEach((item, index) => {
      if (index !== optionIndex) {
        item.isDefault = false
      }
    })
  } else {
    // 如果取消默认，检查是否还有其他默认选项
    const hasDefault = group.options.some((item, index) => 
      index !== optionIndex && item.isDefault
    )
    
    if (!hasDefault && group.isRequired) {
      ElMessage.warning('必选组至少需要一个默认选项')
      option.isDefault = true
    }
  }
}

// 初始化套餐组（供父组件调用）
const initComboGroups = () => {
  // 确保在新增模式下有一个默认的套餐组
  if (groups.value.length === 0) {
    addGroup()
  }
}

// 暴露方法给父组件
defineExpose({
  initComboGroups
})
</script>

<style lang="scss" scoped>
.combo-settings-container {
  padding: 0;
  
  .combo-basic-config {
    background: #fff;
    border-radius: 12px;
    border: 1px solid #e4e7ed;
    padding: 20px 24px;
    margin-bottom: 20px;
    
    .unit {
      margin-left: 8px;
      color: #909399;
      font-size: 14px;
    }
    
    .form-tip {
      font-size: 12px;
      color: #909399;
      margin-top: 4px;
    }
  }
  
  .groups-content-card {
    background: #fff;
    border-radius: 12px;
    border: 1px solid #e4e7ed;
    margin-bottom: 20px;
    overflow: hidden;
    
    .card-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 20px 24px;
      background: #f8f9fa;
      border-bottom: 1px solid #e4e7ed;
      
      .card-title {
        display: flex;
        align-items: center;
        gap: 8px;
        margin: 0;
        font-size: 16px;
        font-weight: 600;
        color: #303133;
      }
    }
  }
  
  .groups-container {
    padding: 24px;
    
    .group-card {
      background: #fafbfc;
      border: 1px solid #e4e7ed;
      border-radius: 8px;
      margin-bottom: 16px;
      
      .group-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 16px 20px;
        background: #f0f2f5;
        border-bottom: 1px solid #e4e7ed;
        
        .group-title {
          display: flex;
          align-items: center;
          gap: 8px;
          font-weight: 600;
          color: #303133;
        }
      }
      
      .group-content {
        padding: 20px;
        
        .group-basic-info {
          margin-bottom: 24px;
          
          :deep(.el-form) {
            .el-form-item {
              margin-bottom: 16px;
              
              .el-form-item__label {
                font-weight: 500;
                color: #606266;
                font-size: 13px;
                line-height: 1.4;
                padding-bottom: 0;
              }
              
              .el-form-item__content {
                line-height: 1.4;
              }
            }
            
            .el-row {
              margin-bottom: 0;
              
              &:last-child {
                margin-bottom: 0;
              }
            }
          }
        }
        
        .group-options-section {
          .section-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 16px;
            
            .form-label {
              font-weight: 500;
              color: #606266;
              font-size: 14px;
              margin: 0;
            }
          }
          
          .options-list {
            .option-item {
              display: flex;
              align-items: center;
              justify-content: space-between;
              padding: 16px;
              background: #fff;
              border: 1px solid #e4e7ed;
              border-radius: 8px;
              margin-bottom: 12px;
              
              &:last-child {
                margin-bottom: 0;
              }
              
              .option-info {
                display: flex;
                align-items: center;
                flex: 1;
                
                .option-image {
                  width: 60px;
                  height: 60px;
                  border-radius: 6px;
                  margin-right: 12px;
                  
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
                }
                
                .option-details {
                  .option-name {
                    font-weight: 500;
                    color: #303133;
                    margin-bottom: 4px;
                  }
                  
                  .option-goods {
                    font-size: 12px;
                    color: #909399;
                  }
                }
              }
              
              .option-config {
                display: flex;
                align-items: center;
                gap: 16px;
                
                .config-item {
                  display: flex;
                  align-items: center;
                  gap: 4px;
                  
                  label {
                    font-size: 12px;
                    color: #606266;
                    white-space: nowrap;
                  }
                  
                  .unit {
                    font-size: 12px;
                    color: #909399;
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}

:deep(.el-empty) {
  padding: 40px 0;
}
</style> 