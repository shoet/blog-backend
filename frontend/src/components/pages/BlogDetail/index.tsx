import { Badge } from '@/components/atoms/Badge'
import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { useBlog } from '@/services/blogs/use-blog'
import { toStringYYYYMMDD_HHMMSS } from '@/utils/date'
import { Responsive, Space, toResponsiveValue } from '@/utils/style'
import { marked } from 'marked'
import { useParams, redirect } from 'react-router-dom'
import styled from 'styled-components'
import { MarkedOptions } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/monokai.css'
import { useEffect } from 'react'

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
  span:not(:last-child) {
    margin-right: 0.5rem;
  }
`

const BadgeWrapper = styled.span.withConfig({
  shouldForwardProp: (prop) => !['paddingTop'].includes(prop),
})<{ paddingTop?: Responsive<Space> }>`
  display: inline-block;
  ${({ paddingTop, theme }) =>
    paddingTop && toResponsiveValue('padding-top', paddingTop, theme)}
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

  useEffect(() => {
    if (blog) {
      hljs.highlightAll()
    }
  }, [blog])

  marked.setOptions({
    langPrefix: '',
    highlight: function (code: string, lang: string) {
      return hljs.highlightAuto(code, [lang]).value
    },
  } as MarkedOptions)

  const markedHtml = marked(blog?.content ?? '')

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
                  <BadgeWrapper key={idx} paddingTop={{ base: '5px', md: '0' }}>
                    <Badge>{tag}</Badge>
                  </BadgeWrapper>
                ))}
              </TagsWrapper>
            )}
          </Flex>
          <ImageWrapper marginTop={2}>
            <img src={blog.thumbnailImageFileName} alt={blog.title} />
          </ImageWrapper>
          <Box marginTop={3}>
            <span
              dangerouslySetInnerHTML={{
                __html: markedHtml,
              }}
            />
          </Box>
        </>
      )}
    </>
  )
}
