import { Meta, StoryFn } from '@storybook/react'
import Dropdown, { DropdownOption } from '.'

export default {
  title: 'atoms/Dropdown',
  component: Dropdown,
  argTypes: {
    xxx: {
      type: 'string',
    },
  },
} as Meta<typeof Dropdown>

const Template: StoryFn<typeof Dropdown> = (args) => <Dropdown {...args} />

export const Normal = Template.bind({})
Normal.args = {
  options: [
    { label: 'label_a', value: 'value_a' },
    { label: 'label_b', value: 'value_b' },
    { label: 'label_c', value: 'value_c' },
  ],
  onChange: (option: DropdownOption) => console.log(option),
}
