import { ApiContext, Blog } from '@/types/api'
import useSWR from 'swr'

export const useBlogList = (context: ApiContext, initial: Blog[] = []) => {
  const url = `${context.apiBaseUrl}/blogs`
  const { data, isLoading, error } = useSWR<Blog[]>(url)
  return {
    blogs: data || initial,
    isLoading,
    error,
  }
}
