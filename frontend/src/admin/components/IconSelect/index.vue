<template>
  <el-dialog
    v-model="dialogVisible"
    title="选择图标"
    width="800px"
  >
    <div class="icon-container">
      <div 
        v-for="(icon, name) in iconList" 
        :key="name"
        class="icon-item"
        @click="handleSelect(name)"
      >
        <el-icon>
          <component :is="icon" />
        </el-icon>
        <span class="icon-name">{{ name }}</span>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: Boolean,
  selectedIcon: String
})

const emit = defineEmits(['update:modelValue', 'select'])

// 图标列表
const iconList = ElementPlusIconsVue

// 对话框显示状态
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// 选择图标
const handleSelect = (iconName) => {
  emit('select', iconName)
  dialogVisible.value = false
}
</script>

<style scoped>
.icon-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 16px;
  padding: 16px;
}

.icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 16px;
  border: 1px solid #eee;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.icon-item:hover {
  background-color: #f5f7fa;
  border-color: #409eff;
  color: #409eff;
}

.icon-name {
  margin-top: 8px;
  font-size: 12px;
}

.el-icon {
  font-size: 24px;
}
</style> 