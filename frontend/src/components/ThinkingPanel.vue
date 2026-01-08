<template>
  <div class="thinking-panel">
    <!-- 思考状态栏 -->
    <div class="thinking-header" @click="toggleThinking">
      <div class="thinking-icon-wrapper">
        <svg class="thinking-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm0 18a8 8 0 1 1 8-8 8 8 0 0 1-8 8z"/>
          <path d="M12 6v6l4 2"/>
        </svg>
      </div>
      <span class="thinking-status">{{ isThinking ? '正在思考' : '已思考' }}</span>
      <svg class="dropdown-icon" :class="{ expanded: isExpanded }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M6 9l6 6 6-6" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </div>

    <!-- 思考内容 -->
    <transition name="thinking-expand">
      <div v-show="isExpanded" class="thinking-content">
        <div class="thinking-markdown">
          <MarkdownRenderer :content="thinking" />
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import MarkdownRenderer from './MarkdownRenderer.vue'

const props = defineProps<{
  thinking: string
  isThinking?: boolean
}>()

const isExpanded = ref(true)

function toggleThinking() {
  isExpanded.value = !isExpanded.value
}
</script>

<style scoped>
.thinking-panel {
  margin-bottom: 16px;
}

.thinking-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background-color: #F9FAFB;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
}

.thinking-header:hover {
  background-color: #F3F4F6;
  border-color: #D1D5DB;
}

.thinking-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
}

.thinking-icon {
  width: 16px;
  height: 16px;
  color: #6B7280;
  flex-shrink: 0;
}

.thinking-status {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.dropdown-icon {
  width: 16px;
  height: 16px;
  color: #6B7280;
  transition: transform 0.2s;
  flex-shrink: 0;
}

.dropdown-icon.expanded {
  transform: rotate(180deg);
}

.thinking-content {
  background-color: #F9FAFB;
  border: 1px solid #E5E7EB;
  border-top: none;
  border-radius: 0 0 8px 8px;
  padding: 12px 16px;
}

.thinking-markdown :deep(.markdown-renderer) {
  color: #6B7280;
}

.thinking-markdown :deep(h1),
.thinking-markdown :deep(h2),
.thinking-markdown :deep(h3),
.thinking-markdown :deep(h4),
.thinking-markdown :deep(h5),
.thinking-markdown :deep(h6) {
  color: #6B7280;
}

.thinking-markdown :deep(p) {
  color: #6B7280;
}

.thinking-markdown :deep(li) {
  color: #6B7280;
}

.thinking-markdown :deep(code) {
  color: #6B7280;
}

.thinking-markdown :deep(strong) {
  color: #6B7280;
}

.thinking-markdown :deep(a) {
  color: #409EFF;
}

.thinking-markdown :deep(blockquote) {
  color: #6B7280;
  background-color: #F3F4F6;
}

.thinking-markdown :deep(.code-block-wrapper) {
  background-color: #F3F4F6;
}

.thinking-markdown :deep(.code-content) {
  color: #6B7280;
}

/* 展开动画 */
.thinking-expand-enter-active,
.thinking-expand-leave-active {
  transition: all 0.2s ease-in-out;
  overflow: hidden;
  max-height: 1000px;
  opacity: 1;
}

.thinking-expand-enter-from,
.thinking-expand-leave-to {
  max-height: 0;
  opacity: 0;
}
</style>