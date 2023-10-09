import { IconSearch } from '@/components/atoms/Icon'
import { Input } from '@/components/atoms/Input'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import styled from 'styled-components'

const Container = styled(Box)`
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 5px;
`

export const SearchForm = () => {
  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    console.log('onSubmit')
  }

  return (
    <form onSubmit={onSubmit}>
      <Container>
        <Flex flexDirection="row" alignItems="center" padding="2px">
          <Input hasBorder={false} placeholder="Search" />
          <IconSearch size={24} />
        </Flex>
      </Container>
    </form>
  )
}
