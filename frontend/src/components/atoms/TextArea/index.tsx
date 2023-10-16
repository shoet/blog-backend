import { useState } from 'react'
import { css, styled } from 'styled-components'

type TextAreaProps = {
  value?: string
  minRows?: number
  maxRows?: number
  onChange?: (text: string) => void
  isError?: boolean
  fontSize?: number
  lineHeight?: number
}
const TextArea = ({
  value,
  minRows = 2,
  maxRows = 10,
  isError = false,
  fontSize = 18,
  lineHeight = 20,
  onChange,
}: TextAreaProps) => {
  const [rows, setRows] = useState<number>(minRows)

  const handleOnChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const text = e.target.value
    e.target.rows = minRows
    const currentRows = Math.floor(e.target.scrollHeight / lineHeight)

    if (currentRows > maxRows) {
      e.target.rows = maxRows
      setRows(maxRows)
    } else {
      e.target.rows = currentRows
      setRows(currentRows)
    }
    onChange && onChange(text)
  }

  return (
    <TextAreaStyled
      value={value}
      isError={isError}
      rows={rows}
      onChange={handleOnChange}
      fontSize={fontSize}
      lineHeight={lineHeight}
    />
  )
}

const TextAreaStyled = styled.textarea<{
  isError: boolean
  fontSize: number
  lineHeight: number
}>`
    width: 100%;
    padding: 5px 5px;
    border: 1px solid ${({ isError, theme }) =>
      isError ? theme.colors.danger : theme.colors.border};
    border-radius: 3px;
    ${({ fontSize, lineHeight }) => css`
      line-height: ${lineHeight}px;
      font-size: ${fontSize}px;
    `}
  `

export default TextArea
