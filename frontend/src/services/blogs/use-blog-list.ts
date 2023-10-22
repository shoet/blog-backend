import { ApiContext, Blog } from '@/types/api'
import useSWR from 'swr'

export const useBlogList = (context: ApiContext, initial: Blog[] = []) => {
  const url = `${context.apiBaseUrl}/blogs`
  const { data, isLoading, error, mutate } = useSWR<Blog[]>(url)
  return {
    blogs: data || initial,
    isLoading,
    error,
    mutate,
  }
}
