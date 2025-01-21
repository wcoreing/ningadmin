<template>
  <el-aside width="250px" class="component-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>表单组件</span>
        </div>
      </template>
      <el-collapse v-model="activeCollapse">
        <el-collapse-item title="基础组件" name="basic">
          <draggable
            :list="basicComponents"
            :group="{ name: 'components', pull: 'clone', put: false }"
            item-key="type"
            :sort="false"
            :clone="cloneComponent"
            @add="handleAdd"
          >
            <template #item="{ element }">
              <div class="component-item">
                <el-icon><component :is="element.icon" /></el-icon>
                <span>{{ element.label }}</span>
              </div>
            </template>
          </draggable>
        </el-collapse-item>
        <el-collapse-item title="高级组件" name="advanced">
          <draggable
            :list="advancedComponents"
            :group="{ name: 'components', pull: 'clone', put: false }"
            item-key="type"
            :sort="false"
            :clone="cloneComponent"
            @add="handleAdd"
          >
            <template #item="{ element }">
              <div class="component-item">
                <el-icon><component :is="element.icon" /></el-icon>
                <span>{{ element.label }}</span>
              </div>
            </template>
          </draggable>
        </el-collapse-item>
      </el-collapse>
    </el-card>
  </el-aside>
</template>

<script setup>
import { ref } from 'vue'
import draggable from 'vuedraggable'
import { basicComponents } from './config/basic-components'
import { advancedComponents } from './config/advanced-components'

const activeCollapse = ref(['basic'])

const emit = defineEmits(['clone'])

const cloneComponent = (item) => {
  const cloned = {
    ...item,
    key: `${item.type}_${Date.now()}`,
    props: JSON.parse(JSON.stringify(item.props))
  }
  return cloned
}

const handleAdd = (evt) => {
  if (evt.added) {
    emit('clone', evt.added.element)
  }
}
</script>

<style scoped>
.component-list {
  border-right: 1px solid #e6e6e6;
  background-color: #fff;
}

.component-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  margin: 4px 0;
  cursor: move;
  border: 1px dashed #d9d9d9;
  border-radius: 4px;
}

.component-item:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.component-item .el-icon {
  margin-right: 8px;
  font-size: 16px;
}
</style> 