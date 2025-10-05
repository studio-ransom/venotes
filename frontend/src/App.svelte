<script>
  import { onMount } from 'svelte';
  import { loadInitialData } from './stores/app.js';
  import Sidebar from './components/Sidebar.svelte';
  import MainPanel from './components/MainPanel.svelte';
  import TopBar from './components/TopBar.svelte';

  onMount(() => {
    loadInitialData();
  });
</script>

<main class="flex h-screen bg-discord-darkest">
  <!-- Sidebar -->
  <Sidebar />
  
  <!-- Main Content -->
  <div class="flex-1 flex flex-col">
    <!-- Top Bar -->
    <TopBar />
    
    <!-- Main Panel -->
    <MainPanel />
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
</style>