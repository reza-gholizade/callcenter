import React, { useState, useEffect } from 'react';
import {
  Box,
  Paper,
  Typography,
  Grid,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  CircularProgress,
  Card,
  CardContent,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from '@mui/material';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

const AdminDashboard = () => {
  const [timeRange, setTimeRange] = useState('24h');
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    totalTickets: 0,
    activeTickets: 0,
    cancelledTickets: 0,
    totalChats: 0,
    activeChats: 0,
    escalatedChats: 0,
  });
  const [recentTickets, setRecentTickets] = useState([]);
  const [chatMetrics, setChatMetrics] = useState([]);

  useEffect(() => {
    fetchDashboardData();
  }, [timeRange]);

  const fetchDashboardData = async () => {
    setLoading(true);
    try {
      // TODO: Implement actual API calls
      // For now, using mock data
      setStats({
        totalTickets: 150,
        activeTickets: 120,
        cancelledTickets: 30,
        totalChats: 200,
        activeChats: 45,
        escalatedChats: 15,
      });

      setRecentTickets([
        {
          id: 1,
          ticketNumber: 'TKT123',
          status: 'active',
          passenger: 'John Doe',
          date: new Date(),
        },
        // Add more mock data as needed
      ]);

      setChatMetrics([
        { time: '00:00', messages: 10, users: 5 },
        { time: '04:00', messages: 15, users: 8 },
        { time: '08:00', messages: 25, users: 12 },
        { time: '12:00', messages: 35, users: 18 },
        { time: '16:00', messages: 30, users: 15 },
        { time: '20:00', messages: 20, users: 10 },
      ]);
    } catch (error) {
      console.error('Error fetching dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          height: '100vh',
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Admin Dashboard
      </Typography>

      {/* Time Range Selector */}
      <FormControl sx={{ mb: 3, minWidth: 200 }}>
        <InputLabel>Time Range</InputLabel>
        <Select
          value={timeRange}
          label="Time Range"
          onChange={(e) => setTimeRange(e.target.value)}
        >
          <MenuItem value="24h">Last 24 Hours</MenuItem>
          <MenuItem value="7d">Last 7 Days</MenuItem>
          <MenuItem value="30d">Last 30 Days</MenuItem>
        </Select>
      </FormControl>

      {/* Stats Cards */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={4}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Total Tickets
              </Typography>
              <Typography variant="h4">{stats.totalTickets}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Active Chats
              </Typography>
              <Typography variant="h4">{stats.activeChats}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Escalated Chats
              </Typography>
              <Typography variant="h4">{stats.escalatedChats}</Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Chat Metrics Chart */}
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom>
          Chat Activity
        </Typography>
        <Box sx={{ height: 300 }}>
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={chatMetrics}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Line
                type="monotone"
                dataKey="messages"
                stroke="#8884d8"
                name="Messages"
              />
              <Line
                type="monotone"
                dataKey="users"
                stroke="#82ca9d"
                name="Active Users"
              />
            </LineChart>
          </ResponsiveContainer>
        </Box>
      </Paper>

      {/* Recent Tickets Table */}
      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Recent Tickets
        </Typography>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Ticket Number</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Passenger</TableCell>
                <TableCell>Date</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {recentTickets.map((ticket) => (
                <TableRow key={ticket.id}>
                  <TableCell>{ticket.ticketNumber}</TableCell>
                  <TableCell>{ticket.status}</TableCell>
                  <TableCell>{ticket.passenger}</TableCell>
                  <TableCell>
                    {new Date(ticket.date).toLocaleString()}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
    </Box>
  );
};

export default AdminDashboard; 