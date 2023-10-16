import TextArea from '.'

import { Meta, StoryFn } from '@storybook/react'

export default {
  title: 'atoms/TextArea',
  component: TextArea,
  argTypes: {},
} as Meta<typeof TextArea>

const Template: StoryFn<typeof TextArea> = (args) => <TextArea {...args} />

export const Normal = Template.bind({})
Normal.args = {}
