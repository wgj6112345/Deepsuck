
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
      const dataLines: string[] = []

      for (const line of lines) {
        if (!line) {
          continue
        }

        if (line.startsWith('event: ')) {
          eventType = line.slice(7).trim()
        } else if (line.startsWith('data: ')) {
          // 收集所有 data 行
          dataLines.push(line.slice(6))
        }
      }

      // 将多个 data 行用换行符连接
      const eventDataJSON = dataLines.join('\n')

      if (eventType && eventDataJSON) {
        // 解析 JSON 编码的数据
        let eventData: string
        try {
          eventData = JSON.parse(eventDataJSON)
        } catch {
          // 如果解析失败，直接使用原始字符串
          eventData = eventDataJSON
        }

        yield { event: eventType as SSEEvent['event'], data: eventData }
      }
    }
  }
}

export async function stopChat(conversationId: string): Promise<void> {
  await fetch(`${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}/api/chat/stop`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ conversationId }),
  })
}