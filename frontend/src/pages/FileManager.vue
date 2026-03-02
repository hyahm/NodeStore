<template>
  <div class="file-manager q-pa-lg">
    <!-- 顶部导航栏 - 玻璃拟态效果 -->
    <div class="nav-header q-mb-lg">
      <q-card class="nav-card" flat bordered>
        <q-card-section class="row items-center no-wrap q-py-sm">
          <!-- 面包屑导航 -->
          <div class="breadcrumb-container col">
            <div class="text-caption text-grey-7 q-mb-xs">当前位置</div>
            <div class="breadcrumb-wrapper row items-center">
              <q-btn flat dense icon="home" color="primary" @click="goHome" size="sm" class="q-mr-sm" />
              <q-icon name="chevron_right" size="16px" color="grey-5" class="q-mr-sm" />
              <template v-for="(segment, index) in pathSegments" :key="index">
                <q-btn 
                  flat 
                  dense 
                  no-caps 
                  :label="segment.name" 
                  color="primary" 
                  class="breadcrumb-item text-weight-medium"
                  @click="navigateTo(index)"
                />
                <q-icon v-if="index < pathSegments.length - 1" name="chevron_right" size="16px" color="grey-5" class="q-mx-xs" />
              </template>
            </div>
          </div>

          <q-space />

          <!-- 操作按钮组 -->
          <div class="action-buttons row items-center gap-sm">
            <q-btn 
              label="新建文件夹" 
              icon="create_new_folder" 
              color="primary" 
              unelevated
              rounded
              class="action-btn"
              @click="createDirDialog = true"
            >
              <q-tooltip>创建新目录</q-tooltip>
            </q-btn>
            
            <q-separator vertical class="q-mx-sm" />
            
            <q-uploader
              ref="uploaderRef"
              class="upload-btn"
              style="max-width: 150px"
              flat
              bordered
              auto-upload
              hide-upload-btn
            
              :headers="uploadHeaders"
              field-name="file"
              @added="handleFileAdded"
              @uploaded="handleUploaded"
              @failed="handleUploadFailed"
            >
              <!-- <template v-slot:header="scope">
                <q-btn 
                  label="上传文件" 
                  icon="cloud_upload" 
                  color="secondary" 
                  unelevated
                  rounded
                  class="full-width"
                  @click="scope.pickFiles"
                />
              </template> -->
            </q-uploader>
          </div>
        </q-card-section>
      </q-card>
    </div>

    <!-- 统计信息卡片 -->
    <div class="stats-row row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="stat-card folder-card" flat bordered>
          <q-card-section class="row items-center">
            <div class="stat-icon bg-amber-1 text-amber-8">
              <q-icon name="folder" size="24px" />
            </div>
            <div class="q-ml-md">
              <div class="text-caption text-grey-6">目录数量</div>
              <div class="text-h6 text-weight-bold">{{ dirList.length }}</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="stat-card file-card" flat bordered>
          <q-card-section class="row items-center">
            <div class="stat-icon bg-blue-1 text-blue-8">
              <q-icon name="insert_drive_file" size="24px" />
            </div>
            <div class="q-ml-md">
              <div class="text-caption text-grey-6">文件数量</div>
              <div class="text-h6 text-weight-bold">{{ fileList.length }}</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="stat-card storage-card" flat bordered>
          <q-card-section class="row items-center">
            <div class="stat-icon bg-green-1 text-green-8">
              <q-icon name="storage" size="24px" />
            </div>
            <div class="q-ml-md">
              <div class="text-caption text-grey-6">总大小</div>
              <div class="text-h6 text-weight-bold">{{ totalSize }}</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <q-card class="stat-card media-card" flat bordered>
          <q-card-section class="row items-center">
            <div class="stat-icon bg-purple-1 text-purple-8">
              <q-icon name="perm_media" size="24px" />
            </div>
            <div class="q-ml-md">
              <div class="text-caption text-grey-6">媒体文件</div>
              <div class="text-h6 text-weight-bold">{{ mediaCount }}</div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- 目录区域 -->
    <q-card v-if="dirList.length > 0" class="section-card q-mb-lg" flat bordered>
      <q-card-section class="section-header row items-center">
        <q-icon name="folder_open" size="20px" color="amber-7" class="q-mr-sm" />
        <span class="text-subtitle1 text-weight-medium">文件夹</span>
        <q-space />
        <q-chip size="sm" color="amber-1" text-color="amber-9">{{ dirList.length }} 个</q-chip>
      </q-card-section>
      
      <q-separator />
      
      <q-card-section>
        <div class="folder-grid">
          <q-intersection 
            v-for="dir in dirList" 
            :key="dir"
            transition="scale"
            class="folder-item-wrapper"
          >
            <q-card 
              class="folder-card-item cursor-pointer"
              flat
              bordered
              @click="enterDir(dir)"
              v-ripple
            >
              <q-card-section class="text-center q-pa-md">
                <div class="folder-icon-wrapper q-mb-sm">
                  <q-icon name="folder" size="48px" color="amber-5" />
                  <q-icon name="folder_open" size="48px" color="amber-6" class="folder-hover-icon" />
                </div>
                <div class="folder-name text-body2 text-weight-medium text-truncate-2-lines" :title="getDirName(dir)">
                  {{ getDirName(dir) }}
                </div>
                <div class="text-caption text-grey-6 q-mt-xs">点击进入</div>
              </q-card-section>
            </q-card>
          </q-intersection>
        </div>
      </q-card-section>
    </q-card>

    <!-- 文件列表 -->
    <q-card class="section-card" flat bordered>
      <q-card-section class="section-header row items-center">
        <q-icon name="description" size="20px" color="blue-7" class="q-mr-sm" />
        <span class="text-subtitle1 text-weight-medium">文件</span>
        <q-space />
        <q-input
          v-model="searchQuery"
          dense
          outlined
          placeholder="搜索文件..."
          class="search-input"
          clearable
        >
          <template v-slot:prepend>
            <q-icon name="search" size="18px" color="grey-5" />
          </template>
        </q-input>
      </q-card-section>

      <q-separator />

      <q-card-section class="q-pa-none">
        <q-table
          :rows="filteredFiles"
          :columns="columns"
          row-key="id"
          flat
          bordered
          class="file-table"
          :pagination="{ rowsPerPage: 10 }"
          :loading="loading"
        >
          <!-- 文件名列 -->
          <template v-slot:body-cell-filename="props">
            <q-td :props="props" class="filename-cell">
              <div class="row items-center no-wrap">
                <q-icon 
                  :name="getFileIcon(props.row.filename)" 
                  :color="getFileColor(props.row.filename)"
                  size="24px"
                  class="q-mr-sm"
                />
                <div>
                  <div class="text-body2 text-weight-medium text-truncate" style="max-width: 200px">
                    {{ props.row.filename }}
                  </div>
                  <div class="text-caption text-grey-6" v-if="isMediaFile(props.row.filename)">
                    媒体文件
                  </div>
                </div>
              </div>
            </q-td>
          </template>

          <!-- 大小列 -->
          <template v-slot:body-cell-size="props">
            <q-td :props="props">
              <q-badge color="grey-3" text-color="grey-8" class="text-weight-medium">
                {{ formatFileSize(props.row.file_size) }}
              </q-badge>
            </q-td>
          </template>

          <!-- MD5列 -->
          <template v-slot:body-cell-md5="props">
            <q-td :props="props">
              <code class="md5-code">{{ props.row.md5.slice(0, 8) }}...</code>
              <q-tooltip>{{ props.row.md5 }}</q-tooltip>
            </q-td>
          </template>

          <!-- 时间列 -->
          <template v-slot:body-cell-uploadTime="props">
            <q-td :props="props">
              <div class="text-body2">{{ formatDate(props.row.upload_time) }}</div>
              <div class="text-caption text-grey-6">{{ formatTime(props.row.upload_time) }}</div>
            </q-td>
          </template>

          <!-- 操作列 -->
          <template v-slot:body-cell-actions="props">
            <q-td :props="props" class="actions-cell">
              <div class="row items-center justify-end no-wrap gap-xs">
                <q-btn
                  round
                  flat
                  size="sm"
                  color="primary"
                  icon="download"
                  @click="handleDownload(props.row)"
                >
                  <q-tooltip>下载</q-tooltip>
                </q-btn>
                
                <q-btn
                  v-if="isMediaFile(props.row.filename)"
                  round
                  flat
                  size="sm"
                  color="green"
                  icon="play_circle"
                  @click="handleStream(props.row)"
                >
                  <q-tooltip>播放</q-tooltip>
                </q-btn>
                
                <q-btn
                  round
                  flat
                  size="sm"
                  color="blue"
                  icon="share"
                  @click="handleShare(props.row)"
                >
                  <q-tooltip>分享</q-tooltip>
                </q-btn>
                
                <q-btn
                  round
                  flat
                  size="sm"
                  color="negative"
                  icon="delete"
                  @click="handleDelete(props.row)"
                >
                  <q-tooltip>删除</q-tooltip>
                </q-btn>
              </div>
            </q-td>
          </template>

          <!-- 空状态 -->
          <template v-slot:no-data>
            <div class="full-width row flex-center q-pa-lg text-grey-6">
              <q-icon name="inbox" size="48px" class="q-mb-sm" />
              <div class="text-h6">暂无文件</div>
              <div class="text-caption">点击上方"上传文件"按钮添加</div>
            </div>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- 创建目录弹窗 - 美化版 -->
    <q-dialog v-model="createDirDialog" persistent>
      <q-card class="dialog-card" style="min-width: 400px">
        <q-card-section class="dialog-header bg-primary text-white">
          <div class="row items-center">
            <q-icon name="create_new_folder" size="28px" class="q-mr-sm" />
            <div class="text-h6">新建文件夹</div>
          </div>
        </q-card-section>
        
        <q-card-section class="q-pt-lg">
          <q-input
            v-model="newDirName"
            label="文件夹名称"
            outlined
            autofocus
            clearable
            :rules="[val => !!val || '请输入名称']"
            @keyup.enter="handleCreateDir"
          >
            <template v-slot:prepend>
              <q-icon name="folder" color="amber" />
            </template>
          </q-input>
        </q-card-section>
        
        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="取消" color="grey-7" v-close-popup />
          <q-btn 
            unelevated 
            label="创建" 
            color="primary" 
            :disable="!newDirName.trim()"
            @click="handleCreateDir"
            icon="add"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- 播放弹窗 - 沉浸式播放器 -->
    <q-dialog v-model="streamDialog" maximized transition-show="slide-up" transition-hide="slide-down">
      <q-card class="media-player-card bg-dark text-white">
        <q-bar class="bg-black">
          <q-icon name="play_circle" />
          <div class="text-weight-bold">{{ currentStreamFile }}</div>
          <q-space />
          <q-btn dense flat icon="close" v-close-popup>
            <q-tooltip>关闭</q-tooltip>
          </q-btn>
        </q-bar>
        
        <q-card-section class="flex flex-center q-pa-none" style="height: calc(100% - 40px)">
          <div class="media-container">
            <video
              v-if="isVideoFile(currentStreamFile) && streamUrl"
              :src="streamUrl"
              controls
              class="media-player"
              autoplay
              playsinline
            />
            <audio
              v-else-if="isAudioFile(currentStreamFile) && streamUrl"
              :src="streamUrl"
              controls
              class="media-player audio-player"
              autoplay
            />
            <div v-else class="text-center">
              <q-spinner size="48px" color="primary" />
              <div class="q-mt-md">加载中...</div>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </q-dialog>

    <!-- 分享弹窗 -->
    <q-dialog v-model="shareDialog">
      <q-card class="dialog-card" style="min-width: 450px">
        <q-card-section class="dialog-header bg-blue text-white">
          <div class="row items-center">
            <q-icon name="share" size="28px" class="q-mr-sm" />
            <div class="text-h6">分享文件</div>
          </div>
        </q-card-section>
        
        <q-card-section class="q-pt-lg">
          <div class="text-subtitle2 text-grey-8 q-mb-sm">分享链接</div>
          <q-input
            v-model="shareLink"
            outlined
            readonly
            type="textarea"
            rows="3"
            class="share-input"
          />
          
          <div class="row q-mt-md q-col-gutter-sm">
            <div class="col">
              <q-btn 
                outline 
                color="primary" 
                icon="content_copy" 
                label="复制链接" 
                class="full-width"
                @click="copyShareLink"
              />
            </div>
            <div class="col">
              <q-btn 
                outline 
                color="secondary" 
                icon="qr_code" 
                label="生成二维码" 
                class="full-width"
                @click="showQR = true"
              />
            </div>
          </div>

          <!-- 二维码展示 -->
          <q-slide-transition>
            <div v-show="showQR" class="q-mt-md flex flex-center">
              <q-card flat bordered class="q-pa-md">
                <!-- 这里可以集成 QRCode 组件 -->
                <div class="text-center text-grey-6">
                  <q-icon name="qr_code_2" size="120px" />
                  <div class="text-caption q-mt-sm">二维码预览</div>
                </div>
              </q-card>
            </div>
          </q-slide-transition>
        </q-card-section>
        
        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="关闭" color="grey-7" v-close-popup @click="showQR = false" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useQuasar, LocalStorage, date } from 'quasar';
