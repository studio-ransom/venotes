<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  export let file;
  export let showContent = false;

  let contentVisible = showContent;
  let fileContent = '';
  let loading = false;

  // File type detection
  $: isImage = file.mime_type?.startsWith('image/');
  $: isText = file.mime_type?.startsWith('text/') || 
             (file.original_name && (
               file.original_name.endsWith('.txt') || 
               file.original_name.endsWith('.md') || 
               file.original_name.endsWith('.json') || 
               file.original_name.endsWith('.js') || 
               file.original_name.endsWith('.css') || 
               file.original_name.endsWith('.html')
             ));

  async function toggleContent() {
    if (!contentVisible && isText && !fileContent) {
      loading = true;
      try {
        const response = await fetch(`/api/files/${file.id}/content`);
        if (response.ok) {
          fileContent = await response.text();
        }
      } catch (err) {
        console.error('Failed to load file content:', err);
      } finally {
        loading = false;
      }
    }
    contentVisible = !contentVisible;
  }

  function downloadFile() {
    const link = document.createElement('a');
    link.href = `/api/files/${file.id}`;
    link.download = file.original_name || 'download';
    link.click();
  }

  function deleteFile() {
    dispatch('delete', { file });
  }

  function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }
</script>

<div class="file-display">
  <div class="file-header">
    <div class="file-info">
      <div class="file-icon">
        {#if isImage}
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z" clip-rule="evenodd"></path>
          </svg>
        {:else if isText}
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clip-rule="evenodd"></path>
          </svg>
        {:else}
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd"></path>
          </svg>
        {/if}
      </div>
      <div class="file-details">
        <div class="file-name">{file.original_name || 'Unknown file'}</div>
        <div class="file-meta">
          {formatFileSize(file.size || 0)} • {file.mime_type || 'Unknown type'}
        </div>
      </div>
    </div>
    
    <div class="file-actions">
      {#if isText}
        <button
          class="action-btn"
          on:click={toggleContent}
          disabled={loading}
        >
          {#if loading}
            <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          {:else}
            {contentVisible ? 'Hide' : 'Show'} content
          {/if}
        </button>
      {/if}
      
      <button class="action-btn" on:click={downloadFile}>
        Download
      </button>
      
      <button class="action-btn text-red-400 hover:text-red-300" on:click={deleteFile}>
        Delete
      </button>
    </div>
  </div>

  {#if isImage}
    <div class="file-preview">
      <img 
        src={`/api/files/${file.id}`} 
        alt={file.original_name || 'Uploaded file'}
        class="max-w-full h-auto rounded"
        loading="lazy"
      />
    </div>
  {/if}

  {#if contentVisible && isText && fileContent}
    <div class="file-content">
      <pre class="text-sm bg-discord-darkest p-3 rounded border border-discord-dark overflow-x-auto"><code>{fileContent}</code></pre>
    </div>
  {/if}
</div>

<style>
  .file-display {
    @apply border border-discord-dark rounded-lg p-3 mb-2 bg-discord-darkest;
  }

  .file-header {
    @apply flex items-center justify-between mb-2;
  }

  .file-info {
    @apply flex items-center space-x-3 flex-1;
  }

  .file-icon {
    @apply text-discord-blue flex-shrink-0;
  }

  .file-details {
    @apply flex-1 min-w-0;
  }

  .file-name {
    @apply text-white font-medium truncate;
  }

  .file-meta {
    @apply text-sm text-discord-light;
  }

  .file-actions {
    @apply flex items-center space-x-2 flex-shrink-0;
  }

  .action-btn {
    @apply px-2 py-1 text-xs text-discord-light hover:text-white transition-colors;
  }

  .file-preview {
    @apply mt-2;
  }

  .file-content {
    @apply mt-2;
  }
</style>
