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

    <!-- 侧边栏 -->
    <q-drawer
      v-model="drawerOpen"
      show-if-above
      bordered
      class="bg-grey-1"
      content-class="bg-white"
    >
      <q-list dense class="pt-0">
        <!-- 侧边栏导航项 -->

        <q-item
          clickable
          v-ripple
          :to="{ name: 'FileManagement' }"
          :active="isActive('/')"
        >
          <q-item-section avatar>
            <q-icon name="storage" />
          </q-item-section>
          <q-item-section>文件管理</q-item-section>
        </q-item>

        
        <q-item
          clickable
          v-ripple
          :to="{ name: 'NodeManagement' }"
          :active="isActive('/node-management')"
        >
          <q-item-section avatar>
            <q-icon name="storage" />
          </q-item-section>
          <q-item-section>节点管理</q-item-section>
        </q-item>


         
        
        <!-- 可添加其他导航项 -->
        <!-- <q-item
          clickable
          v-ripple
          :to="{ name: 'OtherPage' }"
          :active="isActive('/other-page')"
        >
          <q-item-section avatar>
            <q-icon name="folder" />
          </q-item-section>
          <q-item-section>其他页面</q-item-section>
        </q-item> -->
      </q-list>
    </q-drawer>

    <!-- 主内容区 -->
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useQuasar, LocalStorage } from 'quasar';

const $q = useQuasar();
const router = useRouter();
const route = useRoute();

// 侧边栏开关状态
const drawerOpen = ref(true);

/**
 * 判断当前路由是否激活
 * @param path 路由路径
 */
const isActive = (path: string) => {
  return route.path === path;
};

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