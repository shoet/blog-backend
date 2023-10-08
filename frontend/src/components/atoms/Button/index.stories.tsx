import { Meta } from '@storybook/react'
import { Button } from '.'

export default {
  title: 'atoms/Button',
  component: Button,
} as Meta<typeof Button>

export const Default = () => <Button>button</Button>
