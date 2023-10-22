import { Button } from '@/components/atoms/Button'
import { IconGitHub, IconTwitter, IconYoutube } from '@/components/atoms/Icon'
import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { useAuthContext } from '@/contexts/AuthContext'
import { useEffect } from 'react'
import { NavLink, useNavigate } from 'react-router-dom'

export const Profile = () => {
  // TODO: anchor link
  // TODO: icon size
  const {
    authUser,
    signout,
    isLoading,
    mutate: refreshSession,
  } = useAuthContext()
  const handleSignout = async () => {
    await signout()
    refreshSession()
  }
  useEffect(() => {
    refreshSession()
  }, [authUser])
  return (
    <Box>
      <Flex flexDirection="row" alignItems="baseline">
        <Box>
          <Text fontSize="large" fontWeight="bold" letterSpacing="large">
            {authUser ? (
              <NavLink to="/admin">
                <Text color="focusGreen">shoet</Text>
              </NavLink>
            ) : (
              <NavLink to="/admin/login">shoet</NavLink>
            )}
          </Text>
        </Box>
        <Flex flexDirection="row" paddingLeft={1} alignItems="center">
          <Box>
            <IconGitHub size={14} focusColor="focusGreen" />
          </Box>
          <Box paddingLeft={1}>
            <IconTwitter size={14} focusColor="focusGreen" />
          </Box>
          <Box paddingLeft={1}>
            <IconYoutube size={14} focusColor="focusGreen" />
          </Box>
          {authUser && (
            <Box paddingLeft={3}>
              {isLoading && <Text>Loading...</Text>}
              {authUser && (
                <Button variant="primary" onClick={handleSignout}>
                  SignOut
                </Button>
              )}
            </Box>
          )}
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
