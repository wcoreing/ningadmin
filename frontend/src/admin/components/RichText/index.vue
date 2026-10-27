<template>
  <div class="rich-text">
    <Toolbar
      class="rich-text-bar"
      :editor="editorRef"
      :defaultConfig="toolbarConfig"
      mode="default"
    />
    <Editor
      class="rich-text-body"
      :style="{ height: height + 'px' }"
      v-model="html"
      :defaultConfig="editorConfig"
      mode="default"
      @onCreated="handleCreated"
    />
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, shallowRef } from 'vue'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { uploadFile } from '@/admin/api/common/upload'

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '请输入正文' },
  height: { type: Number, default: 420 }
})

const emit = defineEmits(['update:modelValue'])

const editorRef = shallowRef()
const toolbarConfig = {}
const html = computed({
  get: () => props.modelValue || '',
  set: value => emit('update:modelValue', value)
})

const editorConfig = {
  placeholder: props.placeholder,
  autoFocus: false,
  scroll: true,
  MENU_CONF: {
    uploadImage: {
      allowedFileTypes: ['image/*'],
      maxFileSize: 10 * 1024 * 1024,
      base64LimitSize: 5 * 1024,
      withCredentials: false,
      timeout: 10 * 1000,
      async customUpload(file, insertFn) {
        const data = new FormData()
        data.append('file', file)
        const saved = await uploadFile(data)
        const editor = editorRef.value
        editor?.restoreSelection()
        if (editor && editor.selection == null) editor.focus(true)
        insertFn(saved.url, saved.originalName || file.name, saved.url)
      }
    }
  }
}

onBeforeUnmount(() => {
  editorRef.value?.destroy()
  editorRef.value = null
})

function handleCreated(editor) {
  editorRef.value = editor
  if (!props.modelValue) editor.setHtml('<p><br></p>')
}
</script>

<style scoped>
.rich-text {
  border: 1px solid #ccc;
  width: 100%;
  z-index: 2;
}
.rich-text-bar {
  border-bottom: 1px solid #ccc;
  background: #fff;
}
.rich-text-body {
  overflow-y: hidden;
}
:deep(.w-e-text-container) {
  background: #fff;
}
</style>
