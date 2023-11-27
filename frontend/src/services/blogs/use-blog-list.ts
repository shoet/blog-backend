import { ApiContext, Blog } from '@/types/api'
import { parseCookie } from '@/utils/cookie'
import { buildFetchClient, buildFetcher } from '@/utils/fetcher'
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

export const useBlogListAdmin = (context: ApiContext, initial: Blog[] = []) => {
  const url = `${context.apiBaseUrl}/admin/blogs`
  const token = parseCookie(document.cookie)['authToken']
  const client = buildFetchClient({
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })
  const fetcher = buildFetcher(client)
  const { data, isLoading, error, mutate } = useSWR<Blog[]>(url, fetcher)
  return {
    blogs: data || initial,
    isLoading,
    error,
    mutate,
  }
}
