import { Meta } from '@storybook/react'
import { Alert } from '.'

export default {
  title: 'atoms/Alert',
  component: Alert,
} as Meta<typeof Alert>

export const Default = () => <Alert text="Hello" />
