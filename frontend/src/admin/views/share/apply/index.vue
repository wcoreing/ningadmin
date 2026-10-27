<template>
  <div class="app-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :model="queryParams" ref="queryForm" :inline="true">
        <el-form-item label="申请人" prop="realName">
          <el-input
            v-model="queryParams.realName"
            placeholder="请输入申请人姓名"
            clearable
            @keyup.enter.native="handleQuery"
          />
        </el-form-item>
        <el-form-item label="手机号码" prop="phone">
          <el-input
            v-model="queryParams.phone"
            placeholder="请输入手机号码"
            clearable
            @keyup.enter.native="handleQuery"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="queryParams.status" placeholder="请选择状态" clearable>
            <el-option label="待审核" value="0" />
            <el-option label="已通过" value="1" />
            <el-option label="已拒绝" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="el-icon-search" @click="handleQuery">搜索</el-button>
          <el-button icon="el-icon-refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 列表区域 -->
    <el-card class="table-card">
      <el-table
        v-loading="loading"
        :data="applyList"
        style="width: 100%"
      >
        <el-table-column label="申请人" prop="realName" />
        <el-table-column label="手机号码" prop="phone" />
        <el-table-column label="身份证号" prop="idCard" />
        <el-table-column label="申请时间" prop="applyTime" width="160" />
        <el-table-column label="状态" prop="status" width="100">
            <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]">
              {{ statusMap[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
            <el-button
              v-if="row.status === 0"
              size="mini"
              type="success"
              @click="handleApprove(row)"
            >通过</el-button>
            <el-button
              v-if="row.status === 0"
              size="mini"
              type="danger"
              @click="handleReject(row)"
            >拒绝</el-button>
            <el-button
              size="mini"
              type="text"
              @click="handleDetail(row)"
            >详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <pagination
        v-show="total > 0"
        :total="total"
        :page.sync="queryParams.pageNum"
        :limit.sync="queryParams.pageSize"
        @pagination="getList"
      />
    </el-card>

    <!-- 拒绝对话框 -->
    <el-dialog title="拒绝原因" :visible.sync="rejectDialogVisible" width="500px">
      <el-form :model="rejectForm" ref="rejectForm" :rules="rejectRules">
        <el-form-item label="拒绝原因" prop="reason" label-width="80px">
          <el-input
            type="textarea"
            v-model="rejectForm.reason"
            placeholder="请输入拒绝原因"
            :rows="3"
          />
        </el-form-item>
      </el-form>
      <div slot="footer">
        <el-button @click="rejectDialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="confirmReject">确 定</el-button>
      </div>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog 
      title="申请详情" 
      v-model="detailDialogVisible"
      width="600px" 
      append-to-body
      :before-close="handleDetailClose"
    >
      <div v-loading="detailLoading">
        <el-descriptions :column="2" border size="medium" v-if="detailData">
          <el-descriptions-item label="申请人">{{ detailData.realName }}</el-descriptions-item>
          <el-descriptions-item label="手机号码">{{ detailData.phone }}</el-descriptions-item>
          <el-descriptions-item label="身份证号">{{ detailData.idCard }}</el-descriptions-item>
          <el-descriptions-item label="申请时间">{{ detailData.applyTime }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusTypeMap[detailData.status]">
              {{ statusMap[detailData.status] }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="审核时间">{{ detailData.auditTime || '-' }}</el-descriptions-item>
          <el-descriptions-item label="审核人">{{ detailData.auditUser || '-' }}</el-descriptions-item>
          <el-descriptions-item label="拒绝原因" :span="2">
            <el-input
              type="textarea"
              v-model="detailData.rejectReason"
              :rows="2"
              readonly
              placeholder="-"
            />
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <div slot="footer">
        <el-button @click="handleDetailClose">关 闭</el-button>
        <template v-if="detailData && detailData.status === 0">
          <el-button type="success" @click="handleApprove(detailData)">通 过</el-button>
          <el-button type="danger" @click="handleReject(detailData)">拒 绝</el-button>
        </template>
      </div>
    </el-dialog>
  </div>
</template>

<script>
// 导入API请求方法
import { 
  listShareApply, 
  approveShareApply, 
  rejectShareApply,
  getShareApplyDetail
} from '@/admin/api/share/apply'

export default {
  name: 'DistributionApply',
  data() {
    return {
      // 遮罩层
      loading: false,
      // 总条数
      total: 0,
      // 申请列表
      applyList: [],
      // 查询参数
      queryParams: {
        pageNum: 1,
        pageSize: 10,
        realName: undefined,
        phone: undefined,
        status: undefined
      },
      // 状态映射
      statusMap: {
        '0': '待审核',
        '1': '已通过',
        '2': '已拒绝'
      },
      // 状态样式映射
      statusTypeMap: {
        '0': 'warning',
        '1': 'success',
        '2': 'danger'
      },
      // 拒绝对话框
      rejectDialogVisible: false,
      // 拒绝表单
      rejectForm: {
        id: undefined,
        reason: ''
      },
      // 拒绝表单校验规则
      rejectRules: {
        reason: [
          { required: true, message: '请输入拒绝原因', trigger: 'blur' }
        ]
      },
      // 详情对话框
      detailDialogVisible: false,
      detailLoading: false,
      // 详情数据
      detailData: null
    }
  },
  created() {
    this.getList()
  },
  methods: {
    // 获取列表数据
    getList() {
      this.loading = true
      listShareApply(this.queryParams).then(response => {
        console.log('API Response:', response)
        this.applyList = response.list || []
        this.total = response.total || 0
        this.loading = false
      }).catch(() => {
        this.loading = false
      })
    },
    // 搜索按钮操作
    handleQuery() {
      this.queryParams.pageNum = 1
      this.getList()
    },
    // 重置按钮操作
    resetQuery() {
      this.resetForm('queryForm')
      this.handleQuery()
    },
    // 通过操作
    handleApprove(row) {
      this.$confirm('确认通过该申请吗？', '提示', {
        type: 'warning'
      }).then(() => {
        approveShareApply(row.id).then(response => {
          
            this.$message.success('操作成功')
            this.getList()
          
        })
      })
    },
    // 拒绝操作
    handleReject(row) {
      this.rejectForm.id = row.id
      this.rejectForm.reason = ''
      this.rejectDialogVisible = true
    },
    // 确认拒绝
    confirmReject() {
      this.$refs.rejectForm.validate(valid => {
        if (valid) {
          rejectShareApply(this.rejectForm).then(response => {
        
              this.$message.success('操作成功')
              this.rejectDialogVisible = false
              this.getList()
            
          })
        }
      })
    },
    // 查看详情
    handleDetail(row) {
      console.log('Row data:', row)
      this.detailData = null
      this.detailDialogVisible = true
      this.detailLoading = true
      
      getShareApplyDetail(row.id).then(response => {
        console.log('Detail Response:', response)
      
          this.detailData = response
          console.log('Detail data set:', this.detailData)
       
      }).catch(error => {
        console.error('Detail error:', error)
        this.$message.error('获取详情失败')
      }).finally(() => {
        this.detailLoading = false
      })
    },
    
    // 关闭详情弹窗
    handleDetailClose() {
      this.detailDialogVisible = false
      this.detailData = null
    }
  }
}
</script>

<style scoped>
.search-card {
  margin-bottom: 15px;
}
.table-card {
  margin-bottom: 15px;
}
.el-descriptions {
  margin: 20px 0;
}
.el-descriptions-item__label {
  width: 100px;
  text-align: right;
}
</style> 