import { fsApi, fileApi } from 'src/api';
import type { FileItem } from 'src/api';

const $q = useQuasar();
const token = LocalStorage.getItem("token") as string;

// ========== 状态定义 ==========
const currentDir = ref<string>('/');
const dirList = ref<string[]>([]);
const fileList = ref<FileItem[]>([]);
const loading = ref(false);
const searchQuery = ref('');

// 弹窗状态
const createDirDialog = ref(false);
const newDirName = ref('');
const streamDialog = ref(false);
const streamUrl = ref('');
const currentStreamFile = ref('');
const shareDialog = ref(false);
const shareLink = ref('');
const showQR = ref(false);

// 上传相关
// const uploaderRef = ref(null);
// const uploadUrl = computed(() => {
//     console.log(process.env.API_URL)
//     return `${process.env.API_URL}/file/upload`
// });
const uploadHeaders = computed(() => [{ name: 'Authorization', value: `Bearer ${token}` }]);

// ========== 计算属性 ==========

// 路径分段（面包屑用）
const pathSegments = computed(() => {
  if (currentDir.value === '/') return [{ name: '根目录', path: '/' }];
  
  const parts = currentDir.value.split('/').filter(Boolean);
  return parts.map((name, index) => ({
    name,
    path: '/' + parts.slice(0, index + 1).join('/')
  }));
});

