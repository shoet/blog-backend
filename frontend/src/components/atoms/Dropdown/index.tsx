import { css, styled } from 'styled-components'
import { Text } from '@/components/atoms/Text'
import { useDropdown } from './use-dropdown'

export type DropdownOption = {
  label: string
  value: string
}
export const DROPDOWN_DEFAULT_OPTION: DropdownOption = {
  label: '-',
  value: '0',
}

type DropdownProps = {
  value?: DropdownOption
  options: DropdownOption[]
  onChange?: (item: DropdownOption) => void
  isError?: boolean
}
const Dropdown = (props: DropdownProps) => {
  const {
    value = DROPDOWN_DEFAULT_OPTION,
    options,
    onChange,
    isError = false,
  } = props

  const {
    isOpen,
    selectedValue,
    dropdownRootRef,
    dropdownItemRefs,
    openDropdown,
    openDropdownByKeyboard,
    selectDropdownItem,
    selectDropdownItemByKeyboard,
  } = useDropdown(value)

  return (
    <DropdownRoot
      ref={dropdownRootRef}
      isError={isError}
      aria-haspopup="true"
      aria-expanded={isOpen}
      tabIndex={0}
      onClick={openDropdown}
      onKeyDown={openDropdownByKeyboard}
    >
      <DropdownContainer>
        <DropdownSelected>
          <Text variant="medium">{selectedValue.label}</Text>
        </DropdownSelected>
        <DropdownArrowIcon isOpen={isOpen} size={5} />
        {isOpen && (
          <DropdownOptions>
            <ul>
              {[DROPDOWN_DEFAULT_OPTION, ...options].map((option, idx) => (
                <DropdownItem
                  key={option.value}
                  ref={(el) => {
                    dropdownItemRefs.current[idx] = el
                  }}
                  tabIndex={0}
                  onClick={() => {
                    selectDropdownItem(option)
                    onChange && onChange(option)
                  }}
                  onKeyDown={(e) => {
                    selectDropdownItemByKeyboard(e, idx, option)
                    onChange && onChange(option)
                  }}
                >
                  {option.label}
                </DropdownItem>
              ))}
            </ul>
          </DropdownOptions>
        )}
      </DropdownContainer>
    </DropdownRoot>
  )
}

const DropdownRoot = styled.div.withConfig({
  shouldForwardProp: (props) => !['isError'].includes(props),
})<{ isError: boolean }>`
  ${({ theme, isError }) => {
    return css`
        border: 1px solid ${
          isError ? theme.colors.danger : theme.colors.border
        };
      `
  }}
  border-radius: 3px;
  cursor: pointer;
`

const DropdownContainer = styled.div`
  position: relative;
`

const DropdownSelected = styled.div`
  line-height: 30px;
  padding: 5px 10px;
`

const DropdownArrowIcon = styled.div.withConfig({
  shouldForwardProp: (props) => !['isOpen'].includes(props),
})<{ isOpen?: boolean; size?: number }>`
  position: absolute;
  right: 0px;
  top: 50%;
  width: 0;
  height: 0;
  transform: translateX(-100%) translateY(-50%);
  ${({ isOpen = false, size = 5 }) => {
    return isOpen
      ? css`
      border-left: ${size}px solid transparent;
      border-right: ${size}px solid transparent;
      border-bottom: ${size}px solid black;
    `
      : css`
      border-left: ${size}px solid transparent;
      border-right: ${size}px solid transparent;
      border-top: ${size}px solid black;
    `
  }}
`

const DropdownOptions = styled.div`
  position: absolute;
  top: 40px;
  width: 100%;
  max-height: 160px;
  overflow: auto;
  border-radius: 0px 0px 5px 5px;
  border: 1px solid transparent;
  box-shadow: rgba(149, 157, 165, 0.2) 0px 8px 24px;
  background-color: #FFFFFF
`

const DropdownItem = styled.li`
  line-height: 30px;
  padding: 5px 10px;
  &:hover {
    background-color: #CCCCCC;
  }
`

export default Dropdown
