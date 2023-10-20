import { Header } from '@/components/organisms/Header'
import { Footer } from '@/components/organisms/Footer'
import { PropsWithChildren } from 'react'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'

export const VSplit = ({ children }: PropsWithChildren) => {
  return (
    <>
      <Flex flexDirection="column" minHeight="100vh" paddingTop={2}>
        <Box flexGrow={1}>
          <Box maxWidth="1280px" marginLeft="auto" marginRight="auto">
            <Box
              paddingLeft={{ base: 2, md: 4 }}
              paddingRight={{ base: 2, md: 4 }}
            >
              <Box marginBottom={3}>
                <Header />
              </Box>
              <main>{children}</main>
            </Box>
          </Box>
        </Box>
        <Box marginTop={3}>
          <Footer />
        </Box>
      </Flex>
    </>
  )
}
