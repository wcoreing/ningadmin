<template>
  <div class="app-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>任务详情</span>
          <div class="header-actions">
            <el-button 
              type="primary" 
              @click="handleDownload"
              v-if="canDownload"
              :loading="downloadLoading"
              size="small"
            >下载所有文件</el-button>
            <el-button 
              @click="goBack" 
              icon="ArrowLeft" 
              size="small"
            >返回</el-button>
          </div>
        </div>
      </template>

      <div v-if="taskDetail">
        <!-- 任务基本信息 -->
        <el-descriptions 
          title="基本信息" 
          :column="2" 
          border
          class="info-section"
        >
          <el-descriptions-item label="任务名称">{{ taskDetail.name }}</el-descriptions-item>
          <el-descriptions-item label="爬虫类型">{{ taskDetail.typeName }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(taskDetail.status)">
              {{ getStatusText(taskDetail.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ taskDetail.createdAt }}</el-descriptions-item>
          <el-descriptions-item label="完成时间" v-if="taskDetail.endTime">{{ taskDetail.endTime }}</el-descriptions-item>
          <el-descriptions-item label="进度">
            <el-progress 
              :percentage="taskDetail.progress" 
              :status="getProgressStatus(taskDetail)" 
              :format="progressFormat"
            />
          </el-descriptions-item>
        </el-descriptions>

        <!-- 任务配置信息 -->
        <el-collapse class="info-section">
          <el-collapse-item title="任务配置" name="config">
            <el-descriptions :column="2" border>
              <!-- 抖音视频爬虫配置显示 -->
              <template v-if="taskDetail.type === 'douyin_video'">
                <el-descriptions-item label="下载内容">
                  {{ formatDownloadContent(taskDetail.config.downloadContent) }}
                </el-descriptions-item>
                <el-descriptions-item label="视频质量">
                  {{ formatVideoQuality(taskDetail.config.videoQuality) }}
                </el-descriptions-item>
                <el-descriptions-item label="最大抓取数量">
                  {{ taskDetail.config.maxVideos || '不限制' }}
                </el-descriptions-item>
                <el-descriptions-item label="使用代理">
                  {{ taskDetail.config.useProxy ? '是' : '否' }}
                </el-descriptions-item>
                <el-descriptions-item label="超时设置">
                  {{ taskDetail.config.timeout }} 秒
                </el-descriptions-item>
                <el-descriptions-item label="重试次数">
                  {{ taskDetail.config.retries }}
                </el-descriptions-item>
              </template>
              <!-- 这里可以扩展其他爬虫类型的配置展示 -->
            </el-descriptions>
          </el-collapse-item>
          
          <el-collapse-item 
            title="链接信息" 
            name="urls"
            v-if="taskDetail.type === 'douyin_video' && (taskDetail.config.urls || taskDetail.config.authorUrl)"
          >
            <div v-if="taskDetail.config.authorUrl" class="url-item">
              <strong>作者主页：</strong>
              <el-link type="primary" :href="taskDetail.config.authorUrl" target="_blank">
                {{ taskDetail.config.authorUrl }}
              </el-link>
            </div>
            
            <div v-if="taskDetail.config.urls" class="url-list">
              <strong>视频链接：</strong>
              <div v-for="(url, index) in formatUrls(taskDetail.config.urls)" :key="index" class="url-item">
                <el-link type="primary" :href="url" target="_blank">
                  {{ url }}
                </el-link>
              </div>
            </div>
          </el-collapse-item>
        </el-collapse>

        <!-- 爬取结果列表 -->
        <div class="crawler-data">
          <div class="section-header">
            <h3>爬取结果</h3>
            <el-input
              v-model="dataSearch"
              placeholder="搜索标题"
              style="width: 200px"
              @input="filterCrawlerData"
            />
          </div>

          <el-table
            :data="filteredCrawlerData"
            style="width: 100%"
            v-loading="dataLoading"
          >
            <el-table-column type="expand">
              <template #default="props">
                <div class="data-detail">
                  <div v-if="props.row.metadata" class="metadata">
                    <strong>作者：</strong> {{ props.row.metadata.author || '未知' }}<br>
                    <strong>点赞：</strong> {{ props.row.metadata.likes || 0 }}<br>
                    <strong>评论：</strong> {{ props.row.metadata.comments || 0 }}<br>
                    <strong>分享：</strong> {{ props.row.metadata.shares || 0 }}<br>
                    <strong>描述：</strong> {{ props.row.metadata.description || '无' }}
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="预览" width="120">
              <template #default="scope">
                <div class="data-preview" v-if="scope.row.type === 'video'">
                  <el-image 
                    style="width: 100px; height: 100px;"
                    :src="scope.row.coverUrl || '/uploads/default.png'"
                    fit="cover"
                    :preview-src-list="[scope.row.coverUrl]"
                  />
                  <i class="el-icon-video-play preview-icon"></i>
                </div>
                <div class="data-preview" v-else-if="scope.row.type === 'image'">
                  <el-image 
                    style="width: 100px; height: 100px;"
                    :src="scope.row.url || '/uploads/default.png'"
                    fit="cover"
                    :preview-src-list="[scope.row.url]"
                  />
                </div>
                <div v-else>
                  <el-icon class="file-icon"><Document /></el-icon>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="title" label="标题" min-width="200" />
            <el-table-column label="类型" width="80">
              <template #default="scope">
                <el-tag
                  :type="getDataTypeTag(scope.row.type)"
                  size="small"
                >
                  {{ formatDataType(scope.row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="scope">
                <el-tag
                  :type="getDataStatusTag(scope.row.status)"
                  size="small"
                >
                  {{ formatDataStatus(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="大小" width="100">
              <template #default="scope">
                {{ formatFileSize(scope.row.size) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="scope">
                <el-button 
                  size="small" 
                  type="primary" 
                  @click="handleDownloadItem(scope.row)"
                  :disabled="!canDownloadItem(scope.row)"
                >下载</el-button>
                <el-button 
                  size="small" 
                  @click="handlePreview(scope.row)"
                  :disabled="!canPreview(scope.row)"
                >预览</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-container">
            <el-pagination
              background
              layout="total, sizes, prev, pager, next, jumper"
              :total="totalData"
              v-model:page-size="dataQueryParams.pageSize"
              v-model:current-page="dataQueryParams.page"
              @size-change="handleDataSizeChange"
              @current-change="handleDataCurrentChange"
            />
          </div>
        </div>
      </div>

      <div v-else-if="!loading" class="empty-data">
        <el-empty description="未找到任务详情" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
// 导入API
import { getTaskDetail, getTaskData, downloadTask, downloadTaskItem } from '@/api/crawler/task'

const route = useRoute()
const router = useRouter()
const taskId = route.params.id

const loading = ref(true)
const dataLoading = ref(false)
const downloadLoading = ref(false)
const taskDetail = ref(null)
const crawlerData = ref([])
const filteredCrawlerData = ref([])
const dataSearch = ref('')
const totalData = ref(0)

const dataQueryParams = reactive({
  page: 1,
  pageSize: 10
})

// 获取任务详情
const fetchTaskDetail = () => {
  loading.value = true
  getTaskDetail(taskId).then(res => {
    taskDetail.value = res
    loading.value = false
    fetchTaskData()
  }).catch((error) => {
    ElMessage.error('获取任务详情失败')
    loading.value = false
  })
}

// 获取任务数据列表
const fetchTaskData = () => {
  dataLoading.value = true
  getTaskData(taskId, dataQueryParams).then(res => {
    crawlerData.value = res.list || []
    filteredCrawlerData.value = [...crawlerData.value]
    totalData.value = res.total || 0
    dataLoading.value = false
  }).catch(() => {
    dataLoading.value = false
  })
}

// 筛选数据
const filterCrawlerData = () => {
  if (!dataSearch.value) {
    filteredCrawlerData.value = [...crawlerData.value]
  } else {
    filteredCrawlerData.value = crawlerData.value.filter(item => {
      return item.title && item.title.toLowerCase().includes(dataSearch.value.toLowerCase())
    })
  }
}

// 数据分页处理
const handleDataSizeChange = (val) => {
  dataQueryParams.pageSize = val
  fetchTaskData()
}

const handleDataCurrentChange = (val) => {
  dataQueryParams.page = val
  fetchTaskData()
}

// 获取状态样式类型
const getStatusType = (status) => {
  const map = {
    'waiting': 'info',
    'running': 'warning',
    'completed': 'success',
    'failed': 'danger',
    'cancelled': ''
  }
  return map[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    'waiting': '等待中',
    'running': '运行中',
    'completed': '已完成',
    'failed': '失败',
    'cancelled': '已取消'
  }
  return map[status] || '未知'
}

// 获取进度条状态
const getProgressStatus = (task) => {
  if (task.status === 'failed') return 'exception'
  if (task.status === 'completed') return 'success'
  return ''
}

// 进度条格式化
const progressFormat = (percentage) => {
  return percentage === 100 ? '完成' : `${percentage}%`
}

// 格式化下载内容
const formatDownloadContent = (content) => {
  if (!content) return '无'
  
  const map = {
    'video': '视频',
    'cover': '封面',
    'music': '音乐'
  }
  
  return content.map(item => map[item] || item).join('，')
}

// 格式化视频质量
const formatVideoQuality = (quality) => {
  const map = {
    'high': '高清（HD）',
    'medium': '标清（SD）',
    'low': '流畅'
  }
  return map[quality] || quality
}

// 格式化URL列表
const formatUrls = (urls) => {
  if (!urls) return []
  return urls.split('\n').filter(url => url.trim())
}

// 格式化数据类型
const formatDataType = (type) => {
  const map = {
    'video': '视频',
    'image': '图片',
    'music': '音乐',
    'text': '文本',
    'cover': '封面'
  }
  return map[type] || type
}

// 获取数据类型标签样式
const getDataTypeTag = (type) => {
  const map = {
    'video': 'danger',
    'image': 'success',
    'music': 'warning',
    'text': 'info',
    'cover': 'primary'
  }
  return map[type] || ''
}

// 格式化数据状态
const formatDataStatus = (status) => {
  const map = {
    'pending': '待处理',
    'downloading': '下载中',
    'completed': '已完成',
    'failed': '失败'
  }
  return map[status] || status
}

// 获取数据状态标签样式
const getDataStatusTag = (status) => {
  const map = {
    'pending': 'info',
    'downloading': 'warning',
    'completed': 'success',
    'failed': 'danger'
  }
  return map[status] || ''
}

// 格式化文件大小
const formatFileSize = (size) => {
  if (!size) return '0 B'
  
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  while (size >= 1024 && i < units.length - 1) {
    size /= 1024
    i++
  }
  
  return size.toFixed(2) + ' ' + units[i]
}

// 是否可以下载
const canDownload = computed(() => {
  return taskDetail.value && 
         taskDetail.value.status === 'completed' && 
         taskDetail.value.progress === 100
})

// 是否可以下载单个项目
const canDownloadItem = (item) => {
  return item.status === 'completed'
}

// 是否可以预览
const canPreview = (item) => {
  return item.status === 'completed' && (item.type === 'video' || item.type === 'image')
}

// 返回列表
const goBack = () => {
  router.push('/crawler')
}

// 下载所有文件
const handleDownload = () => {
  downloadLoading.value = true
  
  downloadTask(taskId).then(res => {
    ElMessage.success('下载任务已开始，文件将被保存到您的下载目录')
    downloadLoading.value = false
  }).catch(() => {
    downloadLoading.value = false
  })
}

// 下载单个项目
const handleDownloadItem = (item) => {
  downloadTaskItem(taskId, item.id).then(() => {
    ElMessage.success('开始下载文件: ' + item.title)
  })
}

// 预览
const handlePreview = (item) => {
  if (item.type === 'video') {
    // 实现视频预览逻辑
    window.open(item.url, '_blank')
  } else if (item.type === 'image') {
    // 图片直接用el-image组件预览，这里可以不做处理
  }
}

onMounted(() => {
  fetchTaskDetail()
})
</script>

<style lang="scss" scoped>
.app-container {
  padding: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.info-section {
  margin-bottom: 20px;
}

.crawler-data {
  margin-top: 30px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  
  h3 {
    margin: 0;
  }
}

.pagination-container {
  text-align: right;
  margin-top: 20px;
}

.data-preview {
  position: relative;
  
  .preview-icon {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 24px;
    color: white;
    background-color: rgba(0, 0, 0, 0.5);
    border-radius: 50%;
    padding: 8px;
  }
}

.file-icon {
  font-size: 40px;
  color: #909399;
}

.data-detail {
  padding: 16px;
  background-color: #f8f8f8;
  border-radius: 4px;
  
  .metadata {
    line-height: 1.8;
  }
}

.url-list {
  margin-top: 10px;
}

.url-item {
  margin-bottom: 8px;
}

.empty-data {
  padding: 40px 0;
}
</style> 