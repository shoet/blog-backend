import { ApiContext } from '@/types/api'
import { fetcher } from '@/utils/fetcher'

export type GetSignedPutURLPrams = {
  fileName: string
}

type GetSignedPutURLResponse = {
  signedUrl: string
  putUrl: string
}

export const getSignedPutUrl = async (
  context: ApiContext,
  { fileName }: GetSignedPutURLPrams,
  authToken: string,
): Promise<GetSignedPutURLResponse> => {
  const url = `${context.apiBaseUrl}/files/thumbnail/new`
  return await fetcher(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      Authorization: `Bearer ${authToken}`,
    },
    data: JSON.stringify({ fileName: fileName }),
  })
}
