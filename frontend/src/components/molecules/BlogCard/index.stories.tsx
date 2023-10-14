import { Meta } from '@storybook/react'
import { BlogCard } from '.'

export default {
  title: 'molecules/BlogCard',
  component: BlogCard,
} as Meta<typeof BlogCard>

export const Default = () => <BlogCard blogId={1} />
