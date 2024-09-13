export const randomColorFromString = (str: string) => {
  let hash = 0
  if (str.length === 0) return hash.toString()
  for (let i = 0; i < str.length; i += 1) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
    hash = hash & hash
  }
  let color = '#'
  for (let j = 0; j < 3; j += 1) {
    const value = (hash >> (j * 8)) & 255
    color += `00${value.toString(16)}`.slice(-2)
  }
  return color
}
