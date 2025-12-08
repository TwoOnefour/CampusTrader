<script setup lang="ts">
import { h, ref, onMounted, computed, Component, reactive } from 'vue'
import { api, type Product } from '../api'
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu,
  NButton, NCard, NInput, NSpace, NTag, NGrid, NGridItem,
  useMessage, NAvatar, NDropdown, NIcon, NEmpty, NModal, NForm, NFormItem,
  NInputNumber, NSelect, NUpload, NUploadDragger, NText, NImage,NTabs,NTabPane,
  type UploadFileInfo
} from 'naive-ui'
// 引入图标
import {
  BagHandleOutline, PersonOutline, LogOutOutline,
  CartOutline, AddCircleOutline, SearchOutline, CloudUploadOutline
} from '@vicons/ionicons5'
// import { CloudUpload } from '@vicons/fa'

// --- 图标渲染辅助函数 ---
function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const message = useMessage()
const token = ref(localStorage.getItem('token') || '')
const username = ref('User')
const currentView = ref('market')
const products = ref<Product[]>([])

// --- 1. 动态计算菜单 (实现需求一：权限控制) ---
const allMenuOptions = [
  { label: '交易市场', key: 'market', icon: renderIcon(BagHandleOutline) },
  { label: '我发布的', key: 'my-products', icon: renderIcon(CartOutline), requiresAuth: true },
  { label: '发布商品', key: 'create', icon: renderIcon(AddCircleOutline), requiresAuth: true },
  { label: '个人中心', key: 'profile', icon: renderIcon(PersonOutline), requiresAuth: true }
]

const menuOptions = computed(() => {
  return allMenuOptions.filter(option => {
    // 如果选项需要权限，且token不存在，则过滤掉
    if (option.requiresAuth && !token.value) {
      return false
    }
    return true
  })
})

const userDropdownOptions = [
  { label: '个人资料', key: 'profile', icon: renderIcon(PersonOutline) },
  { label: '退出登录', key: 'logout', icon: renderIcon(LogOutOutline) }
]

// --- 2. 发布商品表单相关 (实现需求二：弹出表单+上传) ---
const showCreateModal = ref(false)
const createFormRef = ref(null)
// 表单数据模型
const createForm = reactive({
  name: '',
  description: '',
  price: null as number | null,
  category_id: null as number | null,
  condition: 'good',
  image_url: ''
})
const fileList = ref<UploadFileInfo[]>([])
// --- 3. 登录/注册相关 (新增) ---
const showLoginModal = ref(false)
const activeTab = ref('login') // 控制显示登录还是注册
const loginForm = reactive({
  account: '',
  password: ''
})
const registerForm = reactive({
  username: '',
  password: '',
  re_password: '',
  email: ''
})

// 处理登录提交
const handleLoginSubmit = async () => {
  if (!loginForm.account || !loginForm.password) {
    message.warning('请输入账号和密码')
    return
  }

  try {
    const res = await api.login({
      account: loginForm.account,
      password: loginForm.password
    })

    // @ts-ignore
    const tokenStr = res.token
    localStorage.setItem('token', tokenStr)
    token.value = tokenStr

    message.success('登录成功！')
    showLoginModal.value = false

    // 登录后刷新数据
    loadMarket()

    // (可选) 获取一下用户信息更新右上角名字，这里简单模拟
    username.value = loginForm.account

  } catch (err: any) {
    message.error(err.message || '登录失败，请检查账号密码')
  }
}

// 处理注册提交 (简单实现)
const handleRegisterSubmit = async () => {
  if (registerForm.password !== registerForm.re_password) {
    message.error('两次输入的密码不一致')
    return
  }
  try {
    // 假设你api.ts里有register，如果没有先注释掉这行
    await api.register(registerForm)
    message.success('注册成功，请登录')
    activeTab.value = 'login' // 切换回登录页
  } catch (err: any) {
    message.error(err.message || '注册失败')
  }
}


const conditionOptions = [
  { label: '全新 (New)', value: 'new' },
  { label: '几乎全新 (Like New)', value: 'like_new' },
  { label: '功能完好 (Good)', value: 'good' },
  { label: '有瑕疵 (Fair)', value: 'fair' },
  { label: '功能缺陷 (Poor)', value: 'poor' }
]
// 这里需要你后端提供真实的分类列表，先写死
const categoryOptions = [
  { label: '电子数码', value: 1 },
  { label: '图书教材', value: 2 },
  { label: '宿舍电器', value: 3 }
]

