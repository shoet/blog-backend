import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import { BlogTable } from '@/components/organisms/BlogTable'
import { useBlogList } from '@/services/blogs/use-blog-list'
import { ApiContext } from '@/types/api'

const AdminPage = () => {
  const apiContext: ApiContext = {
    apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
  }
  const { blogs, isLoading, error } = useBlogList(apiContext, [])

  return (
    <Box>
      <Box>
        <Text fontSize="display" fontWeight="bold">
          Admin
        </Text>
      </Box>
      <Box>
        <Text fontSize="extraLarge" fontWeight="bold">
          記事一覧
        </Text>
        {isLoading && <Text>loading...</Text>}
        {error && <Text>{error.message}</Text>}
        {blogs && blogs.length !== 0 && (
          <Text>
            <BlogTable blogs={blogs} />
          </Text>
        )}
      </Box>
    </Box>
  )
}

export default AdminPage
