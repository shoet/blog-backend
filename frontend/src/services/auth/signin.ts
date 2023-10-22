import { ApiContext } from '@/types/api'
import { fetcher } from '@/utils/fetcher'

export type SigninParams = {
  email: string
  password: string
}

type SigninResponse = {
  authToken: string
}

export const signin = async (
  context: ApiContext,
  params: SigninParams,
): Promise<SigninResponse> => {
  const url = `${context.apiBaseUrl}/auth/signin`
  return await fetcher(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    data: JSON.stringify(params),
  })
}

export const signout = async (context: ApiContext): Promise<void> => {
  const url = `${context.apiBaseUrl}/auth/admin/signout`
  return await fetcher(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  })
}
