<template>
  <div class="chat-input-container">
    <div class="input-wrapper">
      <textarea
        v-model="inputContent"
        placeholder="给 DeepSeek 发送消息"
        class="input-field"
        @keydown.enter.prevent="handleEnter"
        rows="1"
        ref="textareaRef"
      />
      <div class="input-actions">
        <div class="left-actions">
          <button
            class="feature-button"
            :class="{ active: thinkingEnabled }"
            title="深度思考"
            @click="toggleThinking"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="3"/>
              <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
            </svg>
            <span>深度思考</span>
          </button>
          <button class="feature-button" title="联网搜索">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <path d="M2 12h20M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
            </svg>
            <span>联网搜索</span>
          </button>
        </div>
        <div class="right-actions">
          <button class="attachment-button" title="附件">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"/>
            </svg>
          </button>
          <button
            class="send-button"
            :class="{ loading: disabled, active: canSend }"
            :disabled="!canSend && !disabled"
            @click="handleClick"
            :title="disabled ? '停止生成' : '发送消息'"
          >
            <svg v-if="!disabled" class="send-icon" viewBox="0 0 24 24" fill="currentColor">
              <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/>
            </svg>
            <svg v-else class="stop-icon" viewBox="0 0 24 24" fill="currentColor">
              <rect x="6" y="6" width="12" height="12" rx="2"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
    
    <!-- 悬浮定位按钮 -->
    <button
      v-if="showScrollButton"
      class="scroll-to-bottom-btn"
      @click="$emit('scrollToBottom')"
      title="查看最新回复"
    >
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M12 5v14M19 12l-7 7M5 12l7 7"/>
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'

const props = defineProps<{
  disabled?: boolean
  showScrollButton?: boolean
  thinkingEnabled?: boolean
}>()

const emit = defineEmits<{
  send: [content: string, thinkingEnabled: boolean]
  stop: []
  scrollToBottom: []
  toggleThinking: [enabled: boolean]
}>()

const inputContent = ref('')
const textareaRef = ref<HTMLTextAreaElement | null>(null)
const internalThinkingEnabled = ref(props.thinkingEnabled ?? false)

const canSend = computed(() => {
  return inputContent.value.trim().length > 0 && inputContent.value.length <= 4096
})

// 监听 props 的变化
watch(() => props.thinkingEnabled, (newVal) => {
  if (newVal !== undefined) {
    internalThinkingEnabled.value = newVal
  }
})

// 自动调整 textarea 高度
watch(inputContent, () => {
  nextTick(() => {
    if (textareaRef.value) {
      textareaRef.value.style.height = 'auto'
      textareaRef.value.style.height = Math.min(textareaRef.value.scrollHeight, 200) + 'px'
    }
  })
})

function handleEnter(event: KeyboardEvent) {
  if (event.shiftKey) {
    return
  }
  if (canSend.value) {
    handleSend()
  }
}

function handleClick() {
  if (props.disabled) {
    // 如果正在输出，先停止
    emit('stop')
    // 如果有输入内容，停止后发送新消息
    if (canSend.value) {
      setTimeout(() => handleSend(), 100)
    }
  } else if (canSend.value) {
    handleSend()
  }
}

function handleSend() {
  if (!canSend.value) return
  
  emit('send', inputContent.value, internalThinkingEnabled.value)
  inputContent.value = ''
  
  // 重置 textarea 高度
  nextTick(() => {
    if (textareaRef.value) {
      textareaRef.value.style.height = 'auto'
    }
  })
}

function toggleThinking() {
  internalThinkingEnabled.value = !internalThinkingEnabled.value
  emit('toggleThinking', internalThinkingEnabled.value)
}
</script>

<style scoped>
.chat-input-container {
  padding: 20px 32px 40px;
  background-color: #FFFFFF;
  width: 100%;
  position: relative;
}

.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 24px;
  width: 100%;
  padding: 24px;
  border-radius: 20px;
  border: 1px solid #EDEDED;
  background-color: #FFFFFF;
  position: relative;
}

.input-field {
  width: 100%;
  background-color: transparent;
  border: none;
  padding: 0;
  color: #333333;
  font-size: 19px;
  resize: none;
  outline: none;
  line-height: 1.5;
  min-height: 40px;
  max-height: 200px;
  overflow-y: auto;
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Helvetica Neue', 'Segoe UI', sans-serif;
}

.input-field:focus {
  outline: none;
  box-shadow: none;
}

.input-field:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.input-field::placeholder {
  color: #999999;
  font-weight: 400;
  font-size: 15px;
}

.input-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.left-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.right-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.feature-button {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: #FFFFFF;
  color: #333333;
  border: 1px solid #E0E0E0;
  border-radius: 20px;
  font-size: 17px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Helvetica Neue', 'Segoe UI', sans-serif;
}

.feature-button:hover {
  background-color: #F5F5F5;
}

.feature-button.active {
  background-color: #EFF6FF;
  border-color: #4A6CF7;
  color: #4A6CF7;
}

.feature-button svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.attachment-button {
  width: 36px;
  height: 36px;
  background-color: transparent;
  color: #1890FF;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.attachment-button:hover {
  background-color: #E6F7FF;
}

.attachment-button svg {
  width: 20px;
  height: 20px;
}

.send-button {
  width: 36px;
  height: 36px;
  background-color: #D9D9D9;
  color: #FFFFFF;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
}

.send-button.active {
  background-color: #4A6CF7;
}

.send-button:hover:not(:disabled) {
  background-color: #BFBFBF;
}

.send-button.active:hover:not(:disabled) {
  background-color: #3B5FCC;
}

.send-button:active:not(:disabled) {
  transform: scale(0.95);
}

.send-button:disabled {
  background-color: #E0E0E0;
  cursor: not-allowed;
  opacity: 0.6;
}

.send-button.loading {
  background-color: #EF4444;
}

.send-button.loading:hover {
  background-color: #DC2626;
}

.send-icon {
  width: 18px;
  height: 18px;
}

.stop-icon {
  width: 18px;
  height: 18px;
}

.scroll-to-bottom-btn {
  position: absolute;
  bottom: calc(100% + 20px);
  left: 50%;
  transform: translateX(-50%);
  width: 48px;
  height: 48px;
  background-color: #FFFFFF;
  color: #6B7280;
  border: 1px solid #E5E7EB;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
  z-index: 50;
}

.scroll-to-bottom-btn:hover {
  background-color: #F9FAFB;
  border-color: #D1D5DB;
  transform: translateX(-50%) scale(1.1);
}

.scroll-to-bottom-btn svg {
  width: 24px;
  height: 24px;
}
</style>