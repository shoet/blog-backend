import Box from '@/components/layout/Box'
import { BlogCard } from '@/components/molecules/BlogCard'
import { useBlogList } from '@/services/blogs/use-blog-list'
import { NavLink } from 'react-router-dom'
import styled from 'styled-components'

const Container = styled.div`
  > div:not(:last-child) {
    margin-bottom: 1rem;
  }
`

export const BlogCardList = () => {
  // TODO: フォームにタグを入力する際、Enterでタグの決定するが、「,」でもタグ決定とみなす。
  // TODO: backend側でタグをカンマ区切りでハンドリングするため
  const { blogs } = useBlogList(
    {
      apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
    },
    [],
  )

  return (
    <>
      <Container>
        {blogs &&
          blogs.map((b, idx) => (
            <Box key={idx}>
              <NavLink to={`/${b.id}`}>
                <BlogCard blog={b} />
              </NavLink>
            </Box>
          ))}
      </Container>
    </>
  )
}
