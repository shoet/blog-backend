import { Meta } from '@storybook/react'
import { Input } from '.'

export default {
  title: 'atoms/Input',
  component: Input,
} as Meta<typeof Input>

export const Border = () => <Input placeholder="Border" />
export const NoBorder = () => <Input hasBorder={false} />
export const HasError = () => <Input hasError={true} />
