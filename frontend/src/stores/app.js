import { writable } from 'svelte/store';
import { guildApi, channelApi, logApi } from '../api/api.js';

// App state stores
export const currentGuild = writable(null);
export const currentChannel = writable(null);
export const sidebarOpen = writable(false);
export const searchQuery = writable('');
export const loading = writable(false);
export const error = writable(null);

// Data stores
export const guilds = writable([]);
export const channels = writable([]);
export const logs = writable([]);

// Load initial data
export async function loadInitialData() {
  loading.set(true);
  error.set(null);
  
  try {
    const guildList = await guildApi.getAll();
    // Ensure guildList is an array
    const guildsArray = Array.isArray(guildList) ? guildList : [];
    guilds.set(guildsArray);
    
    if (guildsArray.length > 0) {
      const firstGuild = guildsArray[0];
      currentGuild.set(firstGuild);
      await loadGuildChannels(firstGuild.id);
    }
  } catch (err) {
    console.error('Failed to load initial data:', err);
    error.set(err.message);
  } finally {
    loading.set(false);
  }
}

// Load channels for a guild
export async function loadGuildChannels(guildId) {
  loading.set(true);
  error.set(null);
  
  try {
    const channelList = await channelApi.getByGuildId(guildId);
    // Ensure channelList is an array
    const channelsArray = Array.isArray(channelList) ? channelList : [];
    channels.set(channelsArray);
    
    if (channelsArray.length > 0) {
      const firstChannel = channelsArray[0];
      currentChannel.set(firstChannel);
      await loadChannelLogs(firstChannel.id);
    } else {
      // Clear logs if no channels
      logs.set([]);
      currentChannel.set(null);
    }
  } catch (err) {
    console.error('Failed to load channels:', err);
    error.set(err.message);
  } finally {
    loading.set(false);
  }
}

// Load logs for a channel
export async function loadChannelLogs(channelId) {
  error.set(null);
  
  try {
    const logList = await logApi.getByChannelId(channelId);
    // Ensure logList is an array and reverse the order so newest logs appear at the bottom
    const logsArray = Array.isArray(logList) ? logList : [];
    logs.set(logsArray.reverse());
  } catch (err) {
    error.set(err.message);
    console.error('Failed to load logs:', err);
    // Set empty array on error to prevent null reference
    logs.set([]);
  }
}

// Set active guild and load its channels
export async function setActiveGuild(guild) {
  currentGuild.set(guild);
  await loadGuildChannels(guild.id);
}

// Set active channel and load its logs
export async function setActiveChannel(channel) {
  loading.set(true);
  currentChannel.set(channel);
  
  try {
    await loadChannelLogs(channel.id);
  } finally {
    loading.set(false);
  }
}

// Create new guild
export async function createGuild(name) {
  loading.set(true);
  error.set(null);
  
  try {
    const newGuild = await guildApi.create(name);
    guilds.update(currentGuilds => [...currentGuilds, newGuild]);
    return newGuild;
  } catch (err) {
    error.set(err.message);
    console.error('Failed to create guild:', err);
    throw err;
  } finally {
    loading.set(false);
  }
}

// Delete guild
export async function deleteGuild(guildId) {
  loading.set(true);
  error.set(null);
  
  try {
    await guildApi.delete(guildId);
    guilds.update(currentGuilds => currentGuilds.filter(g => g.id !== guildId));
    
    // If we deleted the current guild, switch to the first available guild
    if ($currentGuild && $currentGuild.id === guildId) {
      guilds.update(remainingGuilds => {
        if (remainingGuilds.length > 0) {
          setActiveGuild(remainingGuilds[0]);
        } else {
          currentGuild.set(null);
          channels.set([]);
          logs.set([]);
          currentChannel.set(null);
        }
        return remainingGuilds;
      });
    }
  } catch (err) {
    error.set(err.message);
    console.error('Failed to delete guild:', err);
    throw err;
  } finally {
    loading.set(false);
  }
}

// Create new channel
export async function createChannel(guildId, name) {
  loading.set(true);
  error.set(null);
  
  try {
    const newChannel = await channelApi.create(guildId, name);
    channels.update(currentChannels => [...currentChannels, newChannel]);
    return newChannel;
  } catch (err) {
    error.set(err.message);
    console.error('Failed to create channel:', err);
    throw err;
  } finally {
    loading.set(false);
  }
}

// Delete channel
export async function deleteChannel(channelId) {
  loading.set(true);
  error.set(null);
  
  try {
    await channelApi.delete(channelId);
    channels.update(currentChannels => currentChannels.filter(c => c.id !== channelId));
    
    // If we deleted the current channel, switch to the first available channel
    if ($currentChannel && $currentChannel.id === channelId) {
      channels.update(remainingChannels => {
        if (remainingChannels.length > 0) {
          setActiveChannel(remainingChannels[0]);
        } else {
          logs.set([]);
          currentChannel.set(null);
        }
        return remainingChannels;
      });
    }
  } catch (err) {
    error.set(err.message);
    console.error('Failed to delete channel:', err);
    throw err;
  } finally {
    loading.set(false);
  }
}

// Create new log
export async function createLog(channelId, content, files = []) {
  loading.set(true);
  error.set(null);
  
  try {
    const newLog = await logApi.create(channelId, content);
    
    // Upload files if any
    if (files && files.length > 0) {
      try {
        const formData = new FormData();
        files.forEach(file => {
          formData.append('files', file);
        });
        
        const response = await fetch(`/api/logs/${newLog.id}/files`, {
          method: 'POST',
          body: formData
        });
        
        if (response.ok) {
          const fileData = await response.json();
          newLog.files = fileData.files || [];
        } else {
          console.error('File upload failed:', response.status, response.statusText);
          newLog.files = [];
        }
      } catch (fileErr) {
        console.error('File upload error:', fileErr);
        newLog.files = [];
      }
    } else {
      newLog.files = [];
    }
    
    // Get current logs and add new log to the end (bottom)
    logs.update(currentLogs => {
      const logsArray = Array.isArray(currentLogs) ? currentLogs : [];
      return [...logsArray, newLog];
    });
    return newLog;
  } catch (err) {
    error.set(err.message);
    console.error('Failed to create log:', err);
    throw err;
  } finally {
    loading.set(false);
  }
}

// Delete log
export async function deleteLog(logId) {
  loading.set(true);
  error.set(null);
  
  try {
    await logApi.delete(logId);
    logs.update(currentLogs => {
      const logsArray = Array.isArray(currentLogs) ? currentLogs : [];
      return logsArray.filter(l => l.id !== logId);
    });
  } catch (err) {
    error.set(err.message);
    console.error('Failed to delete log:', err);
    throw err;
  } finally {
    loading.set(false);
  }
}

// Toggle sidebar (for mobile)
export function toggleSidebar() {
  sidebarOpen.update(open => !open);
}

// Close sidebar
export function closeSidebar() {
  sidebarOpen.set(false);
}

// Search functionality
export function searchLogs(query) {
  searchQuery.set(query);
  // This could be enhanced to filter logs in real-time
}

// Clear error
export function clearError() {
  error.set(null);
}