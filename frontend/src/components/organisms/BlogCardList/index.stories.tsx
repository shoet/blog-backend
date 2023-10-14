import { Meta } from '@storybook/react'
import { BlogCardList } from '.'

export default {
  title: 'organisms/BlogCardList',
  component: BlogCardList,
} as Meta<typeof BlogCardList>

export const Default = () => <BlogCardList />
