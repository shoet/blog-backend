import { IconSearch } from '@/components/atoms/Icon'
import { Input } from '@/components/atoms/Input'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styled from 'styled-components'

const Container = styled(Box)`
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 5px;
`

const IconWrapper = styled(Box)`
  cursor: pointer;
`

export const SearchForm = () => {
  const navigate = useNavigate()
  const [searchText, setSearchText] = useState('')

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (searchText !== '') {
      navigate(`/search?keyword=${searchText}`)
    }
  }

  const handleClick = () => {
    if (searchText !== '') {
      navigate(`/search?keyword=${searchText}`)
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <Container>
        <Flex flexDirection="row" alignItems="center" padding="2px">
          <Input
            type="text"
            hasBorder={false}
            placeholder="Search"
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
          />
          <IconWrapper paddingTop="3px" onClick={handleClick}>
            <IconSearch size={24} />
          </IconWrapper>
        </Flex>
      </Container>
    </form>
  )
}
