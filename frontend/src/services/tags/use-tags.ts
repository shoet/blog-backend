import { ApiContext, Tag } from '@/types/api'
import useSWR from 'swr'

export const useTags = (context: ApiContext, initial: Tag[] = []) => {
  const url = `${context.apiBaseUrl}/tags`
  const { data, isLoading, error } = useSWR<Tag[]>(url)
  return {
    tags: data || initial,
    isLoading,
    error,
  }
}
