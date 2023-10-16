import { Badge } from '@/components/atoms/Badge'
import { useTags } from '@/services/tags/use-tags'
import { useNavigate } from 'react-router-dom'
import styled from 'styled-components'

const Container = styled.div`
  width: 100%;
  flex-wrap: wrap;
`

const BadgeWrapper = styled.div`
  display: inline-flex;
  padding: 3px 3px;
  cursor: pointer;
`

export const BadgeContainer = () => {
  const navigate = useNavigate()

  const { tags } = useTags(
    {
      apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
    },
    [],
  )

  return (
    <Container>
      {tags &&
        tags.map((t) => (
          <BadgeWrapper
            key={t.id}
            onClick={() => navigate(`/search?tag=${t.name}`)}
          >
            <Badge
              backgroundColor="black"
              color="white"
              focusColor="focusPurple"
            >
              {t.name}
            </Badge>
          </BadgeWrapper>
        ))}
    </Container>
  )
}
