<template>
  <el-input
    v-model="modelValue"
    v-bind="props"
    @update:modelValue="(val) => $emit('update:modelValue', val)"
  />
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  placeholder: String,
  disabled: Boolean
})

defineEmits(['update:modelValue'])
</script> 