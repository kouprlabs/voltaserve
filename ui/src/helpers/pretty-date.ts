// eslint-disable-next-line @typescript-eslint/no-explicit-any
export default function prettyDate(date: any): string {
  return new Date(date).toLocaleString()
}
