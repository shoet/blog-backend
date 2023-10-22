import { AppTheme } from '@/types/style'
import { theme } from '@/themes'

// CSSプロパティ判定用
const ColorPropKeys = new Set(['color', 'background-color'])
const FontSizePropKeys = new Set(['font-size'])
const LetterSpacingPropKeys = new Set(['letter-spacing'])
const SpacePropKeys = new Set([
  'margin',
  'margin-top',
  'margin-right',
  'margin-bottom',
  'margin-left',
  'padding',
  'padding-top',
  'padding-right',
  'padding-bottom',
  'padding-left',
])
const LineHeightPropKeys = new Set(['line-height'])

// Themeキー判定用
type ColorThemeKeys = keyof typeof theme.colors
type FontSizeThemeKeys = keyof typeof theme.fontSizes
type LetterSpacingThemeKeys = keyof typeof theme.letterSpacings
type SpaceThemeKeys = keyof typeof theme.spaces
type LineHeightThemeKeys = keyof typeof theme.lineHeights

// 型推論用
export type Color = ColorThemeKeys | (string & {})
export type FontSize = FontSizeThemeKeys | (string & {})
export type LetterSpacing = LetterSpacingThemeKeys | (string & {})
export type Space = SpaceThemeKeys | (string & {})
export type LineHeight = LineHeightThemeKeys | (string & {})

const BREAKPOINT = {
  sm: '576px',
  md: '768px',
  lg: '992px',
  xl: '1200px',
}

export type ResponsiveProp<T> = {
  base?: T
  sm?: T
  md?: T
  lg?: T
  xl?: T
}
export type Responsive<T> = ResponsiveProp<T> | T

export function toResponsiveValue<T>(
  propKey: string,
  value: Responsive<T>,
  theme: AppTheme,
) {
  if (isResponsiveProp(value)) {
    const styles: string[] = []
    for (const [responsiveKey, responsiveValue] of Object.entries(value)) {
      if (responsiveKey === 'base') {
        styles.push(
          `${propKey}: ${toThemeValue(propKey, responsiveValue, theme)};`,
        )
      } else {
        const breakPoint = BREAKPOINT[responsiveKey as keyof typeof BREAKPOINT]
        const style = `${propKey}: ${toThemeValue(
          propKey,
          responsiveValue,
          theme,
        )};`
        const mediaQuery = `@media screen and (min-width: ${breakPoint}) { ${style} }`
        styles.push(mediaQuery)
      }
    }
    const ret = styles.join('\n')
    return ret
  }
  return `${propKey}: ${toThemeValue(propKey, value, theme)};`
}

function isResponsiveProp<T>(value: any): value is ResponsiveProp<T> {
  return (
    value !== undefined &&
    (value.base !== undefined ||
      value.sm !== undefined ||
      value.md !== undefined ||
      value.lg !== undefined ||
      value.xl !== undefined)
  )
}

export function toThemeValue(propKey: string, value: any, theme: AppTheme) {
  if (ColorPropKeys.has(propKey) && isColorThemeKey(value, theme)) {
    return theme.colors[value]
  } else if (
    FontSizePropKeys.has(propKey) &&
    isFontSizeThemeKey(value, theme)
  ) {
    return theme.fontSizes[value]
  } else if (
    LetterSpacingPropKeys.has(propKey) &&
    isLetterSpacingThemeKey(value, theme)
  ) {
    return theme.letterSpacings[value]
  } else if (SpacePropKeys.has(propKey) && isSpaceThemeKey(value, theme)) {
    return theme.spaces[value]
  } else if (
    LineHeightPropKeys.has(propKey) &&
    isLineHeightThemeKey(value, theme)
  ) {
    return theme.lineHeights[value]
  }
  return value
}

function isColorThemeKey(value: any, theme: AppTheme): value is ColorThemeKeys {
  return Object.keys(theme.colors).filter((k) => k == value).length > 0
}

function isFontSizeThemeKey(
  value: any,
  theme: AppTheme,
): value is FontSizeThemeKeys {
  return Object.keys(theme.fontSizes).filter((k) => k == value).length > 0
}

function isLetterSpacingThemeKey(
  value: any,
  theme: AppTheme,
): value is LetterSpacingThemeKeys {
  return Object.keys(theme.letterSpacings).filter((k) => k == value).length > 0
}

function isSpaceThemeKey(value: any, theme: AppTheme): value is SpaceThemeKeys {
  return Object.keys(theme.spaces).filter((k) => k == value).length > 0
}

function isLineHeightThemeKey(
  value: any,
  theme: AppTheme,
): value is LineHeightThemeKeys {
  return Object.keys(theme.lineHeights).filter((k) => k == value).length > 0
}
