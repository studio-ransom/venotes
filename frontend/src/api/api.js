// API service for communicating with the Go backend
const API_BASE_URL = '/api';

class ApiError extends Error {
  constructor(message, status) {
    super(message);
    this.status = status;
  }
}

async function apiRequest(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
      throw new ApiError(errorData.error || `HTTP ${response.status}`, response.status);
    }

    return await response.json();
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError(`Network error: ${error.message}`, 0);
  }
}

// Guild API
export const guildApi = {
  async getAll() {
    return apiRequest('/guilds');
  },

  async getById(id) {
    return apiRequest(`/guilds/${id}`);
  },

  async create(name) {
    return apiRequest('/guilds', {
      method: 'POST',
      body: JSON.stringify({ name }),
    });
  },

  async delete(id) {
    return apiRequest(`/guilds/${id}`, {
      method: 'DELETE',
    });
  },
};

// Channel API
export const channelApi = {
  async getByGuildId(guildId) {
    return apiRequest(`/guilds/${guildId}/channels`);
  },

  async getById(id) {
    return apiRequest(`/channels/${id}`);
  },

  async create(guildId, name) {
    return apiRequest(`/guilds/${guildId}/channels`, {
      method: 'POST',
      body: JSON.stringify({ name }),
    });
  },

  async delete(id) {
    return apiRequest(`/channels/${id}`, {
      method: 'DELETE',
    });
  },
};

// Log API
export const logApi = {
  async getByChannelId(channelId) {
    return apiRequest(`/channels/${channelId}/logs`);
  },

  async getById(id) {
    return apiRequest(`/logs/${id}`);
  },

  async create(channelId, content) {
    return apiRequest(`/channels/${channelId}/logs`, {
      method: 'POST',
      body: JSON.stringify({ content }),
    });
  },

  async update(id, content) {
    return apiRequest(`/logs/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ content }),
    });
  },

  async delete(id) {
    return apiRequest(`/logs/${id}`, {
      method: 'DELETE',
    });
  },
};

export { ApiError };
