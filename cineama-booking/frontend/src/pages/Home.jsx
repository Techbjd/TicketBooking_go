import { Box, Container, Typography, Grid, Card, CardMedia, CardContent, CardActionArea } from '@mui/material'
import { useNavigate } from 'react-router-dom'
import AccessTimeIcon from '@mui/icons-material/AccessTime'
import StarIcon from '@mui/icons-material/Star'

function Home({ movies }) {
  const navigate = useNavigate()

  if (!movies || movies.length === 0) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Typography variant="h4" sx={{ mb: 4, color: '#fff' }}>
          Now Showing
        </Typography>
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <Typography variant="h6" sx={{ color: '#aaa' }}>
            Loading movies... or backend not running
          </Typography>
          <Typography sx={{ color: '#666', mt: 2 }}>
            Make sure the backend is running on port 8080
          </Typography>
        </Box>
      </Container>
    )
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" sx={{ mb: 4, color: '#fff' }}>
        Now Showing
      </Typography>
      <Grid container spacing={3}>
        {movies.map((movie) => (
          <Grid item xs={12} sm={6} md={4} key={movie.id}>
            <Card sx={{ bgcolor: 'background.paper', borderRadius: 2, overflow: 'hidden' }}>
              <CardActionArea onClick={() => navigate(`/movie/${movie.id}`)}>
                <Box sx={{ 
                  height: 280, 
                  bgcolor: '#333',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  position: 'relative'
                }}>
                  <Typography variant="h1" sx={{ fontSize: 60, color: '#444' }}>
                    🎬
                  </Typography>
                  <Box sx={{ 
                    position: 'absolute', 
                    top: 8, 
                    right: 8,
                    bgcolor: 'rgba(0,0,0,0.7)',
                    px: 1,
                    py: 0.5,
                    borderRadius: 1,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 0.5
                  }}>
                    <StarIcon sx={{ fontSize: 16, color: '#ffd700' }} />
                    <Typography variant="body2" sx={{ color: '#fff' }}>
                      {movie.rating}
                    </Typography>
                  </Box>
                </Box>
                <CardContent>
                  <Typography variant="h6" sx={{ color: '#fff', mb: 1 }} noWrap>
                    {movie.title}
                  </Typography>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="body2" sx={{ color: '#aaa' }}>
                      {movie.genre}
                    </Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                      <AccessTimeIcon sx={{ fontSize: 14, color: '#aaa' }} />
                      <Typography variant="body2" sx={{ color: '#aaa' }}>
                        {movie.duration} min
                      </Typography>
                    </Box>
                  </Box>
                </CardContent>
              </CardActionArea>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Container>
  )
}

export default Home