import { Badge, IsNotPublicBadge } from '@/components/atoms/Badge'
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
    width: 100%;
    display: block;
    object-fit: fit;
  }
`

const TagsWrapper = styled(Box)`
  span:not(:last-child) {
    margin-right: 0.5rem;
  }
`

function addLinkTargetBlank(html: string): string {
  const regex = /<a href="(.*?)"/g
  const replacer = (_: string, p1: string) => `<a href="${p1}" target="_blank"`
  return html.replace(regex, replacer)
}

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
  const markedHtmlWithTargetBlank = addLinkTargetBlank(markedHtml)

  const MarkedWrapper = styled(Box)`
    li {
      list-style: inside;
    }

    a {
      text-decoration: underline;
      color: #b99976;
      target: _blank;
    }

    h2 {
      border-bottom: 1px solid #ccc;
      border-left: 5px solid #ccc;
      padding-left: 10px;
    }

    code:not([class]) {
      margin: 0 2px;
      padding: 0 5px;
      border-radius: 3px;
      background-color: #de6fa1;
      color: white;
    }

    pre > code {
      border-radius: 5px;
    }

    p {
      font-size: 0.9rem;
    }

    h2 {
      margin-top: 2rem;
    }
  `

  return (
    <>
      {isLoading ?? <div>Loading...</div>}
      {blog && (
        <>
          <Box marginTop={2}>
            <Flex alignItems="center">
              {!blog.isPublic && (
                <Box marginRight={1}>
                  <IsNotPublicBadge />
                </Box>
              )}
              <Text fontSize="extraExtraLarge" fontWeight="bold">
                {blog.title}
              </Text>
            </Flex>
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
          <MarkedWrapper marginTop={3}>
            <span
              dangerouslySetInnerHTML={{
                __html: markedHtmlWithTargetBlank,
              }}
            />
          </MarkedWrapper>
        </>
      )}
    </>
  )
}
