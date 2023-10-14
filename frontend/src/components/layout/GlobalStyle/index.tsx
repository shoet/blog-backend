import { createGlobalStyle } from 'styled-components'

export const GlobalStyle = createGlobalStyle`
  html,
  body,
  textarea {
    font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, Oxygen, Ubuntu, Cantarell, Fira Sans, Droid Sans, Helvetica Neue, sans-serif;
    padding: 0;
    margin: 0;
  }
  * {
    box-sizing: border-box;
  }
  a {
    cursor: pointer;
    text-decoration: none;
    color: #000;
  }

  ol, ul {
    margin: 0;
    padding: 0;
    list-style: none;
  }
`
