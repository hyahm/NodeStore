// src/types/shims-window.d.ts
declare global {
  interface Window {
    // 直接声明 $cookies 的核心方法，避免依赖库的导出类型
    $cookies: {
      get: (key: string) => string | null
      set: (key: string, value: string, expire?: string | number) => void
      remove: (key: string) => void
      clear: () => void
    }
  }
}

// 确保文件被识别为模块
export {}