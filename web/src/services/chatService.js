import api from './api';

export const chatService = {
  sendMessage: async (message) => {
    return api.post('/chat/message', message);
  },

  getHistory: async (sessionId) => {
    return api.get(`/chat/history/${sessionId}`);
  },

  createSession: async (platform, userId) => {
    return api.post('/chat/session', { platform, user_id: userId });
  },

  closeSession: async (sessionId) => {
    return api.post(`/chat/session/${sessionId}/close`);
  },

  escalateSession: async (sessionId) => {
    return api.post(`/chat/session/${sessionId}/escalate`);
  },
}; 