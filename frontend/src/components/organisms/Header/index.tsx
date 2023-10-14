import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'

export const Header = () => {
  return (
    <nav>
      <Flex flexDirection="row" alignItems="baseline">
        <Box>
          <Text fontSize="display">shoet Blog</Text>
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
