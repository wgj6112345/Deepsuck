import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUIStore = defineStore('ui', () => {
  const thinkingEnabled = ref(false)
  const sidebarOpen = ref(true)

  function setThinkingEnabled(value: boolean) {
    thinkingEnabled.value = value
    // 持久化到 localStorage
    localStorage.setItem('thinkingEnabled', String(value))
  }

  function toggleSidebar() {
    sidebarOpen.value = !sidebarOpen.value
    // 持久化到 localStorage
    localStorage.setItem('sidebarOpen', String(sidebarOpen.value))
  }

  function setSidebarOpen(value: boolean) {
    sidebarOpen.value = value
    // 持久化到 localStorage
    localStorage.setItem('sidebarOpen', String(value))
  }

  // 从 localStorage 恢复侧边栏状态
  function loadSidebarOpen() {
    const saved = localStorage.getItem('sidebarOpen')
    if (saved !== null) {
      sidebarOpen.value = saved === 'true'
    }
  }

  // 从 localStorage 恢复思考模式状态
  function loadThinkingEnabled() {
    const saved = localStorage.getItem('thinkingEnabled')
    if (saved !== null) {
      thinkingEnabled.value = saved === 'true'
    }
  }

  return {
    thinkingEnabled,
    sidebarOpen,
    setThinkingEnabled,
    toggleSidebar,
    setSidebarOpen,
    loadThinkingEnabled
  }
})