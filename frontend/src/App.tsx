import { ThemeProvider } from 'styled-components'
import { GlobalStyle } from './components/layout/GlobalStyle'
import { theme } from './themes'
import { Outlet } from 'react-router-dom'
import { SWRConfig } from 'swr'
import { fetcher } from './utils/fetcher'
import Layout from './components/templates/Layout'
import Content from './components/templates/Content'
import { SideContent } from './components/organisms/SideContent'

function App() {
  return (
    <>
      <GlobalStyle />
      <ThemeProvider theme={theme}>
        <SWRConfig value={{ fetcher }}>
          <Layout>
            <Content MainContent={<Outlet />} SubContent={<SideContent />} />
          </Layout>
        </SWRConfig>
      </ThemeProvider>
    </>
  )
}

export default App
