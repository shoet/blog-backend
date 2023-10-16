import { Meta } from '@storybook/react'
import { BlogForm } from '.'

export default {
  title: 'organisms/BlogForm',
  component: BlogForm,
} as Meta<typeof BlogForm>

export const Default = () => <BlogForm />
