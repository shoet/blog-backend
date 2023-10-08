import { Meta } from '@storybook/react'
import { SearchForm } from '.'

export default {
  title: 'molecules/SearchForm',
  component: SearchForm,
} as Meta<typeof SearchForm>

export const Default = () => <SearchForm />
