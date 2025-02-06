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

    <el-form-item label="商品品牌" prop="brandId">
      <el-select v-model="formData.brandId" placeholder="请选择商品品牌">
        <el-option
          v-for="item in brandOptions"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
    </el-form-item>

    <el-form-item label="供应商" prop="supplierId">
      <el-select v-model="formData.supplierId" placeholder="请选择供应商">
        <el-option
          v-for="item in supplierOptions"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
    </el-form-item>

    <el-form-item label="销售价" prop="price">
      <el-input-number
        v-model="formData.price"
        :precision="2"
        :step="0.1"
        :min="0"
      />
    </el-form-item>

    <el-form-item label="市场价" prop="originalPrice">
      <el-input-number
        v-model="formData.originalPrice"
        :precision="2"
        :step="0.1"
        :min="0"
      />
    </el-form-item>

    <el-form-item label="商品库存" prop="stock">
      <el-input-number
        v-model="formData.stock"
        :min="0"
        :precision="0"
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
        <el-radio :label="'0'">上架</el-radio>
        <el-radio :label="'1'">下架</el-radio>
      </el-radio-group>
    </el-form-item>

    <el-form-item label="商品标签" prop="tagIds">
      <el-select v-model="formData.tagIds" multiple placeholder="请选择商品标签">
        <el-option
          v-for="item in tagOptions"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
    </el-form-item>
  </div>
</template>

<script setup>
defineProps({
  formData: {
    type: Object,
    required: true
  },
  categoryOptions: {
    type: Array,
    default: () => []
  },
  brandOptions: {
    type: Array,
    default: () => []
  },
  supplierOptions: {
    type: Array,
    default: () => []
  },
  tagOptions: {
    type: Array,
    default: () => []
  }
})
</script> 