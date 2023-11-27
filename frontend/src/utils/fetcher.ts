import axios, { AxiosRequestConfig } from 'axios'

export type FetcherBuildOptions = {
  headers?: any
}

export const buildFetchClient = (options?: FetcherBuildOptions) => {
  let headers = {
    Accept: 'application/json',
    'Content-Type': 'application/json',
  }

  if (options && options.headers) {
    headers = Object.assign(headers, options.headers)
  }

  return axios.create({
    withCredentials: true,
    headers: headers,
  })
}

const client = buildFetchClient()

export const buildFetcher = (client: any) => {
  return async (url: string, config: AxiosRequestConfig = {}): Promise<any> => {
    try {
      config.url = url
      const res = await client.request(config)
      return res.data
    } catch (err) {
      if (axios.isAxiosError(err)) {
        console.log(`Failed request by axios: ${err.message}`)
        throw err
      } else {
        console.log(`Failed request by unknown error: ${err}`)
        throw err
      }
    }
  }
}

export const fetcher = buildFetcher(client)
