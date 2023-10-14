import { IconGitHub, IconTwitter, IconYoutube } from '@/components/atoms/Icon'
import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'

export const Profile = () => {
  // TODO: anchor link
  return (
    <Box>
      <Flex flexDirection="row" alignItems="baseline">
        <Box>
          <Text fontSize="large" fontWeight="bold" letterSpacing="large">
            shoet
          </Text>
        </Box>
        <Flex flexDirection="row" paddingLeft={1} alignItems="center">
          <Box>
            <IconGitHub size={14} />
          </Box>
          <Box paddingLeft={1}>
            <IconTwitter size={14} />
          </Box>
          <Box paddingLeft={1}>
            <IconYoutube size={14} />
          </Box>
        </Flex>
      </Flex>
      <Box paddingTop={1}>
        <Text variant="small">
          エンジニア。
          <br />
          エンジニアリングで価値提供できるよう、日々自己研磨。
        </Text>
      </Box>
    </Box>
  )
}