// 处理文件上传变化
const handleUploadChange = (data: { fileList: UploadFileInfo[] }) => {
  fileList.value = data.fileList
  // 这里是一个简单的模拟，实际需要调用后端上传接口获取 URL
  if (data.fileList.length > 0) {
    const file = data.fileList[0]
    if (file.status === 'finished') {
      // 假设后端返回的 URL 放在 file.url 里
      // createForm.image_url = file.url
      message.success('上传成功 (模拟)')
      // 模拟设置一个图片地址
      createForm.image_url = 'https://via.placeholder.com/300'
    }
  } else {
    createForm.image_url = ''
  }
}

// 提交发布表单
const handleCreateSubmit = () => {
  // 这里应该调用 api.createProduct(createForm)
  console.log('提交表单:', createForm)
  message.info('提交功能后端尚未实现')
  showCreateModal.value = false
}

// --- 业务逻辑 ---

const handleMenuUpdate = (key: string) => {
  const option = allMenuOptions.find(o => o.key === key)
  if (option?.requiresAuth && !token.value) {
    message.warning('请先登录')
    return
  }

  if (key === 'create') {
    showCreateModal.value = true
    return
  }
  currentView.value = key
  if (key === 'market') loadMarket()
  else if (key === 'my-products') loadMyProducts()
}

const handleUserDropdown = (key: string) => {
  if (key === 'logout') {
    token.value = ''
    localStorage.removeItem('token')
    message.success('已退出')
    // 强制刷新页面以更新状态
    window.location.reload()
  } else if (key === 'profile') {
    currentView.value = 'profile'
  }
}

const loadMarket = async () => {
  try {
    const res = await api.getProducts()
    // @ts-ignore
    products.value = res.list || []
  } catch (err) { message.error('加载失败') }
}

const loadMyProducts = async () => {
  if (!token.value) return
  try {
    const res = await api.getMyProducts()
    // @ts-ignore
    products.value = res.list || []
  } catch (err) { message.error('加载失败') }
}

const handleBuy = async (id: number) => {
  try {
    await api.createOrder(id)
    message.success('购买成功')
    loadMarket()
  } catch (err) { message.error('购买失败') }
}

onMounted(() => {
  loadMarket()
})
</script>

