// src/api/node.ts
// import { api } from 'boot/axios' // 基于 Quasar 内置的 axios 封装，若你有自定义请求实例可替换
import service, { type ApiResponse } from './request';

// 定义节点数据类型
export interface Nodes {
    url: string
    alive: boolean
}

export const nodeApi = {
  // 上传文件
  getNodeList(): Promise<ApiResponse<Nodes[]>> {
    return service({
      url: '/node/list',
      method: 'get',
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
  },




  // 按MD5下载文件
  addNode(addr: string): Promise<unknown> {
    return service({
      url: `/node/add`,
      method: 'get',
      params: {addr}
    });
  },


  // 软删除文件
  removeNode(addr: string): Promise<ApiResponse<unknown>> {
    return service({
      url: '/node/remove',
      method: 'get',
      params: {addr}
    });
  },


};

// // 查询节点列表
// export const getNodeList = async (): Promise<string[]> => {
//   try {
//     const { data } = await api.get('/node/list')
//     return data
//   } catch (error) {
//     console.error('获取节点列表失败:', error)
//     throw error // 抛出错误让页面层处理
//   }
// }

// 新增节点
// export const addNode = async (url: string): Promise<boolean> => {
//   try {
//     await api.get('/node/add', { params: url }) // 后端用 Get 接收，故用 params 传参
//     return true
//   } catch (error) {
//     console.error('新增节点失败:', error)
//     throw error
//   }
// }

// 删除节点
// export const removeNode = async (url: string ): Promise<boolean> => {
//   try {
//     await api.get('/node/remove', { params: { url: url } })
//     return true
//   } catch (error) {
//     console.error('删除节点失败:', error)
//     throw error
//   }
// }