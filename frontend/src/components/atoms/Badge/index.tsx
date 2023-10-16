import { Color, Responsive, toResponsiveValue } from '@/utils/style'
import { PropsWithChildren } from 'react'
import { Text } from '../Text'
import styled from 'styled-components'

type BadgeProps = {
  backgroundColor: Responsive<Color>
  color: Responsive<Color>
  focusColor?: Responsive<Color>
  onClicn?: () => void
}

const Container = styled.span<{
  backgroundColor: Responsive<Color>
  color: Responsive<Color>
  focusColor?: Responsive<Color>
}>`
  border-radius: 3px;
  display: inline-flex;
  padding: 2px 6px;
  ${({ backgroundColor, theme }) =>
    toResponsiveValue('background-color', backgroundColor, theme)}
  ${({ color, theme }) => toResponsiveValue('color', color, theme)}
  ${({ focusColor, theme }) =>
    focusColor &&
    `
    cursor: pointer;
    transition: all 0.1s ease-in-out;
    &:hover,
    &:focus {
      ${toResponsiveValue('background-color', focusColor, theme)}
    }
  `}
`

export const Badge = (
  props: PropsWithChildren<React.HTMLAttributes<HTMLDivElement> & BadgeProps>,
) => {
  const { backgroundColor, color, onClicn, focusColor, children } = props

  const handleClick = () => {
    onClicn && onClicn()
  }

  return (
    <Container
      backgroundColor={backgroundColor}
      color={color}
      onClick={handleClick}
      focusColor={focusColor}
    >
      <Text fontSize="small" fontWeight="bold" color={color}>
        {children}
      </Text>
    </Container>
  )
}

Badge.defaultProps = {
  backgroundColor: 'black',
  color: 'white',
}
