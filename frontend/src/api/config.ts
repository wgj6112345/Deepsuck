import client from './client'
import type { AgentConfig } from '../store/config'

export const configAPI = {
  async getConfig(): Promise<AgentConfig> {
    const response = await client.get('/api/config')
    return response.data
  },

  async updateConfig(config: AgentConfig): Promise<{ success: boolean }> {
    const response = await client.put('/api/config', config)
    return response.data
  }
}