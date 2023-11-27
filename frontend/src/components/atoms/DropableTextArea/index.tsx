import { useState } from 'react'
import { css, styled } from 'styled-components'
import Box from '@/components/layout/Box'
import { useDropzone } from '../Dropzone/use-dropzone'

type DropableTextAreaProps = {
  value?: string
  minRows?: number
  maxRows?: number
  onChange?: (text: string) => void
  isError?: boolean
  fontSize?: number
  lineHeight?: number
  onDrop?: (files: File[]) => void
}
const DropableTextArea = ({
  value,
  minRows = 2,
  maxRows = 10,
  isError = false,
  fontSize = 18,
  lineHeight = 20,
  onChange,
  onDrop,
}: DropableTextAreaProps) => {
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

  const {
    isError: isErrorDrop,
    dragOver,
    dropFile,
  } = useDropzone({
    onChange: onDrop,
    acceptFileTypes: ['image/png', 'image/jpeg', 'image/jpg'],
  })

  return (
    <Box padding="5px 5px" onDrop={dropFile} onDragOver={dragOver}>
      <TextAreaStyled
        value={value}
        isError={isError || isErrorDrop}
        rows={rows}
        onChange={handleOnChange}
        fontSize={fontSize}
        lineHeight={lineHeight}
      />
    </Box>
  )
}

const TextAreaStyled = styled.textarea.withConfig({
  shouldForwardProp: (prop) =>
    !['lineHeight', 'isError', 'fontSize'].includes(prop),
})<{
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

export default DropableTextArea
