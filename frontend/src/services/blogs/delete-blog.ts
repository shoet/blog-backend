import { ApiContext } from '@/types/api'
import { fetcher } from '@/utils/fetcher'

export type DeleteBlogParams = {
  blogId: number
}

export const deleteBlog = async (
  context: ApiContext,
  { blogId }: DeleteBlogParams,
  authToken: string,
): Promise<void> => {
  const url = `${context.apiBaseUrl}/blogs/delete`
  return await fetcher(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      Authorization: `Bearer ${authToken}`,
    },
    data: JSON.stringify({ id: blogId }),
  })
}
