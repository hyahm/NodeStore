<template>
  <q-layout view="hHh lpR fFf">
    <!-- 顶部导航栏 -->
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-toolbar-title>文件存储系统</q-toolbar-title>
        <q-space />
        <q-btn label="退出登录" color="white" text-color="primary" @click="handleLogout" />
      </q-toolbar>
    </q-header>

    <!-- 主内容区 -->
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useQuasar, LocalStorage } from 'quasar';
// import Cookies from 'js-cookie';

const $q = useQuasar();
const router = useRouter();

/**
 * 退出登录
 */
const handleLogout = async () => {
  LocalStorage.remove('token');
  LocalStorage.remove('userInfo');
  $q.notify({
    type: 'positive',
    message: '已退出登录'
  });
  await router.push('/login');
};
</script>