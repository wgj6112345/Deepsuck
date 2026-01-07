<template>
  <div class="sidebar" :class="{ closed: !open }">
    <!-- 品牌区域 -->
    <div class="sidebar-header">
      <div class="brand-logo" @click="handleToggleSidebar" :class="{ 'clickable': !open }">
        <svg class="brand-icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
        </svg>
        <span class="brand-text">deepseek</span>
      </div>
      <button
        class="header-action-btn toggle-btn"
        @click="handleToggleSidebar"
        :title="open ? '收起边栏' : '展开边栏'"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline :points="open ? '15 18 9 12 15 6' : '9 18 15 12 9 6'" />
        </svg>
      </button>
      <!-- <button
        class="header-action-btn new-chat-icon-btn"
        @click="handleNewChat"
        title="开启新对话"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 5v14M5 12h14" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button> -->
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
          :class="{ active: conv.id === currentId, pinned: conv.pinned }"
        >
          <div class="conversation-content" @click="handleSelectConversation(conv.id)">
            <svg v-if="conv.pinned" class="pin-icon" viewBox="0 0 24 24" fill="currentColor">
              <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
            </svg>
            <svg v-else class="conversation-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <input
              v-if="editingId === conv.id"
              v-model="editingTitle"
              class="conversation-title-input"
              @blur="saveRename(conv.id)"
              @keydown.enter="saveRename(conv.id)"
              @click.stop
              ref="titleInputRef"
            />
            <div v-else class="conversation-title">{{ conv.title }}</div>
          </div>
          <button
            class="conversation-more-btn"
            @click.stop="toggleConversationMenu(conv.id)"
            title="更多选项"
          >
            <svg viewBox="0 0 24 24" fill="currentColor">
              <circle cx="12" cy="5" r="2"/>
              <circle cx="12" cy="12" r="2"/>
              <circle cx="12" cy="19" r="2"/>
            </svg>
          </button>
          <!-- 对话操作菜单 -->
          <div
            v-if="activeMenuId === conv.id"
            class="conversation-menu"
          >
            <div class="menu-item" @click="handleRename(conv.id)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
              </svg>
              <span>重命名</span>
            </div>
            <div class="menu-item" @click="handlePin(conv.id)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
              </svg>
              <span>置顶</span>
            </div>
            <div class="menu-item" @click="handleShare(conv.id)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="18" cy="5" r="3"/>
                <circle cx="6" cy="12" r="3"/>
                <circle cx="18" cy="19" r="3"/>
                <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/>
                <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/>
              </svg>
              <span>分享</span>
            </div>
            <div class="menu-item delete" @click="handleDelete(conv.id)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              </svg>
              <span>删除</span>
            </div>
          </div>
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

    <!-- 删除确认弹窗 -->
    <div v-if="showDeleteDialog" class="delete-dialog-overlay" @click="closeDeleteDialog">
      <div class="delete-dialog" @click.stop>
        <div class="delete-dialog-header">
          <h3>永久删除对话</h3>
        </div>
        <div class="delete-dialog-body">
          <p>删除后，该对话将不可恢复。确认删除吗？</p>
        </div>
        <div class="delete-dialog-footer">
          <button class="dialog-button cancel-button" @click="closeDeleteDialog">取消</button>
          <button class="dialog-button delete-button" @click="confirmDelete">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import type { Conversation } from '../store/conversation'
import { conversationAPI } from '../api/conversation'

const props = defineProps<{
  conversations: Conversation[]
  currentId: string | null
  open: boolean
}>()

const emit = defineEmits<{
  newChat: []
  selectConversation: [id: string]
  deleteConversation: [id: string]
  toggleSidebar: []
}>()

const router = useRouter()
const showMenu = ref(false)
const activeMenuId = ref<string | null>(null)
const editingId = ref<string | null>(null)
const editingTitle = ref('')
const titleInputRef = ref<HTMLInputElement | null>(null)
const showDeleteDialog = ref(false)
const deleteConversationId = ref<string | null>(null)

