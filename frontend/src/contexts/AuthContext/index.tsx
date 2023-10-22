import { signin, signout } from '@/services/auth/signin'
import { ApiContext, User } from '@/types/api'
import { parseCookie } from '@/utils/cookie'
import { fetcher } from '@/utils/fetcher'
import React, { PropsWithChildren } from 'react'
import useSWR from 'swr'

type AuthContextValue = {
  authUser?: User
  isLoading: boolean
  error?: Error
  singin: (email: string, password: string) => Promise<void>
  signout: () => Promise<void>
  mutate: () => Promise<User | undefined>
}

const AuthContext = React.createContext<AuthContextValue>({
  authUser: undefined,
  isLoading: false,
  error: undefined,
  singin: async () => {},
  signout: async () => {},
  mutate: async () => {
    return {} as User
  },
})

export const useAuthContext = () => React.useContext(AuthContext)

type AuthContextProviderProps = {
  context: ApiContext
}

export const AuthContextProvider = (
  props: PropsWithChildren<AuthContextProviderProps>,
) => {
  const { children, context } = props

  const authFetcher = (url: string) => {
    return fetcher(url, {
      headers: {
        Authorization: `Bearer ${parseCookie(document.cookie)['authToken']}`,
      },
    })
  }

  const {
    data: authUser,
    isLoading,
    error,
    mutate: mutateFunc,
  } = useSWR<User>(
    `${import.meta.env.VITE_API_BASE_URL}/auth/login/me`,
    authFetcher,
  )

  const signinFunc = async (email: string, password: string) => {
    await signin(context, { email, password })
    await mutateFunc()
  }

  const signoutFunc = async () => {
    await signout(context)
    await mutateFunc()
  }

  return (
    <AuthContext.Provider
      value={{
        isLoading,
        error,
        authUser,
        singin: signinFunc,
        signout: signoutFunc,
        mutate: mutateFunc,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}
