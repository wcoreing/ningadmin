<template>
  <el-dialog
    v-model="dialogVisible"
    title="表单预览"
    width="600px"
  >
    <el-form
      :model="formData"
      label-width="100px"
      :style="{ padding: '20px' }"
    >
      <el-form-item
        v-for="item in formItems"
        :key="item.key"
        :label="item.props.label"
        :prop="item.key"
        :rules="item.props.required ? { required: true, message: `请输入${item.props.label}`, trigger: 'blur' } : []"
      >
        <template v-if="item.type === 'el-select'">
          <el-select
            v-model="formData[item.key]"
            :placeholder="item.props.placeholder"
            :style="{ width: '100%' }"
          >
            <el-option
              v-for="option in item.props.options"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </template>
        <template v-else-if="item.type === 'el-checkbox-group'">
          <el-checkbox-group
            v-model="formData[item.key]"
            :style="{ width: '100%' }"
          >
            <el-checkbox
              v-for="option in item.props.options"
              :key="option.value"
              :value="option.value"
            >
              {{ option.label }}
            </el-checkbox>
          </el-checkbox-group>
        </template>
        <template v-else-if="item.type === 'el-date-picker'">
          <el-date-picker
            v-model="formData[item.key]"
            v-bind="item.props"
            :style="{ width: '100%' }"
            :editable="false"
            :clearable="true"
          />
        </template>
        <component
          v-else
          :is="item.type"
          v-model="formData[item.key]"
          v-bind="item.props"
          :style="{ width: '100%' }"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">关闭</el-button>
      <el-button type="primary" @click="handleSubmit">提交</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  formItems: {
    type: Array,
    required: true
  },
  formData: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update:visible', 'submit'])

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const handleSubmit = () => {
  emit('submit')
  dialogVisible.value = false
}
</script> 