// 过滤后的文件列表
const filteredFiles = computed(() => {
  if (!searchQuery.value) return fileList.value;
  const query = searchQuery.value.toLowerCase();
  return fileList.value.filter(f => 
    f.filename.toLowerCase().includes(query)
  );
});

// 总大小
const totalSize = computed(() => {
  const total = fileList.value.reduce((sum, f) => sum + (f.file_size || 0), 0);
  return formatFileSize(total);
});

// 媒体文件数量
const mediaCount = computed(() => 
  fileList.value.filter(f => isMediaFile(f.filename)).length
);

// ========== 表格列配置 ==========
const columns = [
  { 
    name: 'filename', 
    label: '文件名', 
    field: 'filename', 
    align: 'left' as const,
    sortable: true,
    style: 'width: 30%'
  },
  { 
    name: 'size', 
    label: '大小', 
    field: 'file_size', 
    align: 'left' as const,
    sortable: true,
    style: 'width: 15%'
  },
  { 
    name: 'md5', 
    label: '文件标识', 
    field: 'md5', 
    align: 'left' as const,
    style: 'width: 20%'
  },
  { 
    name: 'uploadTime', 
    label: '上传时间', 
    field: 'upload_time', 
    align: 'left' as const,
    sortable: true,
    style: 'width: 20%'
  },
  { 
    name: 'actions', 
    label: '操作', 
    field: 'actions', 
    align: 'right' as const,
    style: 'width: 15%'
  }
];

