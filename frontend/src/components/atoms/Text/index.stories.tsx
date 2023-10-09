import { Meta } from '@storybook/react'
import { Text } from '.'

export default {
  title: 'atoms/Text',
  component: Text,
} as Meta<typeof Text>

export const Small = () => <Text variant="small">Small</Text>
export const Medium = () => <Text variant="medium">Medium</Text>
export const Large = () => <Text variant="large">Large</Text>
