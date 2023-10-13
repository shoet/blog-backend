import { Meta } from '@storybook/react'
import { Badge } from '.'

export default {
  title: 'atoms/Badge',
  component: Badge,
} as Meta<typeof Badge>

export const Default = () => <Badge backgroundColor="black">Docker</Badge>
