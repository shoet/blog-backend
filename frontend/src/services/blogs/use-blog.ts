import { ApiContext, Blog } from '@/types/api'
import { parseCookie } from '@/utils/cookie'
import { buildFetchClient, buildFetcher } from '@/utils/fetcher'
import useSWR from 'swr'

export const useBlog = (context: ApiContext, id: number) => {
  // TODO: paging
  const token = parseCookie(document.cookie)['authToken']
  const client = buildFetchClient({
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })
  const fetcher = buildFetcher(client)
  const url = `${context.apiBaseUrl}/blogs/${id}`
  const { data, isLoading, error } = useSWR<Blog>(url, fetcher)
  return {
    blog: data,
    isLoading,
    error,
  }
}
