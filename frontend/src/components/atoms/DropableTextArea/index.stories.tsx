import DropableTextArea from '.'

import { Meta, StoryFn } from '@storybook/react'

export default {
  title: 'atoms/DropableTextArea',
  component: DropableTextArea,
  argTypes: {},
} as Meta<typeof DropableTextArea>

const Template: StoryFn<typeof DropableTextArea> = (args) => (
  <DropableTextArea {...args} />
)

const onDrop = (files: File[]) => {
  console.log(files[0].name)
}

export const Normal = Template.bind({})
Normal.args = {
  onDrop: onDrop,
}
