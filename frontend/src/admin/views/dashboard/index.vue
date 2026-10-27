<template>
  <div class="dashboard">
    <!-- 顶部欢迎区域 -->
    <div class="welcome-section">
      <div class="welcome-content">
        <h1 class="title">
          <span class="text-gradient">集亿玖泽</span>
          <div class="subtitle">深圳市集亿玖泽贸易有限公司</div>
        </h1>
        <div class="description">
          <el-icon class="icon"><Star /></el-icon>
          <span>首饰包装与箱包，国内贸易和进出口</span>
        </div>
      </div>
      <div class="decoration">
        <div class="circle-1"></div>
        <div class="circle-2"></div>
        <div class="circle-3"></div>
      </div>
    </div>

    <!-- 统计卡片区域 -->
    <div class="stat-cards">
      <el-row :gutter="20">
        <el-col :span="6" v-for="(item, index) in statCards" :key="index">
          <el-card class="stat-card" :class="item.type" shadow="hover">
            <div class="card-content">
              <el-icon class="icon"><component :is="item.icon" /></el-icon>
              <div class="info">
                <div class="value">{{ item.value }}</div>
                <div class="label">{{ item.label }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 底部版权信息 -->
    <div class="footer">
      <div class="copyright">
        © {{ new Date().getFullYear() }} 深圳市集亿玖泽贸易有限公司
        <div class="powered-by">韦桂新 15989589813</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { Star, DataLine, Goods, User, Setting } from '@element-plus/icons-vue'

const statCards = ref([
  { label: '产品目录', value: '在售', icon: 'Goods', type: 'success' },
  { label: '内容资料', value: '可改', icon: 'Document', type: 'primary' },
  { label: '询盘记录', value: '跟进', icon: 'ChatDotRound', type: 'warning' },
  { label: '企业资料', value: '公司', icon: 'OfficeBuilding', type: 'danger' }
])
</script>

<style lang="scss" scoped>
.dashboard {
  padding: 24px;
  min-height: calc(100vh - 84px);
  background: #f5f7fa;
  position: relative;
  overflow: hidden;

  .welcome-section {
    position: relative;
    padding: 40px;
    margin-bottom: 40px;
    background: linear-gradient(135deg, #fff 0%, #f0f2f5 100%);
    border-radius: 16px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
    overflow: hidden;

    .welcome-content {
      position: relative;
      z-index: 2;

      .title {
        margin: 0;
        font-size: 36px;
        font-weight: 600;
        line-height: 1.2;

        .text-gradient {
          background: linear-gradient(45deg, #e1251b, #8d1510);
          -webkit-background-clip: text;
          -webkit-text-fill-color: transparent;
          animation: gradient 6s ease infinite;
        }

        .subtitle {
          font-size: 18px;
          color: #606266;
          margin-top: 12px;
          font-weight: normal;
        }
      }

      .description {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-top: 24px;
        color: #909399;
        font-size: 16px;

        .icon {
          color: #ffd04b;
          animation: rotate 4s linear infinite;
        }
      }
    }

    .decoration {
      position: absolute;
      top: 0;
      right: 0;
      bottom: 0;
      width: 40%;
      overflow: hidden;

      .circle-1, .circle-2, .circle-3 {
        position: absolute;
        border-radius: 50%;
        opacity: 0.1;
      }

      .circle-1 {
        width: 200px;
        height: 200px;
        background: #409eff;
        top: -50px;
        right: -50px;
        animation: float 6s ease-in-out infinite;
      }

      .circle-2 {
        width: 150px;
        height: 150px;
        background: #67c23a;
        bottom: -30px;
        right: 100px;
        animation: float 8s ease-in-out infinite;
      }

      .circle-3 {
        width: 100px;
        height: 100px;
        background: #e6a23c;
        top: 50%;
        right: 50px;
        animation: float 7s ease-in-out infinite;
      }
    }
  }

  .stat-cards {
    .stat-card {
      transition: transform 0.3s ease;
      
      &:hover {
        transform: translateY(-5px);
      }

      .card-content {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 20px;

        .icon {
          font-size: 48px;
          padding: 16px;
          border-radius: 12px;
          
          &.primary {
            color: #409eff;
            background: rgba(64, 158, 255, 0.1);
          }
          
          &.success {
            color: #67c23a;
            background: rgba(103, 194, 58, 0.1);
          }
          
          &.warning {
            color: #e6a23c;
            background: rgba(230, 162, 60, 0.1);
          }
          
          &.danger {
            color: #f56c6c;
            background: rgba(245, 108, 108, 0.1);
          }
        }

        .info {
          .value {
            font-size: 28px;
            font-weight: 600;
            line-height: 1;
            margin-bottom: 8px;
          }

          .label {
            color: #909399;
            font-size: 14px;
          }
        }
      }
    }
  }

  .footer {
    text-align: center;
    margin-top: 60px;
    padding: 24px 0;
    color: #909399;

    .copyright {
      font-size: 14px;
      
      .powered-by {
        margin-top: 8px;
        font-size: 12px;
        opacity: 0.8;
      }
    }
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-20px);
  }
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@keyframes gradient {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}
</style> 