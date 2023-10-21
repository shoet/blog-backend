import { Button } from '@/components/atoms/Button'
import { Blog } from '@/types/api'
import './style.module.css'
import { useNavigate } from 'react-router-dom'

type BlogTableProps = {
  blogs: Blog[]
}

export const BlogTable = (props: BlogTableProps) => {
  const { blogs } = props

  const navigate = useNavigate()

  const onEdit = (id: number) => {
    navigate(`/${id}/edit`)
  }
  const onDelete = (id: number) => {}

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
              <Button variant="danger">Delete</Button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
