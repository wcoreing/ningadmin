<template>
  <el-aside width="300px" class="props-panel">
    <el-card v-if="currentItem">
      <template #header>
        <div class="card-header">
          <span>组件属性</span>
        </div>
      </template>
      <el-form label-position="top">
        <el-form-item label="标签文本">
          <el-input v-model="currentItem.props.label" />
        </el-form-item>
        <el-form-item label="字段名">
          <el-input v-model="currentItem.key" />
        </el-form-item>
        <el-form-item label="占位文本" v-if="currentItem.type === 'el-input'">
          <el-input v-model="currentItem.props.placeholder" />
        </el-form-item>
        <template v-if="['el-select', 'el-checkbox-group'].includes(currentItem.type)">
          <el-form-item label="选项配置">
            <div v-for="(option, index) in currentItem.props.options" :key="index" class="option-item">
              <el-input 
                v-model="option.label" 
                placeholder="选项文本" 
                @change="() => handleOptionChange(index)"
                style="width: 120px; margin-right: 8px" 
              />
              <el-input 
                v-model="option.value" 
                placeholder="选项值" 
                @change="() => handleOptionChange(index)"
                style="width: 120px; margin-right: 8px" 
              />
              <el-button type="danger" link :icon="Delete" @click="removeOption(index)" />
            </div>
            <el-button type="primary" link :icon="Plus" @click="addOption">添加选项</el-button>
          </el-form-item>
        </template>
        <el-form-item label="是否必填">
          <el-switch v-model="currentItem.props.required" />
        </el-form-item>
      </el-form>
    </el-card>
  </el-aside>
</template>

<script setup>
import { Delete, Plus } from '@element-plus/icons-vue'

const props = defineProps({
  currentItem: {
    type: Object,
    required: true
  }
})

const addOption = () => {
  if (!props.currentItem.props.options) {
    props.currentItem.props.options = []
  }
  props.currentItem.props.options.push({
    label: `选项${props.currentItem.props.options.length + 1}`,
    value: `${props.currentItem.props.options.length + 1}`
  })
  props.currentItem.props.options = [...props.currentItem.props.options]
}

const removeOption = (index) => {
  props.currentItem.props.options.splice(index, 1)
  props.currentItem.props.options = [...props.currentItem.props.options]
}

const handleOptionChange = (index) => {
  props.currentItem.props.options = [...props.currentItem.props.options]
}
</script>

<style scoped>
.props-panel {
  border-left: 1px solid #e6e6e6;
  background-color: #fff;
}

.option-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}
</style> 