<template>
  <div class="flex flex-center">
    <q-card class="w-full max-w-md">
      <q-card-section class="q-pa-md">
        <div class="text-center mb-6">
          <h2 class="text-2xl font-bold">文件存储系统</h2>
          <p class="text-gray-500">登录后管理您的文件</p>
        </div>

        <q-form @submit="handleLogin" class="q-gutter-md">
          <q-input
            v-model="form.username"
            label="用户名"
            type="text"
            required
            filled
          />
          <q-input
            v-model="form.password"
            label="密码"
            type="password"
            required
            filled
          />
          <q-btn
            label="登录"
            type="submit"
            color="primary"
            class="full-width"
            :loading="loading"
          />
        </q-form>

        <div class="text-center mt-4 padding-top">
          <router-link to="/register" class="text-primary">
            还没有账号？立即注册
          </router-link>
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar,LocalStorage } from 'quasar';
import { userApi } from 'src/api';
import type { LoginParams } from 'src/api';
// import Cookies from 'js-cookie';

const $q = useQuasar();
const router = useRouter();

// 状态定义（带类型）
const loading = ref<boolean>(false);
const form = ref<LoginParams>({
  username: '',
  password: ''
});

// 登录处理函数
const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    $q.notify({
      type: 'warning',
      message: '请输入用户名和密码'
    });
    return;
  }

  loading.value = true;

    const res = await userApi.login(form.value);
    // 存储Token和用户信息
    LocalStorage.set('token', res.data);
    LocalStorage.set('userInfo', JSON.stringify({
      userId: 0,
      username: form.value.username
    }));
    
    $q.notify({
      type: 'positive',
      message: '登录成功'
    });
    await router.push('/');

    loading.value = false;
  
};
</script>

<style lang="css">
.padding-top {
  margin-top: 10px;
}
</style>