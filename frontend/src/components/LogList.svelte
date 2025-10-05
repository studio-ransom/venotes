<script>
  import { 
    logs, 
    deleteLog, 
    currentChannel,
    loading,
    error,
    clearError
  } from '../stores/app.js';
  import MarkdownRenderer from './MarkdownRenderer.svelte';
  import FileDisplay from './FileDisplay.svelte';
  import { onMount, afterUpdate } from 'svelte';

  let logContainer;
  let shouldScrollToBottom = true;
  let showScrollButton = false;

  // Auto-scroll to bottom when new messages are added
  afterUpdate(() => {
    if (logContainer && shouldScrollToBottom) {
      logContainer.scrollTop = logContainer.scrollHeight;
    }
  });

  // Handle scroll events to detect if user scrolled up
  function handleScroll() {
    const { scrollTop, scrollHeight, clientHeight } = logContainer;
    const isAtBottom = scrollTop + clientHeight >= scrollHeight - 10; // 10px threshold
    shouldScrollToBottom = isAtBottom;
    showScrollButton = !isAtBottom && scrollHeight > clientHeight;
  }

  // Scroll to bottom function
  function scrollToBottom() {
    if (logContainer) {
      logContainer.scrollTop = logContainer.scrollHeight;
      shouldScrollToBottom = true;
      showScrollButton = false;
    }
  }

  async function handleDeleteLog(logId) {
    if (confirm('Delete this log entry?')) {
      try {
        await deleteLog(logId);
      } catch (err) {
        // Error is handled in the store
      }
    }
  }

  function formatTimestamp(timestamp) {
    const date = new Date(timestamp);
    const now = new Date();
    const diffInHours = (now - date) / (1000 * 60 * 60);
    
    if (diffInHours < 24) {
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    } else if (diffInHours < 24 * 7) {
      return date.toLocaleDateString([], { weekday: 'short', hour: '2-digit', minute: '2-digit' });
    } else {
      return date.toLocaleDateString([], { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
    }
  }
</script>

<div class="relative p-4 h-full overflow-y-auto" bind:this={logContainer} on:scroll={handleScroll}>
  <!-- Error Display -->
  {#if $error}
    <div class="mb-4 p-3 bg-red-900 border border-red-700 rounded-md">
      <div class="flex items-center justify-between">
        <span class="text-red-200 text-sm">{$error}</span>
        <button 
          class="text-red-400 hover:text-red-300 text-sm"
          on:click={clearError}
        >
          ×
        </button>
      </div>
    </div>
  {/if}

  {#if !$logs || $logs.length === 0}
    <div class="text-center py-12">
      <div class="text-4xl mb-4">💭</div>
      <h3 class="text-lg font-medium text-white mb-2">No logs yet</h3>
      <p class="text-discord-light">Start by adding your first log entry below</p>
    </div>
  {:else}
    <div class="space-y-1">
      {#each $logs as log (log.id)}
        <div class="log-entry group">
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center mb-1">
                <span class="text-xs text-discord-light mr-2">
                  {formatTimestamp(log.created_at)}
                </span>
              </div>
              
              <!-- Log Content -->
              {#if log.content}
                <div class="mb-2">
                  <MarkdownRenderer content={log.content} />
                </div>
              {/if}
              
              <!-- Files -->
              {#if log.files && log.files.length > 0}
                <div class="files-section">
                  {#each log.files as file}
                    <FileDisplay {file} />
                  {/each}
                </div>
              {/if}
            </div>
            <button
              class="opacity-0 group-hover:opacity-100 text-red-400 hover:text-red-300 text-sm px-2 py-1"
              on:click={() => handleDeleteLog(log.id)}
            >
              ×
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  {#if $loading}
    <div class="text-center py-4">
      <div class="inline-flex items-center">
        <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-discord-light" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span class="text-discord-light text-sm">Loading...</span>
      </div>
    </div>
  {/if}

  <!-- Scroll to Bottom Button -->
  {#if showScrollButton}
    <button
      class="fixed bottom-20 right-6 bg-discord-blue hover:bg-blue-600 text-white p-3 rounded-full shadow-lg transition-all duration-200 z-10"
      on:click={scrollToBottom}
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"></path>
      </svg>
    </button>
  {/if}
</div>