<script>
  import { currentGuild, currentChannel, logs, loading } from '../stores/app.js';
  import LogList from './LogList.svelte';
  import LogInput from './LogInput.svelte';
</script>

<main class="flex-1 flex flex-col bg-discord-darkest min-h-0">
  {#if $loading}
    <!-- Loading Screen -->
    <div class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <div class="inline-flex items-center">
          <svg class="animate-spin -ml-1 mr-3 h-8 w-8 text-discord-blue" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span class="text-discord-light text-lg">Loading...</span>
        </div>
      </div>
    </div>
  {:else if $currentGuild && $currentChannel}
    <!-- Main Content Area with Fixed Layout -->
    <div class="flex-1 flex flex-col min-h-0">
      <!-- Logs Display - Scrollable -->
      <div class="flex-1 min-h-0">
        <LogList />
      </div>
      
      <!-- Log Input - Fixed at bottom -->
      <div class="flex-shrink-0 border-t border-discord-dark p-4 bg-discord-darkest">
        <LogInput />
      </div>
    </div>
  {:else}
    <!-- Welcome Screen -->
    <div class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <div class="text-6xl mb-4">📝</div>
        <h2 class="text-2xl font-semibold text-white mb-2">Welcome to Svelte Notes</h2>
        <p class="text-discord-light mb-4">
          Select a guild and channel from the sidebar to start organizing your thoughts
        </p>
        <div class="text-sm text-discord-light">
          <p>• Create guilds to organize different areas of your life</p>
          <p>• Add channels within guilds for specific topics</p>
          <p>• Start logging your thoughts and ideas</p>
        </div>
      </div>
    </div>
  {/if}
</main>