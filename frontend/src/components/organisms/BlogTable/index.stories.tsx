import { Meta } from '@storybook/react'
import { BlogTable } from '.'
import blogsData from './blogs.json'
import { Blog } from '@/types/api'

export default {
  title: 'organisms/BlogTable',
  component: BlogTable,
} as Meta<typeof BlogTable>

const blogs = blogsData as Blog[]

export const Default = () => <BlogTable blogs={blogs} />
