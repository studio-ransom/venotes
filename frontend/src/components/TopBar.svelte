<script>
  import { currentGuild, currentChannel, sidebarOpen, toggleSidebar } from '../stores/app.js';

  let searchQuery = '';
  let fileInput;

  function handleSearch() {
    // Search functionality can be implemented here
    console.log('Searching for:', searchQuery);
  }

  async function handleExport() {
    try {
      // Ask for password (optional)
      const password = prompt('Enter password for encryption (leave empty for no encryption):');
      
      // Build URL with optional password
      let url = '/api/export';
      if (password && password.trim()) {
        url += `?password=${encodeURIComponent(password.trim())}`;
      }

      // Create download link
      const link = document.createElement('a');
      link.href = url;
          const filename = password && password.trim() 
            ? `venotes-export-${new Date().toISOString().split('T')[0]}.enc`
            : `venotes-export-${new Date().toISOString().split('T')[0]}.zip`;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (error) {
      console.error('Export failed:', error);
      alert('Export failed. Please try again.');
    }
  }

  function handleImport() {
    console.log('Import button clicked, triggering file input');
    // Trigger the hidden file input
    if (fileInput) {
      fileInput.click();
    } else {
      console.error('File input not found');
    }
  }

  async function handleImportFile(event) {
    console.log('File input changed:', event.target.files);
    const file = event.target.files[0];
    if (!file) {
      console.log('No file selected');
      return;
    }

    console.log('Selected file:', file.name, file.size, file.type);

    try {
      // Ask for password if file might be encrypted
      const password = prompt('Enter password if file is encrypted (leave empty if not encrypted):');
      
      // Create form data
      const formData = new FormData();
      formData.append('file', file);
      if (password && password.trim()) {
        formData.append('password', password.trim());
      }

      console.log('Uploading file to /api/import');
      // Upload file
      const response = await fetch('/api/import', {
        method: 'POST',
        body: formData
      });

      console.log('Import response:', response.status, response.statusText);

      if (response.ok) {
        const result = await response.json();
        console.log('Import successful:', result);
        alert('Data imported successfully! Please refresh the page to see the imported data.');
        // Optionally reload the page
        window.location.reload();
      } else {
        const error = await response.json();
        console.error('Import failed:', error);
        alert(`Import failed: ${error.error || 'Unknown error'}`);
      }
    } catch (error) {
      console.error('Import failed:', error);
      alert('Import failed. Please try again.');
    }
  }
</script>

<header class="bg-discord-dark border-b border-discord-darkest px-4 py-3 flex items-center justify-between">
  <!-- Mobile menu button -->
  <button 
    class="lg:hidden text-discord-light hover:text-white mr-3"
    on:click={toggleSidebar}
  >
    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
    </svg>
  </button>

  <!-- Current location -->
  <div class="flex items-center flex-1">
    <h1 class="text-white font-semibold">
      {#if $currentGuild && $currentChannel}
        {#if $currentGuild}
          <span class="text-discord-light">🏢 {$currentGuild.name}</span>
        {/if}
        {#if $currentChannel}
          <span class="text-discord-blue ml-2"># {$currentChannel.name}</span>
        {/if}
      {:else}
        <span class="text-discord-light">Select a channel to start</span>
      {/if}
    </h1>
  </div>

  <!-- Search bar -->
  <div class="hidden md:flex items-center flex-1 max-w-md mx-4">
    <div class="relative w-full">
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Search logs..."
        class="input-field w-full pl-10"
        on:keydown={(e) => e.key === 'Enter' && handleSearch()}
      />
      <svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-discord-light" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
      </svg>
    </div>
  </div>

  <!-- Actions -->
  <div class="flex items-center space-x-2">
    <button 
      class="btn-secondary text-sm"
      on:click={handleExport}
    >
      Export
    </button>
    <button 
      class="btn-secondary text-sm"
      on:click={handleImport}
    >
      Import
    </button>
  </div>
</header>

<!-- Hidden file input for import -->
<input
  bind:this={fileInput}
  type="file"
  accept=".zip,.enc,application/zip,application/x-zip-compressed,application/octet-stream"
  on:change={handleImportFile}
  class="hidden"
/>