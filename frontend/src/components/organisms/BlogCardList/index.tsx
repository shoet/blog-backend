import { BlogCard } from '@/components/molecules/BlogCard'
import { useBlogList } from '@/services/blogs/use-blog-list'
import styled from 'styled-components'

const Container = styled.div`
  > div:not(:last-child) {
    margin-bottom: 1rem;
  }
`

export const BlogCardList = () => {
  const { blogs } = useBlogList(
    {
      apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
    },
    [],
  )

  return (
    <>
      <Container>
        {blogs && blogs.map((b) => <BlogCard key={b.id} blogId={b.id} />)}
      </Container>
    </>
  )
}
