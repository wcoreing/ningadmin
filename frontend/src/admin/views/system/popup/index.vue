<template>
  <div class="popup-management">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :model="queryParams" inline>
        <el-form-item label="标题">
          <el-input
            v-model="queryParams.title"
            placeholder="请输入弹窗标题"
            clearable
            style="width: 200px"
            :prefix-icon="Search"
          />
        </el-form-item>
        <el-form-item label="触发类型">
          <el-select v-model="queryParams.trigger_type" placeholder="请选择触发类型" clearable style="width: 150px">
            <el-option label="手动触发" value="manual" />
            <el-option label="自动触发" value="auto" />
            <el-option label="定时触发" value="schedule" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable style="width: 120px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getList" :icon="Search">搜索</el-button>
          <el-button @click="resetQuery" :icon="Refresh">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作区域 -->
    <el-card class="table-card">
      <div class="table-header">
        <div class="header-left">
          <h3 class="page-title">
            <el-icon><Bell /></el-icon>
            弹窗管理
          </h3>
          <p class="page-desc">管理系统中的各种弹窗消息和通知</p>
        </div>
        <div class="header-right">
          <el-button type="primary" @click="handleAdd" :icon="Plus" size="large">
            新增弹窗
          </el-button>
        </div>
      </div>

      <!-- 表格 -->
      <el-table v-loading="loading" :data="popupList" style="width: 100%" row-key="id">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="150">
          <template #default="{ row }">
            <div class="title-cell">
              <span class="title-text" :class="{ 'no-title': !row.title }">
                {{ row.title || '无标题' }}
              </span>
              <el-tag v-if="row.priority > 0" size="small" type="warning" class="priority-tag">
                优先级 {{ row.priority }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="trigger_type" label="触发类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTriggerTypeTag(row.trigger_type)" :icon="getTriggerTypeIcon(row.trigger_type)">
              {{ getTriggerTypeName(row.trigger_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              active-color="#13ce66"
              inactive-color="#ff4949"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="统计信息" width="160">
          <template #default="{ row }">
            <div class="stats-cell">
              <div class="stat-item">
                <el-icon class="stat-icon"><View /></el-icon>
                <span>{{ row.show_count }}</span>
              </div>
              <div class="stat-item">
                <el-icon class="stat-icon"><Mouse /></el-icon>
                <span>{{ row.click_count }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="时间范围" width="200">
          <template #default="{ row }">
            <div class="time-range">
              <div v-if="row.start_time" class="time-item">
                <el-icon><Timer /></el-icon>
                {{ formatDate(row.start_time) }}
              </div>
              <div v-if="row.end_time" class="time-item">
                <el-icon><Clock /></el-icon>
                {{ formatDate(row.end_time) }}
              </div>
              <div v-if="!row.start_time && !row.end_time" class="time-item permanent">
                <el-icon><Unlock /></el-icon>
                永久有效
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.create_time) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-tooltip content="编辑弹窗" placement="top">
                <el-button type="primary" link @click="handleEdit(row)" :icon="Edit">编辑</el-button>
              </el-tooltip>
              <el-tooltip content="预览效果" placement="top">
                <el-button type="success" link @click="handlePreview(row)" :icon="View">预览</el-button>
              </el-tooltip>
              <el-tooltip content="删除弹窗" placement="top">
                <el-button type="danger" link @click="handleDelete(row)" :icon="Delete">删除</el-button>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.page_size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="getList"
          @current-change="getList"
        />
      </div>
    </el-card>

    <!-- 弹窗表单 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="900px"
      :close-on-click-modal="false"
      class="popup-dialog"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <!-- 基础信息 -->
        <el-card class="form-section compact" shadow="never">
          <template #header>
            <div class="section-header">
              <el-icon><InfoFilled /></el-icon>
              <span>基础信息</span>
            </div>
          </template>
          
          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="弹窗标题" prop="title" label-width="80px">
                <el-input 
                  v-model="form.title" 
                  placeholder="请输入弹窗标题"
                  clearable
                  size="small"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="触发类型" prop="trigger_type" label-width="80px">
                <el-select v-model="form.trigger_type" placeholder="请选择触发类型" style="width: 100%" size="small">
                  <el-option label="手动触发" value="manual" />
                  <el-option label="自动触发" value="auto" />
                  <el-option label="定时触发" value="schedule" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="弹窗内容" prop="content" label-width="80px">
            <el-input
              v-model="form.content"
              type="textarea"
              :rows="3"
              placeholder="请输入弹窗内容，支持HTML标签"
              show-word-limit
              maxlength="500"
              size="small"
            />
          </el-form-item>

          <el-row :gutter="16">
            <el-col :span="8">
              <el-form-item label="优先级" prop="priority" label-width="60px">
                <el-input-number 
                  v-model="form.priority" 
                  :min="0" 
                  :max="999" 
                  style="width: 100%"
                  controls-position="right"
                  size="small"
                />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="自动关闭" label-width="70px">
                <el-input-number 
                  v-model="form.auto_close" 
                  :min="0" 
                  style="width: 100%"
                  controls-position="right"
                  size="small"
                />
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="延迟显示" label-width="70px">
                <el-input-number 
                  v-model="form.delay" 
                  :min="0" 
                  style="width: 100%"
                  controls-position="right"
                  size="small"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="开始时间" label-width="80px">
                <el-date-picker
                  v-model="form.start_time"
                  type="datetime"
                  placeholder="选择开始时间"
                  format="YYYY-MM-DD HH:mm:ss"
                  value-format="YYYY-MM-DD HH:mm:ss"
                  style="width: 100%"
                  size="small"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="结束时间" label-width="80px">
                <el-date-picker
                  v-model="form.end_time"
                  type="datetime"
                  placeholder="选择结束时间"
                  format="YYYY-MM-DD HH:mm:ss"
                  value-format="YYYY-MM-DD HH:mm:ss"
                  style="width: 100%"
                  size="small"
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="遮罩关闭" label-width="80px">
                <el-switch 
                  v-model="form.mask_closable" 
                  :active-value="1" 
                  :inactive-value="0"
                  active-color="#13ce66"
                  inactive-color="#ff4949"
                  size="small"
                />
                <span class="switch-tip-compact">点击遮罩层关闭</span>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="启用状态" label-width="80px">
                <el-switch 
                  v-model="form.status" 
                  :active-value="1" 
                  :inactive-value="0"
                  active-color="#13ce66"
                  inactive-color="#ff4949"
                  size="small"
                />
                <span class="switch-tip-compact">控制弹窗启用</span>
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>

        <!-- 样式配置 -->
        <el-card class="form-section compact" shadow="never">
          <template #header>
            <div class="section-header">
              <el-icon><Brush /></el-icon>
              <span>标题样式</span>
            </div>
          </template>
          
          <div class="style-config-compact">
            <!-- 第一行：基础样式 -->
            <el-row :gutter="12">
              <el-col :span="6">
                <div class="style-item">
                  <span class="style-label">文字颜色</span>
                  <el-color-picker v-model="titleStyle.color" size="small" />
                </div>
              </el-col>
              <el-col :span="6">
                <div class="style-item">
                  <span class="style-label">字体大小</span>
                  <el-select v-model="titleStyle.fontSize" size="small" style="width: 100%">
                    <el-option label="12px" value="12px" />
                    <el-option label="14px" value="14px" />
                    <el-option label="16px" value="16px" />
                    <el-option label="18px" value="18px" />
                    <el-option label="20px" value="20px" />
                    <el-option label="24px" value="24px" />
                  </el-select>
                </div>
              </el-col>
              <el-col :span="6">
                <div class="style-item">
                  <span class="style-label">字体粗细</span>
                  <el-select v-model="titleStyle.fontWeight" size="small" style="width: 100%">
                    <el-option label="正常" value="normal" />
                    <el-option label="粗体" value="bold" />
                    <el-option label="细体" value="lighter" />
                  </el-select>
                </div>
              </el-col>
              <el-col :span="6">
                <div class="style-item">
                  <span class="style-label">文字对齐</span>
                  <el-select v-model="titleStyle.textAlign" size="small" style="width: 100%">
                    <el-option label="左对齐" value="left" />
                    <el-option label="居中" value="center" />
                    <el-option label="右对齐" value="right" />
                  </el-select>
                </div>
              </el-col>
            </el-row>
            
            <!-- 第二行：背景样式 -->
            <el-row :gutter="12" style="margin-top: 12px">
              <el-col :span="6">
                <div class="style-item">
                  <span class="style-label">背景颜色</span>
                  <el-color-picker 
                    v-model="titleStyle.backgroundColor" 
                    size="small" 
                    show-alpha 
                    :predefine="backgroundColors"
                  />
                </div>
              </el-col>
              <el-col :span="10">
                <div class="style-item">
                  <span class="style-label">背景图片</span>
                  <div class="bg-image-input">
                    <el-input 
                      v-model="titleStyle.backgroundImage" 
                      placeholder="输入图片URL或选择上传"
                      size="small"
                      clearable
                    />
                    <el-button type="primary" link size="small" @click="openImageSelector" :icon="Picture">
                      选择
                    </el-button>
                    <el-button 
                      type="danger" 
                      link 
                      size="small" 
                      @click="clearBackground" 
                      v-if="titleStyle.backgroundImage || titleStyle.backgroundColor"
                      title="清除背景"
                    >
                      清除
                    </el-button>
                  </div>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="style-item">
                  <span class="style-label">背景模式</span>
                  <el-select v-model="titleStyle.backgroundSize" size="small" style="width: 100%">
                    <el-option label="覆盖" value="cover" />
                    <el-option label="包含" value="contain" />
                    <el-option label="拉伸" value="100% 100%" />
                    <el-option label="原始" value="auto" />
                  </el-select>
                </div>
              </el-col>
            </el-row>
            
            <!-- 预览区域 -->
            <div class="style-preview-compact" :style="getTitlePreviewStyle()">
              <span class="preview-label">预览：</span>
              <span class="preview-title-text">{{ form.title || '系统维护通知' }}</span>
            </div>
          </div>
        </el-card>

        <!-- 按钮配置 -->
        <el-card class="form-section compact" shadow="never">
          <template #header>
            <div class="section-header">
              <el-icon><Operation /></el-icon>
              <span>按钮配置</span>
              <el-button type="primary" link @click="addButton" :icon="Plus" size="small" style="margin-left: auto;">
                添加按钮
              </el-button>
            </div>
          </template>
          
          <div class="button-config-compact">
            <div v-if="buttonList.length === 0" class="empty-buttons-compact">
              <span class="empty-text">暂无按钮配置</span>
              <el-button type="primary" @click="addButton" :icon="Plus" size="small">添加按钮</el-button>
            </div>
            <div v-else class="button-list">
              <div v-for="(button, index) in buttonList" :key="index" class="button-item-compact">
                <div class="button-row">
                  <span class="button-index">{{ index + 1 }}</span>
                  <el-input v-model="button.text" placeholder="按钮文字" size="small" style="width: 120px" />
                  <el-select v-model="button.type" placeholder="类型" size="small" style="width: 100px">
                    <el-option label="默认" value="default" />
                    <el-option label="主要" value="primary" />
                    <el-option label="成功" value="success" />
                    <el-option label="警告" value="warning" />
                    <el-option label="危险" value="danger" />
                  </el-select>
                  <el-input v-model="button.action" placeholder="动作" size="small" style="width: 120px" />
                  <div class="button-preview-compact">
                    <el-button :type="button.type || 'default'" size="small">
                      {{ button.text || '按钮' }}
                    </el-button>
                  </div>
                  <el-button type="danger" link @click="removeButton(index)" :icon="Delete" size="small">删除</el-button>
                </div>
              </div>
            </div>
          </div>
        </el-card>

        <!-- 触发规则 -->
        <el-card class="form-section compact" shadow="never">
          <template #header>
            <div class="section-header">
              <el-icon><Setting /></el-icon>
              <span>触发规则</span>
            </div>
          </template>
          
          <el-row :gutter="16">
            <el-col :span="24">
              <el-form-item label="触发页面" label-width="80px">
                <el-select
                  v-model="triggerRule.pages"
                  multiple
                  filterable
                  allow-create
                  placeholder="输入页面路径，如 /home, /user，留空表示所有页面"
                  style="width: 100%"
                  size="small"
                >
                  <el-option label="/home" value="/home" />
                  <el-option label="/user" value="/user" />
                  <el-option label="/product" value="/product" />
                  <el-option label="/order" value="/order" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item label="用户类型" label-width="80px">
                <el-select v-model="triggerRule.userType" placeholder="选择用户类型" clearable size="small" style="width: 100%">
                  <el-option label="所有用户" value="all" />
                  <el-option label="新用户" value="new" />
                  <el-option label="老用户" value="old" />
                  <el-option label="VIP用户" value="vip" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="显示频率" label-width="80px">
                <el-select v-model="triggerRule.frequency" placeholder="选择显示频率" clearable size="small" style="width: 100%">
                  <el-option label="每次访问" value="always" />
                  <el-option label="仅一次" value="once" />
                  <el-option label="每日一次" value="daily" />
                  <el-option label="每周一次" value="weekly" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false" :icon="Close">取消</el-button>
          <el-button type="primary" @click="submitForm" :icon="Check">
            {{ form.id ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 预览弹窗 -->
    <el-dialog
      v-model="previewVisible"
      width="500px"
      :close-on-click-modal="previewData.mask_closable === 1"
      class="preview-dialog"
      :show-close="true"
    >
      <template #header>
        <div 
          class="preview-title" 
          :style="previewData.title_style || {}"
          v-if="previewData.title"
        >
          {{ previewData.title }}
        </div>
      </template>
      <div class="preview-content" v-html="previewData.content"></div>
      <template #footer v-if="previewData.buttons && previewData.buttons.length > 0">
        <div class="preview-buttons">
          <el-button
            v-for="(button, index) in previewData.buttons"
            :key="index"
            :type="button.type || 'default'"
            @click="handlePreviewButtonClick(button)"
          >
            {{ button.text }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 文件管理器 -->
    <FileManager
      v-model="fileManagerVisible"
      :multiple="false"
      @select="handleImageSelected"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Search, Refresh, Plus, Edit, Delete, View, Mouse, Timer, Clock, Unlock,
  InfoFilled, Brush, Operation, Setting, EditPen, MagicStick, Bell,
  Close, Check, Picture
} from '@element-plus/icons-vue'
import {
  getPopupList,
  getPopupInfo,
  createPopup,
  updatePopup,
  deletePopup,
  updatePopupStatus
} from '@/admin/api/system/popup'
import FileManager from '@/admin/components/FileManager/index.vue'

// 响应式数据
const loading = ref(false)
const popupList = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const previewVisible = ref(false)
const previewData = ref({})
const formRef = ref()

// 查询参数
const queryParams = reactive({
  page: 1,
  page_size: 10,
  title: '',
  trigger_type: '',
  status: undefined
})

// 表单数据
const form = reactive({
  id: undefined,
  title: '',
  content: '',
  trigger_type: 'manual',
  priority: 0,
  auto_close: 0,
  delay: 0,
  mask_closable: 1,
  start_time: '',
  end_time: '',
  status: 1
})

// 样式配置
const titleStyle = reactive({
  color: '#303133',
  fontSize: '16px',
  fontWeight: 'normal',
  textAlign: 'left',
  backgroundColor: '',
  backgroundImage: '',
  backgroundSize: 'cover'
})

// 预定义背景颜色
const backgroundColors = [
  '#409EFF',
  '#67C23A',
  '#E6A23C',
  '#F56C6C',
  '#909399',
  'rgba(64, 158, 255, 0.1)',
  'rgba(103, 194, 58, 0.1)',
  'rgba(230, 162, 60, 0.1)',
  'rgba(245, 108, 108, 0.1)',
  'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
  'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)'
]

// 按钮配置
const buttonList = ref([])

// 触发规则
const triggerRule = reactive({
  pages: [],
  userType: '',
  frequency: ''
})

// 表单验证规则
const rules = {
  content: [
    { required: true, message: '请输入弹窗内容', trigger: 'blur' }
  ],
  trigger_type: [
    { required: true, message: '请选择触发类型', trigger: 'change' }
  ]
}

// 监听标题变化，更新预览
watch(() => form.title, () => {
  // 实时更新预览
})

// 监听样式变化，确保预览更新
watch(titleStyle, () => {
  // 样式变化时强制更新预览
}, { deep: true })

// 获取列表
const getList = async () => {
  loading.value = true
  try {
    const res = await getPopupList(queryParams)
    popupList.value = res.list || []
    total.value = res.total || 0
  } catch (error) {
    console.error('获取弹窗列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 重置查询
const resetQuery = () => {
  Object.assign(queryParams, {
    page: 1,
    page_size: 10,
    title: '',
    trigger_type: '',
    status: undefined
  })
  getList()
}

// 新增
const handleAdd = () => {
  resetForm()
  dialogTitle.value = '新增弹窗'
  dialogVisible.value = true
}

// 编辑
const handleEdit = async (row) => {
  try {
    const res = await getPopupInfo(row.id)
    Object.assign(form, res)
    
    // 先重置样式配置为默认值
    Object.assign(titleStyle, {
      color: '#303133',
      fontSize: '16px',
      fontWeight: 'normal',
      textAlign: 'left',
      backgroundColor: '',
      backgroundImage: '',
      backgroundSize: 'cover'
    })
    
    // 如果有样式配置，则应用
    if (res.title_style) {
      Object.assign(titleStyle, res.title_style)
    }
    
    // 处理按钮配置
    buttonList.value = res.buttons || []
    
    // 处理触发规则
    Object.assign(triggerRule, {
      pages: [],
      userType: '',
      frequency: ''
    })
    
    if (res.trigger_rule) {
      Object.assign(triggerRule, res.trigger_rule)
    }
    
    dialogTitle.value = '编辑弹窗'
    dialogVisible.value = true
  } catch (error) {
    console.error('获取弹窗详情失败:', error)
  }
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除弹窗"${row.title || '无标题'}"吗？`, '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deletePopup(row.id)
      ElMessage.success('删除成功')
      getList()
    } catch (error) {
      console.error('删除失败:', error)
    }
  })
}

// 状态切换
const handleStatusChange = async (row) => {
  try {
    await updatePopupStatus({
      id: row.id,
      status: row.status
    })
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('状态更新失败:', error)
    // 回滚状态
    row.status = row.status === 1 ? 0 : 1
  }
}

// 预览
const handlePreview = (row) => {
  previewData.value = { ...row }
  
  // 处理标题样式中的背景图片
  if (previewData.value.title_style && previewData.value.title_style.backgroundImage) {
    previewData.value.title_style = {
      ...previewData.value.title_style,
      backgroundImage: `url(${previewData.value.title_style.backgroundImage})`
    }
  }
  
  previewVisible.value = true
}

// 获取标题预览样式
const getTitlePreviewStyle = () => {
  const style = {
    padding: '8px 12px',
    minHeight: '40px',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    borderRadius: '4px'
  }
  
  // 基础文字样式
  style.color = titleStyle.color || '#303133'
  style.fontSize = titleStyle.fontSize || '16px'
  style.fontWeight = titleStyle.fontWeight || 'normal'
  style.textAlign = titleStyle.textAlign || 'left'
  
  // 处理背景颜色
  if (titleStyle.backgroundColor) {
    style.backgroundColor = titleStyle.backgroundColor
  } else {
    style.backgroundColor = '#f8f9fa'
  }
  
  // 处理背景图片
  if (titleStyle.backgroundImage) {
    style.backgroundImage = `url(${titleStyle.backgroundImage})`
    style.backgroundRepeat = 'no-repeat'
    style.backgroundPosition = 'center'
    style.backgroundSize = titleStyle.backgroundSize || 'cover'
  }
  
  return style
}

// 文件管理器相关
const fileManagerVisible = ref(false)

// 打开图片选择器
const openImageSelector = () => {
  fileManagerVisible.value = true
}

// 处理文件选择
const handleImageSelected = (file) => {
  if (file && file.url) {
    titleStyle.backgroundImage = file.url
    ElMessage.success('背景图片设置成功')
  }
}

// 清除背景
const clearBackground = () => {
  titleStyle.backgroundColor = ''
  titleStyle.backgroundImage = ''
  ElMessage.success('背景已清除')
}

// 预览按钮点击
const handlePreviewButtonClick = (button) => {
  ElMessage.info(`点击了按钮: ${button.text}`)
  if (button.action === 'cancel') {
    previewVisible.value = false
  }
}

// 添加按钮
const addButton = () => {
  buttonList.value.push({
    text: '按钮',
    type: 'default',
    action: ''
  })
}

// 删除按钮
const removeButton = (index) => {
  buttonList.value.splice(index, 1)
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 检查是否有有效的样式配置
        const hasValidTitleStyle = titleStyle.color !== '#303133' || 
                                   titleStyle.fontSize !== '16px' || 
                                   titleStyle.fontWeight !== 'normal' || 
                                   titleStyle.textAlign !== 'left' || 
                                   titleStyle.backgroundColor || 
                                   titleStyle.backgroundImage
        
        // 处理时间格式
        const formatTime = (time) => {
          if (!time) return null
          if (typeof time === 'string') {
            // 如果是字符串，转换为Date对象再格式化
            const date = new Date(time)
            return date.toISOString()
          }
          if (time instanceof Date) {
            return time.toISOString()
          }
          return null
        }
        
        const submitData = { 
          ...form,
          title_style: hasValidTitleStyle ? titleStyle : null,
          buttons: buttonList.value.length > 0 ? buttonList.value : null,
          trigger_rule: Object.values(triggerRule).some(v => v && v.length > 0) ? triggerRule : null,
          start_time: formatTime(form.start_time),
          end_time: formatTime(form.end_time)
        }

        if (form.id) {
          await updatePopup(submitData)
          ElMessage.success('更新成功')
        } else {
          await createPopup(submitData)
          ElMessage.success('创建成功')
        }
        
        dialogVisible.value = false
        getList()
      } catch (error) {
        console.error('保存失败:', error)
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    id: undefined,
    title: '',
    content: '',
    trigger_type: 'manual',
    priority: 0,
    auto_close: 0,
    delay: 0,
    mask_closable: 1,
    start_time: '',
    end_time: '',
    status: 1
  })
  
  // 重置样式配置
  Object.assign(titleStyle, {
    color: '#303133',
    fontSize: '16px',
    fontWeight: 'normal',
    textAlign: 'left',
    backgroundColor: '',
    backgroundImage: '',
    backgroundSize: 'cover'
  })
  
  // 重置按钮配置
  buttonList.value = []
  
  // 重置触发规则
  Object.assign(triggerRule, {
    pages: [],
    userType: '',
    frequency: ''
  })
  
  formRef.value?.resetFields()
}

// 工具函数
const getTriggerTypeName = (type) => {
  const map = {
    manual: '手动触发',
    auto: '自动触发',
    schedule: '定时触发'
  }
  return map[type] || type
}

const getTriggerTypeTag = (type) => {
  const map = {
    manual: '',
    auto: 'success',
    schedule: 'warning'
  }
  return map[type] || ''
}

const getTriggerTypeIcon = (type) => {
  const map = {
    manual: Mouse,
    auto: MagicStick,
    schedule: Timer
  }
  return map[type] || Mouse
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 初始化
onMounted(() => {
  getList()
})
</script>

<style scoped>
.popup-management {
  padding: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.header-left .page-title {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-left .page-desc {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.header-right {
  flex-shrink: 0;
}

/* 表格样式 */
.title-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-text {
  flex: 1;
}

.title-text.no-title {
  color: #c0c4cc;
  font-style: italic;
}

.priority-tag {
  flex-shrink: 0;
}

.stats-cell {
  display: flex;
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
}

.stat-icon {
  color: #909399;
}

.time-range .time-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #606266;
  margin-bottom: 4px;
}

.time-item.permanent {
  color: #67c23a;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

/* 弹窗表单样式 */
.popup-dialog :deep(.el-dialog__body) {
  max-height: 70vh;
  overflow-y: auto;
  padding: 0;
}

.form-section {
  margin-bottom: 20px;
  border: 1px solid #ebeef5;
}

.form-section:last-child {
  margin-bottom: 0;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #303133;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.switch-tip-compact {
  margin-left: 8px;
  font-size: 11px;
  color: #909399;
}

/* 紧凑型样式配置 */
.form-section.compact {
  margin-bottom: 16px;
}

.form-section.compact :deep(.el-card__body) {
  padding: 16px;
}

.style-config-compact {
  padding: 0;
}

.style-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.style-label {
  font-size: 12px;
  color: #606266;
  font-weight: 500;
  line-height: 1;
}

.style-preview-compact {
  margin-top: 12px;
  border-radius: 4px;
  border: 1px solid #ebeef5;
  position: relative;
  overflow: hidden;
}

.style-preview-compact .preview-label {
  font-size: 12px;
  color: #909399;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.9);
  padding: 2px 6px;
  border-radius: 3px;
  flex-shrink: 0;
}

.style-preview-compact .preview-title-text {
  flex: 1;
}

/* 紧凑型按钮配置 */
.button-config-compact {
  min-height: auto;
}

.empty-buttons-compact {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #f8f9fa;
  border-radius: 4px;
  border: 1px dashed #dcdfe6;
}

.empty-text {
  color: #909399;
  font-size: 14px;
}

.button-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.button-item-compact {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 12px;
  background: #fafafa;
}

.button-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.button-index {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: #409eff;
  color: white;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.button-preview-compact {
  margin-left: auto;
  flex-shrink: 0;
}

.preview-label {
  font-size: 12px;
  color: #606266;
  font-weight: 500;
}

/* 背景图片输入框 */
.bg-image-input {
  display: flex;
  gap: 8px;
  align-items: center;
}

.bg-image-input .el-input {
  flex: 1;
}

/* 选项样式 */
.option-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 预览弹窗 */
.preview-dialog .preview-content {
  min-height: 60px;
  line-height: 1.6;
}

.preview-dialog .preview-title {
  padding: 8px 12px;
  border-radius: 4px;
  background-repeat: no-repeat;
  background-position: center;
  background-size: cover;
  min-height: 40px;
  display: flex;
  align-items: center;
  font-weight: 600;
}

.preview-buttons {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 20px;
  border-top: 1px solid #ebeef5;
  background: #fafafa;
}

/* 响应式 */
@media (max-width: 768px) {
  .table-header {
    flex-direction: column;
    gap: 15px;
  }
  
  .header-right {
    align-self: stretch;
  }
  
  .stats-cell {
    flex-direction: column;
    gap: 4px;
  }
  
  .button-row {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .button-index {
    align-self: flex-start;
  }
  
  .button-preview-compact {
    margin-left: 0;
    margin-top: 8px;
  }
}
</style> 