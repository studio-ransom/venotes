<script>
  import { marked } from 'marked';
  import hljs from 'highlight.js';
  import 'highlight.js/styles/github.css';

  export let content = '';

  // Configure marked with syntax highlighting
  marked.setOptions({
    highlight: function(code, lang) {
      if (lang && hljs.getLanguage(lang)) {
        try {
          return hljs.highlight(code, { language: lang }).value;
        } catch (err) {
          console.error('Syntax highlighting error:', err);
        }
      }
      return hljs.highlightAuto(code).value;
    },
    breaks: true,
    gfm: true
  });

  // Render markdown to HTML
  $: htmlContent = marked(content);
</script>

<div class="markdown-content">
  {@html htmlContent}
</div>

<style>
  .markdown-content {
    @apply text-discord-lighter;
  }

  .markdown-content :global(h1) {
    @apply text-2xl font-bold text-white mb-4 mt-6;
  }

  .markdown-content :global(h2) {
    @apply text-xl font-semibold text-white mb-3 mt-5;
  }

  .markdown-content :global(h3) {
    @apply text-lg font-medium text-white mb-2 mt-4;
  }

  .markdown-content :global(h4) {
    @apply text-base font-medium text-white mb-2 mt-3;
  }

  .markdown-content :global(p) {
    @apply mb-3 leading-relaxed;
  }

  .markdown-content :global(ul) {
    @apply list-disc list-inside mb-3 ml-4;
  }

  .markdown-content :global(ol) {
    @apply list-decimal list-inside mb-3 ml-4;
  }

  .markdown-content :global(li) {
    @apply mb-1;
  }

  .markdown-content :global(blockquote) {
    @apply border-l-4 border-discord-blue pl-4 italic text-discord-light mb-3;
  }

  .markdown-content :global(code) {
    @apply bg-discord-dark text-discord-blue px-1 py-0.5 rounded text-sm font-mono;
  }

  .markdown-content :global(pre) {
    @apply bg-discord-darkest border border-discord-dark rounded-lg p-4 mb-3 overflow-x-auto;
  }

  .markdown-content :global(pre code) {
    @apply bg-transparent p-0 text-sm;
  }

  .markdown-content :global(a) {
    @apply text-discord-blue hover:text-blue-300 underline;
  }

  .markdown-content :global(strong) {
    @apply font-semibold text-white;
  }

  .markdown-content :global(em) {
    @apply italic;
  }

  .markdown-content :global(hr) {
    @apply border-discord-dark my-4;
  }

  .markdown-content :global(table) {
    @apply w-full border-collapse border border-discord-dark mb-3;
  }

  .markdown-content :global(th) {
    @apply border border-discord-dark bg-discord-dark px-3 py-2 text-left font-semibold;
  }

  .markdown-content :global(td) {
    @apply border border-discord-dark px-3 py-2;
  }
</style>