// ========== 工具函数 ==========

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const formatDate = (timeStr: string): string => {
  return date.formatDate(timeStr, 'YYYY-MM-DD');
};

const formatTime = (timeStr: string): string => {
  return date.formatDate(timeStr, 'HH:mm:ss');
};

const isMediaFile = (filename: string): boolean => {
  const ext = filename.slice(filename.lastIndexOf('.')).toLowerCase();
  const mediaExts = ['.mp3', '.mp4', '.avi', '.mov', '.m4a', '.wav', '.flac', '.webm', '.mkv'];
  return mediaExts.includes(ext);
};

const isVideoFile = (filename: string): boolean => {
  const ext = filename.slice(filename.lastIndexOf('.')).toLowerCase();
  const videoExts = ['.mp4', '.avi', '.mov', '.webm', '.mkv'];
  return videoExts.includes(ext);
};

const isAudioFile = (filename: string): boolean => {
  const ext = filename.slice(filename.lastIndexOf('.')).toLowerCase();
  const audioExts = ['.mp3', '.m4a', '.wav', '.flac', '.ogg'];
  return audioExts.includes(ext);
};

const getDirName = (path: string): string => {
  return path.split('/').filter(Boolean).pop() || '/';
};

// 根据文件类型返回图标
const getFileIcon = (filename: string): string => {
  const ext = filename.slice(filename.lastIndexOf('.')).toLowerCase();
  const iconMap: Record<string, string> = {
    '.mp4': 'movie', '.avi': 'movie', '.mov': 'movie', '.webm': 'movie',
    '.mp3': 'audiotrack', '.wav': 'audiotrack', '.flac': 'audiotrack',
    '.jpg': 'image', '.jpeg': 'image', '.png': 'image', '.gif': 'image',
    '.pdf': 'picture_as_pdf',
    '.zip': 'folder_zip', '.rar': 'folder_zip', '.7z': 'folder_zip',
    '.doc': 'description', '.docx': 'description',
    '.xls': 'table_chart', '.xlsx': 'table_chart',
    '.ppt': 'slideshow', '.pptx': 'slideshow'
  };
  return iconMap[ext] || 'insert_drive_file';
};

