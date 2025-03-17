import api from './api';

export const authService = {
  login: async (credentials) => {
    return api.post('/auth/login', credentials);
  },

  logout: async () => {
    return api.post('/auth/logout');
  },

  getCurrentUser: async () => {
    return api.get('/auth/me');
  },

  refreshToken: async () => {
    return api.post('/auth/refresh');
  },

  forgotPassword: async (email) => {
    return api.post('/auth/forgot-password', { email });
  },

  resetPassword: async (token, password) => {
    return api.post('/auth/reset-password', { token, password });
  },

  changePassword: async (currentPassword, newPassword) => {
    return api.post('/auth/change-password', {
      current_password: currentPassword,
      new_password: newPassword,
    });
  },
}; 