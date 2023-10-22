import { Alert } from '@/components/atoms/Alert'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { LoginForm, LoginFormData } from '@/components/organisms/LoginForm'
import { useAuthContext } from '@/contexts/AuthContext'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

export const LoginPage = () => {
  const [error, setError] = useState<string>()
  const navigate = useNavigate()
  const { signin } = useAuthContext()
  const onSubmit = async (data: LoginFormData) => {
    try {
      await signin(data.email, data.password)
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
