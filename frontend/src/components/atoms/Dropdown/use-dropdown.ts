import { useState, useRef, useEffect, useCallback } from 'react'
import { DropdownOption } from '.'

export const useDropdown = (initialValue: DropdownOption) => {
  const [isOpen, setIsOpen] = useState(false)
  const [selectedValue, setSelectedValue] =
    useState<DropdownOption>(initialValue)
  const dropdownRootRef = useRef<HTMLDivElement | null>(null)
  const dropdownItemRefs = useRef<(HTMLElement | null)[]>([])

  // 欄外クリック
  const handleClickOutside = useCallback((e: MouseEvent) => {
    if (!dropdownRootRef.current?.contains(e.target as Node)) {
      // DropdownRootに外側が含まれていない
      setIsOpen(false)
    }
  }, [])

  useEffect(() => {
    if (!dropdownRootRef.current) return
    document.addEventListener('click', (e) => handleClickOutside(e))
    return () => document.removeEventListener('click', handleClickOutside) // cleanup
  }, [dropdownRootRef, handleClickOutside])

  const openDropdown = () => {
    setIsOpen(!isOpen)
  }

  const selectDropdownItem = (item: DropdownOption) => {
    setSelectedValue(item)
    setIsOpen(false)
  }

  const openDropdownByKeyboard = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      setIsOpen(!isOpen)
    } else if (e.key === 'Escape') {
      setIsOpen(false)
    } else if (e.key === 'ArrowDown') {
      e.preventDefault()
      if (isOpen) {
        dropdownItemRefs.current[0]?.focus()
      }
    }
  }

  const selectDropdownItemByKeyboard = (
    e: React.KeyboardEvent,
    index: number,
    item: DropdownOption,
  ) => {
    e.stopPropagation() // 親(DropdownRoot)でのイベントを無視する
    if (e.key === 'Enter') {
      e.preventDefault()
      setSelectedValue(item)
      setIsOpen(!isOpen)
      dropdownRootRef.current?.focus()
    } else if (e.key === 'Escape') {
      setIsOpen(false)
    } else if (e.key === 'ArrowDown') {
      e.preventDefault()
      dropdownItemRefs.current[index + 1]?.focus()
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      dropdownItemRefs.current[index - 1]?.focus()
    }
  }

  return {
    isOpen,
    selectedValue,
    dropdownRootRef,
    dropdownItemRefs,
    openDropdown,
    openDropdownByKeyboard,
    selectDropdownItem,
    selectDropdownItemByKeyboard,
  }
}
