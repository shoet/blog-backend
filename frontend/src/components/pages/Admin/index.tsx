import { Button } from '@/components/atoms/Button'
import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { BlogTable } from '@/components/organisms/BlogTable'
import { useAuthGuard } from '@/services/auth/auth-guard'
import { useBlogList } from '@/services/blogs/use-blog-list'
import { ApiContext } from '@/types/api'
import { NavLink } from 'react-router-dom'

const AdminPage = () => {
  useAuthGuard()

  const apiContext: ApiContext = {
    apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
  }
  const {
    blogs,
    isLoading,
    error,
    mutate: refreshBlogList,
  } = useBlogList(apiContext, [])

  const handleOnClickDelete = () => {
    refreshBlogList()
  }

  return (
    <Box>
      <Box>
        <Text fontSize="display" fontWeight="bold">
          Admin
        </Text>
      </Box>
      <Box>
        <Flex
          flexDirection="row"
          justifyContent="space-between"
          alignItems="center"
        >
          <Text fontSize="extraLarge" fontWeight="bold">
            記事一覧
          </Text>
          <NavLink to="/new">
            <Button variant="primary">New</Button>
          </NavLink>
        </Flex>
        {isLoading && <Text>loading...</Text>}
        {error && <Text>{error.message}</Text>}
        {blogs && blogs.length !== 0 && (
          <Text>
            <BlogTable blogs={blogs} onClickDelete={handleOnClickDelete} />
          </Text>
        )}
      </Box>
    </Box>
  )
}

export default AdminPage
