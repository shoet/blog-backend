import { Alert } from '@/components/atoms/Alert'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { LoginForm, LoginFormData } from '@/components/organisms/LoginForm'
import { signin } from '@/services/auth/signin'
import { ApiContext } from '@/types/api'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

export const LoginPage = () => {
  const [error, setError] = useState<string>()
  const apiContext: ApiContext = {
    apiBaseUrl: import.meta.env.VITE_API_BASE_URL,
  }
  const navigate = useNavigate()
  const onSubmit = async (data: LoginFormData) => {
    try {
      await signin(apiContext, {
        email: data.email,
        password: data.password,
      })
    } catch {
      setError('Invalid email or password')
      return
    }
    navigate('/admin')
  }
  return (
    <Flex height="100%" flexDirection="column" width="100%" alignItems="center">
      <Box width="50%" minWidth="300px" height="150px" marginTop="50px">
        {error && (
          <Box>
            <Alert text={error} onClick={() => setError('')} />
          </Box>
        )}
      </Box>
      <Box width="300px" minWidth="300px">
        <LoginForm onSubmit={onSubmit} />
      </Box>
    </Flex>
  )
}
