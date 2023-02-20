export function byteToMegabyte(value: number): number {
  return value / 1e6
}

export function byteToGigabyte(value: number): number {
  return value / 1e9
}

export function byteToTerabyte(value: number): number {
  return value / 1e12
}

export function terabyteToByte(value: number): number {
  return value * 1e12
}

export function gigabyteToByte(value: number): number {
  return value * 1e9
}

export function megabyteToByte(value: number): number {
  return value * 1e6
}
