import { Button } from '@/components/atoms/Button'
import { ApiContext, Blog } from '@/types/api'
import { NavLink, useNavigate } from 'react-router-dom'
import { deleteBlog } from '@/services/blogs/delete-blog'
import { parseCookie } from '@/utils/cookie'
import { IsPublicBadge, IsNotPublicBadge } from '@/components/atoms/Badge'
import styled from 'styled-components'
import { toResponsiveValue } from '@/utils/style'

type BlogTableProps = {
  blogs: Blog[]
  onClickDelete?: () => void
}

const TitleLink = styled(NavLink)`
  text-decoration: underline;
  color: ${({ theme }) =>
    toResponsiveValue('color', theme.colors.primaryDark, theme)};
`

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

  const TableWrapper = styled.div`
    th,
    td {
      border: solid 1px;
      padding: 10px;
    }

    table {
      border-collapse: collapse;
    }
  `

  return (
    <TableWrapper>
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
              <td>
                <TitleLink to={`/${blog.id}`} target="_blank">
                  {blog.title}
                </TitleLink>
              </td>
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
    </TableWrapper>
  )
}
