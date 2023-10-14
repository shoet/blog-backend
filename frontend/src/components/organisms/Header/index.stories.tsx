import { Meta } from '@storybook/react'
import { Header } from '.'

export default {
  title: 'organisms/Header',
  component: Header,
} as Meta<typeof Header>

export const Default = () => <Header />
