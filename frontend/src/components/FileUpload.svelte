<script>
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  export let disabled = false;
  export let multiple = false;

  let fileInput;
  let dragOver = false;

  function handleFileSelect(event) {
    const files = Array.from(event.target.files);
    if (files.length > 0) {
      dispatch('filesSelected', { files });
    }
  }

  function handleDrop(event) {
    event.preventDefault();
    dragOver = false;
    
    const files = Array.from(event.dataTransfer.files);
    if (files.length > 0) {
      dispatch('filesSelected', { files });
    }
  }

  function handleDragOver(event) {
    event.preventDefault();
    dragOver = true;
  }

  function handleDragLeave(event) {
    event.preventDefault();
    dragOver = false;
  }

  function triggerFileInput() {
    if (!disabled) {
      fileInput.click();
    }
  }
</script>

<div 
  class="file-upload-area"
  class:drag-over={dragOver}
  class:disabled={disabled}
  on:drop={handleDrop}
  on:dragover={handleDragOver}
  on:dragleave={handleDragLeave}
  on:click={triggerFileInput}
  role="button"
  tabindex="0"
  on:keydown={(e) => e.key === 'Enter' && triggerFileInput()}
>
  <input
    bind:this={fileInput}
    type="file"
    multiple={multiple}
    on:change={handleFileSelect}
    class="hidden"
    {disabled}
    accept="*/*"
  />
  
  <div class="upload-content">
    <svg class="upload-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
    </svg>
    <p class="upload-text">
      {#if dragOver}
        Drop files here
      {:else}
        Click to upload or drag and drop
      {/if}
    </p>
    <p class="upload-hint">Images, text files, and more</p>
  </div>
</div>

<style>
  .file-upload-area {
    @apply border-2 border-dashed border-discord-dark rounded-lg p-6 text-center cursor-pointer transition-colors;
  }

  .file-upload-area:hover:not(.disabled) {
    @apply border-discord-blue bg-discord-darkest;
  }

  .file-upload-area.drag-over {
    @apply border-discord-blue bg-discord-darkest;
  }

  .file-upload-area.disabled {
    @apply opacity-50 cursor-not-allowed;
  }

  .upload-content {
    @apply flex flex-col items-center space-y-2;
  }

  .upload-icon {
    @apply w-8 h-8 text-discord-light;
  }

  .upload-text {
    @apply text-discord-light font-medium;
  }

  .upload-hint {
    @apply text-sm text-discord-light;
  }
</style>
