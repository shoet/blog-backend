import Flex from '@/components/layout/Flex'
import styled from 'styled-components'
import { PropsWithChildren } from 'react'
import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import { useDropzone } from './use-dropzone'

const ACCEPT_FILE_TYPES_DEFAULT = [
  'image/jpeg',
  'image/jpg',
  'image/png',
  'image/gif',
]

type DropzoneProps = {
  value?: File[]
  onChange?: (value: File[]) => void
  isError?: boolean
  acceptFileTypes?: string[]
}
const Dropzone = (props: PropsWithChildren<DropzoneProps>) => {
  const {
    onChange,
    isError: isErrorInit = false,
    acceptFileTypes = ACCEPT_FILE_TYPES_DEFAULT,
    children,
  } = props

  const {
    isError,
    errorMessages,
    openFileDialog,
    selectFile,
    dropFile,
    dragOver,
    inputRef,
  } = useDropzone({ isError: isErrorInit, onChange, acceptFileTypes })

  return (
    <>
      <DropzoneContainer
        isError={isError}
        onClick={openFileDialog}
        onDrop={dropFile}
        onDragOver={dragOver}
      >
        {children || 'ファイルをドロップまたはクリックして下さい'}
        <input
          type="file"
          multiple
          hidden
          accept={acceptFileTypes.join(',')}
          onChange={selectFile}
          ref={inputRef}
        />
      </DropzoneContainer>
      {isError &&
        errorMessages.length > 0 &&
        errorMessages.map((m, idx) => (
          <Box key={idx}>
            <Text as="label" variant="small" color="danger">
              {m}
            </Text>
          </Box>
        ))}
    </>
  )
}

const DropzoneContainer = styled(Flex)<{ isError?: boolean }>`
  border: dashed 1px ${({ isError, theme }) =>
    isError ? theme.colors.danger : theme.colors.border};
  border-radius: 10px;
  height: 200px;
  width: 100%;
  align-items: center;
  justify-content: center;
  cursor: pointer;
`

export default Dropzone
