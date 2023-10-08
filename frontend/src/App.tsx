import { ThemeProvider } from 'styled-components'
import { GlobalStyle } from './components/layout/GlobalStyle'
import { theme } from './themes'
import { Outlet } from 'react-router-dom'

function App() {
  return (
    <>
      <GlobalStyle />
      <ThemeProvider theme={theme}>
        <Outlet />
      </ThemeProvider>
    </>
  )
}

export default App
