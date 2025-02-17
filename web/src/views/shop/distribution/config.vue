<template>
  <div class="distribution-config">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="分销模式" name="mode">
        <!-- 分销模式配置 -->
        <el-form :model="configForm" label-width="120px">
          <el-form-item label="分销启用">
            <el-switch v-model="configForm.enabled" />
          </el-form-item>
          
          <el-form-item label="分销层级">
            <el-radio-group v-model="configForm.level">
              <el-radio :label="1">一级分销</el-radio>
              <el-radio :label="2">二级分销</el-radio>
            </el-radio-group>
          </el-form-item>
          
          <el-form-item label="分销模式">
            <el-radio-group v-model="configForm.mode">
              <el-radio :label="1">指定分销</el-radio>
              <el-radio :label="2">人人分销</el-radio>
              <el-radio :label="3">满额分销</el-radio>
            </el-radio-group>
          </el-form-item>
          
          <!-- 其他配置项... -->
        </el-form>
      </el-tab-pane>
      
      <!-- 其他配置标签页... -->
    </el-tabs>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const activeTab = ref('mode')
const configForm = ref({
  enabled: false,
  level: 1,
  mode: 1,
  applyEnabled: false,
  bindType: 1,
  bindMode: 1,
  centerType: 1
})
</script> 