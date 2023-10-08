import { ThemeProvider } from 'styled-components'
import { GlobalStyle } from './components/layout/GlobalStyle'
import { ExamplePage } from './components/pages/exmaple'
import { theme } from './themes'

function App() {
  return (
    <>
      <GlobalStyle />
      <ThemeProvider theme={theme}>
        <ExamplePage />
      </ThemeProvider>
    </>
  )
}

export default App
