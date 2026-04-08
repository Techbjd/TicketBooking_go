import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Box, Container, Typography, Button, Grid, Paper, Alert, Snackbar, Chip } from '@mui/material'
import CheckCircleIcon from '@mui/icons-material/CheckCircle'
import WarningIcon from '@mui/icons-material/Warning'

function SeatSelection({ movies }) {
  const { id, showTime } = useParams()
  const navigate = useNavigate()
  const [seats, setSeats] = useState([])
  const [selectedSeats, setSelectedSeats] = useState([])
  const [loading, setLoading] = useState(true)
  const [userId, setUserId] = useState('')
  const [showUserInput, setShowUserInput] = useState(false)
  const [booking, setBooking] = useState(null)
  const [error, setError] = useState(null)
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' })

  const movie = movies.find(m => m.id === id)
  const rows = movie?.rows || 5
  const seatsPerRow = movie?.seat_per_row || 8
  const decodedShowTime = showTime ? decodeURIComponent(showTime) : ''

  useEffect(() => {
    fetchSeats()
  }, [id, showTime])

  const fetchSeats = async () => {
    setLoading(true)
    try {
      const res = await fetch(`http://localhost:8080/movies/${id}/seats?rows=${rows}&seats_per_row=${seatsPerRow}&show_time=${encodeURIComponent(decodedShowTime)}`)
      const data = await res.json()
      setSeats(data)
    } catch (err) {
      setError('Failed to load seats')
    } finally {
      setLoading(false)
    }
  }

  const toggleSeat = (seatId) => {
    const seat = seats.find(s => s.seat_id === seatId)
    if (seat?.booked) return
    
    setSelectedSeats(prev => 
      prev.includes(seatId) 
        ? prev.filter(s => s !== seatId)
        : [...prev, seatId]
    )
  }

  const handleBookSeats = async () => {
    if (!userId.trim()) {
      setSnackbar({ open: true, message: 'Please enter your name', severity: 'error' })
      return
    }
    
    if (selectedSeats.length === 0) {
      setSnackbar({ open: true, message: 'Please select at least one seat', severity: 'error' })
      return
    }

    try {
      const seatId = selectedSeats[0]
      const res = await fetch(`http://localhost:8080/movies/${id}/seats/${seatId}?show_time=${encodeURIComponent(decodedShowTime)}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId })
      })
      
      const data = await res.json()
      
      if (!res.ok) {
        throw new Error(data.message || 'Failed to book seat')
      }

      setBooking(data)
      setShowUserInput(false)
      setSnackbar({ open: true, message: 'Seat booked successfully!', severity: 'success' })
      fetchSeats()
    } catch (err) {
      setSnackbar({ open: true, message: err.message, severity: 'error' })
    }
  }

  const handleConfirm = async () => {
    if (!booking) return
    
    try {
      const res = await fetch(`http://localhost:8080/session/${booking.session_id}/confirm`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId })
      })
      
      if (!res.ok) throw new Error('Failed to confirm booking')
      
      setSnackbar({ open: true, message: 'Booking confirmed! Enjoy your movie!', severity: 'success' })
      setTimeout(() => navigate('/'), 2000)
    } catch (err) {
      setSnackbar({ open: true, message: err.message, severity: 'error' })
    }
  }

  const handleRelease = async () => {
    if (!booking) return
    
    try {
      await fetch(`http://localhost:8080/session/${booking.session_id}`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId })
      })
      setBooking(null)
      fetchSeats()
    } catch (err) {
      setSnackbar({ open: true, message: 'Failed to release seat', severity: 'error' })
    }
  }

  if (!movie) {
    return <Typography sx={{ color: '#fff', p: 4 }}>Movie not found</Typography>
  }

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Button onClick={() => navigate(`/movie/${id}`)} sx={{ color: '#aaa', mb: 2 }}>
        ← Back to Movie
      </Button>

      <Paper sx={{ p: 3, bgcolor: '#1f1f1f', borderRadius: 2 }}>
        <Typography variant="h5" sx={{ color: '#fff', mb: 1 }}>
          {movie.title}
        </Typography>
        <Chip 
          label={`Show Time: ${decodedShowTime}`} 
          sx={{ bgcolor: '#e50914', mb: 2 }}
        />
        <Typography sx={{ color: '#aaa', mb: 3 }}>
          Select your seats
        </Typography>

        <Box sx={{ 
          bgcolor: '#2a2a2a', 
          py: 4, 
          px: 2, 
          borderRadius: 2, 
          mb: 3,
          position: 'relative'
        }}>
          <Box sx={{ 
            width: '60%', 
            height: 40, 
            bgcolor: '#444', 
            borderRadius: '0 0 50% 50% / 0 0 100% 100%',
            mx: 'auto',
            mb: 4,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}>
            <Typography sx={{ color: '#666', fontSize: 12 }}>SCREEN</Typography>
          </Box>

          {loading ? (
            <Typography sx={{ color: '#aaa' }}>Loading seats...</Typography>
          ) : (
            <Grid container spacing={1} justifyContent="center">
              {seats.map((seat) => {
                const isSelected = selectedSeats.includes(seat.seat_id)
                const isBooked = seat.booked
                
                return (
                  <Grid item key={seat.seat_id}>
                    <Box
                      onClick={() => toggleSeat(seat.seat_id)}
                      sx={{
                        width: 40,
                        height: 40,
                        borderRadius: 1,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        cursor: isBooked ? 'not-allowed' : 'pointer',
                        bgcolor: isBooked 
                          ? '#e50914' 
                          : isSelected 
                            ? '#ffd700' 
                            : '#333',
                        color: isBooked || isSelected ? '#000' : '#666',
                        fontWeight: 600,
                        fontSize: 12,
                        transition: 'all 0.2s',
                        '&:hover': !isBooked && {
                          transform: 'scale(1.1)',
                          bgcolor: isSelected ? '#ffd700' : '#444'
                        }
                      }}
                    >
                      {seat.seat_id.replace('R', '').replace('S', '')}
                    </Box>
                  </Grid>
                )
              })}
            </Grid>
          )}
        </Box>

        <Box sx={{ display: 'flex', gap: 3, mb: 3, justifyContent: 'center' }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Box sx={{ width: 20, height: 20, bgcolor: '#333', borderRadius: 0.5 }} />
            <Typography sx={{ color: '#aaa', fontSize: 12 }}>Available</Typography>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Box sx={{ width: 20, height: 20, bgcolor: '#ffd700', borderRadius: 0.5 }} />
            <Typography sx={{ color: '#aaa', fontSize: 12 }}>Selected</Typography>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Box sx={{ width: 20, height: 20, bgcolor: '#e50914', borderRadius: 0.5 }} />
            <Typography sx={{ color: '#aaa', fontSize: 12 }}>Booked</Typography>
          </Box>
        </Box>

        {selectedSeats.length > 0 && !booking && (
          <Box sx={{ mt: 3, textAlign: 'center' }}>
            <Typography sx={{ color: '#fff', mb: 2 }}>
              Selected: <strong>{selectedSeats.join(', ')}</strong>
            </Typography>
            {!showUserInput ? (
              <Button 
                variant="contained" 
                onClick={() => setShowUserInput(true)}
                sx={{ bgcolor: '#e50914', px: 4 }}
              >
                Book Selected Seats
              </Button>
            ) : (
              <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center', alignItems: 'center' }}>
                <input
                  type="text"
                  placeholder="Enter your name"
                  value={userId}
                  onChange={(e) => setUserId(e.target.value)}
                  style={{
                    padding: '10px 15px',
                    borderRadius: 4,
                    border: '1px solid #444',
                    background: '#333',
                    color: '#fff',
                    fontSize: 16
                  }}
                />
                <Button 
                  variant="contained" 
                  onClick={handleBookSeats}
                  sx={{ bgcolor: '#e50914', px: 4 }}
                >
                  Confirm Booking
                </Button>
              </Box>
            )}
          </Box>
        )}

        {booking && (
          <Box sx={{ mt: 3, textAlign: 'center' }}>
            <Alert severity="success" sx={{ mb: 2 }}>
              <Typography variant="body1">
                <CheckCircleIcon sx={{ mr: 1, verticalAlign: 'middle' }} />
                Seat held! Session expires at: {booking.expires_at}
              </Typography>
            </Alert>
            <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center' }}>
              <Button 
                variant="contained" 
                onClick={handleConfirm}
                sx={{ bgcolor: '#4caf50', px: 4 }}
              >
                Confirm Booking
              </Button>
              <Button 
                variant="outlined" 
                onClick={handleRelease}
                sx={{ borderColor: '#f44336', color: '#f44336' }}
              >
                Release Seat
              </Button>
            </Box>
          </Box>
        )}
      </Paper>

      <Snackbar 
        open={snackbar.open} 
        autoHideDuration={4000} 
        onClose={() => setSnackbar({ ...snackbar, open: false })}
      >
        <Alert severity={snackbar.severity} onClose={() => setSnackbar({ ...snackbar, open: false })}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Container>
  )
}

export default SeatSelection