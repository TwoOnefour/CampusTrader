<script setup lang="ts">
import { h, ref, onMounted, computed, Component, reactive,  } from 'vue'
import {api, type Category, type Product} from '../api'
import {
  NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu,
  NButton, NCard, NInput, NSpace, NTag, NGrid, NGridItem, NList, NListItem, NThing, NResult,
  useMessage, NAvatar, NDropdown, NIcon, NEmpty, NModal, NForm, NFormItem,
  NInputNumber, NSelect, NUpload, NUploadDragger, NText, NImage,NTabs,NTabPane,
  NAutoComplete,
  type AutoCompleteOption,
  type UploadFileInfo,
} from 'naive-ui'

import {
  BagHandleOutline, PersonOutline, LogOutOutline,
  CartOutline, AddCircleOutline, SearchOutline, CloudUploadOutline,FlameOutline,MenuOutline
} from '@vicons/ionicons5'
function renderIcon(icon: Component) {
  return () => h(NIcon, null, { default: () => h(icon) })
}
const submitting = ref<boolean>(false)
const message = useMessage()
const token = ref(localStorage.getItem('token') || '')
const username = ref(localStorage.getItem('user') || '')
const currentView = ref('market')
const products = ref<Product[]>([])
const searchKeyword = ref('')
const category_id = ref('')

// 存放下拉框的选项，格式必须是 { label: '显示文字', value: '选中后的值' }
const searchOptions = ref<AutoCompleteOption[]>([])
const collapsed = ref(false)
// 简单的防抖定时器，防止请求太频繁
let searchTimer: any = null

// 当用户输入内容变化时触发
const handleSearchInput = (value: string) => {
  searchKeyword.value = value

  // 1. 如果输入为空
  if (!value || !value.trim()) {
    searchOptions.value = []
    loadMarket()
    return
  }

  // 2. 清除上一次没执行的定时器
  if (searchTimer) clearTimeout(searchTimer)

  // 3. 开启防抖定时器 (保持不变)
  searchTimer = setTimeout(async () => {
    try {
      const res = await api.getSuggestions(value)
      const list = res.list || []
      searchOptions.value = list.map(item => ({
        label: item,
        value: item
      }))
    } catch (e) {
      console.error('获取建议失败', e)
    }
  }, 300)
}

// 选中建议或者按回车时触发搜索
const handleSearch = async (value?: string) => {
  // 如果是选中建议传进来的 value，就用它；否则用输入框当前的 searchKeyword
  const keyword = (typeof value === 'string' ? value : searchKeyword.value).trim()

  // 更新一下输入框显示的值（如果是点选建议的话）
  searchKeyword.value = keyword

  // 关闭下拉建议（清空选项即可）
  searchOptions.value = []

  if (!keyword) {
    // 清空搜索时，重置回到市场第一页
    loadMarket(false, true)
    return
  }

  try {
    message.loading('搜索中...')
    const res = await api.searchProducts(keyword)
    // @ts-ignore
    products.value = res.list || []
    if (products.value.length === 0) {
      message.info('未找到相关商品')
    }
  } catch (err) {
    message.error('搜索出错')
  }
}


// --- 1. 动态计算菜单 (实现需求一：权限控制) ---
const allMenuOptions = [
  { label: '交易市场', key: 'market', icon: renderIcon(BagHandleOutline) },
  { label: '热门分类', key: 'hot-categories', icon: renderIcon(FlameOutline) }, // 新增
  { label: '我发布的', key: 'my-products', icon: renderIcon(CartOutline), requiresAuth: true },
  { label: '发布商品', key: 'create', icon: renderIcon(AddCircleOutline), requiresAuth: true },
  { label: '个人中心', key: 'profile', icon: renderIcon(PersonOutline), requiresAuth: true }
]

const hotCategories = ref<Category[]>([])

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

