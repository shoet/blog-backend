import * as base32 from 'hi-base32'
import { v4 as uuidv4 } from 'uuid'

export function generateBase32EncodedUuid(): string {
  const uuid = uuidv4().replace(/-/g, '')
  const base32UUID = base32.encode(uuid).replace(/=/g, '')
  return base32UUID
}