// 根据文件类型返回颜色
const getFileColor = (filename: string): string => {
  const ext = filename.slice(filename.lastIndexOf('.')).toLowerCase();
  const colorMap: Record<string, string> = {
    '.mp4': 'red', '.avi': 'red', '.mov': 'red',
    '.mp3': 'orange', '.wav': 'orange',
    '.jpg': 'purple', '.jpeg': 'purple', '.png': 'purple',
    '.pdf': 'red-6',
    '.zip': 'amber',
    '.doc': 'blue', '.docx': 'blue',
    '.xls': 'green', '.xlsx': 'green'
  };
  return colorMap[ext] || 'grey-6';
};

// ========== 导航功能 ==========

const goHome = async () => {
  currentDir.value = '/';
  await loadFsList();
};

const navigateTo = async (index: number) => {
  const segments = pathSegments.value;
  if (index < 0 || index >= segments.length) return;
  currentDir.value = (segments[index] as {name: string, path: string}).path;
  await loadFsList();
};

// ========== 数据加载 ==========

const loadFsList = async () => {
  loading.value = true;
  try {
    
    const res = await fsApi.getFsList(currentDir.value);
    dirList.value = res.data.dirs || [];
    fileList.value = res.data.files || [];
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: '加载失败：' + (err as Error).message,
      position: 'top'
    });
  } finally {
    loading.value = false;
  }
};
const enterDir = async (dirPath: string) => {
    // console.log(dirPath)
    currentDir.value = dirPath
    // console.log( currentDir.value)
    // if (currentDir.value.slice(-1) == '/') {
    //     currentDir.value = currentDir.value + '/' + dirPath;
    // } else {
    //     currentDir.value = currentDir.value +  dirPath;
    // }
    // console.log( currentDir.value)
  await loadFsList();
};

// ========== 文件操作 ==========

const handleCreateDir = async () => {
  if (!newDirName.value.trim()) return;

//   const fullPath = currentDir.value === '/' 
//     ? `/${newDirName.value}` 
//     : `${currentDir.value}/${newDirName.value}`;

  try {
    await fsApi.createDir({path: currentDir.value});
    $q.notify({ 
      type: 'positive', 
      message: '文件夹创建成功',
      position: 'top-right',
      timeout: 2000
    });
    createDirDialog.value = false;
    newDirName.value = '';
    await loadFsList();
  } catch (err) {
    $q.notify({ 
      type: 'negative', 
      message: '创建失败：' + (err as Error).message,
      position: 'top'
    });
  }
};

const handleFileAdded = async (files: readonly File[]) => {
  if (!files.length) return;
  
  // 这里使用自定义上传逻辑，而不是依赖 q-uploader 的自动上传
  for (const file of files) {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('dir', currentDir.value);
    
      await fileApi.upload(formData).then(() => {
        $q.notify({
        type: 'positive',
        message: `${file.name} 上传成功`,
        position: 'top-right',
        timeout: 2000
      });
      });
      
  }
  
  await loadFsList();
//   uploaderRef.value?.reset();
};

const handleUploaded = async () => {
  await loadFsList();
};

const handleUploadFailed = () => {
  $q.notify({
    type: 'negative',
    message: '上传失败：未知错误',
    position: 'top'
  });
};

const handleDownload = async (file: FileItem) => {
  try {
    const blob = await fileApi.downloadByMd5(file.md5);
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = file.filename;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
    
    $q.notify({
      type: 'positive',
      message: '开始下载...',
      position: 'top-right',
      timeout: 1500
    });
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: '下载失败：' + (err as Error).message,
      position: 'top'
    });
  }
};

const handleStream = (file: FileItem) => {
  currentStreamFile.value = file.filename;
  streamUrl.value = fileApi.streamByMd5(file.md5);
  streamDialog.value = true;
};

