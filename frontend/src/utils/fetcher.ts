import axios, { AxiosRequestConfig } from 'axios'

const client = axios.create({
  withCredentials: true,
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json',
  },
})

export const fetcher = async (
  url: string,
  config: AxiosRequestConfig = {},
): Promise<any> => {
  try {
    config.url = url
    const res = await client.request(config)
    return res.data
  } catch (err) {
    if (axios.isAxiosError(err)) {
      console.log(`Failed request by axios: ${err.message}`)
      throw new Error('Failed request by fetcher')
    } else {
      console.log(`Failed request by unknown error: ${err}`)
      throw new Error('Failed request by unknown error')
    }
  }
}
