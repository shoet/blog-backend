import React from 'react'
import styled from 'styled-components'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {
  faGithub,
  faYoutube,
  faTwitter,
} from '@fortawesome/free-brands-svg-icons'
import { Link } from 'react-router-dom'

type IconProps = {
  size?: number
  href?: string
}

const IconStyle = styled.div<IconProps>`
  display: 'inline-block';
  width: ${({ size }) => `${size}px`};
  height: ${({ size }) => `${size}px`};
`

const withIconStyle = (Icon: React.ReactNode) => {
  return (props: IconProps) => {
    return (
      <IconStyle size={props.size}>
        <Link to={props.href ?? ''} target="_blank">
          {Icon}
        </Link>
      </IconStyle>
    )
  }
}

export const IconGitHub = withIconStyle(<FontAwesomeIcon icon={faGithub} />)
export const IconYoutube = withIconStyle(<FontAwesomeIcon icon={faYoutube} />)
export const IconTwitter = withIconStyle(<FontAwesomeIcon icon={faTwitter} />)
