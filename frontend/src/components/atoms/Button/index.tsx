import { Color, FontSize, Responsive, toResponsiveValue } from '@/utils/style'
import styled from 'styled-components'

type ButtonVariant = 'primary' | 'secondary' | 'danger'

const variants = {
  primary: {
    color: 'white',
    backgroundColor: 'primary',
    pseudo: {
      hover: {
        backgroundColor: 'primaryDark',
      },
    },
  },
  secondary: {
    color: '',
    backgroundColor: '',
    pseudo: {
      hover: {
        backgroundColor: '',
      },
    },
  },
  danger: {
    color: 'white',
    backgroundColor: 'danger',
    pseudo: {
      hover: {
        backgroundColor: 'dangerDark',
      },
    },
  },
}

export const Button = styled.button<{
  variant?: ButtonVariant
  color?: Responsive<Color>
  backgroundColor?: Responsive<Color>
  pseudo?: {
    hover?: {
      backgroundColor?: Responsive<Color>
    }
  }
  fontSize?: Responsive<FontSize>
}>`
  ${({ variant, color, backgroundColor, pseudo, theme }) => {
    const styles: string[] = []
    if (variant && variants[variant]) {
      const style = variants[variant]
      !color && styles.push(toResponsiveValue('color', style.color, theme))
      !backgroundColor &&
        styles.push(
          toResponsiveValue('background-color', style.backgroundColor, theme),
        )
      !pseudo &&
        styles.push(
          `&:hover { ${toResponsiveValue(
            'background-color',
            style.pseudo.hover.backgroundColor,
            theme,
          )} }`,
        )
    }
    return styles.join('\n')
  }}
  ${({ fontSize, theme }) => toResponsiveValue('font-size', fontSize, theme)}
  padding: 5px 10px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
`

Button.defaultProps = {
  variant: 'primary',
  fontSize: 'extraLarge',
}
