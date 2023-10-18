import { ApiContext, Blog } from '@/types/api'
import { fetcher } from '@/utils/fetcher'

export type PutBlogParams = {
  blog: Omit<Omit<Omit<Blog, 'id'>, 'created'>, 'modified'>
}

export const putBlog = async (
  context: ApiContext,
  { blog }: PutBlogParams,
): Promise<Blog> => {
  const url = `${context.apiBaseUrl}/blogs/update`
  return await fetcher(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    data: JSON.stringify(blog),
  })
}
