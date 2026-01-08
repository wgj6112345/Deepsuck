<template>
  <div class="markdown-renderer" v-html="renderedContent" ref="containerRef"></div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import MarkdownIt from 'markdown-it'
import Prism from 'prismjs'
import 'prismjs/themes/prism.css'
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-java'
import 'prismjs/components/prism-c'
import 'prismjs/components/prism-cpp'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-rust'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-css'
import 'prismjs/components/prism-scss'
import 'prismjs/components/prism-markdown'
import 'prismjs/components/prism-yaml'
import 'prismjs/components/prism-xml-doc'

// Props
const props = defineProps<{
  content: string
}>()

const containerRef = ref<HTMLElement | null>(null)

// MarkdownIt 实例
const md: MarkdownIt = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: false,
  breaks: true,
  highlight: (str: string, lang: string): string => {
    if (lang && Prism.languages[lang]) {
      try {
        const highlighted = Prism.highlight(str, Prism.languages[lang], lang)
        const languageName = lang.toLowerCase()
        const codeId = `code-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
        return `<div class="code-block-wrapper"><div class="code-block-header"><span class="code-language">${languageName}</span><button class="copy-button" onclick="copyCode('${codeId}')" title="复制代码"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg><span class="copy-text">复制</span></button></div><code id="${codeId}" class="code-content language-${lang}">${highlighted}</code></div>`
      } catch {
        // 如果高亮失败，回退到普通渲染
      }
    }
    return `<div class="code-block-wrapper"><code>${md.utils.escapeHtml(str)}</code></div>`
  }
})

// 渲染后的 HTML
const renderedContent = ref<string>('')

// 监听 content prop
watch(
  () => props.content,
  (newContent) => {
    const rendered = md.render(newContent)
    renderedContent.value = rendered
  },
  { immediate: true }
)

// 全局复制函数
onMounted(() => {
  // 将复制函数挂载到 window 对象
  ;(window as any).copyCode = (codeId: string) => {
    const codeElement = document.getElementById(codeId)
    if (codeElement) {
      const code = codeElement.textContent || ''
      navigator.clipboard.writeText(code).then(() => {
        // 更新按钮文本
        const button = document.querySelector(`button[onclick="copyCode('${codeId}')"]`)
        if (button) {
          const textSpan = button.querySelector('.copy-text')
          if (textSpan) {
            textSpan.textContent = '已复制'
            setTimeout(() => {
              textSpan.textContent = '复制'
            }, 2000)
          }
        }
      }).catch((err) => {
        console.error('复制失败:', err)
      })
    }
  }
})
</script>

<style scoped>
.markdown-renderer :deep(h1),
.markdown-renderer :deep(h2),
.markdown-renderer :deep(h3),
.markdown-renderer :deep(h4),
.markdown-renderer :deep(h5),
.markdown-renderer :deep(h6) {
  margin-top: var(--ds-spacing-2xl);
  margin-bottom: var(--ds-spacing-md);
  font-weight: var(--ds-font-semibold);
  color: var(--ds-text-primary);
  line-height: var(--ds-leading-tight);
}

.markdown-renderer :deep(h1) { 
  font-size: var(--ds-font-4xl); 
  border-bottom: 2px solid var(--ds-border-default); 
  padding-bottom: var(--ds-spacing-md);
}

.markdown-renderer :deep(h2) { 
  font-size: var(--ds-font-3xl); 
  border-bottom: 1px solid var(--ds-border-default); 
  padding-bottom: var(--ds-spacing-sm);
}

.markdown-renderer :deep(h3) { font-size: var(--ds-font-2xl); }
.markdown-renderer :deep(h4),
.markdown-renderer :deep(h5),
.markdown-renderer :deep(h6) { font-size: var(--ds-font-xl); }

.markdown-renderer :deep(p) {
  margin: var(--ds-spacing-md) 0;
  line-height: var(--ds-leading-relaxed);
  color: var(--ds-text-primary);
}

.markdown-renderer :deep(br) {
  display: block;
  margin: var(--ds-spacing-sm) 0;
  content: '';
}

.markdown-renderer :deep(ul),
.markdown-renderer :deep(ol) { 
  margin: var(--ds-spacing-md) 0; 
  padding-left: var(--ds-spacing-2xl); 
}

.markdown-renderer :deep(li) { 
  margin: var(--ds-spacing-xs) 0; 
  line-height: var(--ds-leading-normal);
  color: var(--ds-text-primary);
}

.markdown-renderer :deep(ul) { list-style-type: disc; }
.markdown-renderer :deep(ol) { list-style-type: decimal; }

/* 行内代码样式 - 排除代码块中的 code */
.markdown-renderer :deep(p code),
.markdown-renderer :deep(li code),
.markdown-renderer :deep(td code),
.markdown-renderer :deep(th code) {
  font-family: 'JetBrains Mono', 'SF Mono', 'Fira Code', 'Courier New', monospace;
  font-size: 0.9em;
  background-color: #F9FAFB;
  color: var(--ds-primary);
  padding: 0.2em 0.5em;
  border-radius: var(--ds-radius-sm);
  font-weight: var(--ds-font-medium);
}

/* 代码块容器 - DeepSeek 样式 */

.markdown-renderer :deep(.code-block-wrapper) {

  margin: 8px 0;

  background-color: #F9FAFB;

  border: none;

  border-radius: 8px;

  overflow: hidden;

  display: flex;

  flex-direction: column;

}



/* 代码块头部 - DeepSeek 样式 */

.markdown-renderer :deep(.code-block-header) {

  display: flex;

  justify-content: space-between;

  align-items: center;

  padding: 8px 12px;

  background-color: transparent;

  border-bottom: none;

  margin: 0;

  flex-shrink: 0;

}

/* 代码语言标签 */

.markdown-renderer :deep(.code-language) {

  font-size: 0.75em;

  font-weight: 500;

  color: var(--ds-text-secondary);

  text-transform: lowercase;

}



/* 代码内容 - DeepSeek 样式 */

.markdown-renderer :deep(.code-content) {

  font-family: 'JetBrains Mono', 'SF Mono', 'Fira Code', 'Courier New', monospace;

  font-size: var(--ds-font-base);

  font-weight: 600;

  white-space: pre;

  padding: 12px;

  margin: 0;

  overflow-x: auto;

  background-color: transparent;

  line-height: 1.6;

  display: block;

}



/* 其他 Markdown 元素样式 */

.markdown-renderer :deep(a) { 

  color: var(--ds-primary); 

  text-decoration: none;

  font-weight: var(--ds-font-medium);

  border-bottom: 1px solid transparent;

  transition: all var(--ds-transition-fast);

}



.markdown-renderer :deep(a:hover) { 

  color: var(--ds-primary-hover);

  border-bottom-color: var(--ds-primary-hover);

}



.markdown-renderer :deep(blockquote) { 

  border-left: 4px solid var(--ds-primary); 

  padding: var(--ds-spacing-md) var(--ds-spacing-lg); 

  margin: var(--ds-spacing-lg) 0; 

  background-color: var(--ds-bg-secondary); 

  color: var(--ds-text-secondary);

  border-radius: 0 var(--ds-radius-md) var(--ds-radius-md) 0;

}



.markdown-renderer :deep(blockquote p) { margin: 0; }



.markdown-renderer :deep(table) { 

  border-collapse: collapse; 

  width: 100%; 

  margin: var(--ds-spacing-lg) 0; 

  overflow-x: auto;

  border-radius: var(--ds-radius-md);

  overflow: hidden;

}



.markdown-renderer :deep(th),

.markdown-renderer :deep(td) { 

  border: 1px solid var(--ds-border-default); 

  padding: var(--ds-spacing-md) var(--ds-spacing-lg); 

  text-align: left; 

}



.markdown-renderer :deep(th) { 

  background-color: var(--ds-bg-tertiary); 

  font-weight: var(--ds-font-semibold); 

  color: var(--ds-text-primary);

}



.markdown-renderer :deep(tr:nth-child(even)) { 

  background-color: var(--ds-bg-secondary); 

}



.markdown-renderer :deep(strong) { 

  font-weight: var(--ds-font-semibold); 

  color: var(--ds-text-primary); 

}



.markdown-renderer :deep(em) { 

  font-style: italic; 

  color: var(--ds-text-secondary); 

}



.markdown-renderer :deep(hr) { 

  border: none; 

  border-top: 1px solid var(--ds-border-default); 

  margin: var(--ds-spacing-3xl) 0; 

}



.markdown-renderer :deep(img) { 

  max-width: 100%; 

  height: auto; 

  border-radius: var(--ds-radius-lg); 

  margin: var(--ds-spacing-md) 0;

  box-shadow: var(--ds-shadow-md);

}

/* 复制按钮 - DeepSeek 样式 */
.markdown-renderer :deep(.copy-button) {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: transparent;
  border: none;
  color: #666666;
  padding: 0;
  border-radius: 0;
  font-size: 14px;
  font-family: inherit;
  font-weight: 400;
  cursor: pointer;
  transition: color var(--ds-transition-fast);
}

.markdown-renderer :deep(.copy-button:hover) {
  color: #333333;
}

.markdown-renderer :deep(.copy-button:active) {
  color: #000000;
}

.markdown-renderer :deep(.copy-button svg) {
  width: 16px;
  height: 16px;
  stroke: currentColor;
  stroke-width: 1.5;
}

.markdown-renderer :deep(.copy-button .copy-text) {
  font-weight: 400;
}
</style>