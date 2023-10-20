import { ThemeProvider } from 'styled-components'
import { GlobalStyle } from './components/layout/GlobalStyle'
import { theme } from './themes'
import { Outlet } from 'react-router-dom'
import { SWRConfig } from 'swr'
import { fetcher } from './utils/fetcher'
import { SideContent } from './components/organisms/SideContent'
import { BaseLayout } from './components/templates/BaseLayout'
import { VSplit } from './components/templates/VSplit'

function App() {
  return (
    <>
      <GlobalStyle />
      <ThemeProvider theme={theme}>
        <SWRConfig value={{ fetcher }}>
          <BaseLayout>
            <VSplit MainContent={<Outlet />} SubContent={<SideContent />} />
          </BaseLayout>
        </SWRConfig>
      </ThemeProvider>
    </>
  )
}

export default App
