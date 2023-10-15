import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'

export const Header = () => {
  return (
    <nav>
      <Flex flexDirection="row" alignItems="baseline">
        <Box>
          <Box display="inline-flex">
            <Text fontSize="display" fontWeight="bold" letterSpacing="large">
              shoet
            </Text>
          </Box>
          <Box display="inline-flex" marginLeft={1}>
            <Text fontSize="display" letterSpacing="large">
              Blog
            </Text>
          </Box>
        </Box>
        <Box marginLeft={2}>
          <Text fontSize="small" color="gray">
            技術や好きなことについて発信しています。
          </Text>
        </Box>
      </Flex>
    </nav>
  )
}
