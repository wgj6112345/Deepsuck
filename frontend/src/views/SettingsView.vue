<template>
  <div class="settings-view">
    <div class="settings-header">
      <router-link to="/" class="back-link">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M19 12H5M12 19l-7-7 7-7" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span>返回</span>
      </router-link>
      <h1>设置</h1>
    </div>
    <div class="settings-content">
      <div class="settings-section">
        <div class="section-header">
          <svg class="section-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <h2>Agent 配置</h2>
        </div>
        <form @submit.prevent="handleSave">
          <div class="form-group">
            <label for="apiKey">
              <span class="label-text">API Key</span>
              <span class="label-hint">用于身份验证的密钥</span>
            </label>
            <input
              id="apiKey"
              v-model="config.apiKey"
              type="password"
              placeholder="输入你的 API Key"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label for="modelName">
              <span class="label-text">模型名称</span>
              <span class="label-hint">要使用的 AI 模型</span>
            </label>
            <input
              id="modelName"
              v-model="config.modelName"
              type="text"
              placeholder="deepseek-chat"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label for="baseUrl">
              <span class="label-text">Base URL</span>
              <span class="label-hint">API 服务的地址</span>
            </label>
            <input
              id="baseUrl"
              v-model="config.baseUrl"
              type="text"
              placeholder="https://api.deepseek.com"
              class="form-input"
            />
          </div>
          <div class="form-actions">
            <button type="submit" class="save-button" :disabled="loading">
              <svg v-if="!loading" class="save-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M17 21v-8H7v8" stroke-linecap="round" stroke-linejoin="round"/>
                <path d="M7 3v5h8" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              <span v-if="loading">保存中...</span>
              <span v-else>保存配置</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useConfigStore } from '../store/config'
import { configAPI } from '../api/config'
import type { AgentConfig } from '../store/config'

const configStore = useConfigStore()

const loading = ref(false)
const config = ref<AgentConfig>({
  apiKey: '',
  modelName: import.meta.env.VITE_DEFAULT_AGENT_MODEL || 'mimo-v2-flash',
  baseUrl: import.meta.env.VITE_DEFAULT_AGENT_BASE_URL || 'https://api.xiaomimimo.com/v1'
})

onMounted(async () => {
  await loadConfig()
})

async function loadConfig() {
  try {
    loading.value = true
    const data = await configAPI.getConfig()
    config.value = data
  } catch (error) {
    // Error loading config
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  try {
    loading.value = true
    await configAPI.updateConfig(config.value)
    configStore.setConfig(config.value)
    alert('配置已保存')
  } catch (error) {
    alert('保存失败，请重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.settings-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--ds-bg-primary);
  color: var(--ds-text-primary);
}

.settings-header {
  display: flex;
  align-items: center;
  gap: var(--ds-spacing-md);
  padding: var(--ds-spacing-lg) var(--ds-spacing-xl);
  background-color: var(--ds-bg-primary);
  border-bottom: 1px solid var(--ds-border-default);
  height: 60px;
  flex-shrink: 0;
}

.back-link {
  display: flex;
  align-items: center;
  gap: var(--ds-spacing-xs);
  color: var(--ds-primary);
  text-decoration: none;
  font-size: var(--ds-font-base);
  font-weight: var(--ds-font-medium);
  padding: var(--ds-spacing-xs) var(--ds-spacing-sm);
  border-radius: var(--ds-radius-sm);
  transition: all var(--ds-transition-base);
}

.back-link:hover {
  background-color: var(--ds-primary-light);
  text-decoration: none;
}

.back-link svg {
  width: 18px;
  height: 18px;
}

.settings-header h1 {
  flex: 1;
  font-size: var(--ds-font-2xl);
  font-weight: var(--ds-font-semibold);
  margin: 0;
}

.settings-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--ds-spacing-2xl);
}

.settings-section {
  max-width: 600px;
  margin: 0 auto;
}

.section-header {
  display: flex;
  align-items: center;
  gap: var(--ds-spacing-md);
  margin-bottom: var(--ds-spacing-xl);
  padding-bottom: var(--ds-spacing-lg);
  border-bottom: 1px solid var(--ds-border-default);
}

.section-icon {
  width: 24px;
  height: 24px;
  color: var(--ds-primary);
}

.settings-section h2 {
  font-size: var(--ds-font-xl);
  font-weight: var(--ds-font-semibold);
  margin: 0;
  color: var(--ds-text-primary);
}

.form-group {
  margin-bottom: var(--ds-spacing-xl);
}

.form-group label {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  margin-bottom: var(--ds-spacing-sm);
}

.label-text {
  font-size: var(--ds-font-base);
  font-weight: var(--ds-font-medium);
  color: var(--ds-text-primary);
}

.label-hint {
  font-size: var(--ds-font-sm);
  color: var(--ds-text-tertiary);
  font-weight: var(--ds-font-normal);
}

.form-input {
  width: 100%;
  padding: var(--ds-spacing-md) var(--ds-spacing-lg);
  background-color: var(--ds-bg-primary);
  border: 1px solid var(--ds-border-default);
  border-radius: var(--ds-radius-md);
  color: var(--ds-text-primary);
  font-size: var(--ds-font-base);
  outline: none;
  transition: all var(--ds-transition-base);
  font-family: inherit;
}

.form-input:focus {
  border-color: var(--ds-primary);
  box-shadow: 0 0 0 3px rgba(77, 107, 254, 0.1);
}

.form-input::placeholder {
  color: var(--ds-text-tertiary);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: var(--ds-spacing-2xl);
  padding-top: var(--ds-spacing-lg);
  border-top: 1px solid var(--ds-border-default);
}

.save-button {
  display: flex;
  align-items: center;
  gap: var(--ds-spacing-sm);
  padding: var(--ds-spacing-md) var(--ds-spacing-xl);
  background-color: var(--ds-primary);
  color: white;
  border: none;
  border-radius: var(--ds-radius-md);
  font-size: var(--ds-font-base);
  font-weight: var(--ds-font-semibold);
  cursor: pointer;
  transition: all var(--ds-transition-base);
  box-shadow: var(--ds-shadow-sm);
}

.save-button:hover:not(:disabled) {
  background-color: var(--ds-primary-hover);
  transform: translateY(-1px);
  box-shadow: var(--ds-shadow-md);
}

.save-button:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: var(--ds-shadow-sm);
}

.save-button:disabled {
  background-color: var(--ds-gray-300);
  cursor: not-allowed;
  opacity: 0.6;
  box-shadow: none;
}

.save-icon {
  width: 18px;
  height: 18px;
}
</style>