import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { Color } from '@/utils/style'
import styled from 'styled-components'
import { IconXmark } from '../Icon'

const Container = styled.div<{ borderColor: string; backgroundColor: string }>`
  background-color: ${({ backgroundColor }) => backgroundColor};
  border: 1px solid ${({ borderColor }) => borderColor};
  color: ${({ borderColor }) => borderColor};
  border-radius: 5px;
  opacity: 0.8;
`

type AlertProps = {
  text?: string
  borderColor?: Color
  backgroundColor?: Color
  onClick?: () => void
}

export const Alert = (props: AlertProps) => {
  // TODO: icon size
  const {
    text,
    borderColor = '#ff4545',
    backgroundColor = '#ffc2c2',
    onClick,
  } = props
  return (
    <Container borderColor={borderColor} backgroundColor={backgroundColor}>
      <Flex
        flexDirection="row"
        justifyContent="space-between"
        alignItems="center"
      >
        <Box padding={2}>{text}</Box>
        <Box padding={2} onClick={onClick}>
          <IconXmark />
        </Box>
      </Flex>
    </Container>
  )
}
