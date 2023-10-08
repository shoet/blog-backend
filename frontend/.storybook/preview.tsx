import React from 'react'
import type { Preview } from '@storybook/react'
import { ThemeProvider } from 'styled-components'
import { theme } from '../src/themes'
import { withConsole } from '@storybook/addon-console'
import { BrowserRouter } from 'react-router-dom'
import { GlobalStyle } from '../src/components/layout/GlobalStyle'

const preview: Preview = {
  parameters: {
    actions: { argTypesRegex: '^on[A-Z].*' },
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/,
      },
    },
  },
  decorators: [
    (Story) => (
      <>
        <GlobalStyle />
        <ThemeProvider theme={theme}>
          <BrowserRouter>
            <Story />
          </BrowserRouter>
        </ThemeProvider>
      </>
    ),
    (storyFn, context) => withConsole()(storyFn)(context),
  ],
}

export default preview