<template>
  <n-layout position="absolute" has-sider>
    <n-layout-sider bordered width="240" content-style="padding: 24px;" :native-scrollbar="false">
      <div style="margin-bottom: 30px; display: flex; align-items: center; gap: 10px;">
        <n-icon size="30" color="#18a058"><CartOutline /></n-icon>
        <span style="font-size: 18px; font-weight: bold; color: #333;">CampusTrader</span>
      </div>
      <n-menu :options="menuOptions" :value="currentView" @update:value="handleMenuUpdate" />
    </n-layout-sider>

    <n-layout>
      <n-layout-header bordered style="height: 64px; display: flex; align-items: center; padding: 0 24px; justify-content: space-between;">
        <div style="display: flex; align-items: center;">
          <h2 style="margin: 0; font-size: 16px;">
            {{ currentView === 'market' ? '🛒 交易市场' : currentView === 'my-products' ? '📦 我的商品' : '个人中心' }}
          </h2>
        </div>
        <div style="display: flex; align-items: center; gap: 20px;">
          <div v-if="token">
            <n-dropdown :options="userDropdownOptions" @select="handleUserDropdown">
              <div style="display: flex; align-items: center; cursor: pointer; gap: 10px;">
                <n-avatar round size="small" src="https://api.dicebear.com/7.x/avataaars/svg?seed=Felix" />
                <span>{{ username }}</span>
              </div>
            </n-dropdown>
          </div>
          <div v-else>
            <n-button type="primary" size="small" @click="showLoginModal = true">
              登录 / 注册
            </n-button>
          </div>
        </div>
      </n-layout-header>

      <n-layout-content content-style="padding: 24px; background-color: #f5f7f9; min-height: 100vh;">
        <div v-if="currentView === 'market' || currentView === 'my-products'">
          <n-grid x-gap="16" y-gap="16" cols="1 600:2 900:3 1200:4">
            <template v-for="item in products" :key="item.id">
              <n-grid-item v-if="item.image_url">
                <n-card hoverable content-style="padding: 0;">
                  <template #cover>
                    <n-image
                        width="100%"
                        height="180"
                        :src="item.image_url"
                        object-fit="cover"
                        preview-disabled
                    />
                  </template>

                  <div style="padding: 15px;">
                    <div style="font-size: 16px; font-weight: bold; margin-bottom: 10px;">{{ item.name }}</div>
                    <n-space justify="space-between">
                      <n-tag size="small" :type="item.status === 'available' ? 'success' : 'default'">{{ item.status }}</n-tag>
                      <n-tag size="small" :bordered="false">{{ item.condition }}</n-tag>
                    </n-space>

                    <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 15px;">
                      <span style="color: #f59e0b; font-size: 18px; font-weight: bold;">¥ {{ item.price }}</span>
                      <n-button
                          v-if="currentView !== 'my-products'"
                          type="primary"
                          size="small"
                          :disabled="item.status !== 'available'"
                          @click="handleBuy(item.id)"
                      >
                        Buy Now
                      </n-button>
                    </div>
                  </div>
                </n-card>
              </n-grid-item>
            </template>
          </n-grid>
          <n-empty v-if="products.length === 0" description="这里空空如也" style="margin-top: 100px" />
        </div>

        <div v-else-if="currentView === 'profile'">
          <n-card title="个人中心">
            <p>这里放用户的个人信息修改表单...</p>
          </n-card>
        </div>
      </n-layout-content>
    </n-layout>

    <n-modal v-model:show="showCreateModal" preset="card" title="发布新商品" style="width: 600px;">
      <n-form ref="createFormRef" :model="createForm" label-placement="left" label-width="auto">
        <n-form-item label="商品名称" path="name">
          <n-input v-model:value="createForm.name" placeholder="例如：MacBook Pro M1" />
        </n-form-item>
        <n-form-item label="商品描述" path="description">
          <n-input v-model:value="createForm.description" type="textarea" placeholder="描述一下商品的细节、新旧程度等" />
        </n-form-item>
        <n-grid cols="2" x-gap="12">
          <n-form-item label="价格 (¥)" path="price">
            <n-input-number v-model:value="createForm.price" :min="0" placeholder="0.00" style="width: 100%"/>
          </n-form-item>
          <n-form-item label="分类" path="category_id">
            <n-select v-model:value="createForm.category_id" :options="categoryOptions" placeholder="选择分类" />
          </n-form-item>
        </n-grid>
        <n-form-item label="成色" path="condition">
          <n-select v-model:value="createForm.condition" :options="conditionOptions" />
        </n-form-item>

        <n-form-item label="商品图片">
          <n-upload
              multiple
              directory-dnd
              :max="1"
              list-type="image"
              :file-list="fileList"
              @update:file-list="handleUploadChange"
          action="/api/v1/upload/image"
          >
          <n-upload-dragger>
            <div style="margin-bottom: 12px">
              <n-icon size="48" :depth="3"><CloudUploadOutline /></n-icon>
            </div>
            <n-text style="font-size: 16px">点击或者拖动图片到该区域来上传</n-text>
            <n-p depth="3" style="margin: 8px 0 0 0">支持 JPG、PNG 格式，请勿上传敏感图片</n-p>
          </n-upload-dragger>
          </n-upload>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showCreateModal = false">取消</n-button>
          <n-button type="primary" @click="handleCreateSubmit">确认发布</n-button>
        </n-space>
      </template>
    </n-modal>
    <n-modal v-model:show="showLoginModal" preset="card" style="width: 400px;">
      <n-tabs v-model:value="activeTab" justify-content="space-evenly" animated>

        <n-tab-pane name="login" tab="登录">
          <n-form>
            <n-form-item label="账号">
              <n-input v-model:value="loginForm.account" placeholder="用户名 / 邮箱" />
            </n-form-item>
            <n-form-item label="密码">
              <n-input
                  v-model:value="loginForm.password"
                  type="password"
                  show-password-on="click"
                  placeholder="请输入密码"
                  @keydown.enter="handleLoginSubmit"
              />
            </n-form-item>
            <n-button type="primary" block @click="handleLoginSubmit">
              立即登录
            </n-button>
          </n-form>
        </n-tab-pane>

        <n-tab-pane name="register" tab="注册新账号">
          <n-form>
            <n-form-item label="用户名">
              <n-input v-model:value="registerForm.username" placeholder="设置用户名" />
            </n-form-item>
            <n-form-item label="邮箱">
              <n-input v-model:value="registerForm.email" placeholder="用于找回密码" />
            </n-form-item>
            <n-form-item label="密码">
              <n-input v-model:value="registerForm.password" type="password" show-password-on="click" />
            </n-form-item>
            <n-form-item label="确认密码">
              <n-input v-model:value="registerForm.re_password" type="password" show-password-on="click" />
            </n-form-item>
            <n-button type="success" block @click="handleRegisterSubmit">
              注册并登录
            </n-button>
          </n-form>
        </n-tab-pane>

      </n-tabs>
    </n-modal>
  </n-layout>
</template>