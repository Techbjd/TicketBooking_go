import { Box, Container, Typography, Grid, Button, Chip, Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material'
import { useParams, useNavigate } from 'react-router-dom'
import AccessTimeIcon from '@mui/icons-material/AccessTime'
import StarIcon from '@mui/icons-material/Star'
import CalendarMonthIcon from '@mui/icons-material/CalendarMonth'
import { useState } from 'react'

function MovieDetails({ movies }) {
  const { id } = useParams()
  const navigate = useNavigate()
  const [selectedTime, setSelectedTime] = useState(null)
  const [openDialog, setOpenDialog] = useState(false)

  const movie = movies.find(m => m.id === id)

  if (!movie) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography variant="h5" sx={{ color: '#fff' }}>Movie not found</Typography>
      </Container>
    )
  }

  const handleTimeSelect = (time) => {
    setSelectedTime(time)
    setOpenDialog(true)
  }

  const handleConfirm = () => {
    setOpenDialog(false)
    navigate(`/movie/${id}/seats/${encodeURIComponent(selectedTime)}`)
  }

  return (
    <Box sx={{ py: 4 }}>
      <Container maxWidth="lg">
        <Grid container spacing={4}>
          <Grid item xs={12} md={4}>
            <Box sx={{ 
              height: 400, 
              bgcolor: '#333', 
              borderRadius: 2,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <Typography variant="h1" sx={{ fontSize: 100, color: '#444' }}>
                🎬
              </Typography>
            </Box>
          </Grid>
          <Grid item xs={12} md={8}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
              <Typography variant="h3" sx={{ color: '#fff' }}>
                {movie.title}
              </Typography>
              <Chip 
                icon={<StarIcon sx={{ color: '#ffd700 !important' }} />} 
                label={movie.rating} 
                sx={{ bgcolor: 'rgba(255,215,0,0.1)', color: '#ffd700' }}
              />
            </Box>
            
            <Box sx={{ display: 'flex', gap: 2, mb: 3 }}>
              <Chip label={movie.genre} sx={{ bgcolor: '#e50914' }} />
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                <AccessTimeIcon sx={{ color: '#aaa' }} />
                <Typography sx={{ color: '#aaa' }}>{movie.duration} min</Typography>
              </Box>
            </Box>

            <Typography sx={{ color: '#ccc', mb: 4, lineHeight: 1.8 }}>
              {movie.description}
            </Typography>

            <Typography variant="h6" sx={{ color: '#fff', mb: 2 }}>
              Select Show Time
            </Typography>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2, mb: 4 }}>
              {movie.show_times.map((time) => (
                <Button 
                  key={time}
                  variant="outlined"
                  onClick={() => handleTimeSelect(time)}
                  sx={{ 
                    borderColor: '#444',
                    color: '#fff',
                    px: 3,
                    '&:hover': { borderColor: '#e50914', bgcolor: 'rgba(229,9,20,0.1)' }
                  }}
                >
                  {time}
                </Button>
              ))}
            </Box>
          </Grid>
        </Grid>
      </Container>

      <Dialog open={openDialog} onClose={() => setOpenDialog(false)}>
        <DialogTitle>Confirm Booking</DialogTitle>
        <DialogContent>
          <Typography>
            Movie: <strong>{movie.title}</strong><br/>
            Time: <strong>{selectedTime}</strong>
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
          <Button variant="contained" onClick={handleConfirm} sx={{ bgcolor: '#e50914' }}>
            Continue to Seats
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}

export default MovieDetails