import { boot } from 'quasar/wrappers';
import axios from 'axios';
import type  { AxiosInstance, AxiosError, AxiosResponse } from 'axios'
// import { useCookies } from 'vue3-cookies';
import { Notify, LocalStorage } from 'quasar';

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $api: AxiosInstance;
  }
}

const api: AxiosInstance = axios.create({
  baseURL: 'http://localhost:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // const { cookies } = LocalStorage.getItem("token");
    const token = LocalStorage.getItem('token');
    if (token) {
      config.headers['X-Token'] = token;
    }
    return config;
  },
  (error: AxiosError) => {
    Notify.create({
      type: 'negative',
      message: `请求异常：${error.message}`
    });
    return Promise.reject(new Error(`请求拦截器错误：${error.message}`));
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response.data;
  },
  (error: AxiosError) => {
    let msg = '请求失败';
    if (error.response) {
      const status = error.response.status;
      switch (status) {
        case 401:
          msg = '登录已过期，请重新登录';
          // const { cookies } = useCookies();
          LocalStorage.remove('token');
          LocalStorage.remove('user');
          window.location.href = '/#/login';
          break;
        case 403:
          msg = '无权限操作';
          break;
        case 404:
          msg = '资源不存在';
          break;
        default:
          msg = (error.response.data as string) || `请求错误(${status})`;
          break;
      }
    }
    Notify.create({
      type: 'negative',
      message: msg
    });
    const rejectError = new Error(msg);
    rejectError.name = 'AxiosError';
    // rejectError.message = error;
    return Promise.reject(rejectError);
  }
);

export default boot(({ app }) => {
  app.config.globalProperties.$api = api;
});

export { api };