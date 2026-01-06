import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface AgentConfig {
  apiKey: string
  modelName: string
  baseUrl: string
}

export const useConfigStore = defineStore('config', () => {
  const config = ref<AgentConfig>({
    apiKey: '',
    modelName: import.meta.env.VITE_DEFAULT_AGENT_MODEL || 'mimo-v2-flash',
    baseUrl: import.meta.env.VITE_DEFAULT_AGENT_BASE_URL || 'https://api.xiaomimimo.com/v1'
  })
  const loading = ref(false)

  function setConfig(data: AgentConfig) {
    config.value = data
  }

  function setLoading(value: boolean) {
    loading.value = value
  }

  return {
    config,
    loading,
    setConfig,
    setLoading
  }
})