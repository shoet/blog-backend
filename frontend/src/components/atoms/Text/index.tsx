import {
  Color,
  FontSize,
  LetterSpacing,
  Responsive,
  toResponsiveValue,
} from '@/utils/style'
import styled from 'styled-components'

type TextVariant = 'small' | 'medium' | 'large'

const variants = {
  small: {
    fontSize: 'small',
    letterSpacing: 'small',
  },
  medium: {
    fontSize: 'medium',
    letterSpacing: 'small',
  },
  large: {
    fontSize: 'large',
    letterSpacing: 'medium',
  },
}

export const Text = styled.text.withConfig({
  shouldForwardProp: (prop) =>
    !['fontSize', 'fontWeight', 'letterSpacing'].includes(prop),
})<{
  variant?: TextVariant
  fontSize?: Responsive<FontSize>
  fontWeight?: Responsive<string>
  letterSpacing?: Responsive<LetterSpacing>
  color?: Responsive<Color>
}>`
  ${({ variant, fontSize, letterSpacing, theme }) => {
    const styles: string[] = []
    if (variant && variants[variant]) {
      const style = variants[variant]
      !fontSize &&
        styles.push(toResponsiveValue('font-size', style.fontSize, theme))
      !letterSpacing &&
        styles.push(
          toResponsiveValue('letter-spacing', style.letterSpacing, theme),
        )
    }
    return styles.join('\n')
  }}
  ${({ fontSize, theme }) => toResponsiveValue('font-size', fontSize, theme)}
  ${({ letterSpacing, theme }) =>
    toResponsiveValue('letter-spacing', letterSpacing, theme)}
  ${({ color, theme }) => toResponsiveValue('color', color, theme)}
  ${({ fontWeight, theme }) =>
    toResponsiveValue('font-weight', fontWeight, theme)}
  `

Text.defaultProps = {
  variant: 'medium',
  color: 'text',
}
