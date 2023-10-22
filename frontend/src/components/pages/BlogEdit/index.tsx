import { BlogForm, BlogFormData } from '@/components/organisms/BlogForm'
import { useAuthGuard } from '@/services/auth/auth-guard'
import { putBlog } from '@/services/blogs/put-blog'
import { useBlog } from '@/services/blogs/use-blog'
import { parseCookie } from '@/utils/cookie'
import { redirect, useNavigate, useParams } from 'react-router-dom'

type BlogEditPageParams = {
  id: string
}

export const BlogEditPage = () => {
  useAuthGuard()
  const navigate = useNavigate()

  const { id } = useParams<BlogEditPageParams>()
  if (!id) {
    redirect('/404')
  }

  const { blog, isLoading } = useBlog(
    {
      apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
    },
    Number(id),
  )

  const onSubmit = async (data: BlogFormData) => {
    const newBlog = {
      id: data.id,
      title: data.title,
      description: data.description,
      content: data.content,
      authorId: data.authorId,
      isPublic: data.isPublic,
      thumbnailImageFileName: data.thumbnailImageFileName ?? '',
      tags: data.tags,
    }
    const token = parseCookie(document.cookie)['authToken']
    await putBlog(
      {
        apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
      },
      { blog: newBlog },
      token,
    )
    navigate(`/`)
  }

  return isLoading ? (
    <div>loading...</div>
  ) : (
    <BlogForm data={blog} onSubmit={onSubmit} />
  )
}
