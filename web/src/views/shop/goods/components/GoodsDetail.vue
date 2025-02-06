<template>
  <div class="goods-detail">
    <el-form-item label="商品详情">
      <div class="editor-wrapper">
        <Toolbar
          class="editor-toolbar"
          :editor="editorRef"
          :defaultConfig="toolbarConfig"
          mode="default"
        />
        <Editor
          class="editor-content"
          v-model="formData.detailDesc"
          :defaultConfig="editorConfig"
          mode="default"
          @onCreated="handleCreated"
        />
      </div>
    </el-form-item>
  </div>
</template>

<script setup>
import { ref, shallowRef, onBeforeUnmount } from 'vue'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  formData: {
    type: Object,
    required: true
  }
})

// 编辑器实例，必须用 shallowRef
const editorRef = shallowRef()

// 工具栏配置
const toolbarConfig = {}

// 编辑器配置
const editorConfig = {
  placeholder: '请输入商品详情...',
  // 基础配置
  autoFocus: false,
  scroll: false,
  maxLength: 99999,
  MENU_CONF: {
    uploadImage: {
      // 上传配置
      allowedFileTypes: ['image/*'],
      maxFileSize: 10 * 1024 * 1024,
      base64LimitSize: 5 * 1024,
      withCredentials: false,
      timeout: 10 * 1000,
      // 自定义上传
      customUpload(file) {
        const formData = new FormData()
        formData.append('file', file)
        
        fetch(import.meta.env.VITE_API_BASE_URL + '/files/upload', {
          method: 'POST',
          headers: {
            Authorization: 'Bearer ' + localStorage.getItem('token')
          },
          body: formData
        })
        .then(res => res.json())
        .then(res => {
          if (res.code === 200) {
            const imageUrl = import.meta.env.VITE_IMG_URL + res.data.url
            console.log('上传成功:', imageUrl)
            // 使用编辑器实例插入图片
            const editor = editorRef.value
            editor.dangerouslyInsertHtml(`<img src="${imageUrl}" alt="${file.name}">`)
          } else {
            ElMessage.error(res.msg || '上传失败')
          }
        })
      }
    }
  }
}

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()
})

const handleCreated = (editor) => {
  editorRef.value = editor
  // 设置默认内容
  if (!props.formData.detailDesc) {
    editor.setHtml('<p><br></p>')
  }
}
</script>

<style lang="scss" scoped>
.editor-wrapper {
  position: relative;
  border: 1px solid #ccc;
  width: 100%;
}

.editor-toolbar {
  position: sticky;
  top: 0;
  border-bottom: 1px solid #ccc;
  background: #fff;
  z-index: 101;
}

.editor-content {
  height: 800px !important;
  overflow-y: auto;
}

:deep(.w-e-text-container) {
  background-color: #fff;
  height: 800px !important;
}

:deep(.w-e-text-container[data-mode=default]) {
  height: 800px !important;
}
</style> 