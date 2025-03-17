import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { ticketService } from '../../services/ticketService';

export const getTicketDetails = createAsyncThunk(
  'tickets/getDetails',
  async (ticketNumber, { rejectWithValue }) => {
    try {
      const response = await ticketService.getTicketDetails(ticketNumber);
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const cancelTicket = createAsyncThunk(
  'tickets/cancel',
  async ({ ticketNumber, reason }, { rejectWithValue }) => {
    try {
      const response = await ticketService.cancelTicket(ticketNumber, reason);
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const getRefundStatus = createAsyncThunk(
  'tickets/getRefundStatus',
  async (ticketNumber, { rejectWithValue }) => {
    try {
      const response = await ticketService.getRefundStatus(ticketNumber);
      return response.data;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

const initialState = {
  currentTicket: null,
  refundStatus: null,
  loading: false,
  error: null,
};

const ticketSlice = createSlice({
  name: 'tickets',
  initialState,
  reducers: {
    clearTicket: (state) => {
      state.currentTicket = null;
      state.refundStatus = null;
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(getTicketDetails.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(getTicketDetails.fulfilled, (state, action) => {
        state.loading = false;
        state.currentTicket = action.payload.ticket;
      })
      .addCase(getTicketDetails.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload?.error || 'Failed to get ticket details';
      })
      .addCase(cancelTicket.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(cancelTicket.fulfilled, (state, action) => {
        state.loading = false;
        if (state.currentTicket) {
          state.currentTicket.status = 'cancelled';
        }
      })
      .addCase(cancelTicket.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload?.error || 'Failed to cancel ticket';
      })
      .addCase(getRefundStatus.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(getRefundStatus.fulfilled, (state, action) => {
        state.loading = false;
        state.refundStatus = action.payload.refund_request;
      })
      .addCase(getRefundStatus.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload?.error || 'Failed to get refund status';
      });
  },
});

export const { clearTicket } = ticketSlice.actions;
export default ticketSlice.reducer; 