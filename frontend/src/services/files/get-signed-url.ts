import { ApiContext } from '@/types/api'
import { fetcher } from '@/utils/fetcher'

type FileType = 'thumbnail' | 'content'

export type GetSignedPutURLPrams = {
  fileName: string
  fileType: FileType
}

type GetSignedPutURLResponse = {
  signedUrl: string
  putUrl: string
}

export const getSignedPutUrl = async (
  context: ApiContext,
  { fileName, fileType = 'thumbnail' }: GetSignedPutURLPrams,
  authToken: string,
): Promise<GetSignedPutURLResponse> => {
  const url = `${context.apiBaseUrl}/files/${fileType}/new`
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
