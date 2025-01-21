import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ArrowDown } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

// 处理下拉菜单命令
const handleCommand = async (command) => {
  if (command === 'logout') {
    try {
      const success = await userStore.logoutFn()
      if (success) {
        router.push('/login')
      }
    } catch (error) {
      console.error('退出登录失败:', error)
    }
  } else if (command === 'profile') {
    // 处理个人中心
    router.push('/profile')
  }
} 