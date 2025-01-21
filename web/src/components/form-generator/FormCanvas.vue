<template>
  <el-main class="form-canvas">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>表单画布</span>
          <div class="button-group">
            <slot name="header-actions" />
            <el-button type="primary" @click="emit('preview')">预览</el-button>
            <el-button type="success" @click="emit('generate')">生成代码</el-button>
            <el-button @click="handleClear">清空</el-button>
          </div>
        </div>
      </template>
      <el-form
        :model="formData"
        label-width="100px"
        :disabled="preview"
      >
        <draggable
          v-model="formItemsValue"
          group="components"
          item-key="key"
          class="form-items-container"
          @start="dragStart"
          @end="dragEnd"
          @add="handleAdd"
        >
          <template #item="{ element, index }">
            <div 
              class="form-item-wrapper"
              :class="{ active: currentItemValue?.key === element.key }"
              @click="selectItem(element)"
            >
              <el-form-item
                :label="element.props.label"
                :prop="element.key"
              >
                <component
                  :is="element.type"
                  v-model="formData[element.key]"
                  v-bind="element.props"
                >
                  <template v-if="element.type === 'el-date-picker'">
                    <!-- 日期选择器不需要额外的模板 -->
                  </template>
                  <template v-if="element.type === 'el-select'">
                    <el-option
                      v-for="option in element.props.options"
                      :key="option.value"
                      :label="option.label"
                      :value="option.value"
                    />
                  </template>
                  <template v-if="element.type === 'el-checkbox-group'">
                    <el-checkbox
                      v-for="option in element.props.options"
                      :key="option.value"
                      :label="option.value"
                    >
                      {{ option.label }}
                    </el-checkbox>
                  </template>
                </component>
              </el-form-item>
              <div class="form-item-actions" v-show="currentItemValue?.key === element.key">
                <el-button-group>
                  <el-button 
                    link 
                    type="primary" 
                    :icon="ArrowUp"
                    @click.stop="moveItem(index, -1)"
                    :disabled="index === 0"
                  />
                  <el-button 
                    link 
                    type="primary" 
                    :icon="ArrowDown"
                    @click.stop="moveItem(index, 1)"
                    :disabled="index === formItemsValue.length - 1"
                  />
                  <el-button 
                    link 
                    type="danger" 
                    :icon="Delete"
                    @click.stop="removeItem(index)"
                  />
                </el-button-group>
              </div>
            </div>
          </template>
        </draggable>
      </el-form>
    </el-card>
  </el-main>
</template>

<script setup>
import { ref, computed } from 'vue'
import draggable from 'vuedraggable'
import { ArrowUp, ArrowDown, Delete } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  formItems: {
    type: Array,
    required: true
  },
  currentItem: {
    type: Object,
    default: null
  },
  formData: {
    type: Object,
    required: true
  },
  preview: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:formItems', 'update:currentItem', 'preview', 'generate'])

const formItemsValue = computed({
  get: () => props.formItems,
  set: (val) => emit('update:formItems', val)
})

const currentItemValue = computed({
  get: () => props.currentItem,
  set: (val) => emit('update:currentItem', val)
})

const selectItem = (item) => {
  currentItemValue.value = item
}

const moveItem = (index, direction) => {
  const newIndex = index + direction
  const items = [...formItemsValue.value]
  const item = items[index]
  items.splice(index, 1)
  items.splice(newIndex, 0, item)
  formItemsValue.value = items
}

const removeItem = (index) => {
  formItemsValue.value.splice(index, 1)
  if (currentItemValue.value?.key === formItemsValue.value[index]?.key) {
    currentItemValue.value = null
  }
}

const handleClear = () => {
  formItemsValue.value = []
  currentItemValue.value = null
}

const dragStart = () => {
  currentItemValue.value = null
}

const dragEnd = () => {
  // 拖拽结束后的处理
}

const handleAdd = (evt) => {
  emit('clone', evt.item.__draggable_context.element)
}
</script>

<style scoped>
.form-canvas {
  background-color: #f0f2f5;
  padding: 20px;
}

.form-items-container {
  min-height: 200px;
  padding: 20px;
  border: 1px dashed #d9d9d9;
  border-radius: 4px;
}

.form-item-wrapper {
  position: relative;
  padding: 12px;
  border: 1px dashed transparent;
  border-radius: 4px;
}

.form-item-wrapper:hover {
  border-color: #1890ff;
}

.form-item-wrapper.active {
  border-color: #1890ff;
  background-color: #e6f7ff;
}

.form-item-actions {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.button-group {
  display: flex;
  gap: 8px;
}
</style> 