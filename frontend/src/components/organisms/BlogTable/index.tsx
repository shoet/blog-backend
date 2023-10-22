import { Button } from '@/components/atoms/Button'
import { ApiContext, Blog } from '@/types/api'
import './style.module.css'
import { useNavigate } from 'react-router-dom'
import { deleteBlog } from '@/services/blogs/delete-blog'
import { parseCookie } from '@/utils/cookie'

type BlogTableProps = {
  blogs: Blog[]
  onClickDelete?: () => void
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
            <td>{blog.isPublic}</td>
            <td>
              <Button variant="secondary" onClick={() => onEdit(blog.id)}>
                Edit
              </Button>
            </td>
            <td>
              <Button variant="danger" onClick={() => onDelete(blog.id)}>
                Delete
              </Button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
