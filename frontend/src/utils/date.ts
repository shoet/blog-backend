export function toStringYYYYMMDD_HHMMSS(date: string): string {
  const d = new Date(date)
  return `${d.toLocaleDateString()} ${d.toLocaleTimeString()}`
}
