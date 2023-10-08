import Box from '@/components/layout/Box'
import { useParams } from 'react-router-dom'

type BlogDetailPageParams = {
  id: string
}

export const BlogDetailPage = () => {
  const { id } = useParams<BlogDetailPageParams>()
  return (
    <Box>
      <Box>BlogDetailPage</Box>
      <Box>{id}</Box>
    </Box>
  )
}
