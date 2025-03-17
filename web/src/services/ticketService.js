import api from './api';

export const ticketService = {
  getTicketDetails: async (ticketNumber) => {
    return api.get(`/tickets/${ticketNumber}`);
  },

  cancelTicket: async (ticketNumber, reason) => {
    return api.post(`/tickets/${ticketNumber}/cancel`, { reason });
  },

  getRefundStatus: async (ticketNumber) => {
    return api.get(`/tickets/${ticketNumber}/refund-status`);
  },

  updateRefundStatus: async (ticketNumber, status, processedBy) => {
    return api.put(`/tickets/${ticketNumber}/refund-status`, {
      status,
      processed_by: processedBy,
    });
  },

  searchTickets: async (params) => {
    return api.get('/tickets/search', { params });
  },

  getTicketHistory: async (ticketNumber) => {
    return api.get(`/tickets/${ticketNumber}/history`);
  },
}; 