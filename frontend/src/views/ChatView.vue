<template>
  <div class="chat-view">
    <Sidebar
      :conversations="conversationStore.conversations"
      :current-id="conversationStore.currentConversationId"
      :open="uiStore.sidebarOpen"
      @new-chat="handleNewChat"
      @select-conversation="handleSelectConversation"
      @delete-conversation="handleDeleteConversation"
      @toggle-sidebar="handleToggleSidebar"
    />
    
    <div class="chat-content" :class="{ 'sidebar-closed': !uiStore.sidebarOpen }">
      <MessageList 
        ref="messageListRef"
        :messages="currentMessages" 
        :conversation-title="currentConversation?.title"
        @scroll="handleScroll"
      />
      <ChatInput
        :key="conversationStore.currentConversationId"
        :disabled="loading"
        :show-scroll-button="showScrollButton"
        :thinking-enabled="thinkingEnabled"
        @send="handleSendMessage"
        @stop="handleStop"
        @scroll-to-bottom="scrollToBottom"
        @toggle-thinking="handleToggleThinking"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useConversationStore } from '../store/conversation'
import { useUIStore } from '../store/ui'
import { conversationAPI } from '../api/conversation'
import { streamChat, stopChat } from '../api/chat'
import type { Message } from '../store/conversation'
import MessageList from '../components/MessageList.vue'
import ChatInput from '../components/ChatInput.vue'
import Sidebar from '../components/Sidebar.vue'

const conversationStore = useConversationStore()
const uiStore = useUIStore()

const loading = ref(false)
const thinkingEnabled = ref(false)
let abortController: AbortController | null = null

const currentConversation = computed(() => conversationStore.currentConversation)
const currentMessages = computed(() => currentConversation.value?.messages || [])

const messageListRef = ref<InstanceType<typeof MessageList> | null>(null)
const showScrollButton = ref(false)

// 监听 loading 状态，当 agent 输出时显示按钮
watch(loading, (newLoading) => {
  showScrollButton.value = newLoading
})

onMounted(async () => {
  uiStore.loadThinkingEnabled()
  uiStore.loadSidebarOpen()
  await loadConversations()
  
  // 从 localStorage 读取上次保存的对话ID
  const savedConversationId = localStorage.getItem('currentConversationId')
  
  // 如果没有对话或没有当前对话，创建并进入默认对话
  if (conversationStore.conversations.length === 0) {
    await handleNewChat()
  } else if (savedConversationId) {
    // 如果有保存的对话ID，尝试恢复
    const savedConv = conversationStore.conversations.find(c => c.id === savedConversationId)
    if (savedConv) {
      conversationStore.setCurrentConversation(savedConv.id)
    } else {
      // 保存的对话不存在，切换到第一个对话
      const firstConv = conversationStore.conversations[0]
      if (firstConv) {
        conversationStore.setCurrentConversation(firstConv.id)
      }
    }
  } else if (!conversationStore.currentConversationId && conversationStore.conversations.length > 0) {
    // 没有保存的对话，切换到第一个对话
    const firstConv = conversationStore.conversations[0]
    if (firstConv) {
      conversationStore.setCurrentConversation(firstConv.id)
    }
  }
})

function handleScroll() {
  // 检查是否滚动到底部
  if (messageListRef.value) {
    const isNearBottom = messageListRef.value.isNearBottom()
    // 如果 agent 正在输出且不在底部，显示按钮
    showScrollButton.value = loading.value && !isNearBottom
  }
}

function scrollToBottom() {
  if (messageListRef.value) {
    messageListRef.value.scrollToBottom()
    // 点击按钮后立即隐藏按钮
    showScrollButton.value = false
  }
}

async function loadConversations() {
  try {
    conversationStore.setLoading(true)
    const conversations = await conversationAPI.getConversations()
    conversationStore.setConversations(conversations)
  } catch (error) {
    console.error('Error loading conversations:', error)
  } finally {
    conversationStore.setLoading(false)
  }
}

async function handleNewChat() {
  // 如果当前对话没有消息，直接使用当前对话，不创建新对话
  if (currentConversation.value && currentConversation.value.messages.length === 0) {
    return
  }

  try {
    const newConv = await conversationAPI.createConversation()
    conversationStore.addConversation(newConv)
    conversationStore.setCurrentConversation(newConv.id)
    // 保存当前对话ID到 localStorage
    localStorage.setItem('currentConversationId', newConv.id)
  } catch {
    // Error creating conversation
  }
}

