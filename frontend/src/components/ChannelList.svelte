<script>
  import { 
    channels, 
    currentChannel, 
    setActiveChannel, 
    createChannel, 
    deleteChannel,
    currentGuild,
    loading,
    error,
    clearError
  } from '../stores/app.js';

  let newChannelName = '';
  let showAddChannel = false;

  function handleChannelClick(channel) {
    setActiveChannel(channel);
  }

  async function addChannel() {
    if (newChannelName.trim() && $currentGuild) {
      try {
        await createChannel($currentGuild.id, newChannelName.trim());
        newChannelName = '';
        showAddChannel = false;
      } catch (err) {
        // Error is handled in the store
      }
    }
  }

  async function handleDeleteChannel(channelId) {
    if (confirm('Delete this channel and all its logs?')) {
      try {
        await deleteChannel(channelId);
      } catch (err) {
        // Error is handled in the store
      }
    }
  }

  function toggleAddChannel() {
    showAddChannel = !showAddChannel;
    if (!showAddChannel) {
      newChannelName = '';
    }
    clearError();
  }
</script>

<div class="px-4 pb-4">
  <div class="flex items-center justify-between mb-4">
    <h3 class="text-md font-medium text-discord-light">Channels</h3>
    <button 
      class="text-discord-light hover:text-white text-lg"
      on:click={toggleAddChannel}
    >
      +
    </button>
  </div>

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

  <!-- Add Channel Form -->
  {#if showAddChannel}
    <div class="mb-4 p-3 bg-discord-dark rounded-md">
      <input 
        type="text"
        bind:value={newChannelName}
        placeholder="Channel name..."
        class="input-field w-full mb-2"
        on:keydown={(e) => e.key === 'Enter' && addChannel()}
      />
      <div class="flex gap-2">
        <button 
          class="btn-primary text-sm" 
          on:click={addChannel}
          disabled={$loading || !newChannelName.trim()}
        >
          {#if $loading}
            <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          {/if}
          Add
        </button>
        <button class="btn-secondary text-sm" on:click={toggleAddChannel}>
          Cancel
        </button>
      </div>
    </div>
  {/if}

  <!-- Channel List -->
  <div class="space-y-1">
    {#each $channels as channel}
      <div class="flex items-center group">
        <button
          class="sidebar-item flex-1 text-left {$currentChannel && $currentChannel.id === channel.id ? 'active' : ''}"
          on:click={() => handleChannelClick(channel)}
        >
          <span class="text-sm mr-2">#</span>
          {channel.name}
        </button>
        <button
          class="opacity-0 group-hover:opacity-100 text-red-400 hover:text-red-300 text-sm px-2"
          on:click={() => handleDeleteChannel(channel.id)}
        >
          ×
        </button>
      </div>
    {/each}
  </div>

  {#if (!$channels || $channels.length === 0) && !$loading}
    <div class="text-center py-4">
      <p class="text-discord-light text-sm">No channels yet</p>
    </div>
  {/if}
</div>