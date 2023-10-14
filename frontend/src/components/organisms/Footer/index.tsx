import Flex from '@/components/layout/Flex'
import { IconGitHub, IconTwitter, IconYoutube } from '@/components/atoms/Icon'
import Box from '@/components/layout/Box'
import { Text } from '@/components/atoms/Text'

export const Footer = () => {
  return (
    <Flex
      flexDirection="column"
      alignItems="center"
      justifyContent="space-between"
      height="200px"
      width="100%"
      padding={4}
    >
      <Flex flexDirection="row" alignItems="center">
        <Box>
          <IconGitHub size={14} />
        </Box>
        <Box marginLeft={2}>
          <IconTwitter size={14} />
        </Box>
        <Box marginLeft={2}>
          <IconYoutube size={14} />
        </Box>
      </Flex>
      <Box>
        <Text>
          &copy;{` ${new Date().getFullYear()} shoet. All rights reserved.`}
        </Text>
      </Box>
    </Flex>
  )
}
