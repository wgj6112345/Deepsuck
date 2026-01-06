import type { Message } from '../store/conversation'

export interface ChatRequest {
  conversationId: string
  content: string
  thinkingEnabled: boolean
}

export interface SSEEvent {
  event: 'thinking' | 'content' | 'done' | 'error' | 'title_update'
  data: string
}

export async function* streamChat(request: ChatRequest, signal?: AbortSignal): AsyncGenerator<SSEEvent> {
  const response = await fetch(`${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}/api/chat`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
    signal,
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  const reader = response.body?.getReader()
  if (!reader) {
    throw new Error('No response body')
  }

  const decoder = new TextDecoder()
  let buffer = ''
  let currentEvent = ''

  while (true) {
    if (signal?.aborted) {
      reader.cancel()
      throw new DOMException('Aborted', 'AbortError')
    }

    const { done, value } = await reader.read()
    if (done) break

    buffer += decoder.decode(value, { stream: true })

    // 按双换行符分割 SSE 事件
    const events = buffer.split('\n\n')
    buffer = events.pop() || ''

    for (const event of events) {
      if (!event || event.trim() === '') {
        continue
      }

      const lines = event.split('\n')
      let eventType = ''
      let eventData = ''

      for (const line of lines) {
        if (!line) {
          continue
        }

        if (line.startsWith('event: ')) {
          eventType = line.slice(7).trim()
        } else if (line.startsWith('data: ')) {
          eventData = line.slice(6)
          console.log(`SSE Data: "${eventData.replace(/\n/g, '\\n')}"`)
          console.log(`Chars:`, Array.from(eventData).map(c => {
            if (c === '\n') return '\\n';
            return c;
          }).join(''))
        }
      }

      if (eventType && eventData) {
        yield { event: eventType as SSEEvent['event'], data: eventData }
      }
    }
  }
}