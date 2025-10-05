// Data structures
export class LogData {
  constructor(id, timestamp, content) {
    this.id = id;
    this.timestamp = timestamp;
    this.content = content;
  }
}

export class GuildData {
  constructor(name) {
    this.name = name;
    this.channels = {};
  }
}

// Storage key for localStorage
const STORAGE_KEY = 'venotes-data';

// Generate unique ID
function generateId() {
  return 'log_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
}

// Get all data from localStorage
function getAllData() {
  const data = localStorage.getItem(STORAGE_KEY);
  return data ? JSON.parse(data) : { guilds: {} };
}

// Save all data to localStorage
function saveAllData(data) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
}

// Guild API
export function get_guilds() {
  const data = getAllData();
  return Object.keys(data.guilds);
}

export function load_guild(guildName) {
  const data = getAllData();
  return data.guilds[guildName] || new GuildData(guildName);
}

export function save_guild(guildName, guildData) {
  const data = getAllData();
  data.guilds[guildName] = guildData;
  saveAllData(data);
}

export function create_guild(guildName) {
  const guildData = new GuildData(guildName);
  save_guild(guildName, guildData);
  return guildData;
}

export function delete_guild(guildName) {
  const data = getAllData();
  delete data.guilds[guildName];
  saveAllData(data);
}

// Channel API
export function get_channels(guild) {
  const guildData = load_guild(guild);
  return Object.keys(guildData.channels);
}

export function create_channel(guild, channelName) {
  const guildData = load_guild(guild);
  guildData.channels[channelName] = [];
  save_guild(guild, guildData);
}

export function delete_channel(guild, channelName) {
  const guildData = load_guild(guild);
  delete guildData.channels[channelName];
  save_guild(guild, guildData);
}

// Logs API
export function get_logs(guild, channel) {
  const guildData = load_guild(guild);
  return guildData.channels[channel] || [];
}

export function channel_post_log(guild, channel, content) {
  const guildData = load_guild(guild);
  if (!guildData.channels[channel]) {
    guildData.channels[channel] = [];
  }
  
  const log = new LogData(
    generateId(),
    new Date().toISOString(),
    content
  );
  
  guildData.channels[channel].push(log);
  save_guild(guild, guildData);
  return log;
}

export function channel_delete_log(guild, channel, logId) {
  const guildData = load_guild(guild);
  if (guildData.channels[channel]) {
    guildData.channels[channel] = guildData.channels[channel].filter(
      log => log.id !== logId
    );
    save_guild(guild, guildData);
  }
}

// Storage API
export function load_all() {
  return getAllData();
}

export function save_all() {
  return getAllData();
}

export function export_data() {
  const data = getAllData();
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'venotes-backup.json';
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

export function import_data(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const data = JSON.parse(e.target.result);
        saveAllData(data);
        resolve(data);
      } catch (error) {
        reject(error);
      }
    };
    reader.readAsText(file);
  });
}

// Initialize with default data if empty
export function initialize_default_data() {
  const data = getAllData();
  if (Object.keys(data.guilds).length === 0) {
    // Create default guilds
    const workGuild = new GuildData('Work');
    workGuild.channels = {
      'Coding': [],
      'Ideas': [],
      'Meetings': []
    };
    
    const personalGuild = new GuildData('Personal');
    personalGuild.channels = {
      'Thoughts': [],
      'Ideas': [],
      'Journal': []
    };
    
    data.guilds = {
      'Work': workGuild,
      'Personal': personalGuild
    };
    
    saveAllData(data);
  }
}
