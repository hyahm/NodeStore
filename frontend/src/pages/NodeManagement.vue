<!-- src/pages/NodeManagement.vue -->
<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- 页面头部 -->
    <div class="row justify-between items-center q-mb-lg">
      <div>
        <h2 class="text-h4 text-primary q-mt-none q-mb-sm">节点管理</h2>
        <p class="text-subtitle1 text-grey-7">管理存储集群中的节点地址与状态</p>
      </div>
      <q-btn
        label="新增节点"
        color="primary"
        icon="add"
        @click="openAddDialog"
        rounded
        class="shadow-1"
      />
    </div>

    <!-- 节点列表卡片 -->
    <q-card class="shadow-1 rounded-borders">
      <q-card-section >
        <div class="text-h6">节点列表</div>
      </q-card-section>
      <q-table
        :rows="nodeList"
        :columns="columns"
        row-key="url"
        dense
        flat
        :pagination="{ rowsPerPage: 5 }"
        class="rounded-borders"
      >
        <!-- 状态列 -->
        <template v-slot:body-cell-alive="props">
          <q-td :props="props" class="text-center">
            <q-badge
              :color="props.row.alive ? 'green-6' : 'red-6'"
              rounded
              class="text-white"
            >
              {{ props.row.alive ? '在线' : '离线' }}
            </q-badge>
          </q-td>
        </template>

        <!-- 操作列 -->
        <template v-slot:body-cell-actions="props">
          <q-td :props="props" class="text-center">
            <q-btn
              size="sm"
              color="primary"
              icon="edit"
              @click="openEditDialog(props.row)"
              dense
              flat
              class="q-mr-xs"
            />
            <q-btn
              size="sm"
              color="negative"
              icon="delete"
              @click="handleRemoveNode(props.row)"
              dense
              flat
            />
          </q-td>
        </template>
      </q-table>

      <!-- 空数据提示 -->
      <q-empty-state
        v-if="nodeList.length === 0"
        icon="cloud_off"
        label="暂无节点数据"
        icon-color="grey-5"
        class="q-py-xl"
      >
        <q-btn
          label="添加节点"
          color="primary"
          icon="add"
          @click="openAddDialog"
          class="q-mt-md"
        />
      </q-empty-state>
    </q-card>

    <!-- 新增/编辑弹窗 -->
    <q-dialog v-model="addDialogOpen" persistent>
      <q-card class="q-pa-lg" style="width: 500px; max-width: 90vw;">
        <q-card-section>
          <h4 class="text-h5 q-mt-none q-mb-md">
            {{ isEdit ? '编辑节点' : '新增节点' }}
          </h4>
          <q-form ref="addFormRef" @submit="handleAddOrEditNode">
            <q-input
              v-model="nodeUrl"
              label="节点地址"
              placeholder="例如：http://192.168.1.100:8080"
              required
              :rules="[val => val && val.length > 0 || '请输入节点URL']"
              class="q-mb-md"
            />
          </q-form>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn label="取消" flat @click="closeAddDialog" />
          <q-btn
            label="确认"
            color="primary"
            @click="handleAddOrEditNode"
            :loading="submitLoading"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import type { QForm } from 'quasar'
import { nodeApi } from 'src/api/node'
import type { Nodes } from 'src/api/node'

const $q = useQuasar()
const addFormRef = ref<QForm | null>(null)
const submitLoading = ref(false)
const addDialogOpen = ref(false)
const isEdit = ref(false)
const nodeUrl = ref('')
const currentNode = ref<Nodes | null>(null)

const columns = [
  { name: 'url', label: '节点地址', field: 'url', align: 'left'as const },
  { name: 'alive', label: '是否在线', field: 'alive', align: 'center' as const},
  { name: 'actions', label: '操作', field: 'actions', align: 'center' as const }
]

const nodeList = ref<Nodes[]>([])

const openAddDialog = () => {
  isEdit.value = false
  nodeUrl.value = ''
  currentNode.value = null
  addDialogOpen.value = true
}

const openEditDialog = (row: Nodes) => {
  isEdit.value = true
  nodeUrl.value = row.url
  currentNode.value = row
  addDialogOpen.value = true
}

const closeAddDialog = () => {
  addDialogOpen.value = false
  addFormRef.value?.resetValidation()
}

const loadNodeList = async () => {
    const { data } = await nodeApi.getNodeList()
    nodeList.value = data
}

const handleAddOrEditNode = async () => {
  if (!addFormRef.value?.validate()) return
  submitLoading.value = true

    if (isEdit.value && currentNode.value) {
      // 编辑逻辑：先删除旧节点，再添加新节点（适配你的后端接口）
      await nodeApi.removeNode(currentNode.value.url)
      await nodeApi.addNode(nodeUrl.value)
      $q.notify({ type: 'positive', message: '编辑节点成功', position: 'top-right' })
    } else {
      await nodeApi.addNode(nodeUrl.value)
      $q.notify({ type: 'positive', message: '新增节点成功', position: 'top-right' })
    }
    closeAddDialog()
    await loadNodeList()
    submitLoading.value = false
}

const handleRemoveNode = async (row: Nodes) => {
  const confirm = $q.dialog({
    title: '确认删除',
    message: `确定要删除节点 ${row.url} 吗？`,
    cancel: true,
    ok: { label: '删除', color: 'negative' }
  })
  if (!confirm) return

    await nodeApi.removeNode(row.url)
    $q.notify({ type: 'positive', message: '删除节点成功', position: 'top-right' })
    await loadNodeList()
}

onMounted(async () => {
  await loadNodeList()
})
</script>

<style scoped>
.q-page {
  background-color: #f5f7fa;
  min-height: calc(100vh - 64px);
}
</style>