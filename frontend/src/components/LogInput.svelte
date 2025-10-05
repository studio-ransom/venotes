<script>
  import { 
    createLog, 
    currentChannel,
    loading,
    error,
    clearError
  } from '../stores/app.js';
  import { tick } from 'svelte';
  import FileUpload from './FileUpload.svelte';

  let newLogContent = '';
  let isSubmitting = false;
  let selectedFiles = [];
  let showFileUpload = false;
  let textarea;

  async function submitLog() {
    if ((newLogContent.trim() || selectedFiles.length > 0) && !isSubmitting && $currentChannel) {
      isSubmitting = true;
      
      try {
        await createLog($currentChannel.id, newLogContent.trim(), selectedFiles);
        newLogContent = '';
        selectedFiles = [];
        showFileUpload = false;
        
        // Wait for DOM to update, then refocus the textarea
        await tick();
        setTimeout(() => {
          if (textarea) {
            textarea.focus();
          }
        }, 10);
      } catch (err) {
        // Error is handled in the store
      } finally {
        isSubmitting = false;
      }
    }
  }

  function handleKeyDown(event) {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      submitLog();
    }
  }

  function handleFilesSelected(event) {
    selectedFiles = [...selectedFiles, ...event.detail.files];
  }

  function removeFile(index) {
    selectedFiles = selectedFiles.filter((_, i) => i !== index);
  }

  function toggleFileUpload() {
    showFileUpload = !showFileUpload;
  }
</script>

<div class="space-y-3">
  <!-- File Upload Area -->
  {#if showFileUpload}
    <FileUpload 
      on:filesSelected={handleFilesSelected}
      disabled={isSubmitting || $loading}
      multiple={true}
    />
  {/if}

  <!-- Selected Files -->
  {#if selectedFiles.length > 0}
    <div class="selected-files">
      <div class="text-sm text-discord-light mb-2">Selected files:</div>
      {#each selectedFiles as file, index}
        <div class="file-item">
          <span class="file-name">{file.name}</span>
          <span class="file-size">({Math.round(file.size / 1024)}KB)</span>
          <button 
            class="remove-file"
            on:click={() => removeFile(index)}
          >
            ×
          </button>
        </div>
      {/each}
    </div>
  {/if}

  <!-- Input Area -->
  <div class="flex items-end space-x-3">
    <div class="flex-1">
      <textarea
        bind:this={textarea}
        bind:value={newLogContent}
        placeholder="What's on your mind? (Supports Markdown)"
        class="input-field w-full resize-none min-h-[60px] max-h-[200px]"
        on:keydown={handleKeyDown}
        disabled={isSubmitting || $loading}
      ></textarea>
      <div class="flex items-center justify-between mt-1">
        <div class="text-xs text-discord-light">
          Press Enter to send, Shift+Enter for new line
        </div>
        <button
          class="text-xs text-discord-blue hover:text-discord-blue-light"
          on:click={toggleFileUpload}
        >
          {showFileUpload ? 'Hide' : 'Add'} files
        </button>
      </div>
    </div>
    
    <button
      class="btn-primary px-6 py-3 disabled:opacity-50 disabled:cursor-not-allowed"
      on:click={submitLog}
      disabled={(!newLogContent.trim() && selectedFiles.length === 0) || isSubmitting || $loading}
    >
      {#if isSubmitting}
        <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      {/if}
      Send
    </button>
  </div>
</div>

<style>
  .selected-files {
    @apply bg-discord-darkest border border-discord-dark rounded p-3;
  }

  .file-item {
    @apply flex items-center justify-between py-1;
  }

  .file-name {
    @apply text-discord-light text-sm;
  }

  .file-size {
    @apply text-discord-light text-xs;
  }

  .remove-file {
    @apply text-red-400 hover:text-red-300 text-sm px-1;
  }
</style>