// src/api.ts
import axios from 'axios'

// 1. 配置 Axios 实例
const request = axios.create({
    baseURL: '/api/v1', // 配合 vite.config.ts 的代理
    timeout: 5000,
})

// 2. 请求拦截器：自动把 Token 带上
request.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        if (config.params) {
            const params = config.params;
            Object.keys(params).forEach(key => {
                if (params[key] === null || params[key] === undefined || params[key] === '') {
                    delete params[key];
                }
            });
        }
        return config
    },

)

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
        if (error.response) {
            console.error(error.response.data)
        }
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
    user_rating_stat: UserRatingStat // 新增字段
}

export interface UserRatingStat {
    target_user_id: number
    avg_rating:    number
    review_count: number
}

export interface Category {
    id: number
    name: string
}

// 对应 internal/controller/user.go 中的 LoginReq
export interface LoginReq {
    account: string
    password: string
}

// 对应 internal/controller/product.go 中的 ListProductSearchResult
export interface ProductListResp {
    list: Product[]
    has_more: boolean
    last_id: number
}

export interface CreateProductReq {
    name: string
    description: string
    price: number
    category_id: number | null // 允许初始为空
    condition: string
    image_url: string
}

export interface UserMeResp {
    nickname: string
    email: string
    phone: string
    username: string

}

// --- API 方法导出 ---

export const api = {
    // ----------------- Auth Group (/api/v1/auth) -------网络异常----------

    // 登录
    // 对应 Go: authGroup.POST("/login", ...) -> /api/v1/auth/login
    login: (data: LoginReq) => request.post<{ token: string }>('/auth/login', data),

    // 注册
    // 对应 Go: authGroup.POST("/register", ...) -> /api/v1/auth/register
    register: (data: RegisterReq) => request.post('/auth/register', data),

    // ----------------- User/Private Group (/api/v1) -----------------

    // 获取当前用户信息
    // 对应 Go: privateGroup.GET("/users/me", ...) -> /api/v1/users/me
    getMe: () => request.get<UserMeResp>('/users/me'),

    // 获取我发布的商品
    // 对应 Go: privateGroup.GET("/users/me/products", ...) -> /api/v1/users/me/products
    getMyProducts: (lastId: number = 0, pageSize: number = 8) => request.get<ProductListResp>('/users/me/products', {
        params: { last_id: lastId, page_size: pageSize }
    }),

    // ----------------- Product Group (/api/v1/products) -----------------

    // 获取商品列表
    // 对应 Go: productGroup.GET("", ...) -> /api/v1/products
    getProducts: (lastId: number = 0, pageSize: number = 8, category: string) =>
        request.get<ProductListResp>('/products', {
            params: { last_id: lastId, page_size: pageSize, category: category }
        }),

    // 搜索商品
    // 对应 Go: productGroup.GET("/search", ...) -> /api/v1/products/search
    searchProducts: (keyword: string, last_id: number, page_size: number) =>
        request.get<ProductListResp>('/products/search', {
            params: {
                keyword: keyword,
                last_id: last_id,
                page_size: page_size,
            }
        }),
    getHotCategories: () => request.get<{ list: Category[] }>('/categories/popular'),
    // 搜索建议
    // 对应 Go: productGroup.GET("/suggestion", ...) -> /api/v1/products/suggestion
    getSuggestions: (prefix: string) =>
        request.get<{ list: string[] }>('/products/suggestion', {
            params: { prefix }
        }),

    // ----------------- Order Group (/api/v1/orders) -----------------

    // 创建订单
    // 对应 Go: orderGroup.POST("", ...) -> /api/v1/orders
    // 注意：后端这里是 privateGroup 下的 /orders 组
    createOrder: (itemId: number) => request.post('/orders', { item_id: itemId }),

    createProduct: (data: CreateProductReq) => request.post('/products', data),

    uploadImage: (file: File) => {
        const formData = new FormData()
        formData.append('file', file)
        return request.post<{ url: string }>('/images', formData)
    }

}

export interface RegisterReq {
    username: string
    password: string
    re_password: string
    email: string
    // phone, nickname 可选
}