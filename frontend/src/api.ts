// src/api.ts
import axios from 'axios'

// 1. 配置 Axios 实例
const request = axios.create({
    baseURL: '/api/v1', // 配合 vite.config.ts 的代理
    timeout: 5000,
})

// 2. 请求拦截器：自动把 Token 带上
request.interceptors.request.use((config) => {
    const token = localStorage.getItem('token')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

// 3. 响应拦截器：简化数据返回，处理报错
request.interceptors.response.use(
    (response) => {
        const res = response.data
        // 你的后端约定：code === 0 代表成功
        if (res.code === 0) {
            return res.data
        } else {
            // 这里可以对接 Naive UI 的 message 报错
            console.error(res.msg)
            return Promise.reject(new Error(res.msg))
        }
    },
    (error) => {
        console.error('网络异常', error)
        return Promise.reject(error)
    }
)

// --- 类型定义 (对应你的 Go Struct) ---

// 对应 internal/model/product.go
export interface Product {
    id: number
    name: string
    description: string
    price: number
    category_id: number
    seller_id: number
    status: 'available' | 'sold' | 'removed'
    condition: 'new' | 'like_new' | 'good' | 'fair' | 'poor'
    image_url: string
    created_at: string
}

// 对应 internal/controller/user.go 中的 LoginReq
export interface LoginReq {
    account: string
    password: string
}

// 对应 internal/controller/product.go 中的 ListProductSearchResult
export interface ProductListResp {
    list: Product[]
    total: number
    page: number
    size: number
}

// --- API 方法导出 ---

export const api = {
    // 登录
    login: (data: LoginReq) => request.post<{ token: string }>('/login', data),

    // 获取当前用户信息
    getMe: () => request.get('/me'),

    // 获取商品列表 (对应 ProductController.ListProducts)
    getProducts: (lastId: number = 0) =>
        request.get<ProductListResp>('/product/list', { params: { last_id: lastId, page_size: 10 } }),

    // 创建订单 (对应 OrderController.Order)
    createOrder: (itemId: number) => request.post('/order/create', { item_id: itemId }),
    getMyProducts: () => request.get<ProductListResp>('/product/my'),
    register: (data: RegisterReq) => request.post('/register', data),
}

export interface RegisterReq {
    username: string
    password: string
    re_password: string
    email: string
    // phone, nickname 可选
}