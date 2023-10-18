import { BlogForm, BlogFormData } from '@/components/organisms/BlogForm'
import { addBlog } from '@/services/blogs/add-blog'
import { useNavigate } from 'react-router-dom'

export const BlogPostPage = () => {
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
    await addBlog(
      {
        apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
      },
      { blog: newBlog },
    )
    navigate(`/`)
  }

  return <BlogForm onSubmit={onSubmit} />
}
