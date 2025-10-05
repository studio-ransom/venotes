<script>
  import { 
    guilds, 
    currentGuild, 
    setActiveGuild, 
    createGuild, 
    deleteGuild,
    loading,
    error,
    clearError
  } from '../stores/app.js';

  let newGuildName = '';
  let showAddGuild = false;

  function handleGuildClick(guild) {
    setActiveGuild(guild);
  }

  async function addGuild() {
    if (newGuildName.trim()) {
      try {
        await createGuild(newGuildName.trim());
        newGuildName = '';
        showAddGuild = false;
      } catch (err) {
        // Error is handled in the store
      }
    }
  }

  async function handleDeleteGuild(guildId) {
    if (confirm('Delete this guild and all its channels and logs?')) {
      try {
        await deleteGuild(guildId);
      } catch (err) {
        // Error is handled in the store
      }
    }
  }

  function toggleAddGuild() {
    showAddGuild = !showAddGuild;
    if (!showAddGuild) {
      newGuildName = '';
    }
    clearError();
  }
</script>

<div class="p-4">
  <div class="flex items-center justify-between mb-4">
    <h2 class="text-lg font-semibold text-white">Guilds</h2>
    <button 
      class="text-discord-light hover:text-white text-xl"
      on:click={toggleAddGuild}
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

  <!-- Add Guild Form -->
  {#if showAddGuild}
    <div class="mb-4 p-3 bg-discord-dark rounded-md">
      <input 
        type="text"
        bind:value={newGuildName}
        placeholder="Guild name..."
        class="input-field w-full mb-2"
        on:keydown={(e) => e.key === 'Enter' && addGuild()}
      />
      <div class="flex gap-2">
        <button 
          class="btn-primary text-sm" 
          on:click={addGuild}
          disabled={$loading || !newGuildName.trim()}
        >
          {#if $loading}
            <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          {/if}
          Add
        </button>
        <button class="btn-secondary text-sm" on:click={toggleAddGuild}>
          Cancel
        </button>
      </div>
    </div>
  {/if}

  <!-- Guild List -->
  <div class="space-y-1">
    {#each $guilds as guild}
      <div class="flex items-center group">
        <button
          class="sidebar-item flex-1 text-left {$currentGuild && $currentGuild.id === guild.id ? 'active' : ''}"
          on:click={() => handleGuildClick(guild)}
        >
          <span class="text-lg mr-2">🏢</span>
          {guild.name}
        </button>
        <button
          class="opacity-0 group-hover:opacity-100 text-red-400 hover:text-red-300 text-sm px-2"
          on:click={() => handleDeleteGuild(guild.id)}
        >
          ×
        </button>
      </div>
    {/each}
  </div>

  {#if (!$guilds || $guilds.length === 0) && !$loading}
    <div class="text-center py-4">
      <p class="text-discord-light text-sm">No guilds yet</p>
    </div>
  {/if}
</div>