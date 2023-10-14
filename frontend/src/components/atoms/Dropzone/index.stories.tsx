import { Meta, StoryFn } from '@storybook/react'
import Dropzone from '.'

export default {
  title: 'atoms/Dropzone',
  component: Dropzone,
  argTypes: {},
} as Meta<typeof Dropzone>

const Template: StoryFn<typeof Dropzone> = (args) => <Dropzone {...args} />

export const Normal = Template.bind({})
Normal.args = {}
