<template>
  <div class="settings-view">
    <div class="settings-header">
      <router-link to="/" class="back-link">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M19 12H5M12 19l-7-7 7-7" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span>返回</span>
      </router-link>
      <h1>Agent Providers</h1>
    </div>
    <div class="settings-content">
      <div class="providers-grid">
        <div
          v-for="provider in providers"
          :key="provider.providerType"
          class="provider-card"
          :class="{ active: provider.enabled, editing: editingProvider === provider.providerType }"
          @click="handleCardClick(provider)"
        >
          <div class="provider-header">
            <div class="provider-icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
              </svg>
            </div>
            <div class="provider-info">
              <h3 class="provider-name">{{ getProviderName(provider.providerType) }}</h3>
              <p class="provider-status">
                <span v-if="provider.enabled" class="status-badge active">已激活</span>
                <span v-else class="status-badge">未激活</span>
              </p>
            </div>
          </div>

          <div v-if="editingProvider === provider.providerType" class="provider-form" @click.stop>
            <div class="form-field">
              <label>API Key</label>
              <input
                v-model="provider.apiKey"
                type="password"
                placeholder="输入 API Key"
                class="form-input"
              />
            </div>
            <div class="form-field">
              <label>模型名称</label>
              <input
                v-model="provider.modelName"
                type="text"
                placeholder="模型名称"
                class="form-input"
              />
            </div>
            <div class="form-field">
              <label>Base URL</label>
              <input
                v-model="provider.baseUrl"
                type="text"
                placeholder="Base URL"
                class="form-input"
              />
            </div>
            <div class="form-actions">
              <button class="btn btn-secondary" @click="cancelEdit">取消</button>
              <button class="btn btn-primary" @click="saveProvider(provider)">保存</button>
            </div>
          </div>

          <div v-else class="provider-actions">
            <button class="action-btn" @click.stop="activateProvider(provider.providerType)" :disabled="provider.enabled">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18.36 6.64a9 9 0 1 1-12.73 0"/>
                <line x1="12" y1="2" x2="12" y2="12"/>
              </svg>
              激活
            </button>
            <button class="action-btn" @click.stop="editProvider(provider.providerType)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
              </svg>
              配置
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { configAPI, type ProviderConfig } from '../api/config'

const providers = ref<ProviderConfig[]>([])
const editingProvider = ref<string | null>(null)
const loading = ref(false)

onMounted(async () => {
  await loadProviders()
})

async function loadProviders() {
  try {
    loading.value = true
    const response = await configAPI.getAllProviders()
    providers.value = response.providers
  } catch (error) {
    console.error('Failed to load providers:', error)
  } finally {
    loading.value = false
  }
}

function getProviderName(type: string): string {
  const names: Record<string, string> = {
    'mimo': 'Mimo',
    'iflow': 'IFlow',
    'openai': 'OpenAI',
    'claude': 'Claude'
  }
  return names[type] || type
}

function handleCardClick(provider: ProviderConfig) {
  // 点击卡片不直接激活，需要点击激活按钮
}

function editProvider(providerType: string) {
  editingProvider.value = providerType
}

function cancelEdit() {
  editingProvider.value = null
  // 重新加载以恢复原始值
  loadProviders()
}

async function saveProvider(provider: ProviderConfig) {
  try {
    loading.value = true
    await configAPI.saveProviderConfig(provider)
    editingProvider.value = null
    alert('配置已保存')
    await loadProviders()
  } catch (error) {
    alert('保存失败，请重试')
  } finally {
    loading.value = false
  }
}

async function activateProvider(providerType: string) {
  try {
    loading.value = true
    await configAPI.activateProvider(providerType)
    alert('已激活 ' + getProviderName(providerType))
    await loadProviders()
  } catch (error) {
    alert('激活失败，请重试')
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
  background-color: #FFFFFF;
  color: #1F2937;
}

.settings-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  background-color: #FFFFFF;
  border-bottom: 1px solid #E5E7EB;
  height: 60px;
  flex-shrink: 0;
}

.back-link {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #4A6CF7;
  text-decoration: none;
  font-size: 16px;
  font-weight: 500;
  padding: 6px 12px;
  border-radius: 6px;
  transition: all 0.2s;
}

.back-link:hover {
  background-color: #EFF6FF;
  text-decoration: none;
}

.back-link svg {
  width: 18px;
  height: 18px;
}

.settings-header h1 {
  flex: 1;
  font-size: 24px;
  font-weight: 600;
  margin: 0;
}

.settings-content {
  flex: 1;
  overflow-y: auto;
  padding: 32px;
}

.providers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.provider-card {
  background-color: #FFFFFF;
  border: 2px solid #E5E7EB;
  border-radius: 16px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s;
}

.provider-card:hover {
  border-color: #4A6CF7;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(74, 108, 247, 0.1);
}

.provider-card.active {
  border-color: #4A6CF7;
  background-color: #EFF6FF;
}

.provider-card.editing {
  border-color: #4A6CF7;
  cursor: default;
}

.provider-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.provider-icon {
  width: 48px;
  height: 48px;
  background-color: #4A6CF7;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.provider-icon svg {
  width: 24px;
  height: 24px;
  color: #FFFFFF;
}

.provider-info {
  flex: 1;
}

.provider-name {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 4px 0;
  color: #1F2937;
}

.provider-status {
  margin: 0;
}

.status-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
}

.status-badge.active {
  background-color: #4A6CF7;
  color: #FFFFFF;
}

.status-badge:not(.active) {
  background-color: #E5E7EB;
  color: #6B7280;
}

.provider-form {
  margin-top: 16px;
}

.form-field {
  margin-bottom: 16px;
}

.form-field label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

.form-input {
  width: 100%;
  padding: 10px 14px;
  background-color: #FFFFFF;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  color: #1F2937;
  font-size: 16px;
  outline: none;
  transition: all 0.2s;
  font-family: inherit;
}

.form-input:focus {
  border-color: #4A6CF7;
  box-shadow: 0 0 0 3px rgba(74, 108, 247, 0.1);
}

.form-input::placeholder {
  color: #9CA3AF;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 20px;
}

.btn {
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
  display: flex;
  align-items: center;
  gap: 6px;
}

.btn-primary {
  background-color: #4A6CF7;
  color: #FFFFFF;
  border: none;
}

.btn-primary:hover {
  background-color: #3B5FCC;
}

.btn-secondary {
  background-color: #F3F4F6;
  color: #1F2937;
  border: 1px solid #D1D5DB;
}

.btn-secondary:hover {
  background-color: #E5E7EB;
}

.provider-actions {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  background-color: #F3F4F6;
  color: #1F2937;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
}

.action-btn:hover:not(:disabled) {
  background-color: #E5E7EB;
  border-color: #9CA3AF;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn svg {
  width: 16px;
  height: 16px;
}
</style>