import { ThemeProvider } from 'styled-components'
import { GlobalStyle } from './components/layout/GlobalStyle'
import { theme } from './themes'
import { Outlet } from 'react-router-dom'
import { SWRConfig } from 'swr'
import { fetcher } from './utils/fetcher'

function App() {
  return (
    <>
      <GlobalStyle />
      <ThemeProvider theme={theme}>
        <SWRConfig value={{ fetcher }}>
          <Outlet />
        </SWRConfig>
      </ThemeProvider>
    </>
  )
}

export default App
