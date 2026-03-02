import axios from 'axios';
import type { AxiosInstance,  AxiosError } from 'axios'
import { LocalStorage, useQuasar } from 'quasar'
// import  Cookies from 'js-cookie';

const $q = useQuasar()
// 定义响应数据类型
export interface ApiResponse<T> {
  code: number;
  msg: string;
  data: T;
//   [key: string]: any;
}

// 创建axios实例
const service: AxiosInstance = axios.create({
  baseURL: 'http://localhost:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    const token = LocalStorage.getItem('token') as string;
    if (token && config.headers) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => {
    console.error('Request error:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器
service.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error: AxiosError) => {
    console.error('Response error:', error);
    // const msg = error.response?.data?.msg || error.message || 'Request failed';
    
    // 全局错误提示
      $q.notify({
        type: 'negative',
        message: error.message,
        position: 'top-right'
      });

    // Token过期处理
    if (error.response?.status === 401) {
      LocalStorage.remove('token');
      LocalStorage.remove('userInfo');
      window.location.href = '/login';
    }

    return Promise.reject(error);
  }
);

export default service;