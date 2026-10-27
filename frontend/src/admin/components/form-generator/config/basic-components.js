import * as ElementPlusIconsVue from '@element-plus/icons-vue'

export const basicComponents = [
  { 
    type: 'el-input',
    label: '单行文本',
    icon: ElementPlusIconsVue.Edit,
    props: { label: '单行文本' }
  },
  {
    type: 'el-date-picker',
    label: '日期选择',
    icon: ElementPlusIconsVue.Calendar,
    props: {
      label: '日期选择',
      type: 'date',
      placeholder: '请选择日期时间',
      format: 'YYYY-MM-DD',
      valueFormat: 'YYYY-MM-DD',
      editable: false,
      clearable: true,
      disabled: false,
      shortcuts: [
        {
          text: '今天',
          value: new Date(),
        },
        {
          text: '昨天',
          value: () => {
            const date = new Date()
            date.setTime(date.getTime() - 3600 * 1000 * 24)
            return date
          },
        },
        {
          text: '一周前',
          value: () => {
            const date = new Date()
            date.setTime(date.getTime() - 3600 * 1000 * 24 * 7)
            return date
          },
        },
      ],
      style: {
        width: '100%'
      }
    }
  },
  {
    type: 'el-input',
    label: '多行文本',
    icon: ElementPlusIconsVue.Document,
    props: { 
      label: '多行文本', 
      type: 'textarea' 
    }
  },
  {
    type: 'el-select',
    label: '下拉选择',
    icon: ElementPlusIconsVue.Select,
    props: {
      label: '下拉选择',
      placeholder: '请选择',
      options: [
        { label: '选项1', value: '1' },
        { label: '选项2', value: '2' }
      ]
    }
  },
  {
    type: 'el-checkbox-group',
    label: '多选框组',
    icon: ElementPlusIconsVue.Check,
    props: {
      label: '多选框组',
      options: [
        { label: '选项1', value: '1' },
        { label: '选项2', value: '2' },
        { label: '选项3', value: '3' }
      ]
    }
  },
  {
    type: 'el-switch',
    component: 'SwitchField',
    label: '开关',
    icon: ElementPlusIconsVue.Switch,
    props: { label: '开关' }
  },
  {
    type: 'el-input-number',
    component: 'NumberField',
    label: '数字输入',
    icon: ElementPlusIconsVue.Plus,
    props: {
      label: '数字输入',
      min: 0,
      max: 100
    }
  },
  {
    type: 'el-upload',
    component: 'UploadField',
    label: '图片上传',
    icon: ElementPlusIconsVue.Upload,
    props: {
      label: '图片上传',
      action: '/api/upload',
      listType: 'picture-card',
      limit: 5,
      multiple: true,
      accept: 'image/*',
      autoUpload: false,
      showFileList: true,
      beforeUpload: (file) => {
        const isImage = file.type.startsWith('image/')
        const isLt2M = file.size / 1024 / 1024 < 2
        
        if (!isImage) {
          ElMessage.error('只能上传图片文件!')
          return false
        }
        if (!isLt2M) {
          ElMessage.error('图片大小不能超过 2MB!')
          return false
        }
        
        const reader = new FileReader()
        reader.readAsDataURL(file)
        return false
      }
    }
  }
] 