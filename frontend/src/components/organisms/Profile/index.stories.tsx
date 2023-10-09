import { Meta } from '@storybook/react'
import { Profile } from '.'

export default {
  title: 'organisms/Profile',
  component: Profile,
} as Meta<typeof Profile>

export const Default = () => <Profile />
