import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { useRouteError } from 'react-router-dom'

function isError(error: any): error is { statusText: string; status: string } {
  return 'statusText' in error && 'status' in error
}

export const ErrorPage = () => {
  const error = useRouteError()
  return (
    <Box width="100%" height="100vh" backgroundColor="">
      <h3>Occurd error</h3>
      <Flex flexDirection="column" alignItems="center">
        {isError(error) && (
          <Box display="inline" backgroundColor="">
            <Flex flexDirection="column" alignItems="center" backgroundColor="">
              <h1>{error.status}</h1>
            </Flex>
            <div>{error.statusText}</div>
          </Box>
        )}
      </Flex>
    </Box>
  )
}
