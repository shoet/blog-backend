import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'

export const Profile = () => {
  return (
    <Box>
      <Text fontSize="large" fontWeight="bold" letterSpacing="large">
        shoet
      </Text>
      <Box paddingTop={1}>
        <Text variant="small">
          Webエンジニア。
          <br />
          エンジニアリングで価値提供できるよう、日々自己研磨。
        </Text>
      </Box>
    </Box>
  )
}
