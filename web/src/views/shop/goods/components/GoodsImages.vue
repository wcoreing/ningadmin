<template>
  <div class="goods-images">
    <!-- 商品封面 -->
    <el-form-item label="商品封面" prop="coverImg" required>
      <div class="image-list">
        <!-- 显示主图 -->
        <div
          v-if="mainImage"
          class="image-item"
        >
          <el-image :src="mainImage.url" fit="cover" />
          <div class="image-actions">
            <el-button
              circle
              type="danger"
              @click="handleRemoveImage('cover')"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
        <div v-if="!mainImage" class="image-upload" @click="openFileManager('cover')">
          <el-icon><Plus /></el-icon>
        </div>
      </div>
      <div class="el-form-item__tip">建议尺寸：800x800px，支持jpg、png格式</div>
    </el-form-item>

    <!-- 商品图片 -->
    <el-form-item label="商品图片" prop="images">
      <div class="image-list">
        <div
          v-for="(img, index) in normalImages"
          :key="index"
          class="image-item"
        >
          <el-image :src="img.url" fit="cover" />
          <div class="image-actions">
            <el-button
              circle
              type="danger"
              @click="handleRemoveImage('images', index)"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
        <div class="image-upload" @click="openFileManager('images')">
          <el-icon><Plus /></el-icon>
        </div>
      </div>
    </el-form-item>

    <!-- 附件管理器 -->
    <file-manager
      v-model="fileManagerVisible"
      :multiple="currentSelectType === 'images'"
      @select="handleImageSelect"
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Plus, Delete } from '@element-plus/icons-vue'
import FileManager from '@/components/FileManager/index.vue'

const props = defineProps({
  formData: {
    type: Object,
    required: true
  }
})

const fileManagerVisible = ref(false)
const currentSelectType = ref('')

// 图片位置常量
const IMAGE_POSITIONS = {
  NORMAL: 0,
  MAIN: 1,
  LIST: 2,
  DETAIL: 3,
  SHARE: 4
}

// 计算主图
const mainImage = computed(() => {
  if (props.formData.coverImg) {
    return {
      url: props.formData.coverImg,
      position: IMAGE_POSITIONS.MAIN
    }
  }
  return null
})

// 计算普通图片
const normalImages = computed(() => {
  // 确保 images 数组存在
  if (!props.formData.images) {
    props.formData.images = []
  }
  // 打印调试信息
  console.log('商品图片数据:', props.formData.images)
  // 只返回 position 为 0 的普通图片
  return props.formData.images.filter(img => img.position === IMAGE_POSITIONS.NORMAL)
})

const openFileManager = (type) => {
  currentSelectType.value = type
  fileManagerVisible.value = true
}

const handleImageSelect = (files) => {
  const fileData = Array.isArray(files) ? files : [files]
  
  if (!props.formData.images) {
    props.formData.images = []
  }

  if (currentSelectType.value === 'cover') {
    // 设置封面图
    if (fileData.length > 0) {
      props.formData.coverImg = fileData[0].url
    }
  } else {
    // 添加普通图片，position 设置为 0
    const newImages = fileData.map(file => ({
      url: file.url,
      position: IMAGE_POSITIONS.NORMAL, // 这里 NORMAL 是 0
      sort: props.formData.images.length
    }))
    props.formData.images.push(...newImages)
  }
  
  fileManagerVisible.value = false
}

const handleRemoveImage = (type, index) => {
  if (!props.formData.images) {
    props.formData.images = []
    return
  }

  if (type === 'cover') {
    props.formData.coverImg = '' // 清除封面图
  } else {
    props.formData.images.splice(index, 1)
  }
}
</script>

<style lang="scss" scoped>
.image-list {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;

  .image-item {
    position: relative;
    width: 120px;
    height: 120px;
    border: 1px solid var(--el-border-color);
    border-radius: 4px;
    overflow: hidden;

    .el-image {
      width: 100%;
      height: 100%;
    }

    .image-actions {
      position: absolute;
      top: 8px;
      right: 8px;
      opacity: 0;
      transition: opacity 0.2s;
    }

    &:hover .image-actions {
      opacity: 1;
    }
  }

  .image-upload {
    width: 120px;
    height: 120px;
    border: 1px dashed var(--el-border-color);
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: border-color 0.2s;
    color: var(--el-text-color-secondary);

    &:hover {
      border-color: var(--el-color-primary);
      color: var(--el-color-primary);
    }

    .el-icon {
      font-size: 24px;
    }
  }
}
</style> 