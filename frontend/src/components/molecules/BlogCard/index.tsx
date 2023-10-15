import { Badge } from '@/components/atoms/Badge'
import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { Blog } from '@/types/api'
import { toStringYYYYMMDD_HHMMSS } from '@/utils/date'
import styled from 'styled-components'

type BlogCardProps = {
  blog: Blog
}

const Container = styled.div`
  display: flex;
  overflow: hidden;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 5px;
  padding: 20px;
`

const ImageWrapper = styled(Box)`
  flex: 1;
  width: 400px;
  img {
    width: 100%;
    height: 100%;
    display: block;
    object-fit: cover;
  }
`

const ContentWrapper = styled(Box)`
  flex: 2;
  padding-left: 1rem;
`

const TagsWrapper = styled(Box)`
  div:not(:last-child) {
    margin-right: 0.5rem;
  }
`

export const BlogCard = (props: BlogCardProps) => {
  // TODO: anchor link
  const { blog } = props
  return (
    <>
      <Container>
        <Flex flexDirection="row" alignItems="start">
          <ImageWrapper>
            <img src={blog.thumbnailImageFileName} alt={blog.title} />
          </ImageWrapper>
          <ContentWrapper>
            <Text fontSize="extraExtraLarge" fontWeight="bold">
              {blog.title}
            </Text>
            {blog.tags && (
              <TagsWrapper marginTop={1}>
                {blog.tags.map((tag) => (
                  <Badge>{tag}</Badge>
                ))}
              </TagsWrapper>
            )}
            <Box marginTop={2}>
              <Text fontSize="medium">{blog.description}</Text>
            </Box>
            <Box marginTop={2}>
              <Text fontSize="small" fontWeight="bold" color="gray">
                {toStringYYYYMMDD_HHMMSS(blog.created)}
              </Text>
            </Box>
          </ContentWrapper>
        </Flex>
      </Container>
    </>
  )
}
