// 任务状态枚举
const TaskStatus = {
  WAITING: 'waiting',
  RUNNING: 'running',
  COMPLETED: 'completed',
  FAILED: 'failed',
  CANCELLED: 'cancelled'
}

// 数据状态枚举
const DataStatus = {
  PENDING: 'pending',
  DOWNLOADING: 'downloading',
  COMPLETED: 'completed',
  FAILED: 'failed'
}

// 爬虫类型列表
const crawlerTypes = [
  { code: 'douyin_video', name: '抖音视频' },
  { code: 'bilibili_video', name: 'B站视频' },
  { code: 'weibo_image', name: '微博图片' }
]

// 模拟任务列表数据
const mockTaskList = [
  {
    id: 1,
    name: '抖音热门视频抓取',
    type: 'douyin_video',
    typeName: '抖音视频',
    status: TaskStatus.COMPLETED,
    progress: 100,
    createdAt: '2023-10-01 12:30:45',
    updatedAt: '2023-10-01 13:20:10',
    endTime: '2023-10-01 13:20:10'
  },
  {
    id: 2,
    name: 'B站舞蹈区视频',
    type: 'bilibili_video',
    typeName: 'B站视频',
    status: TaskStatus.RUNNING,
    progress: 45,
    createdAt: '2023-10-02 09:15:22',
    updatedAt: '2023-10-02 09:35:30'
  },
  {
    id: 3,
    name: '微博热门图片',
    type: 'weibo_image',
    typeName: '微博图片',
    status: TaskStatus.WAITING,
    progress: 0,
    createdAt: '2023-10-02 14:50:11',
    updatedAt: '2023-10-02 14:50:11'
  },
  {
    id: 4,
    name: '抖音账号视频备份',
    type: 'douyin_video',
    typeName: '抖音视频',
    status: TaskStatus.FAILED,
    progress: 67,
    createdAt: '2023-09-28 18:22:33',
    updatedAt: '2023-09-28 19:10:05',
    endTime: '2023-09-28 19:10:05'
  }
]

// 模拟任务详情数据
const mockTaskDetails = {
  1: {
    id: 1,
    name: '抖音热门视频抓取',
    type: 'douyin_video',
    typeName: '抖音视频',
    status: TaskStatus.COMPLETED,
    progress: 100,
    createdAt: '2023-10-01 12:30:45',
    updatedAt: '2023-10-01 13:20:10',
    endTime: '2023-10-01 13:20:10',
    config: {
      urls: 'https://v.douyin.com/abcdef/\nhttps://v.douyin.com/123456/',
      authorUrl: 'https://v.douyin.com/author123/',
      downloadContent: ['video', 'cover'],
      videoQuality: 'high',
      maxVideos: 100,
      useProxy: true,
      timeout: 30,
      retries: 3
    }
  },
  2: {
    id: 2,
    name: 'B站舞蹈区视频',
    type: 'bilibili_video',
    typeName: 'B站视频',
    status: TaskStatus.RUNNING,
    progress: 45,
    createdAt: '2023-10-02 09:15:22',
    updatedAt: '2023-10-02 09:35:30',
    config: {
      urls: 'https://www.bilibili.com/video/BV123456\nhttps://www.bilibili.com/video/BV654321',
      downloadContent: ['video'],
      videoQuality: 'high',
      maxVideos: 50,
      useProxy: false,
      timeout: 60,
      retries: 2
    }
  }
}

