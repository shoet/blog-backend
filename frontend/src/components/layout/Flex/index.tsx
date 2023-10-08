import { styled } from 'styled-components'
import Box from '../Box'
import {
  CSSPropertyAlignContent,
  CSSPropertyAlignItems,
  CSSPropertyAlignSelf,
  CSSPropertyFlexDirection,
  CSSPropertyFlexWrap,
  CSSPropertyJustifyContent,
  CSSPropertyJustifyItems,
  CSSPropertyJustifySelf,
} from '@/types/style'
import { Responsive, toResponsiveValue } from '@/utils/style'
import { theme } from '@/themes'

type FlexProps = {
  justifyContent?: Responsive<CSSPropertyJustifyContent>
  justifyItems?: Responsive<CSSPropertyJustifyItems>
  justifySelf?: Responsive<CSSPropertyJustifySelf>
  alignContent?: Responsive<CSSPropertyAlignContent>
  alignItems?: Responsive<CSSPropertyAlignItems>
  alignSelf?: Responsive<CSSPropertyAlignSelf>
  flexDirection?: Responsive<CSSPropertyFlexDirection>
  flexWrap?: Responsive<CSSPropertyFlexWrap>
}

const Flex = styled(Box).withConfig({
  shouldForwardProp: (props) =>
    ![
      'justifyContent',
      'justifyItems',
      'justifySelf',
      'alignContent',
      'alignItems',
      'alignSelf',
      'flexDirection',
      'flexWrap',
    ].includes(props),
})<FlexProps>`
  display: flex;
  ${(props) =>
    toResponsiveValue('justify-content', props.justifyContent, theme)}
  ${(props) => toResponsiveValue('justify-items', props.justifyItems, theme)}
  ${(props) => toResponsiveValue('justify-self', props.justifySelf, theme)}
  ${(props) => toResponsiveValue('align-content', props.alignContent, theme)}
  ${(props) => toResponsiveValue('align-items', props.alignItems, theme)}
  ${(props) => toResponsiveValue('align-self', props.alignSelf, theme)}
  ${(props) => toResponsiveValue('flex-direction', props.flexDirection, theme)}
  ${(props) => toResponsiveValue('flex-wrap', props.flexWrap, theme)}
`

export default Flex
