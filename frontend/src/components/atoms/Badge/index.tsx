import { Color, toResponsiveValue } from '@/utils/style'
import { PropsWithChildren } from 'react'
import { Text } from '../Text'
import styled from 'styled-components'

type BadgeProps = {
  backgroundColor: Color
  color: Color
  onClicn?: () => void
}

const Container = styled.div<{ backgroundColor: Color; color: Color }>`
  border-radius: 3px;
  display: inline-flex;
  padding: 2px 6px;
  ${({ backgroundColor, theme }) =>
    toResponsiveValue('background-color', backgroundColor, theme)}
  ${({ color, theme }) => toResponsiveValue('color', color, theme)}
`

export const Badge = (props: PropsWithChildren<BadgeProps>) => {
  const { backgroundColor, color, onClicn, children } = props

  const handleClick = () => {
    onClicn && onClicn()
  }

  return (
    <Container
      backgroundColor={backgroundColor}
      color={color}
      onClick={handleClick}
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
