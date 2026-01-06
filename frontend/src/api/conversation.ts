import client from './client'
import type { Conversation } from '../store/conversation'

export const conversationAPI = {
  async createConversation(title?: string): Promise<Conversation> {
    const response = await client.post('/api/conversations', { title })
    return response.data
  },

  async getConversations(): Promise<Conversation[]> {
    const response = await client.get('/api/conversations')
    return response.data.conversations
  },

  async getConversation(id: string): Promise<Conversation> {
    const response = await client.get(`/api/conversations/${id}`)
    return response.data
  },

  async deleteConversation(id: string): Promise<{ success: boolean }> {
    const response = await client.delete(`/api/conversations/${id}`)
    return response.data
  }
}