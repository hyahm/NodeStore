<template>
  <div class="flex flex-center">
    <q-card class="w-full max-w-md">
      <q-card-section class="q-pa-md">
        <div class="text-center mb-6">
          <h2 class="text-2xl font-bold">用户注册</h2>
          <p class="text-gray-500">创建您的文件存储账号</p>
        </div>

        <q-form @submit="handleRegister" class="q-gutter-md">
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
            label="注册"
            type="submit"
            color="primary"
            class="full-width"
            :loading="loading"
          />
        </q-form>

        <div class="text-center mt-4">
          <router-link to="/login" class="text-primary">
            已有账号？立即登录
          </router-link>
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import { userApi } from 'src/api';
import type { LoginParams } from 'src/api';

const $q = useQuasar();
const router = useRouter();

// 状态定义
const loading = ref<boolean>(false);
const form = ref<LoginParams>({
  username: '',
  password: ''
});

/**
 * 处理注册
 */
const handleRegister = async (): Promise<void> => {
  if (!form.value.username || !form.value.password) {
    $q.notify({
      type: 'warning',
      message: '请输入用户名和密码'
    });
    return;
  }

  loading.value = true;
    await userApi.register(form.value);
    $q.notify({
      type: 'positive',
      message: '注册成功，请登录'
    });
    await router.push('/login');
    loading.value = false;
};
</script>