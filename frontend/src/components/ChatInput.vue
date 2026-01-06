<template>
  <div class="chat-input-container">
    <div class="chat-input-wrapper">
      <div class="input-area">
        <textarea
          v-model="inputContent"
          placeholder="给 DeepSeek 发送消息"
          class="input-textarea"
          :disabled="disabled"
          @keydown.enter.prevent="handleEnter"
          rows="1"
          ref="textareaRef"
        />
        <div class="input-actions">
          <button class="feature-button" title="深度思考">
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
      </div>
      <div class="input-right">
        <button class="attachment-button" title="附件">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"/>
          </svg>
        </button>
        <button
          class="send-button"
          :class="{ loading: disabled }"
          :disabled="!canSend && !disabled"
          @click="handleClick"
          :title="disabled ? '停止生成' : '发送消息'"
        >
          <svg v-if="!disabled" class="send-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 2L11 13M22 2l-7 20-4-9-9-4 20-7z" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <svg v-else class="stop-icon" viewBox="0 0 24 24" fill="currentColor">
            <rect x="6" y="6" width="12" height="12" rx="2"/>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'

const props = defineProps<{
  disabled?: boolean
}>()

const emit = defineEmits<{
  send: [content: string]
  stop: []
}>()

const inputContent = ref('')
const textareaRef = ref<HTMLTextAreaElement | null>(null)

const canSend = computed(() => {
  return inputContent.value.trim().length > 0 && inputContent.value.length <= 4096
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
    emit('stop')
  } else if (canSend.value) {
    handleSend()
  }
}

function handleSend() {
  if (!canSend.value) return
  
  emit('send', inputContent.value)
  inputContent.value = ''
  
  // 重置 textarea 高度
  nextTick(() => {
    if (textareaRef.value) {
      textareaRef.value.style.height = 'auto'
    }
  })
}
</script>

<style scoped>
.chat-input-container {
  padding: 20px 32px;
  background-color: #FFFFFF;
}

.chat-input-wrapper {
  display: flex;
  gap: 12px;
  align-items: flex-end;
  max-width: 900px;
  margin: 0 auto;
}

.input-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
}

.input-textarea {
  flex: 1;
  background-color: #FFFFFF;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 12px 16px;
  color: #1F2937;
  font-size: 16px;
  resize: none;
  outline: none;
  transition: all 0.2s;
  line-height: 1.5;
  min-height: 48px;
  max-height: 200px;
  overflow-y: auto;
}

.input-textarea:focus {
  border-color: #4A6CF7;
  box-shadow: 0 0 0 3px rgba(74, 108, 247, 0.1);
}

.input-textarea:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background-color: #F9FAFB;
}

.input-textarea::placeholder {
  color: #9CA3AF;
}

.input-actions {
  display: flex;
  gap: 8px;
  padding: 0 4px;
}

.feature-button {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background-color: #F3F4F6;
  color: #1F2937;
  border: none;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.feature-button:hover {
  background-color: #E5E7EB;
}

.feature-button svg {
  width: 16px;
  height: 16px;
}

.input-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.attachment-button {
  width: 40px;
  height: 40px;
  background-color: transparent;
  color: #6B7280;
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
  background-color: #F3F4F6;
  color: #1F2937;
}

.attachment-button svg {
  width: 20px;
  height: 20px;
}

.send-button {
  width: 40px;
  height: 40px;
  background-color: #4A6CF7;
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
}

.send-button:hover:not(:disabled) {
  background-color: #3B5BD8;
  transform: translateY(-1px);
}

.send-button:active:not(:disabled) {
  transform: translateY(0);
}

.send-button:disabled {
  background-color: #D1D5DB;
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
  width: 20px;
  height: 20px;
  margin-left: 1px;
  margin-top: 1px;
}

.stop-icon {
  width: 20px;
  height: 20px;
}
</style>