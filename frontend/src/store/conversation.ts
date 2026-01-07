import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Message {
  id: string
  conversationId: string
  role: 'user' | 'assistant'
  content: string
  thinking: string
  thinkingEnabled: boolean
  timestamp: string
}

export interface Conversation {
  id: string
  title: string
  messages: Message[]
  createdAt: string
  updatedAt: string
  pinned: boolean
}

export const useConversationStore = defineStore('conversation', () => {
  const conversations = ref<Conversation[]>([])
  const currentConversationId = ref<string | null>(null)
  const loading = ref(false)

  const currentConversation = computed(() => {
    return conversations.value.find(c => c.id === currentConversationId.value) || null
  })

  function setConversations(data: Conversation[]) {
    conversations.value = data
  }

  function setCurrentConversation(id: string | null) {
    currentConversationId.value = id
  }

  function addConversation(conv: Conversation) {
    conversations.value = [conv, ...conversations.value]
  }

  function updateConversation(conv: Conversation) {
    const index = conversations.value.findIndex(c => c.id === conv.id)
    if (index !== -1) {
      conversations.value = [
        ...conversations.value.slice(0, index),
        conv,
        ...conversations.value.slice(index + 1)
      ]
    }
  }

  function removeConversation(id: string) {
    conversations.value = conversations.value.filter(c => c.id !== id)
    if (currentConversationId.value === id) {
      currentConversationId.value = null
    }
  }

  function addMessage(convId: string, message: Message) {
    const conv = conversations.value.find(c => c.id === convId)
    if (conv) {
      const messages = conv.messages || []
      updateConversation({
        ...conv,
        messages: [...messages, message],
        updatedAt: new Date().toISOString()
      })
    }
  }

  function updateMessage(convId: string, messageId: string, updates: Partial<Message>) {
    const conv = conversations.value.find(c => c.id === convId)
    if (conv) {
      const messages = conv.messages || []
      const updatedMessages = messages.map(msg =>
        msg.id === messageId ? { ...msg, ...updates } : msg
      )
      updateConversation({
        ...conv,
        messages: updatedMessages
      })
    }
  }

  function updateConversationTitle(convId: string, title: string) {
    const conv = conversations.value.find(c => c.id === convId)
    if (conv) {
      updateConversation({
        ...conv,
        title: title,
        updatedAt: new Date().toISOString()
      })
    }
  }

  function setLoading(value: boolean) {
    loading.value = value
  }

  return {
    conversations,
    currentConversationId,
    currentConversation,
    loading,
    setConversations,
    setCurrentConversation,
    addConversation,
    updateConversation,
    removeConversation,
    addMessage,
    updateMessage,
    updateConversationTitle,
    setLoading
  }
})