// 模拟爬取数据
const mockCrawlerData = {
  1: [
    {
      id: 1,
      taskId: 1,
      title: '这是一个抖音视频标题1',
      type: 'video',
      originalUrl: 'https://v.douyin.com/abcdef/',
      url: 'https://example.com/videos/video1.mp4',
      coverUrl: 'https://example.com/covers/cover1.jpg',
      size: 15728640, // 15MB
      status: DataStatus.COMPLETED,
      metadata: {
        author: '抖音用户1',
        likes: 1200,
        comments: 300,
        shares: 150,
        description: '这是视频描述：非常有趣的一段视频，大家快来看'
      },
      createdAt: '2023-10-01 12:35:22'
    },
    {
      id: 2,
      taskId: 1,
      title: '这是一个抖音视频标题2',
      type: 'video',
      originalUrl: 'https://v.douyin.com/123456/',
      url: 'https://example.com/videos/video2.mp4',
      coverUrl: 'https://example.com/covers/cover2.jpg',
      size: 20971520, // 20MB
      status: DataStatus.COMPLETED,
      metadata: {
        author: '抖音用户2',
        likes: 5500,
        comments: 800,
        shares: 320,
        description: '这是视频描述：这段视频太精彩了，不敢相信自己的眼睛'
      },
      createdAt: '2023-10-01 12:42:33'
    },
    {
      id: 3,
      taskId: 1,
      title: '这是一个抖音视频封面1',
      type: 'cover',
      originalUrl: 'https://v.douyin.com/abcdef/',
      url: 'https://example.com/covers/cover1.jpg',
      size: 524288, // 512KB
      status: DataStatus.COMPLETED,
      metadata: {
        author: '抖音用户1'
      },
      createdAt: '2023-10-01 12:35:24'
    },
    {
      id: 4,
      taskId: 1,
      title: '这是一个抖音视频封面2',
      type: 'cover',
      originalUrl: 'https://v.douyin.com/123456/',
      url: 'https://example.com/covers/cover2.jpg',
      size: 655360, // 640KB
      status: DataStatus.COMPLETED,
      metadata: {
        author: '抖音用户2'
      },
      createdAt: '2023-10-01 12:42:35'
    }
  ],
  2: [
    {
      id: 5,
      taskId: 2,
      title: 'B站舞蹈视频1',
      type: 'video',
      originalUrl: 'https://www.bilibili.com/video/BV123456',
      url: 'https://example.com/videos/bilibili1.mp4',
      coverUrl: 'https://example.com/covers/bilibili1.jpg',
      size: 31457280, // 30MB
      status: DataStatus.DOWNLOADING,
      metadata: {
        author: 'B站用户1',
        likes: 3300,
        comments: 520,
        shares: 210,
        description: 'B站舞蹈视频描述'
      },
      createdAt: '2023-10-02 09:20:11'
    },
    {
      id: 6,
      taskId: 2,
      title: 'B站舞蹈视频2',
      type: 'video',
      originalUrl: 'https://www.bilibili.com/video/BV654321',
      url: 'https://example.com/videos/bilibili2.mp4',
      coverUrl: 'https://example.com/covers/bilibili2.jpg',
      size: 26214400, // 25MB
      status: DataStatus.PENDING,
      metadata: {
        author: 'B站用户2',
        likes: 1800,
        comments: 320,
        shares: 150,
        description: '又一个B站舞蹈视频描述'
      },
      createdAt: '2023-10-02 09:22:05'
    }
  ]
}

// 获取任务列表
export function getTaskList(params) {
  // 筛选
  let filteredList = [...mockTaskList]
  
  if (params?.keyword) {
    filteredList = filteredList.filter(item => item.name.includes(params.keyword))
  }
  
  if (params?.status) {
    filteredList = filteredList.filter(item => item.status === params.status)
  }
  
  if (params?.type) {
    filteredList = filteredList.filter(item => item.type === params.type)
  }
  
  // 总数
  const total = filteredList.length
  
  // 分页
  let result = [...filteredList]
  if (params?.page && params?.pageSize) {
    const start = (params.page - 1) * params.pageSize
    const end = start + params.pageSize
    result = filteredList.slice(start, end)
  }
  
  return {
    list: result,
    total: total,
    page: params?.page || 1,
    pageSize: params?.pageSize || 10
  }
}

// 获取任务详情
export function getTaskDetail(id) {
  const detail = mockTaskDetails[id]
  
  if (!detail) {
    return {
      code: 404,
      message: '任务不存在'
    }
  }
  
  return {
    code: 200,
    data: detail
  }
}

// 获取任务数据
export function getTaskData(taskId, params) {
  const dataList = mockCrawlerData[taskId] || []
  
  // 筛选
  let filteredList = [...dataList]
  
  // 总数
  const total = filteredList.length
  
  // 分页
  let result = [...filteredList]
  if (params?.page && params?.pageSize) {
    const start = (params.page - 1) * params.pageSize
    const end = start + params.pageSize
    result = filteredList.slice(start, end)
  }
  
  return {
    list: result,
    total: total,
    page: params?.page || 1,
    pageSize: params?.pageSize || 10
  }
}

// 创建任务
export function createTask(data) {
  // 生成新ID
  const newId = mockTaskList.length > 0 ? Math.max(...mockTaskList.map(item => item.id)) + 1 : 1
  
  // 创建新任务
  const newTask = {
    id: newId,
    name: data.name,
    type: data.type,
    typeName: crawlerTypes.find(item => item.code === data.type)?.name || data.type,
    status: TaskStatus.WAITING,
    progress: 0,
    createdAt: new Date().toLocaleString(),
    updatedAt: new Date().toLocaleString(),
    config: data.config
  }
  
  // 添加到列表
  mockTaskList.unshift(newTask)
  
  // 添加到详情
  mockTaskDetails[newId] = newTask
  
  // 初始化空数据列表
  mockCrawlerData[newId] = []
  
  return {
    code: 200,
    message: '创建成功',
    data: newTask
  }
}

// 删除任务
export function deleteTask(id) {
  const index = mockTaskList.findIndex(item => item.id == id)
  
  if (index === -1) {
    return {
      code: 404,
      message: '任务不存在'
    }
  }
  
  // 从列表中移除
  mockTaskList.splice(index, 1)
  
  // 从详情中移除
  delete mockTaskDetails[id]
  
  // 从数据中移除
  delete mockCrawlerData[id]
  
  return {
    code: 200,
    message: '删除成功'
  }
}

// 获取爬虫类型列表
export function getCrawlerTypes() {
  return {
    code: 200,
    data: crawlerTypes
  }
} 