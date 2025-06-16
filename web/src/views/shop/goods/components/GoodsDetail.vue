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
import { ref, shallowRef, onBeforeUnmount, watch } from 'vue'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  formData: {
    type: Object,
    required: true
  }
})

// 处理 HTML 内容，确保所有图片都有最大宽度属性
const processHtml = (html) => {
  const div = document.createElement('div')
  div.innerHTML = html
  
  // 处理所有图片
  div.querySelectorAll('img').forEach(img => {
    if (!img.style.maxWidth) {
      img.style.maxWidth = '100%'
    }
  })
  
  return div.innerHTML
}

// 监听 detailDesc 的变化
watch(() => props.formData.detailDesc, (newVal) => {
  if (newVal) {
    // 处理内容并更新
    props.formData.detailDesc = processHtml(newVal)
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
  // 添加 parseHtml 配置
  parseHtml: (elem) => {
    // 处理图片元素
    if (elem.tagName === 'IMG') {
      elem.style.maxWidth = '100%'
    }
    return elem
  },
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
            // res.data 在这里是正确的，因为这是直接的 fetch 响应，不是通过 request.js
            const imageUrl = import.meta.env.VITE_IMG_URL + res.data.url
            console.log('上传成功:', imageUrl)
            // 修改插入图片的方式
            const editor = editorRef.value
            editor.dangerouslyInsertHtml(`<img src="${imageUrl}" alt="${file.name}" style="max-width:100%">`)
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