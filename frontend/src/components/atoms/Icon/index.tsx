import React from 'react'
import styled from 'styled-components'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {
  faGithub,
  faYoutube,
  faTwitter,
} from '@fortawesome/free-brands-svg-icons'
import { faMagnifyingGlass, faXmark } from '@fortawesome/free-solid-svg-icons'
import { Color, Responsive, toResponsiveValue } from '@/utils/style'

type IconProps = {
  size?: number
  href?: string
  focusColor?: Responsive<Color>
}

const IconStyle = styled.div.withConfig({
  shouldForwardProp: (prop) => !['focusColor'].includes(prop),
})<IconProps>`
  display: 'inline-block';
  width: ${({ size }) => `${size}px`};
  height: ${({ size }) => `${size}px`};
  cursor: pointer;
  transition: all 0.1s ease-in-out;
  ${({ focusColor, theme }) =>
    focusColor &&
    `
    &:hover,
    &:focus {
      ${toResponsiveValue('color', focusColor, theme)}
    }
  `}
`

const withIconStyle = (Icon: React.ReactNode) => {
  return (props: IconProps) => {
    return (
      <IconStyle size={props.size} focusColor={props.focusColor}>
        {Icon}
      </IconStyle>
    )
  }
}

export const IconGitHub = withIconStyle(<FontAwesomeIcon icon={faGithub} />)
export const IconYoutube = withIconStyle(<FontAwesomeIcon icon={faYoutube} />)
export const IconTwitter = withIconStyle(<FontAwesomeIcon icon={faTwitter} />)
export const IconSearch = withIconStyle(
  <FontAwesomeIcon icon={faMagnifyingGlass} />,
)
export const IconXmark = withIconStyle(<FontAwesomeIcon icon={faXmark} />)
