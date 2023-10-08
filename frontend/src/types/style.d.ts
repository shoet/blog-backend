import { theme } from '@/themes'

export type AppTheme = typeof theme

declare module 'styled-components' {
  export interface DefaultTheme extends AppTheme {}
}

// CSSProperties
// https://developer.mozilla.org/ja/docs/Web/CSS
// > Propertyes

type CSSPropertyGlobal = 'inherit' | 'initial' | 'unset' | 'revert'

type CSSPropertyPosition =
  | 'start'
  | 'end'
  | 'right'
  | 'center'
  | 'left'
  | 'flex-start'
  | 'flex-end'
  | 'normal'
  | 'streatch'
  | 'safe center'
  | 'unsafe center'

type CSSPropertyContentPosition =
  | 'start'
  | 'end'
  | 'center'
  | 'flex-start'
  | 'flex-end'

// justify-*
export type CSSPropertyJustifySelf =
  | CSSPropertyGlobal
  | CSSPropertyPosition
  | 'auto'
  | 'baseline'
  | 'first baseline'
  | 'last baseline'
  | 'self-end'
  | 'self-start'
  | (string & {})

export type CSSPropertyJustifyItems =
  | CSSPropertyGlobal
  | CSSPropertyPosition
  | CSSPropertyJustifySelf
  | 'legacy center'
  | 'legacy left'
  | 'legacy right'
  | (string & {})

export type CSSPropertyJustifyContent =
  | CSSPropertyGlobal
  | CSSPropertyPosition
  | 'space-between'
  | 'space-evenly'
  | 'space-around'
  | (string & {})

// align-*
export type CSSPropertyAlignSelf =
  | CSSPropertyGlobal
  | CSSPropertyPosition
  | 'auto'
  | 'baseline'
  | 'first baseline'
  | 'last baseline'
  | 'self-end'
  | 'self-start'

export type CSSPropertyAlignItems =
  | CSSPropertyGlobal
  | CSSPropertyPosition
  | 'baseline'
  | 'first baseline'
  | 'last baseline'

export type CSSPropertyAlignContent =
  | CSSPropertyGlobal
  | CSSPropertyContentPosition
  | 'space-between'
  | 'space-evenly'
  | 'space-around'
  | 'baseline'
  | 'first baseline'
  | 'last baseline'
  | 'normal'
  | 'streatch'
  | 'safe center'
  | 'unsafe center'

// flex-*
export type CSSPropertyFlexDirection =
  | CSSPropertyGlobal
  | 'row'
  | 'row-reverse'
  | 'column'
  | 'column-reverse'

export type CSSPropertyFlexWrap =
  | CSSPropertyGlobal
  | 'nowrap'
  | 'wrap'
  | 'wrap-reverse'