// 点击外部关闭菜单
function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('.user-section')) {
    showMenu.value = false
  }
  if (!target.closest('.conversation-item')) {
    activeMenuId.value = null
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
    '置顶': [],
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
    // 置顶的对话单独分组
    if (conv.pinned) {
      groups['置顶'].push(conv)
      return
    }

    // 使用 updatedAt 而不是 createdAt 来判断分组
    const convDate = new Date(conv.updatedAt)
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
    if (groups[key as keyof typeof groups].length > 0) {
      result[key] = groups[key]!
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

function handleToggleSidebar() {
  emit('toggleSidebar')
}

function handleSelectConversation(id: string) {
  emit('selectConversation', id)
}

function handleDeleteConversation(id: string) {
  emit('deleteConversation', id)
}

// 对话菜单相关方法
function toggleConversationMenu(convId: string) {
  if (activeMenuId.value === convId) {
    activeMenuId.value = null
  } else {
    activeMenuId.value = convId
  }
}

function handleRename(convId: string) {
  activeMenuId.value = null
  const conv = props.conversations.find(c => c.id === convId)
  if (conv) {
    editingId.value = convId
    editingTitle.value = conv.title
    // 聚焦输入框
    nextTick(() => {
      if (titleInputRef.value) {
        titleInputRef.value.focus()
        titleInputRef.value.select()
      }
    })
  }
}

async function saveRename(convId: string) {
  if (editingTitle.value.trim() && editingTitle.value !== props.conversations.find(c => c.id === convId)?.title) {
    try {
      // 调用 API 更新标题
      await conversationAPI.updateConversation(convId, { title: editingTitle.value.trim() })
      // 更新本地状态
      const conv = props.conversations.find(c => c.id === convId)
      if (conv) {
        conv.title = editingTitle.value.trim()
      }
    } catch (error) {
      alert('重命名失败')
    }
  }
  editingId.value = null
  editingTitle.value = ''
}

async function handlePin(convId: string) {
  activeMenuId.value = null
  try {
    const updatedConv = await conversationAPI.togglePin(convId)
    // 更新本地状态
    const conv = props.conversations.find(c => c.id === convId)
    if (conv) {
      conv.pinned = updatedConv.pinned
    }
  } catch (error) {
    alert('置顶失败')
  }
}

function handleShare(convId: string) {
  activeMenuId.value = null
  // TODO: 实现分享功能
  alert('分享功能开发中')
}

function handleDelete(convId: string) {
  activeMenuId.value = null
  deleteConversationId.value = convId
  showDeleteDialog.value = true
}

function closeDeleteDialog() {
  showDeleteDialog.value = false
  deleteConversationId.value = null
}

function confirmDelete() {
  if (deleteConversationId.value) {
    emit('deleteConversation', deleteConversationId.value)
  }
  closeDeleteDialog()
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
  transition: width 0.3s ease;
  flex-shrink: 0;
}

.sidebar.closed {
  width: 56px;
  min-width: 56px;
  max-width: 56px;
  overflow: hidden;
  border-right: none;
}

.sidebar.closed .new-chat-section,
.sidebar.closed .conversation-list,
.sidebar.closed .toggle-btn,
.sidebar.closed .new-chat-icon-btn {
  display: none;
}

/* 品牌区域 */
.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #E5E7EB;
}

.sidebar.closed .sidebar-header {
  justify-content: center;
  padding: 12px 8px;
  border-bottom: none;
}

.sidebar.closed .brand-logo {
  width: 40px;
  height: 40px;
  justify-content: center;
}

.sidebar.closed .brand-icon {
  width: 24px;
  height: 24px;
}

.sidebar.closed .brand-text {
  display: none;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.brand-logo.clickable {
  cursor: pointer;
  transition: all 0.2s;
}

.brand-logo.clickable:hover {
  transform: scale(1.05);
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

.sidebar.closed .brand-text {
  display: none;
}

/* Header中的收起/展开按钮 */
.header-action-btn {
  width: 32px;
  height: 32px;
  background-color: #F9FAFB;
  border: 1px solid #E5E7EB;
  border-radius: 6px;
  color: #6B7280;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.header-action-btn:hover {
  background-color: #F3F4F6;
  color: #1F2937;
  border-color: #D1D5DB;
}

.header-action-btn svg {
  width: 18px;
  height: 18px;
}

/* 新对话图标按钮（仅在收起状态下显示） - 已注释 */
/*
.new-chat-icon-btn {
  display: none;
}

.sidebar.closed .new-chat-icon-btn {
  display: flex;
  width: 40px;
  height: 40px;
  background-color: #F3F4F6;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
}

.sidebar.closed .new-chat-icon-btn svg {
  width: 20px;
  height: 20px;
}
*/

/* 收起/展开按钮在收起状态下的样式 - 已注释 */
/*
.sidebar.closed .toggle-btn {
  width: 40px;
  height: 40px;
  background-color: #F3F4F6;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
}

.sidebar.closed .toggle-btn svg {
  width: 20px;
  height: 20px;
}
*/

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
  position: relative;
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

.conversation-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.conversation-more-btn {
  width: 24px;
  height: 24px;
  background-color: transparent;
  border: none;
  color: #6B7280;
  cursor: pointer;
  padding: 0;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  flex-shrink: 0;
  opacity: 0;
}

.conversation-item:hover .conversation-more-btn {
  opacity: 1;
}

.conversation-more-btn:hover {
  background-color: #F3F4F6;
  color: #1F2937;
}

.conversation-more-btn svg {
  width: 16px;
  height: 16px;
}

.conversation-icon {
  width: 16px;
  height: 16px;
  color: #6B7280;
  flex-shrink: 0;
  transition: color 0.2s;
}

.pin-icon {
  width: 16px;
  height: 16px;
  color: #F59E0B;
  flex-shrink: 0;
}

.conversation-item.pinned .conversation-title {
  font-weight: 600;
  color: #1F2937;
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

.conversation-title-input {
  flex: 1;
  font-size: 14px;
  color: #1F2937;
  padding: 2px 4px;
  border: 1px solid #4A6CF7;
  border-radius: 4px;
  outline: none;
  background-color: #FFFFFF;
  line-height: 1.4;
}

.conversation-title-input:focus {
  box-shadow: 0 0 0 2px rgba(74, 108, 247, 0.1);
}

.conversation-menu-button {
  width: 24px;
  height: 24px;
  background-color: transparent;
  border: none;
  color: #6B7280;
  cursor: pointer;
  padding: 0;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  flex-shrink: 0;
  opacity: 0;
}

.conversation-item:hover .conversation-menu-button {
  opacity: 1;
}

.conversation-menu-button:hover {
  background-color: #F3F4F6;
  color: #1F2937;
}

.conversation-menu-button svg {
  width: 16px;
  height: 16px;
}

/* 对话操作菜单 */
.conversation-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 4px;
  background-color: #F5F5F5;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  min-width: 160px;
  z-index: 100;
  overflow: hidden;
  padding: 8px 0;
}

.conversation-menu .menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  cursor: pointer;
  transition: background-color 0.2s;
  color: #333333;
  font-size: 15px;
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Helvetica Neue', 'Segoe UI', sans-serif;
}

.conversation-menu .menu-item:hover {
  background-color: #E0E0E0;
}

.conversation-menu .menu-item svg {
  width: 16px;
  height: 16px;
  color: #333333;
  flex-shrink: 0;
}

.conversation-menu .menu-item.delete {
  color: #E53935;
}

.conversation-menu .menu-item.delete svg {
  color: #E53935;
}

.conversation-menu .menu-item.delete:hover {
  background-color: #FFEBEE;
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

/* 删除确认弹窗 */
.delete-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.delete-dialog {
  background-color: #FFFFFF;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  min-width: 320px;
  max-width: 90%;
  padding: 24px;
}

.delete-dialog-header {
  margin-bottom: 16px;
}

.delete-dialog-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1F2937;
}

.delete-dialog-body {
  margin-bottom: 24px;
}

.delete-dialog-body p {
  margin: 0;
  font-size: 14px;
  color: #6B7280;
  line-height: 1.5;
}

.delete-dialog-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.dialog-button {
  padding: 8px 20px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Helvetica Neue', 'Segoe UI', sans-serif;
}

.cancel-button {
  background-color: #F3F4F6;
  color: #111827;
  border: 1px solid #D1D5DB;
}

.cancel-button:hover {
  background-color: #E5E7EB;
}

.delete-button {
  background-color: #FFFFFF;
  color: #EF4444;
  border: 2px solid #EF4444;
}

.delete-button:hover {
  background-color: #FEF2F2;
}
</style>