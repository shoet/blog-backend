import Box from '@/components/layout/Box'
import { Profile } from '../Profile'
import { SearchForm } from '@/components/molecules/SearchForm'
import { BadgeContainer } from '../BadgeContainer'
import styled from 'styled-components'

const Container = styled.div`
  position: sticky;
  top: 40px;
`

export const SideContent = () => {
  return (
    <Container>
      <Box>
        <Profile />
      </Box>
      <Box marginTop={3}>
        <SearchForm />
      </Box>
      <Box marginTop={3}>
        <BadgeContainer />
      </Box>
    </Container>
  )
}
