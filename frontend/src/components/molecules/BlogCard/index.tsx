import { Badge } from '@/components/atoms/Badge'
import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { useBlog } from '@/services/blogs/use-blog'
import styled from 'styled-components'

type BlogCardProps = {
  blogId: number
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
  const { blogId } = props
  const { blog, isLoading } = useBlog(
    {
      apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
    },
    blogId,
  )

  return (
    <>
      <Container>
        {isLoading && <div>loading...</div>}
        {blog && (
          <>
            <Flex flexDirection="row" alignItems="start">
              <ImageWrapper>
                <img src={blog.thumbnailImageFileName} alt={blog.title} />
              </ImageWrapper>
              <ContentWrapper>
                <Text fontSize="extraExtraLarge" fontWeight="bold">
                  {blog.title}
                </Text>
                <TagsWrapper marginTop={1}>
                  {blog.tags.map((tag) => (
                    <Badge>{tag}</Badge>
                  ))}
                </TagsWrapper>
                <Box marginTop={2}>
                  <Text fontSize="medium">{blog.description}</Text>
                </Box>
                <Box marginTop={2}>
                  <Text fontSize="small" fontWeight="bold" color="gray">
                    {blog.created}
                  </Text>
                </Box>
              </ContentWrapper>
            </Flex>
          </>
        )}
      </Container>
    </>
  )
}
