<template>
  <div class="goods-display">
    <div class="selected-goods" v-if="modelValue && goodsInfo">
      <div class="goods-info">
        <el-image
          :src="goodsInfo.coverImg"
          fit="contain"
          class="goods-image"
        >
          <template #error>
            <div class="image-placeholder">
              <el-icon><Picture /></el-icon>
            </div>
          </template>
        </el-image>
        <div class="goods-detail">
          <div class="goods-name" :title="goodsInfo.name">{{ goodsInfo.name }}</div>
          <div class="goods-meta">
            <span class="goods-id">ID: {{ goodsInfo.goodsId }}</span>
            <span class="goods-price">¥{{ goodsInfo.price }}</span>
          </div>
        </div>
      </div>
      <div class="actions">
        <el-button link type="primary" @click="handleChange">更换</el-button>
        <el-button link type="danger" @click="handleClear">清除</el-button>
      </div>
    </div>
    <el-button v-else type="primary" @click="handleChange">
      <el-icon><Plus /></el-icon>选择商品
    </el-button>
  </div>
</template>

<script setup>
import { Picture, Plus } from '@element-plus/icons-vue'

defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  goodsInfo: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'change', 'clear'])

const handleChange = () => {
  emit('change')
}

const handleClear = () => {
  emit('update:modelValue', '')
  emit('clear')
}
</script>

<style lang="scss" scoped>
.goods-display {
  display: flex;
  align-items: center;
  gap: 8px;

  .selected-goods {
    flex: 1;
    border: 1px solid var(--el-border-color);
    border-radius: 4px;
    padding: 8px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    background-color: var(--el-fill-color-blank);
    
    .goods-info {
      display: flex;
      align-items: center;
      gap: 12px;
      flex: 1;
      min-width: 0;
      
      .goods-image {
        width: 60px;
        height: 60px;
        border-radius: 4px;
        border: 1px solid var(--el-border-color-lighter);
        padding: 4px;
        background-color: #fff;
      }
      
      .goods-detail {
        flex: 1;
        min-width: 0;
        
        .goods-name {
          font-size: 14px;
          color: var(--el-text-color-primary);
          margin-bottom: 4px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        
        .goods-meta {
          display: flex;
          align-items: center;
          gap: 12px;
          color: var(--el-text-color-secondary);
          font-size: 13px;
          
          .goods-id {
            color: var(--el-text-color-secondary);
          }
          
          .goods-price {
            color: #f56c6c;
            font-weight: bold;
          }
        }
      }
    }
    
    .actions {
      display: flex;
      gap: 8px;
      margin-left: 12px;
    }
  }
}
</style> 