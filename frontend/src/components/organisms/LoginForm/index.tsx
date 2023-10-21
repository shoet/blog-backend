import { Alert } from '@/components/atoms/Alert'
import { Button } from '@/components/atoms/Button'
import { Input } from '@/components/atoms/Input'
import { Text } from '@/components/atoms/Text'
import Box from '@/components/layout/Box'
import Flex from '@/components/layout/Flex'
import { useForm } from 'react-hook-form'

export type LoginFormData = {
  email: string
  password: string
}

type LoginFormProps = {
  errorMessgae?: string
  onSubmit?: (data: LoginFormData) => void
}

export const LoginForm = (props: LoginFormProps) => {
  const { onSubmit } = props
  const {
    control: _,
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>()

  const handleOnSubmit = async (data: LoginFormData) => {
    onSubmit && onSubmit(data)
  }

  return (
    <form>
      <Box>
        <Box>
          <Box>
            <Input
              placeholder="Email"
              {...register('email', {
                required: 'メールアドレスを入力してください。',
              })}
              name="email"
              hasError={!!errors.email}
            />
            {errors.email && (
              <Text as="label" variant="small" color="danger">
                {errors.email.message}
              </Text>
            )}
          </Box>
          <Box marginTop={2}>
            <Input
              type="password"
              placeholder="Password"
              {...register('password', {
                required: 'パスワードを入力してください。',
              })}
              name="password"
              hasError={!!errors.password}
            />
            {errors.password && (
              <Text as="label" variant="small" color="danger">
                {errors.password.message}
              </Text>
            )}
          </Box>
        </Box>
        <Flex justifyContent="end" marginTop={3}>
          <Button
            variant="primary"
            type="button"
            onClick={handleSubmit(handleOnSubmit)}
          >
            Sign In
          </Button>
        </Flex>
      </Box>
    </form>
  )
}
