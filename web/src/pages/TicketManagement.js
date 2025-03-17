import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Box,
  Paper,
  TextField,
  Button,
  Typography,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Alert,
} from '@mui/material';
import { getTicketDetails, cancelTicket, getRefundStatus } from '../store/slices/ticketSlice';

const TicketManagement = () => {
  const dispatch = useDispatch();
  const { currentTicket, refundStatus, loading, error } = useSelector((state) => state.tickets);
  const [ticketNumber, setTicketNumber] = useState('');
  const [cancelDialogOpen, setCancelDialogOpen] = useState(false);
  const [cancelReason, setCancelReason] = useState('');

  const handleSearch = async (e) => {
    e.preventDefault();
    if (!ticketNumber.trim()) return;

    await dispatch(getTicketDetails(ticketNumber));
    await dispatch(getRefundStatus(ticketNumber));
  };

  const handleCancelClick = () => {
    setCancelDialogOpen(true);
  };

  const handleCancelConfirm = async () => {
    if (!cancelReason.trim()) return;

    await dispatch(cancelTicket({ ticketNumber, reason: cancelReason }));
    setCancelDialogOpen(false);
    setCancelReason('');
  };

  return (
    <Box sx={{ p: 3 }}>
      <Paper elevation={3} sx={{ p: 3 }}>
        <Typography variant="h5" gutterBottom>
          Ticket Management
        </Typography>

        {/* Search Form */}
        <Box component="form" onSubmit={handleSearch} sx={{ mb: 3 }}>
          <TextField
            fullWidth
            label="Ticket Number"
            value={ticketNumber}
            onChange={(e) => setTicketNumber(e.target.value)}
            sx={{ mb: 2 }}
          />
          <Button
            variant="contained"
            type="submit"
            disabled={loading || !ticketNumber.trim()}
          >
            {loading ? <CircularProgress size={24} /> : 'Search'}
          </Button>
        </Box>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {/* Ticket Details */}
        {currentTicket && (
          <Box>
            <Typography variant="h6" gutterBottom>
              Ticket Details
            </Typography>
            <TableContainer>
              <Table>
                <TableBody>
                  <TableRow>
                    <TableCell component="th">Ticket Number</TableCell>
                    <TableCell>{currentTicket.ticket_number}</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell component="th">Status</TableCell>
                    <TableCell>{currentTicket.status}</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell component="th">Airline</TableCell>
                    <TableCell>{currentTicket.airline}</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell component="th">Flight Number</TableCell>
                    <TableCell>{currentTicket.flight_number}</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell component="th">Departure Date</TableCell>
                    <TableCell>
                      {new Date(currentTicket.departure_date).toLocaleString()}
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell component="th">Arrival Date</TableCell>
                    <TableCell>
                      {new Date(currentTicket.arrival_date).toLocaleString()}
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell component="th">Passenger Name</TableCell>
                    <TableCell>{currentTicket.passenger_name}</TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell component="th">Price</TableCell>
                    <TableCell>
                      {currentTicket.price} {currentTicket.currency}
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </TableContainer>

            {/* Actions */}
            <Box sx={{ mt: 2 }}>
              {currentTicket.status === 'active' && (
                <Button
                  variant="contained"
                  color="error"
                  onClick={handleCancelClick}
                >
                  Cancel Ticket
                </Button>
              )}
            </Box>

            {/* Refund Status */}
            {refundStatus && (
              <Box sx={{ mt: 3 }}>
                <Typography variant="h6" gutterBottom>
                  Refund Status
                </Typography>
                <TableContainer>
                  <Table>
                    <TableBody>
                      <TableRow>
                        <TableCell component="th">Status</TableCell>
                        <TableCell>{refundStatus.status}</TableCell>
                      </TableRow>
                      <TableRow>
                        <TableCell component="th">Amount</TableCell>
                        <TableCell>
                          {refundStatus.amount} {refundStatus.currency}
                        </TableCell>
                      </TableRow>
                      <TableRow>
                        <TableCell component="th">Requested At</TableCell>
                        <TableCell>
                          {new Date(refundStatus.created_at).toLocaleString()}
                        </TableCell>
                      </TableRow>
                    </TableBody>
                  </Table>
                </TableContainer>
              </Box>
            )}
          </Box>
        )}
      </Paper>

      {/* Cancel Dialog */}
      <Dialog open={cancelDialogOpen} onClose={() => setCancelDialogOpen(false)}>
        <DialogTitle>Cancel Ticket</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Reason for Cancellation"
            fullWidth
            multiline
            rows={4}
            value={cancelReason}
            onChange={(e) => setCancelReason(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCancelDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleCancelConfirm} color="error">
            Confirm Cancellation
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default TicketManagement; 