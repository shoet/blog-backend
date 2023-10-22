import { BlogForm, BlogFormData } from '@/components/organisms/BlogForm'
import { useAuthGuard } from '@/services/auth/auth-guard'
import { addBlog } from '@/services/blogs/add-blog'
import { parseCookie } from '@/utils/cookie'
import { useNavigate } from 'react-router-dom'

export const BlogPostPage = () => {
  useAuthGuard()
  const navigate = useNavigate()

  const onSubmit = async (data: BlogFormData) => {
    const newBlog = {
      title: data.title,
      description: data.description,
      content: data.content,
      authorId: data.authorId,
      isPublic: data.isPublic,
      thumbnailImageFileName: data.thumbnailImageFileName ?? '',
      tags: data.tags,
    }
    const token = parseCookie(document.cookie)['authToken']
    await addBlog(
      {
        apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
      },
      { blog: newBlog },
      token,
    )
    navigate(`/`)
  }

  return <BlogForm onSubmit={onSubmit} />
}
