<template>
  <div class="sidebar" :class="{ closed: !open }">
    <!-- 品牌区域 -->
    <div class="sidebar-header">
      <div class="brand-logo">
        <svg class="brand-icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
        </svg>
        <span class="brand-text">deepseek</span>
      </div>
    </div>

    <!-- 新对话按钮 -->
    <div class="new-chat-section">
      <button class="new-chat-button" @click="handleNewChat">
        <svg class="new-chat-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 5v14M5 12h14" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span>开启新对话</span>
      </button>
    </div>

    <!-- 对话历史列表 -->
    <div class="conversation-list">
      <div v-for="(group, groupName) in groupedConversations" :key="groupName" class="conversation-group">
        <div class="group-title">{{ groupName }}</div>
        <div
          v-for="conv in group"
          :key="conv.id"
          class="conversation-item"
          :class="{ active: conv.id === currentId }"
          @click="handleSelectConversation(conv.id)"
        >
          <svg class="conversation-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <div class="conversation-title">{{ conv.title }}</div>
        </div>
      </div>
    </div>

    <!-- 用户信息 -->
    <div class="user-section">
      <div class="user-email">227*****32@qq.com</div>
      <button class="menu-button" @click="toggleMenu" title="更多选项">
        <svg viewBox="0 0 24 24" fill="currentColor">
          <circle cx="12" cy="5" r="2"/>
          <circle cx="12" cy="12" r="2"/>
          <circle cx="12" cy="19" r="2"/>
        </svg>
      </button>
      <!-- 用户菜单弹窗 -->
      <div v-if="showMenu" class="user-menu">
        <div class="menu-item" @click="handleDownloadApp">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="5" y="2" width="14" height="20" rx="2" ry="2"/>
            <line x1="12" y1="18" x2="12.01" y2="18"/>
          </svg>
          <span>下载手机应用</span>
        </div>
        <div class="menu-item" @click="handleSettings">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
          <span>系统设置</span>
        </div>
        <div class="menu-item" @click="handleContact">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 2L11 13"/>
            <path d="M22 2l-7 20-4-9-9-4 20-7z"/>
          </svg>
          <span>联系我们</span>
        </div>
        <div class="menu-item" @click="handleLogout">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
            <polyline points="16 17 21 12 16 7"/>
            <line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
          <span>退出登录</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import type { Conversation } from '../store/conversation'

const props = defineProps<{
  conversations: Conversation[]
  currentId: string | null
  open: boolean
}>()

const emit = defineEmits<{
  newChat: []
  selectConversation: [id: string]
  deleteConversation: [id: string]
}>()

const router = useRouter()
const showMenu = ref(false)

// 点击外部关闭菜单
function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('.user-section')) {
    showMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

// 按时间分组对话
const groupedConversations = computed(() => {
  const groups: Record<string, Conversation[]> = {
    '今天': [],
    '7天内': [],
    '30天内': [],
    '更早': []
  }

  const now = new Date()
  const oneDay = 24 * 60 * 60 * 1000
  const sevenDays = 7 * oneDay
  const thirtyDays = 30 * oneDay

  props.conversations.forEach(conv => {
    const convDate = new Date(conv.createdAt)
    const diff = now.getTime() - convDate.getTime()

    if (diff < oneDay) {
      groups['今天'].push(conv)
    } else if (diff < sevenDays) {
      groups['7天内'].push(conv)
    } else if (diff < thirtyDays) {
      groups['30天内'].push(conv)
    } else {
      groups['更早'].push(conv)
    }
  })

  // 移除空分组
  const result: Record<string, Conversation[]> = {}
  Object.keys(groups).forEach(key => {
    if (groups[key].length > 0) {
      result[key] = groups[key]
    }
  })

  return result
})

function toggleMenu() {
  showMenu.value = !showMenu.value
}

function handleDownloadApp() {
  showMenu.value = false
  alert('下载手机应用功能开发中')
}

function handleSettings() {
  showMenu.value = false
  router.push('/settings')
}

function handleContact() {
  showMenu.value = false
  alert('联系我们功能开发中')
}

function handleLogout() {
  showMenu.value = false
  if (confirm('确定要退出登录吗？')) {
    alert('退出登录功能开发中')
  }
}

function handleNewChat() {
  emit('newChat')
}

function handleSelectConversation(id: string) {
  emit('selectConversation', id)
}

function handleDeleteConversation(id: string) {
  emit('deleteConversation', id)
}
</script>

<style scoped>
.sidebar {
  width: 280px;
  min-width: 280px;
  max-width: 280px;
  background-color: #FFFFFF;
  border-right: 1px solid #E5E7EB;
  display: flex;
  flex-direction: column;
  transition: width var(--ds-transition-slow) ease;
  flex-shrink: 0;
}

.sidebar.closed {
  width: 0;
  min-width: 0;
  max-width: 0;
  overflow: hidden;
}

/* 品牌区域 */
.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  padding: 16px 20px;
  border-bottom: 1px solid #E5E7EB;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.brand-icon {
  width: 24px;
  height: 24px;
  color: #4A6CF7;
}

.brand-text {
  font-size: 18px;
  font-weight: 600;
  color: #4A6CF7;
  letter-spacing: -0.3px;
}

/* 新对话按钮区域 */
.new-chat-section {
  padding: 16px 20px;
}

.new-chat-button {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 16px;
  background-color: #F3F4F6;
  color: #1F2937;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.new-chat-button:hover {
  background-color: #E5E7EB;
  transform: translateY(-1px);
}

.new-chat-button:active {
  transform: translateY(0);
}

.new-chat-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

/* 对话列表 */
.conversation-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 12px;
}

.conversation-group {
  margin-bottom: 16px;
}

.group-title {
  font-size: 12px;
  font-weight: 500;
  color: #6B7280;
  padding: 8px 8px 4px;
  letter-spacing: 0.3px;
}

.conversation-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 2px;
}

.conversation-item:hover {
  background-color: #F3F4F6;
}

.conversation-item.active {
  background-color: #EFF6FF;
  color: #4A6CF7;
}

.conversation-item.active .conversation-icon {
  color: #4A6CF7;
}

.conversation-icon {
  width: 16px;
  height: 16px;
  color: #6B7280;
  flex-shrink: 0;
  transition: color 0.2s;
}

.conversation-title {
  flex: 1;
  font-size: 14px;
  color: #1F2937;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
}

/* 用户信息区域 */
.user-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-top: 1px solid #E5E7EB;
  background-color: #FFFFFF;
  position: relative;
  margin-top: auto;
  flex-shrink: 0;
}

.user-email {
  font-size: 14px;
  color: #1F2937;
  font-weight: 500;
}

.menu-button {
  width: 32px;
  height: 32px;
  background-color: transparent;
  border: none;
  color: #6B7280;
  cursor: pointer;
  padding: 0;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.menu-button:hover {
  background-color: #F3F4F6;
  color: #1F2937;
}

.menu-button svg {
  width: 20px;
  height: 20px;
}

/* 用户菜单弹窗 */
.user-menu {
  position: absolute;
  bottom: 100%;
  left: 0;
  margin-bottom: 8px;
  background-color: #FFFFFF;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  min-width: 200px;
  z-index: 100;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
  color: #1F2937;
  font-size: 14px;
}

.menu-item:hover {
  background-color: #F3F4F6;
}

.menu-item svg {
  width: 18px;
  height: 18px;
  color: #6B7280;
  flex-shrink: 0;
}

.menu-item span {
  flex: 1;
}
</style>