const initialFormState = {
   name: '',
   description: '',
   price: null as number | null,
   category_id: null as number | null,
   condition: 'good',
   image_url: ''
}

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
    const meresp = await api.getMe()
    token.value = tokenStr
    localStorage.setItem('user', meresp.nickname)


    message.success('登录成功！')
    showLoginModal.value = false

    // 登录后刷新数据
    loadMarket()

    // (可选) 获取一下用户信息更新右上角名字，这里简单模拟
    username.value = meresp.nickname

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

const handleCustomRequest = async ({ file, onFinish, onError }: UploadCustomRequestOptions) => {
  try {
    // 1. 防御性检查：确保拿到了原生 JS File 对象
    if (!file.file) {
      message.error('未找到文件流')
      return
    }

    // 2. 调用我们之前封装好的 API (它会自动带 Token)
    // 注意：这里 api.uploadImage 返回的是 Promise<{ url: string }>
    const res = await api.uploadImage(file.file)

    // 3. 上传成功，拿到 URL 填入表单
    createForm.image_url = res.url

    // 4. 告诉 Naive UI 组件：“完事了，把进度条变绿”
    onFinish()
    message.success('上传成功')

  } catch (error) {
    // 5. 告诉组件：“挂了，把进度条变红”
    onError()
    message.error('上传失败，请检查网络或图片大小')
  }
}

// 【修改】这个函数只负责同步文件列表的显示状态，不需要做业务逻辑了
const handleUploadChange = (data: { fileList: UploadFileInfo[] }) => {
  fileList.value = data.fileList

  // 如果用户点击了删除图片（列表为空），要把表单里的 URL 也清空
  if (data.length === 0) {
    createForm.image_url = ''
  }
}

// --- 分页状态管理 ---
const pageSize = 8             // 每页显示多少条
const currentPage = ref(1)     // 当前第几页（仅用于显示）
const cursorHistory = ref([0]) // 游标历史栈：第1页对应0，第2页对应上一页最后一条ID...
const hasMore = ref(true)      // 是否还有下一页数据

// --- 修改后的加载函数 ---
// useCursor: 是否使用当前记录的游标去加载（用于翻页）
// reset: 是否重置分页（用于切换菜单或搜索时）
const loadMarket = async (useCursor = false, reset = false, type = 'market') => {
  try {
    if (reset) {
      currentPage.value = 1
      cursorHistory.value = [0]
      hasMore.value = true
    }

    // 获取当前页对应的 last_id
    // cursorHistory 比如是 [0, 100, 92]，currentPage=2，则取 index=1 的 100
    const cursorIndex = currentPage.value - 1
    const lastId = useCursor ? cursorHistory.value[cursorIndex] : 0
    let res;
    // 调用 API
    if (type === 'market') {
      res = await api.getProducts(lastId, pageSize, category_id.value)
    } else if (type === 'my-products') {
      res = await api.getMyProducts(lastId, pageSize)
    }

    // @ts-ignore
    const list = res.list || []
    products.value = list

    // 判断是否还有下一页：如果返回的数量 < pageSize，说明数据取完了
    // 注意：后端返回 total 仅仅是本次查询的数量，不是总数，所以不能用 total 判断
    if (!res.has_more) {
      hasMore.value = false
    } else {
      hasMore.value = true
      // 预存下一页的游标：取当前列表最后一条数据的 ID
      const lastItem = list[list.length - 1]
      // 只有当是最新页时才 push，防止在“上一页”操作时重复添加
      if (cursorHistory.value.length <= currentPage.value) {
        cursorHistory.value.push(lastItem.id)
      }
    }
  } catch (err) {
    message.error('加载失败')
  }
}

// --- 翻页操作 ---
const handlePrevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    switch (currentView.value) {
      case ('market'):
          loadMarket(true) // true 表示使用历史游标
          break
      case ('my-products'):
          loadMarket(true, false, 'my-products')
    }
  }
}

