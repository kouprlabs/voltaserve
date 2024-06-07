export default function truncateMiddle(text: string, maxLength: number) {
  if (text.length <= maxLength) {
    return text
  }
  const half = Math.floor((maxLength - 3) / 2)
  return text.slice(0, half) + '…' + text.slice(text.length - half)
}
