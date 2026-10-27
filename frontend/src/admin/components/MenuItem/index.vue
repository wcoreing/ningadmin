<template>
  <template v-for="menu in menuList" :key="menu.id">
    <el-sub-menu v-if="menu.children && menu.children.length" :index="String(menu.id)">
      <template #title>
        <el-icon v-if="menu.icon"><component :is="menu.icon" /></el-icon>
        <span>{{ menu.name }}</span>
      </template>
      <recursive-menu :menu-list="menu.children" />
    </el-sub-menu>

    <el-menu-item v-else :index="adminPath(menu.path)">
      <el-icon v-if="menu.icon"><component :is="menu.icon" /></el-icon>
      <template #title>{{ menu.name }}</template>
    </el-menu-item>
  </template>
</template>

<script setup>
import { adminPath } from '@/router/base'

defineProps({
  menuList: {
    type: Array,
    required: true
  }
})

// 关键：为递归组件定义名称
defineOptions({
  name: 'recursive-menu'
})
</script> 