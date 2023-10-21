import { Meta } from '@storybook/react'
import { LoginForm } from '.'

export default {
  title: 'organisms/LoginForm',
  component: LoginForm,
} as Meta<typeof LoginForm>

export const Default = () => <LoginForm />