const handleShare = async (file: FileItem) => {
  try {
    const res = await fileApi.createShare({ md5: file.md5, days: '7', token });
    shareLink.value = res.share_link;
    shareDialog.value = true;
    showQR.value = false;
  } catch (err) {
    $q.notify({
      type: 'negative',
      message: '分享失败：' + (err as Error).message,
      position: 'top'
    });
  }
};

const copyShareLink = () => {
  navigator.clipboard.writeText(shareLink.value)
    .then(() => {
      $q.notify({ 
        type: 'positive', 
        message: '链接已复制到剪贴板',
        position: 'top-right',
        timeout: 2000
      });
    })
    .catch(() => {
      $q.notify({ 
        type: 'warning', 
        message: '复制失败，请手动复制',
        position: 'top'
      });
    });
};

const handleDelete = async (file: FileItem) => {

    await  fileApi.deleteFile({ md5: file.md5, token });
      $q.notify({ 
        type: 'positive', 
        message: '文件已删除',
        position: 'top-right',
        timeout: 2000
      });
  await    loadFsList();
 
};

onMounted(async () => {
  await loadFsList();
});
</script>

<style scoped>
.file-manager {
  max-width: 1400px;
  margin: 0 auto;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);
  min-height: 100vh;
}

/* 导航栏样式 */
.nav-header {
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.breadcrumb-container {
  min-width: 0;
}

.breadcrumb-wrapper {
  flex-wrap: nowrap;
  overflow-x: auto;
  scrollbar-width: none;
}

.breadcrumb-wrapper::-webkit-scrollbar {
  display: none;
}

.breadcrumb-item {
  font-size: 14px;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.3s;
}

.breadcrumb-item:hover {
  background: rgba(25, 118, 210, 0.1);
}

.action-btn {
  padding: 8px 20px;
  font-weight: 500;
  letter-spacing: 0.5px;
  transition: all 0.3s;
}

.action-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(25, 118, 210, 0.3);
}

/* 统计卡片 */
.stats-row {
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 12px;
  transition: all 0.3s;
  background: white;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 区域卡片 */
.section-card {
  border-radius: 12px;
  background: white;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.section-header {
  background: linear-gradient(90deg, rgba(0,0,0,0.02) 0%, rgba(0,0,0,0) 100%);
  padding: 16px 20px;
}

/* 文件夹网格 */
.folder-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 16px;
  padding: 8px;
}

.folder-item-wrapper {
  aspect-ratio: 1;
}

.folder-card-item {
  height: 100%;
  border-radius: 12px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 2px solid transparent;
}

.folder-card-item:hover {
  transform: translateY(-4px) scale(1.02);
  border-color: #ffc107;
  box-shadow: 0 12px 32px rgba(255, 193, 7, 0.2);
}

.folder-icon-wrapper {
  position: relative;
  display: inline-block;
}

.folder-hover-icon {
  position: absolute;
  top: 0;
  left: 0;
  opacity: 0;
  transition: opacity 0.3s;
}

.folder-card-item:hover .folder-hover-icon {
  opacity: 1;
}

.folder-name {
  line-height: 1.4;
  max-height: 2.8em;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* 文件表格 */
.file-table {
  background: transparent;
}

.filename-cell {
  padding: 12px 16px;
}

.md5-code {
  background: #f5f5f5;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-family: 'Courier New', monospace;
  color: #666;
}

.actions-cell {
  padding: 8px;
}

.actions-cell .q-btn {
  opacity: 0.7;
  transition: all 0.2s;
}

.actions-cell .q-btn:hover {
  opacity: 1;
  transform: scale(1.1);
}

.gap-xs {
  gap: 4px;
}

.search-input {
  width: 280px;
  transition: all 0.3s;
}

.search-input:focus-within {
  width: 320px;
}

/* 弹窗样式 */
.dialog-card {
  border-radius: 16px;
  overflow: hidden;
}

.dialog-header {
  padding: 20px;
}

/* 媒体播放器 */
.media-player-card {
  border-radius: 0;
}

.media-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #000;
}

.media-player {
  max-width: 100%;
  max-height: 100%;
  width: auto;
  height: auto;
}

.audio-player {
  width: 80%;
  max-width: 600px;
}

.share-input :deep(.q-field__control) {
  background: #f5f7fa;
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .folder-grid {
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: 12px;
  }
  
  .search-input {
    width: 100%;
    margin-top: 12px;
  }
  
  .action-buttons {
    flex-wrap: wrap;
  }
}

/* 动画 */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.file-manager > * {
  animation: fadeIn 0.4s ease-out forwards;
}
</style>