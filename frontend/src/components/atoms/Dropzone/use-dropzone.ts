import { useRef, useState } from 'react'

type useDropzoneProps = {
  acceptFileTypes: string[]
  isError?: boolean
  onChange?: (files: File[]) => void
}
export const useDropzone = ({
  isError: isErrorInit = false,
  onChange,
  acceptFileTypes,
}: useDropzoneProps) => {
  const [isError, setIsError] = useState(isErrorInit)
  const [errorMessages, setErrorMessages] = useState<string[]>([])

  function checkFileTypes(files: File[] | FileList, acceptFileTypes: string[]) {
    const validFiles = []
    const inValidFiles = []

    for (let i = 0, numFiles = files.length; i < numFiles; i++) {
      if (!acceptFileTypes.includes(files[i].type)) {
        inValidFiles.push(files[i])
      } else {
        validFiles.push(files[i])
      }
    }
    return { valid: validFiles, invalid: inValidFiles }
  }

  const fileTypeErrorMessage = (file: File) =>
    `無効なファイル形式です。[${file.name}: ${file.type}]`

  const selectFile = (e: React.ChangeEvent<HTMLInputElement>) => {
    setErrorMessages([])
    if (e.target.files) {
      const files = e.target.files
      const { valid, invalid } = checkFileTypes(files, acceptFileTypes)
      if (invalid.length > 0) {
        setIsError(true)
        setErrorMessages(invalid.map((f) => fileTypeErrorMessage(f)))
      }
      onChange && onChange(valid)
    }
  }

  const dropFile = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setErrorMessages([])
    if (e.dataTransfer.files) {
      const files = e.dataTransfer.files
      const { valid, invalid } = checkFileTypes(files, acceptFileTypes)
      if (invalid.length > 0) {
        setIsError(true)
        setErrorMessages(invalid.map((f) => fileTypeErrorMessage(f)))
      }
      onChange && onChange(valid)
    }
  }

  const dragOver = (e: React.DragEvent) => {
    e.preventDefault()
  }

  const inputRef = useRef<HTMLInputElement>(null)

  const openFileDialog = () => {
    if (inputRef.current) {
      inputRef.current.click()
    }
  }

  return {
    isError,
    errorMessages,
    selectFile,
    dropFile,
    dragOver,
    inputRef,
    openFileDialog,
  }
}
