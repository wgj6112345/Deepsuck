<template>
  <div class="message-item" :class="message.role">
    <div class="message-content">
      <div v-if="message.thinking && message.thinkingEnabled" class="thinking-wrapper">
        <ThinkingPanel
          :thinking="message.thinking"
          :is-thinking="!message.content"
        />
      </div>
      <div class="message-text">
        <!-- 如果是 assistant 且 content 为空，显示思考中提示 -->
        <div v-if="message.role === 'assistant' && !message.content && (!message.thinking || !message.thinkingEnabled)" class="thinking-indicator">
          <span class="thinking-text">正在思考中</span>
          <div class="thinking-dots">
            <span class="dot"></span>
            <span class="dot"></span>
            <span class="dot"></span>
          </div>
        </div>
        <!-- 否则显示实际内容 -->
        <MarkdownRenderer v-else :content="message.content" />
      </div>
      <div v-if="message.role === 'assistant'" class="message-actions">
        <button class="action-button" @click="handleCopy" title="复制">
          <svg v-if="!copied" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
        </button>
        <button class="action-button" @click="handleRetry" title="重试">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="23 4 23 10 17 10"/>
            <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
          </svg>
        </button>
        <button class="action-button" @click="handleLike" title="点赞">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"/>
          </svg>
        </button>
        <button class="action-button" @click="handleDislike" title="点踩">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10 15v4a3 3 0 0 0 3 3l4-9V2H5.72a2 2 0 0 0-2 1.7l-1.38 9a2 2 0 0 0 2 2.3zm7-13h2.67A2.31 2.31 0 0 1 22 4v7a2.31 2.31 0 0 1-2.33 2H17"/>
          </svg>
        </button>
        <button class="action-button" @click="handleShare" title="分享">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="18" cy="5" r="3"/>
            <circle cx="6" cy="12" r="3"/>
            <circle cx="18" cy="19" r="3"/>
            <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/>
            <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import MarkdownRenderer from './MarkdownRenderer.vue'
import ThinkingPanel from './ThinkingPanel.vue'
import type { Message } from '../store/conversation'

const props = defineProps<{
  message: Message
}>()

const copied = ref(false)

function handleCopy() {
  navigator.clipboard.writeText(props.message.content)
    .then(() => {
      copied.value = true
      setTimeout(() => {
        copied.value = false
      }, 1000)
    })
    .catch(() => {
      // 复制失败，不做任何提示
    })
}

function handleRetry() {
  alert('重试功能开发中')
}

function handleLike() {
  alert('点赞功能开发中')
}

function handleDislike() {
  alert('点踩功能开发中')
}

function handleShare() {
  alert('分享功能开发中')
}
</script>

<style scoped>
.message-item {
  display: flex;
  gap: 12px;
  padding: 12px 32px;
  border-radius: 0;
  margin-bottom: 0;
  transition: background-color 0.2s;
}

.message-item.user {
  flex-direction: row-reverse;
  justify-content: flex-start;
}

.message-item.assistant {
  background-color: #FFFFFF;
}

.message-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: #4A6CF7;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: white;
}

.message-avatar svg {
  width: 18px;
  height: 18px;
}

.message-content {
  flex: 1;
  min-width: 0;
}

.message-item.user .message-content {
  flex: 0 1 auto;
  max-width: 70%;
  background-color: #E3F0FF;
  border-radius: 16px;
  padding: 10px 16px;
}

.message-item.assistant .message-content {
  background-color: transparent;
  padding: 0;
  width: 100%;
  max-width: 100%;
}

.thinking-wrapper {
  margin-bottom: 16px;
}

.message-text {
  line-height: 1.7;
  color: #1F2937;
  font-size: 17px;
}

.message-text :deep(p) {
  margin: 12px 0;
}

.message-text :deep(p:first-child) {
  margin-top: 0;
}

.message-text :deep(p:last-child) {
  margin-bottom: 0;
}

.message-text :deep(h1),
.message-text :deep(h2),
.message-text :deep(h3),
.message-text :deep(h4),
.message-text :deep(h5),
.message-text :deep(h6) {
  margin-top: 24px;
  margin-bottom: 12px;
  font-weight: 600;
  color: #1F2937;
}

.message-text :deep(h1) {
  font-size: 25px;
}

.message-text :deep(h2) {
  font-size: 21px;
}

.message-text :deep(h3) {
  font-size: 19px;
}

.message-actions {
  display: flex;
  gap: 16px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #E5E7EB;
}

.action-button {
  width: 32px;
  height: 32px;
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
}

.action-button:hover {
  background-color: #F3F4F6;
  color: #1F2937;
}

.action-button svg {
  width: 18px;
  height: 18px;
}

.thinking-indicator {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
}

.thinking-dots {
  display: flex;
  gap: 6px;
  align-items: center;
}

.dot {
  width: 8px;
  height: 8px;
  background-color: #9CA3AF;
  border-radius: 50%;
  animation: bounce 1.4s infinite ease-in-out both;
}

.dot:nth-child(1) {
  animation-delay: -0.32s;
}

.dot:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes bounce {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

.thinking-text {
  font-size: 15px;
  color: #6B7280;
  font-weight: 500;
}
</style>