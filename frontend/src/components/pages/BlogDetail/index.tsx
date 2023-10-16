import { Badge } from '@/components/atoms/Badge'
import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { useBlog } from '@/services/blogs/use-blog'
import { toStringYYYYMMDD_HHMMSS } from '@/utils/date'
import { marked } from 'marked'
import { useParams, redirect } from 'react-router-dom'
import styled from 'styled-components'

type BlogDetailPageParams = {
  id: string
}

const ImageWrapper = styled(Box)`
  img {
    display: block;
    object-fit: fit;
  }
`

const TagsWrapper = styled(Box)`
  div:not(:last-child) {
    margin-right: 0.5rem;
  }
`

export const BlogDetailPage = () => {
  const { id } = useParams<BlogDetailPageParams>()
  if (!id) {
    redirect('/404')
  }
  const { blog, isLoading } = useBlog(
    {
      apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
    },
    Number(id),
  )

  const options = {
    gfm: true,
    breaks: true,
    pedantic: false,
    smartLists: true,
    smartypants: true,
  }

  const renderMarkdown = (text: string) => {
    const __html = marked(text, options)
    return { __html }
  }

  return (
    <>
      {isLoading ?? <div>Loading...</div>}
      {blog && (
        <>
          <Box marginTop={2}>
            <Text fontSize="extraExtraLarge" fontWeight="bold">
              {blog.title}
            </Text>
          </Box>
          <Flex flexDirection="row" alignItems="center" marginTop={2}>
            <Box>
              <Text fontSize="medium" fontWeight="bold" color="gray">
                {toStringYYYYMMDD_HHMMSS(blog.created)}
              </Text>
            </Box>
            {blog.tags && (
              <TagsWrapper marginLeft={2}>
                {blog.tags.map((tag, idx) => (
                  <Badge key={idx}>{tag}</Badge>
                ))}
              </TagsWrapper>
            )}
          </Flex>
          <ImageWrapper>
            <img src={blog.thumbnailImageFileName} alt={blog.title} />
          </ImageWrapper>
          <Box marginTop={3}>
            <div dangerouslySetInnerHTML={renderMarkdown(blog.content)}></div>
            <div>{blog.content}</div>
          </Box>
        </>
      )}
    </>
  )
}