const handleNextPage = () => {
  if (hasMore.value) {
    currentPage.value++
    switch (currentView.value) {
      case ('market'):
        loadMarket(true)
        break
      case ('my-products'):
        loadMarket(true, false, 'my-products')
    }
  }
}
const handleCreateSubmit = async () => {
  // 1. 基础校验
  if (!createForm.name || !createForm.price || !createForm.category_id) {
    message.warning('请补全必填项 (名称、价格、分类)')
    return
  }

  if (!createForm.image_url) {
    message.warning('请先上传一张商品图片')
    return
  }

  try {
    submitting.value = true

    // 2. 调用 API
    await api.createProduct({
        name: createForm.name,
        description: createForm.description,
        price: createForm.price,
        category_id: createForm.category_id,
        condition: createForm.condition,
        image_url: createForm.image_url
    })

    message.success('发布成功！')
    // 3. 收尾工作
    showCreateModal.value = false // 关闭弹窗
    resetForm() // 清空表单

    // TODO: 这里通常需要触发一个事件，通知父组件刷新商品列表
    // emit('refresh')
    loadMarket()
  } catch (error) {
    // 错误处理交给了 api.ts 里的拦截器打印，这里只需要通过 catch 停住 loading 状态
    // message.error('发布失败')
  } finally {
    submitting.value = false
  }
}

// 重置表单辅助函数
const resetForm = () => {
  Object.assign(createForm, initialFormState)
  fileList.value = [] // 清空上传组件的显示
}

// --- 业务逻辑 ---

const handleMenuUpdate = (key: string) => {
  const option = allMenuOptions.find(o => o.key === key)

  if (option?.requiresAuth && !token.value) {
    message.warning('请先登录')
    return
  }
  if (window.innerWidth < 768) {
    collapsed.value = true // 手机端点击菜单后自动收回
  }
  if (key === 'create') {
    showCreateModal.value = true
    return
  }
  currentView.value = key

  if (key === 'market') {
    // 切换回市场时，重置分页 (reset = true)
    loadMarket(false, true)
  } else if (key === 'my-products') {
    loadMyProducts()
  } else if (key === 'hot-categories') {
    loadHotCategories()
  }
}

const loadHotCategories = async () => {
  try {
    const res = await api.getHotCategories()
    // @ts-ignore
    hotCategories.value = res.list || []
  } catch (e) {
    message.error('加载热门分类失败')
  }
}

const handleUserDropdown = (key: string) => {
  if (key === 'logout') {
    token.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    message.success('已退出')
    // 强制刷新页面以更新状态
    window.location.reload()
  } else if (key === 'profile') {
    currentView.value = 'profile'
  }
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
    const resp = await api.createOrder(id)

    message.success('购买成功')
    loadMarket()
  } catch (err) {
    if (err.response) {
      message.error(err.response.data.msg)
      return
    }
    message.error(err.toString())
  }
}

const handleUserSession = async () => {
  try {
    await api.getMe()
  } catch (error) {

    if (error?.response.data.statusCode === 401 && localStorage.getItem('token') != null)  {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      message.info("登陆已过期")
    }
  }
}

onMounted(() => {
  if (window.innerWidth < 768) {
    collapsed.value = true
  }
  handleUserSession()
  loadMarket()
})
</script>

