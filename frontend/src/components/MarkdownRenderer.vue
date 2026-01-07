<template>
  <div class="markdown-renderer" v-html="renderedContent"></div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

// Props
const props = defineProps<{
  content: string
}>()

// MarkdownIt 实例
const md: MarkdownIt = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true, // 单换行转 <br>
  highlight: (str: string, lang: string): string => {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code>${hljs.highlight(str, { language: lang }).value}</code></pre>`
      } catch {
        // 如果高亮失败，回退到普通渲染
      }
    }
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`
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

.markdown-renderer :deep(code) {
  font-family: 'SF Mono', 'Fira Code', 'Courier New', monospace;
  font-size: 0.9em;
  background-color: var(--ds-primary-light);
  color: var(--ds-primary);
  padding: 0.2em 0.5em;
  border-radius: var(--ds-radius-sm);
  font-weight: var(--ds-font-medium);
}

.markdown-renderer :deep(pre) {
  background-color: #1E1E1E;
  border: 1px solid var(--ds-border-default);
  border-radius: var(--ds-radius-lg);
  padding: var(--ds-spacing-lg);
  overflow-x: auto;
  margin: var(--ds-spacing-lg) 0;
}

.markdown-renderer :deep(pre code) { 
  background-color: transparent; 
  padding: 0; 
  border-radius: 0; 
  font-size: var(--ds-font-sm);
  color: #D4D4D4;
}

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

.markdown-renderer :deep(.hljs) { 
  background-color: transparent; 
  padding: 0; 
}

/* Code block syntax highlighting colors override */
.markdown-renderer :deep(.hljs-keyword) { color: #569CD6; }
.markdown-renderer :deep(.hljs-string) { color: #CE9178; }
.markdown-renderer :deep(.hljs-comment) { color: #6A9955; }
.markdown-renderer :deep(.hljs-function) { color: #DCDCAA; }
.markdown-renderer :deep(.hljs-number) { color: #B5CEA8; }
.markdown-renderer :deep(.hljs-operator) { color: #D4D4D4; }
.markdown-renderer :deep(.hljs-class) { color: #4EC9B0; }
.markdown-renderer :deep(.hljs-variable) { color: #9CDCFE; }
</style>
