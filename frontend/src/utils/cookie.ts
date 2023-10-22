export function parseCookie(cookie: string): { [key: string]: string } {
  const cookieObj: { [key: string]: string } = {}
  cookie.split(';').forEach((c) => {
    const [key, value] = c.split('=')
    cookieObj[key.trim()] = value
  })
  return cookieObj
}
