import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'

type BaseLayoutProps = {
  MainContent: React.ReactNode
  SubContent: React.ReactNode
}

export const BaseLayout = (props: BaseLayoutProps) => {
  const { MainContent, SubContent } = props
  return (
    <Flex
      flexDirection={{ base: 'column', sm: 'row' }}
      justifyContent="space-between"
    >
      <Box
        width={{ base: '100%', sm: '70%' }}
        marginBottom={{ base: 2, sm: 0 }}
      >
        {MainContent}
      </Box>
      <Box width={{ base: '100%', sm: '25%' }}>{SubContent}</Box>
    </Flex>
  )
}
