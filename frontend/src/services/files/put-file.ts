import { fetcher } from '@/utils/fetcher'

export type PutSignedURLPrams = {
  signedPutUrl: string
  file: File
  contentType: string
}

export const putSignedUrl = async (
  params: PutSignedURLPrams,
): Promise<void> => {
  const { signedPutUrl, file, contentType } = params
  try {
    await fetcher(signedPutUrl, {
      method: 'PUT',
      headers: {
        Accept: 'application/json',
        'Content-Type': contentType,
      },
      data: file,
    })
  } catch {
    throw new Error('Failed to upload file')
  }
}