async function handleSelectConversation(id: string) {
  conversationStore.setCurrentConversation(id)
  // 保存当前对话ID到 localStorage
  localStorage.setItem('currentConversationId', id)
}

async function handleDeleteConversation(id: string) {
  // 如果只有一个对话，不允许删除
  if (conversationStore.conversations.length <= 1) {
    alert('至少保留一个对话')
    return
  }
  
  try {
    await conversationAPI.deleteConversation(id)
    conversationStore.removeConversation(id)
    
    // 如果删除的是当前对话，切换到第一个对话
    if (conversationStore.currentConversationId === id && conversationStore.conversations.length > 0) {
      const firstConv = conversationStore.conversations[0]
      if (firstConv) {
        conversationStore.setCurrentConversation(firstConv.id)
        // 保存新的当前对话ID到 localStorage
        localStorage.setItem('currentConversationId', firstConv.id)
      }
    } else if (conversationStore.currentConversationId === id) {
      // 删除的是当前对话且没有其他对话了
      localStorage.removeItem('currentConversationId')
    }
  } catch {
    // Error deleting conversation
  }
}

async function handleSendMessage(content: string, enabled: boolean) {
  if (!currentConversation.value) {
    await handleNewChat()
    return
  }

  loading.value = true
  abortController = new AbortController()
  const convId = currentConversation.value.id

  // 创建临时用户消息
  const userMsg: Message = {
    id: `temp-${Date.now()}`,
    conversationId: convId,
    role: 'user',
    content,
    thinking: '',
    thinkingEnabled: false,
    timestamp: new Date().toISOString()
  }
  conversationStore.addMessage(convId, userMsg)

  // 创建临时助手消息
  const assistantMsg: Message = {
    id: `temp-assistant-${Date.now()}`,
    conversationId: convId,
    role: 'assistant',
    content: '',
    thinking: '',
    thinkingEnabled: enabled,
    timestamp: new Date().toISOString()
  }
  conversationStore.addMessage(convId, assistantMsg)

  try {
    for await (const event of streamChat({
      conversationId: convId,
      content,
      thinkingEnabled: enabled
    }, abortController.signal)) {
      if (event.event === 'thinking') {
        conversationStore.updateMessage(convId, assistantMsg.id, {
          thinking: assistantMsg.thinking + event.data
        })
        assistantMsg.thinking += event.data
      } else if (event.event === 'content') {
        conversationStore.updateMessage(convId, assistantMsg.id, {
          content: assistantMsg.content + event.data
        })
        assistantMsg.content += event.data
      } else if (event.event === 'done') {
        conversationStore.updateMessage(convId, assistantMsg.id, {
          id: event.data
        })
      } else if (event.event === 'title_update') {
        // 更新对话标题
        conversationStore.updateConversationTitle(convId, event.data)
      } else if (event.event === 'error') {
        alert('发送失败: ' + event.data)
      }
    }
  } catch (error: unknown) {
    if (error instanceof Error && error.name === 'AbortError') {
      // Chat stopped by user
    } else {
      alert('发送失败，请检查网络和配置')
    }
  } finally {
    loading.value = false
    abortController = null
  }
}

function handleToggleThinking(enabled: boolean) {
  thinkingEnabled.value = enabled
}

function handleStop() {
  if (abortController) {
    abortController.abort()
  }

  // 调用 stop 接口通知后端取消 Agent 请求
  if (currentConversation.value) {
    stopChat(currentConversation.value.id)
      .then(async () => {
        // 停止后重新获取当前对话，获取更新后的标题
        const updatedConv = await conversationAPI.getConversation(currentConversation.value.id)
        conversationStore.updateConversation(updatedConv)
      })
      .catch(error => {
        console.error('Failed to stop chat:', error)
      })
  }
}

function handleToggleSidebar() {
  uiStore.toggleSidebar()
}
</script>

<style scoped>
.chat-view {
  display: flex;
  flex-direction: row;
  width: 100%;
  height: 100vh;
  background-color: #FFFFFF;
  color: #1F2937;
}

.chat-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: #FFFFFF;
  min-width: 0;
  width: 100%;
  padding-left: 280px;
  box-sizing: border-box;
}

.chat-content.sidebar-closed {
  padding-left: 56px;
}
</style>