import { Badge } from '@/components/atoms/Badge'
import Box from '@/components/layout/Box'
import { useTags } from '@/services/tags/use-tags'
import styled from 'styled-components'

const Container = styled.div`
  width: 100%;
  flex-wrap: wrap;
`

export const BadgeContainer = () => {
  // TODO: anchor
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
          <Box key={t.id} display="inline-flex" padding="3px 3px">
            <Badge backgroundColor="black" color="white">
              {t.name}
            </Badge>
          </Box>
        ))}
    </Container>
  )
}
