import styled from 'styled-components'

export const Input = styled.input.withConfig({
  shouldForwardProp: (prop) => !['hasBorder', 'hasError'].includes(prop),
})<{ hasBorder?: boolean; hasError?: boolean }>`
  ${({ hasBorder, hasError, theme }) => {
    if (hasBorder) {
      return `
        border: 1px solid ${
          hasError ? theme.colors.danger : theme.colors.border
        };
        border-radius: 5px;
      `
    } else {
      return `
        border: none;
      `
    }
  }}

  width: 100%;
  padding: 5px 10px;
  height: 36px;
  font-size: 18px;
  outline: none;
`

Input.defaultProps = {
  hasBorder: true,
  hasError: false,
}
