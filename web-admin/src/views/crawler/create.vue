<template>
  <div class="crawler-task-create">
    <el-form
      ref="taskFormRef"
      :model="taskForm"
      :rules="rules"
      label-width="100px"
      label-position="right"
    >
      <!-- 基本信息部分 -->
      <el-form-item label="任务名称" prop="name">
        <el-input v-model="taskForm.name" placeholder="请输入任务名称" />
      </el-form-item>

      <el-form-item label="爬虫类型" prop="type">
        <el-select v-model="taskForm.type" placeholder="请选择爬虫类型" style="width: 100%">
          <el-option
            v-for="item in crawlerTypes"
            :key="item.code"
            :label="item.name"
            :value="item.code"
          />
        </el-select>
      </el-form-item>

      <!-- 抖音视频爬虫配置 -->
      <template v-if="taskForm.type === 'douyin_video'">
        <el-divider content-position="left">抖音视频配置</el-divider>
        
        <el-form-item label="视频链接" prop="config.urls">
          <el-tabs v-model="inputMode">
            <el-tab-pane label="直接输入" name="direct">
              <el-input
                v-model="taskForm.config.urls"
                type="textarea"
                :rows="5"
                placeholder="输入抖音视频链接，一行一个链接"
              />
            </el-tab-pane>
            <el-tab-pane label="批量上传" name="file">
              <el-upload
                class="link-file-uploader"
                action="#"
                :auto-upload="false"
                :on-change="handleFileChange"
                :limit="1"
              >
                <el-button type="primary">选择文本文件</el-button>
                <template #tip>
                  <div class="el-upload__tip">
                    请上传文本文件，每行一个抖音链接
                  </div>
                </template>
              </el-upload>
            </el-tab-pane>
          </el-tabs>
        </el-form-item>

        <el-form-item label="作者主页" prop="config.authorUrl">
          <el-input 
            v-model="taskForm.config.authorUrl" 
            placeholder="请输入抖音作者主页链接（可选，将爬取该作者的所有视频）" 
          />
        </el-form-item>

        <el-form-item label="下载内容" prop="config.downloadContent">
          <el-checkbox-group v-model="taskForm.config.downloadContent">
            <el-checkbox label="video">视频</el-checkbox>
            <el-checkbox label="cover">封面</el-checkbox>
            <el-checkbox label="music">音乐</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="视频质量" prop="config.videoQuality">
          <el-radio-group v-model="taskForm.config.videoQuality">
            <el-radio value="high">高清（HD）</el-radio>
            <el-radio value="medium">标清（SD）</el-radio>
            <el-radio value="low">流畅</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="抓取数量" prop="config.maxVideos">
          <el-input-number 
            v-model="taskForm.config.maxVideos" 
            :min="1" 
            :max="1000" 
            placeholder="最大抓取视频数量" 
          />
          <span class="form-item-tip">（0表示不限制）</span>
        </el-form-item>
      </template>

      <!-- B站视频爬虫配置 -->
      <template v-else-if="taskForm.type === 'bilibili_video'">
        <el-divider content-position="left">B站视频配置</el-divider>
        <!-- B站特定配置 -->
        <!-- 这里是B站爬虫的特定配置，暂时不实现 -->
        <div class="placeholder-text">B站视频爬虫配置（功能开发中）</div>
      </template>

      <!-- 微博图片爬虫配置 -->
      <template v-else-if="taskForm.type === 'weibo_image'">
        <el-divider content-position="left">微博图片配置</el-divider>
        <!-- 微博特定配置 -->
        <!-- 这里是微博爬虫的特定配置，暂时不实现 -->
        <div class="placeholder-text">微博图片爬虫配置（功能开发中）</div>
      </template>

      <!-- 高级选项（所有爬虫通用） -->
      <el-divider content-position="left">高级选项</el-divider>

      <el-form-item label="代理设置" prop="config.useProxy">
        <el-switch v-model="taskForm.config.useProxy" />
        <span class="form-item-tip">（启用代理可能会提高成功率，但会降低爬取速度）</span>
      </el-form-item>

      <el-form-item label="超时设置" prop="config.timeout">
        <el-input-number
          v-model="taskForm.config.timeout"
          :min="5"
          :max="120"
          :step="5"
        />
        <span class="form-item-tip">秒</span>
      </el-form-item>

      <el-form-item label="重试次数" prop="config.retries">
        <el-input-number
          v-model="taskForm.config.retries"
          :min="0"
          :max="10"
          :step="1"
        />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="submitForm">创建任务</el-button>
        <el-button @click="cancelCreate">取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { createTask } from '@/api/crawler/task'
import { getCrawlerTypes } from '@/api/crawler/type'

const emit = defineEmits(['created', 'cancel'])

const taskFormRef = ref(null)
const inputMode = ref('direct')

// 爬虫类型
const crawlerTypes = ref([])

// 表单数据
const taskForm = reactive({
  name: '',
  type: 'douyin_video',
  config: {
    urls: '',
    authorUrl: '',
    downloadContent: ['video'],
    videoQuality: 'high',
    maxVideos: 100,
    useProxy: false,
    timeout: 30,
    retries: 3
  }
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择爬虫类型', trigger: 'change' }
  ],
  'config.urls': [
    { 
      validator: (rule, value, callback) => {
        // 如果作者主页填写了，则链接可以为空
        if (!value && !taskForm.config.authorUrl) {
          callback(new Error('请输入视频链接或作者主页'))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ],
  'config.authorUrl': [
    { 
      validator: (rule, value, callback) => {
        // 如果链接填写了，则作者主页可以为空
        if (!value && !taskForm.config.urls) {
          callback(new Error('请输入视频链接或作者主页'))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ]
}

// 处理文件上传改变
const handleFileChange = (file) => {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = (e) => {
      taskForm.config.urls = e.target.result
    }
    reader.readAsText(file.raw)
  }
}

// 提交表单
const submitForm = async () => {
  if (!taskFormRef.value) return
  
  await taskFormRef.value.validate((valid) => {
    if (valid) {
      createTask(taskForm).then(res => {
        console.log('创建的任务数据:', res)
        emit('created')
      }).catch(err => {
        ElMessage.error(err.message || '创建任务失败')
      })
    }
  })
}

// 取消创建
const cancelCreate = () => {
  emit('cancel')
}

// 获取爬虫类型列表
const fetchCrawlerTypes = () => {
  getCrawlerTypes().then(res => {
    crawlerTypes.value = res || []
    // 如果有类型，默认选第一个
    if (crawlerTypes.value.length > 0) {
      taskForm.type = crawlerTypes.value[0].code
    }
  })
}

onMounted(() => {
  fetchCrawlerTypes()
})
</script>

<style lang="scss" scoped>
.crawler-task-create {
  padding: 10px;
}

.form-item-tip {
  margin-left: 10px;
  color: #909399;
  font-size: 12px;
}

.link-file-uploader {
  width: 100%;
}

.placeholder-text {
  color: #909399;
  font-style: italic;
  padding: 20px 0;
  text-align: center;
  border: 1px dashed #d9d9d9;
  border-radius: 4px;
}
</style> 