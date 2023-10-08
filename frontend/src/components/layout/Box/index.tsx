import { Color, Responsive, Space, toResponsiveValue } from '@/utils/style'
import { styled } from 'styled-components'
import { theme } from '@/themes'

type BoxProps = {
  color?: Responsive<Color>
  backgroundColor?: Responsive<Color>
  margin?: Responsive<Space>
  marginTop?: Responsive<Space>
  marginRight?: Responsive<Space>
  marginBottom?: Responsive<Space>
  marginLeft?: Responsive<Space>
  padding?: Responsive<Space>
  paddingTop?: Responsive<Space>
  paddingRight?: Responsive<Space>
  paddingBottom?: Responsive<Space>
  paddingLeft?: Responsive<Space>
  width?: Responsive<string>
  maxWidth?: Responsive<string>
  minWidth?: Responsive<string>
  height?: Responsive<string>
  maxHeight?: Responsive<string>
  minHeight?: Responsive<string>
  flexGrow?: Responsive<number>
  position?: Responsive<string>
  top?: Responsive<string>
  right?: Responsive<string>
  bottom?: Responsive<string>
  left?: Responsive<string>
  textAlign?: Responsive<string>
  display?: Responsive<string>
}

const Box = styled.div.withConfig({
  shouldForwardProp: (props) =>
    ![
      'backgroundColor',
      'marginTop',
      'marginRight',
      'marginBottom',
      'marginLeft',
      'paddingTop',
      'paddingRight',
      'paddingBottom',
      'paddingLeft',
      'maxWidth',
      'minWidth',
      'maxHeight',
      'minHeight',
      'flexGrow',
    ].includes(props),
})<BoxProps>`

  ${(props) => toResponsiveValue('color', props.color, theme)}
  ${(props) =>
    toResponsiveValue('background-color', props.backgroundColor, theme)}
  ${(props) => toResponsiveValue('margin', props.margin, theme)}
  ${(props) => toResponsiveValue('margin-top', props.marginTop, theme)}
  ${(props) => toResponsiveValue('margin-right', props.marginRight, theme)}
  ${(props) => toResponsiveValue('margin-bottom', props.marginBottom, theme)}
  ${(props) => toResponsiveValue('margin-left', props.marginLeft, theme)}
  ${(props) => toResponsiveValue('padding', props.padding, theme)}
  ${(props) => toResponsiveValue('padding-top', props.paddingTop, theme)}
  ${(props) => toResponsiveValue('padding-right', props.paddingRight, theme)}
  ${(props) => toResponsiveValue('padding-bottom', props.paddingBottom, theme)}
  ${(props) => toResponsiveValue('padding-left', props.paddingLeft, theme)}
  ${(props) => toResponsiveValue('width', props.width, theme)}
  ${(props) => toResponsiveValue('max-width', props.maxWidth, theme)}
  ${(props) => toResponsiveValue('min-width', props.minWidth, theme)}
  ${(props) => toResponsiveValue('height', props.height, theme)}
  ${(props) => toResponsiveValue('max-height', props.maxHeight, theme)}
  ${(props) => toResponsiveValue('min-height', props.minHeight, theme)}
  ${(props) => toResponsiveValue('flex-grow', props.flexGrow, theme)}
  ${(props) => toResponsiveValue('position', props.position, theme)}
  ${(props) => toResponsiveValue('top', props.top, theme)}
  ${(props) => toResponsiveValue('right', props.right, theme)}
  ${(props) => toResponsiveValue('bottom', props.bottom, theme)}
  ${(props) => toResponsiveValue('left', props.left, theme)}
  ${(props) => toResponsiveValue('text-align', props.textAlign, theme)}
  ${(props) => toResponsiveValue('display', props.display, theme)}
`

export default Box
