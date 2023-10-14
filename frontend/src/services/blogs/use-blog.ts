import { ApiContext, Blog } from '@/types/api'
import useSWR from 'swr'

export const useBlog = (context: ApiContext, id: number) => {
  // TODO: paging
  const url = `${context.apiBaseUrl}/blogs/${id}`
  const { data, isLoading, error } = useSWR<Blog>(url)
  return {
    blog: data,
    isLoading,
    error,
  }
}
