import { useAuthContext } from '@/contexts/AuthContext'
import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'

export const useAuthGuard = () => {
  const { authUser, isLoading } = useAuthContext()
  const navigate = useNavigate()

  useEffect(() => {
    if (!isLoading && !authUser) {
      navigate('/')
    }
  })
}
