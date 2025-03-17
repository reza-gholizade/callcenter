import { configureStore } from '@reduxjs/toolkit';
import chatReducer from './slices/chatSlice';
import ticketReducer from './slices/ticketSlice';
import authReducer from './slices/authSlice';

export const store = configureStore({
  reducer: {
    chat: chatReducer,
    tickets: ticketReducer,
    auth: authReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false,
    }),
}); 