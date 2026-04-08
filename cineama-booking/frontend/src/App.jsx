import { useState, useEffect } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { ThemeProvider, createTheme, CssBaseline, Box, AppBar, Toolbar, Typography, Button } from '@mui/material'
import Home from './pages/Home'
import MovieDetails from './pages/MovieDetails'
import SeatSelection from './pages/SeatSelection'

const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: { main: '#e50914' },
    secondary: { main: '#ffd700' },
    background: { default: '#141414', paper: '#1f1f1f' },
  },
  typography: {
    fontFamily: '"Netflix Sans", "Helvetica", "Arial", sans-serif',
    h4: { fontWeight: 700 },
    h5: { fontWeight: 600 },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: { textTransform: 'none', borderRadius: 4 },
      },
    },
  },
})

function App() {
  const [movies, setMovies] = useState([])

  useEffect(() => {
    fetch('http://localhost:8080/movies')
      .then(res => res.json())
      .then(data => setMovies(data))
      .catch(err => console.error('Failed to fetch movies:', err))
  }, [])

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <BrowserRouter>
        <Box sx={{ minHeight: '100vh', bgcolor: 'background.default' }}>
          <AppBar position="static" sx={{ bgcolor: '#000' }}>
            <Toolbar>
              <Typography variant="h5" sx={{ flexGrow: 1, color: '#e50914', fontWeight: 700 }}>
                CINEMA
              </Typography>
              <Button color="inherit" href="/">Movies</Button>
            </Toolbar>
          </AppBar>
          <Routes>
            <Route path="/" element={<Home movies={movies} />} />
            <Route path="/movie/:id" element={<MovieDetails movies={movies} />} />
            <Route path="/movie/:id/seats/:showTime" element={<SeatSelection movies={movies} />} />
          </Routes>
        </Box>
      </BrowserRouter>
    </ThemeProvider>
  )
}

export default App