<template>
  <n-layout position="absolute" has-sider>
    <div class="mobile-menu-btn" @click="collapsed = !collapsed">
      <n-icon size="28"><MenuOutline /></n-icon>
    </div>

    <n-layout-sider
        bordered
        collapse-mode="transform"
        :collapsed-width="0"
        :width="240"
        :collapsed="collapsed"
        show-trigger="bar"
        @collapse="collapsed = true"
        @expand="collapsed = false"
        :native-scrollbar="false"
    style="z-index: 100"
    >
    <div style="padding: 24px;">
      <div style="margin-bottom: 30px; display: flex; align-items: center; gap: 10px;">
        <n-icon size="30" color="#18a058"><CartOutline /></n-icon>
        <span style="font-size: 18px; font-weight: bold; color: #333;">CampusTrader</span>
      </div>
      <n-menu :options="menuOptions" :value="currentView" @update:value="handleMenuUpdate" />
    </div>
    </n-layout-sider>

    <n-layout>
      <n-layout-header bordered style="height: 64px; display: flex; align-items: center; padding: 0 24px; justify-content: space-between;">

        <div style="display: flex; align-items: center; min-width: 120px;">
          <h2 style="margin: 0; font-size: 16px;">
            {{ currentView === 'market' ? '🛒 交易市场' : currentView === 'my-products' ? '📦 我的商品' : '个人中心' }}
          </h2>
        </div>

        <div v-if="currentView === 'market'" style="flex: 1; max-width: 400px; margin: 0 20px;">
          <n-auto-complete
              v-model:value="searchKeyword"
              :options="searchOptions"
              placeholder="输入关键词搜索 (例如: Mac)"
              clearable
              @update:value="handleSearchInput"
              @select="handleSearch"
              @keydown.enter="handleSearch"
          >
            <template #prefix>
              <n-icon :component="SearchOutline" />
            </template>
          </n-auto-complete>
        </div>

        <div style="display: flex; align-items: center; gap: 20px; min-width: 120px; justify-content: flex-end;">
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
                    <n-tag v-if="item.user_rating_stat.review_count > 0" size="small" type="warning" :bordered="false">
                      卖家评分 ⭐ {{ item.user_rating_stat.avg_rating.toFixed(1) }}  共 {{ item.user_rating_stat.review_count }} 条评价
                    </n-tag>
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
          <div v-if="(currentView === 'market' || currentView === 'my-products') && products.length > 0"
               style="display: flex; justify-content: center; align-items: center; margin-top: 30px; gap: 20px;">
            <n-button :disabled="currentPage === 1" @click="handlePrevPage">
              上一页
            </n-button>

            <span style="font-weight: bold; color: #666;">
        第 {{ currentPage }} 页
    </span>

            <n-button :disabled="!hasMore" @click="handleNextPage">
              下一页
            </n-button>
          </div>
        </div>
        <div v-else-if="currentView === 'hot-categories'">
          <n-card title="🔥 最火商品分类 (Top 3)">
            <n-list hoverable clickable>
              <n-list-item v-for="(cat, index) in hotCategories" :key="cat.id">
                <n-thing :title="cat.name" content-style="margin-top: 10px;">
                  <template #description>
                    <n-tag type="error" size="small">No. {{ index + 1 }}</n-tag>
                  </template>
                  这里是 {{ cat.name }} 专区，包含了大量热门交易商品。
                </n-thing>
                <template #suffix>
                  <n-button size="small">查看详情</n-button>
                </template>
              </n-list-item>
            </n-list>
            <n-empty v-if="hotCategories.length === 0" description="暂无热门数据" />
          </n-card>
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
        <n-grid cols="10" x-gap="12">
          <n-form-item-gi :span="5" label="价格 (¥)" path="price">
            <n-input-number
                v-model:value="createForm.price"
                :min="0"
                placeholder="0.00"
                style="width: 100%"
            />
          </n-form-item-gi>

            <n-form-item-gi :span="5" label="分类" path="category_id">
              <n-select v-model:value="createForm.category_id" :options="categoryOptions" />
            </n-form-item-gi>
        </n-grid>
        <n-form-item label="成色" path="condition">
          <n-select v-model:value="createForm.condition" :options="conditionOptions" />
        </n-form-item>

        <n-form-item label="商品图片">
          <n-upload
            :custom-request="handleCustomRequest"
            @update:file-list="handleUploadChange"
            :file-list="fileList"

            directory-dnd
            :max="1"
            list-type="image"
          >
            <n-upload-dragger>
              <div style="margin-bottom: 12px">
                <n-icon size="48" :depth="3">
                  <CloudUploadOutline />
                </n-icon>
              </div>
              <n-text style="font-size: 16px">点击或者拖动图片到该区域来上传</n-text>
              <p depth="3" style="margin: 8px 0 0 0">
                支持 JPG、PNG 格式，建议小于 5MB
              </p>
            </n-upload-dragger>
          </n-upload>
        </n-form-item>

      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="() => { showCreateModal = false;resetForm() }">取消</n-button>
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