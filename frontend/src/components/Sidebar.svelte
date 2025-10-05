<script>
  import { 
    currentGuild, 
    currentChannel, 
    sidebarOpen,
    closeSidebar
  } from '../stores/app.js';
  import GuildList from './GuildList.svelte';
  import ChannelList from './ChannelList.svelte';
</script>

<!-- Mobile overlay -->
{#if $sidebarOpen}
  <div 
    class="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
    role="button"
    tabindex="0"
    on:click={closeSidebar}
    on:keydown={(e) => e.key === 'Enter' && closeSidebar()}
  ></div>
{/if}

<!-- Sidebar -->
<aside 
  class="w-64 bg-discord-darker flex flex-col h-full fixed lg:relative z-50 transform transition-transform duration-300 ease-in-out
         {$sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}"
>
  <!-- Guild List -->
  <div class="flex-1 overflow-y-auto">
    <GuildList />
    
    <!-- Channel List -->
    {#if $currentGuild}
      <ChannelList />
    {/if}
  </div>
  
  <!-- Mobile close button -->
  <div class="lg:hidden p-4 border-t border-discord-dark">
    <button 
      class="w-full btn-secondary"
      on:click={closeSidebar}
    >
      Close
    </button>
  </div>
</aside>