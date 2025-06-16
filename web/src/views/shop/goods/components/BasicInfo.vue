<template>
  <div class="basic-info">
    <el-form-item label="商品名称" prop="goodsName">
      <el-input v-model="formData.goodsName" placeholder="请输入商品名称" />
    </el-form-item>
    
    <el-form-item label="商品描述" prop="description">
      <el-input v-model="formData.description" type="textarea" placeholder="请输入商品描述" />
    </el-form-item>

    <el-form-item label="商品分类" prop="categoryId">
      <el-tree-select
        v-model="formData.categoryId"
        :data="categoryOptions"
        :props="{
          value: 'id',
          label: 'name',
          children: 'children'
        }"
        placeholder="请选择商品分类"
        clearable
        check-strictly
      />
    </el-form-item>

    <el-form-item label="商品品牌" prop="brand">
      <el-input v-model="formData.brand" placeholder="请输入商品品牌" />
    </el-form-item>

    <el-form-item label="供应商" prop="supplier">
      <el-input v-model="formData.supplier" placeholder="请输入供应商名称" />
    </el-form-item>

    <el-form-item label="商品类型" prop="type">
      <el-radio-group v-model="formData.type" @change="handleTypeChange">
        <el-radio :value="1">普通商品</el-radio>
        <el-radio :value="2">虚拟商品</el-radio>
        
      </el-radio-group>
    </el-form-item>

    <!-- 价格库存设置提示 -->
    <el-form-item label="价格库存">
      <el-alert
        title="商品的价格和库存请在规格设置页面中配置"
        type="info"
        :closable="false"
        show-icon
      />
    </el-form-item>

    <el-form-item label="排序" prop="sort">
      <el-input-number
        v-model="formData.sort"
        :min="0"
        :precision="0"
      />
    </el-form-item>

    <el-form-item label="商品状态" prop="status">
      <el-radio-group v-model="formData.status">
        <el-radio :value="1">上架</el-radio>
        <el-radio :value="0">下架</el-radio>
      </el-radio-group>
    </el-form-item>

    <el-form-item label="商品标签" prop="tagNames">
      <el-select 
        v-model="formData.tagNames" 
        multiple 
        filterable 
        allow-create 
        default-first-option 
        placeholder="请选择或输入商品标签"
        @change="handleTagChange"
      >
        <el-option
          v-for="item in tagOptions"
          :key="item.name"
          :label="item.name"
          :value="item.name"
        />
      </el-select>
      <div class="form-tip">可以选择已有标签，也可以输入新标签名称</div>
    </el-form-item>
  </div>
</template>

<script setup>
const emit = defineEmits(['type-change'])

defineProps({
  formData: {
    type: Object,
    required: true
  },
  categoryOptions: {
    type: Array,
    default: () => []
  },
  tagOptions: {
    type: Array,
    default: () => []
  }
})

// 商品类型变化处理
const handleTypeChange = (value) => {
  emit('type-change', value)
}

// 标签变化处理
const handleTagChange = (value) => {
  // 可以在这里添加标签变化的处理逻辑
  console.log('标签变化:', value)
}
</script>

<style lang="scss" scoped>
.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>

<style lang="scss" scoped>
:deep(.el-alert) {
  margin: 0;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style> 