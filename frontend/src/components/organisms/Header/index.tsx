import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { Color, Responsive, toResponsiveValue } from '@/utils/style'
import { NavLink } from 'react-router-dom'
import styled from 'styled-components'

const Div = styled.div.withConfig({
  shouldForwardProp: (prop) => !['space'].includes(prop),
})<{ space: string }>`
  width: 1px;
  border: 1px solid gray;
  height: 20px;
  ${({ space }) => space && `margin: 0px ${space};`}
`

const Divider = ({ space = 0 }: { space?: number }) => {
  const half = space / 2
  return <Div space={`${half}px`} />
}

const NavItem = styled.div.withConfig({
  shouldForwardProp: (prop) => !['focusColor'].includes(prop),
})<{
  focusColor?: Responsive<Color>
}>`
  ${({ focusColor, theme }) =>
    focusColor &&
    `
    cursor: pointer;
    transition: all 0.1s ease-in-out;
    &:hover,
    &:focus {
      ${toResponsiveValue('color', focusColor, theme)}
    }
  `}
`

export const Header = () => {
  return (
    <nav>
      <Flex
        flexDirection={{ base: 'column', md: 'row' }}
        justifyContent="space-between"
      >
        <Flex flexDirection="row" alignItems="baseline">
          <NavLink to="/">
            <Box>
              <Box display="inline-flex">
                <Text
                  fontSize="display"
                  fontWeight="bold"
                  letterSpacing="large"
                >
                  shoet
                </Text>
              </Box>
              <Box display="inline-flex" marginLeft={1}>
                <Text fontSize="display" letterSpacing="large">
                  Blog
                </Text>
              </Box>
            </Box>
          </NavLink>
          <Box marginLeft={2}>
            <Text fontSize="small" color="gray">
              技術や好きなことについて発信しています。
            </Text>
          </Box>
        </Flex>
        <Flex
          marginTop={{ base: 2, md: 0 }}
          marginRight={{ base: 0, md: 4 }}
          flexDirection="row"
          alignItems="center"
          justifyContent="space-evenly"
        >
          <Box>
            <NavItem focusColor="placeholder">
              <Text variant="large" color="inherit">
                Blog
              </Text>
            </NavItem>
          </Box>
          <Divider space={20} />
          <Box>
            <NavItem focusColor="placeholder">
              <Text variant="large" color="inherit">
                Portfolio
              </Text>
            </NavItem>
          </Box>
          <Divider space={20} />
          <Box>
            <NavItem focusColor="placeholder">
              <Text variant="large" color="inherit">
                About
              </Text>
            </NavItem>
          </Box>
        </Flex>
      </Flex>
    </nav>
  )
}
