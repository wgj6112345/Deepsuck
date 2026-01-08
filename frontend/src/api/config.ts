import client from './client'
import type { AgentConfig } from '../store/config'

export interface ProviderConfig {
  id: number
  providerType: string
  apiKey: string
  modelName: string
  baseUrl: string
  enabled: boolean
}

export const configAPI = {
  async getConfig(): Promise<AgentConfig> {
    const response = await client.get('/api/config')
    return response.data
  },

  async updateConfig(config: AgentConfig): Promise<{ success: boolean }> {
    const response = await client.put('/api/config', config)
    return response.data
  },

  async getAllProviders(): Promise<{ success: boolean; providers: ProviderConfig[] }> {
    const response = await client.get('/api/config/providers')
    return response.data
  },

  async saveProviderConfig(config: ProviderConfig): Promise<{ success: boolean }> {
    const response = await client.put('/api/config/providers', config)
    return response.data
  },

  async activateProvider(providerType: string): Promise<{ success: boolean }> {
    const response = await client.post('/api/config/activate', { providerType })
    return response.data
  }
}