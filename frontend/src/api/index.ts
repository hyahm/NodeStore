import service, { type ApiResponse } from './request';
// import Cookies from 'js-cookie';
import { LocalStorage } from 'quasar'
const token = localStorage.getItem("token")
// ========== 类型定义 ==========
// 用户信息类型
export interface UserInfo {
  user_id: number;
  username: string;
  token: string;
  expire: string;
}

// 登录/注册参数类型
export interface LoginParams {
  username: string;
  password: string;
}

// 文件信息类型
export interface FileItem {
  id: number;
  user_id: number;
  filename: string;
  encoded_name: string;
  dir: string;
  upload_time: string;
  file_size: number;
  md5: string;
  content_id: number;
  is_deleted: number;
}

// 目录/文件列表响应类型
export interface FsListResponse {
  dirs: string[];
  files: FileItem[];
  dir: string;
  user_id: number;
}

export interface UploadResponse {
  msg: string;
  file: string;
  dir: string;
  md5: string;
  is_new: boolean;
}

// 分享链接响应类型
export interface ShareResponse {
  share_key: string;
  expire_at: string;
  share_link: string;
}

// ========== 接口封装 ==========
// 用户相关
export const userApi = {
  // 注册
  register(data: LoginParams): Promise<ApiResponse<unknown>> {
    return service({
      url: '/user/register',
      method: 'post',
      data
    });
  },

  // 登录
  login(data: LoginParams): Promise<ApiResponse<string>> {
    return service({
      url: '/user/login',
      method: 'post',
      data
    });
  }
};

// 文件相关
export const fileApi = {
  // 上传文件
  upload(formData: FormData): Promise<ApiResponse<UploadResponse>> {
    return service({
      url: '/file/upload',
      method: 'post',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
  },

  // 按MD5下载文件
  downloadByMd5(md5: string): Promise<Blob> {
    return service({
      url: `/file/download_by_md5?token=${token}&md5=${md5}`,
      method: 'get',
      responseType: 'blob'
    });
  },

  // 在线播放（获取文件流地址）
  streamByMd5(md5: string): string {
    const token = LocalStorage.getItem('token') as string|| '';
    return `${service.defaults.baseURL}/file/stream?md5=${md5}&token=${token}`;
  },

  // 软删除文件
  deleteFile(params: { md5?: string; filename?: string; dir?: string,token: string }): Promise<ApiResponse<unknown>> {
    return service({
      url: '/file/delete',
      method: 'get',
      params
    });
  },

  // 创建分享链接
  createShare(params: { md5: string; days?: string, token: string }): Promise<ShareResponse> {
    return service({
      url: '/share/create',
      method: 'get',
      params
    });
  }
};

// 目录/文件列表相关
export const fsApi = {
  // 创建目录
  createDir(data: {path: string}): Promise<ApiResponse<unknown>> {
    return service({
      url: `/dir/create`,
      method: 'post',
      data,
    });
  },

  // 获取文件/目录列表
  getFsList(dir: string): Promise<ApiResponse<FsListResponse>> {
    
    return service({
      url: `/fs/list?dir=${encodeURIComponent(dir || '/')}`,
      method: 'get'
    });
  }
};

// 节点管理
export const nodeApi = {
  addNode(addr: string): Promise<ApiResponse<unknown>> {
    return service({
      url: `/node/add?token=${token}&addr=${encodeURIComponent(addr)}`,
      method: 'get'
    });
  },

  listNodes(): Promise<ApiResponse<string[]>> {
    return service({
      url: '/node/list',
      method: 'get',
      params: {token}
    });
  }
};