<template>
  <div class="message-list">
    <div v-if="messages.length === 0" class="empty-state">
      <div class="empty-icon">
        <svg viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
        </svg>
      </div>
      <h3 class="empty-title">今天有什么可以帮到你？</h3>
    </div>
    <template v-else>
      <div class="conversation-title">
        {{ conversationTitle }}
      </div>
      <MessageItem
        v-for="message in messages"
        :key="message.id"
        :message="message"
      />
      <div class="ai-disclaimer">
        内容由 AI 生成，请仔细甄别
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import MessageItem from './MessageItem.vue'
import type { Message } from '../store/conversation'

defineProps<{
  messages: Message[]
  conversationTitle?: string
}>()
</script>

<style scoped>
.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 0;
  min-height: 0;
  max-width: 900px;
  margin: 0 auto;
  width: 100%;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  min-height: 400px;
  color: #1F2937;
  text-align: center;
  max-width: 900px;
  margin: 0 auto;
  width: 100%;
}

.empty-icon {
  width: 64px;
  height: 64px;
  color: #3B82F6;
}

.empty-icon svg {
  width: 100%;
  height: 100%;
}

.empty-title {
  font-size: 24px;
  font-weight: 600;
  color: #1F2937;
  letter-spacing: -0.3px;
  line-height: 1.4;
}

.conversation-title {
  font-size: 24px;
  font-weight: 700;
  color: #1F2937;
  text-align: center;
  padding: 24px 32px 16px;
  border-bottom: 1px solid #E5E7EB;
  margin-bottom: 24px;
}

.ai-disclaimer {
  text-align: center;
  padding: 16px 32px 24px;
  font-size: 12px;
  color: #9CA3AF;
}
</style>