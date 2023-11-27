import { Button } from '@/components/atoms/Button'
import { ApiContext, Blog } from '@/types/api'
import './style.module.css'
import { useNavigate } from 'react-router-dom'
import { deleteBlog } from '@/services/blogs/delete-blog'
import { parseCookie } from '@/utils/cookie'
import { Badge } from '@/components/atoms/Badge'

type BlogTableProps = {
  blogs: Blog[]
  onClickDelete?: () => void
}

const IsPublicBadge = () => {
  return <Badge backgroundColor="primary">公開</Badge>
}

const IsNotPublicBadge = () => {
  return <Badge>非公開</Badge>
}

export const BlogTable = (props: BlogTableProps) => {
  const { blogs, onClickDelete } = props

  const navigate = useNavigate()

  const onEdit = (id: number) => {
    navigate(`/${id}/edit`)
  }

  const context: ApiContext = {
    apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
  }

  const token = parseCookie(document.cookie)['authToken']
  const onDelete = async (id: number) => {
    await deleteBlog(context, { blogId: id }, token)
    onClickDelete && onClickDelete()
  }

  return (
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Title</th>
          <th>Created</th>
          <th>IsPublic</th>
          <th>Edit</th>
          <th>Delete</th>
        </tr>
      </thead>
      <tbody>
        {blogs.map((blog, idx) => (
          <tr key={idx}>
            <td>{blog.id}</td>
            <td>{blog.title}</td>
            <td>{blog.created}</td>
            <td style={{ textAlign: 'center' }}>
              {blog.isPublic ? <IsPublicBadge /> : <IsNotPublicBadge />}
            </td>
            <td>
              <Button variant="secondary" onClick={() => onEdit(blog.id)}>
                Edit
              </Button>
            </td>
            <td>
              <Button
                backgroundColor="dangerSoft"
                onClick={() => onDelete(blog.id)}
              >
                Delete
              </Button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
