import { createRouter, createWebHistory,  createWebHashHistory } from 'vue-router';
import type {RouteRecordRaw} from 'vue-router';
// import Cookies from 'js-cookie';
import Layout from 'src/layouts/MainLayout.vue';
import Login from 'src/pages/LoginPage.vue';
import Register from 'src/pages/RegisterPage.vue';
import FileManager from 'src/pages/FileManager.vue';
import { LocalStorage } from 'quasar'

// 路由类型定义
const routes: Array<RouteRecordRaw> = [
  {
    path: '/login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    component: Register,
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: Layout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        component: FileManager
      }
    ]
  }
];

// 创建路由实例
const router = createRouter({
  history: process.env.VUE_ROUTER_MODE === 'history' 
    ? createWebHistory(process.env.BASE_URL) 
    : createWebHashHistory(process.env.BASE_URL),
  routes
});

// 路由守卫
router.beforeEach((to, from, next) => {
  const isAuth = !!LocalStorage.getItem('token');
  if (to.meta.requiresAuth && !isAuth) {
    next('/login');
  } else {
    next();
  }
});

export default router;