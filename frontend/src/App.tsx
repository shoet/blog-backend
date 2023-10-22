import { ThemeProvider } from 'styled-components'
import { GlobalStyle } from '@/components/layout/GlobalStyle'
import { theme } from '@/themes'
import { Outlet } from 'react-router-dom'
import { SWRConfig } from 'swr'
import { fetcher } from '@/utils/fetcher'
import { ApiContext } from '@/types/api'
import { AuthContextProvider } from '@/contexts/AuthContext'

function App() {
  const context: ApiContext = {
    apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
  }
  return (
    <>
      <GlobalStyle />
      <ThemeProvider theme={theme}>
        <SWRConfig value={{ fetcher, shouldRetryOnError: false }}>
          <AuthContextProvider context={context}>
            <Outlet />
          </AuthContextProvider>
        </SWRConfig>
      </ThemeProvider>
    </>
  )
}

export default App
