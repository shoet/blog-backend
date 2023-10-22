import { ApiContext, Blog } from '@/types/api'
import { fetcher } from '@/utils/fetcher'

export type AddBlogParams = {
  blog: Omit<Omit<Omit<Blog, 'id'>, 'created'>, 'modified'>
}

export const addBlog = async (
  context: ApiContext,
  { blog }: AddBlogParams,
  authToken: string,
): Promise<Blog> => {
  const url = `${context.apiBaseUrl}/blogs`
  return await fetcher(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      Authorization: `Bearer ${authToken}`,
    },
    data: JSON.stringify(blog),
  })
}
