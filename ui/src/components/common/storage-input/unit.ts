import {
  byteToGigabyte,
  byteToMegabyte,
  byteToTerabyte,
  gigabyteToByte,
  megabyteToByte,
  terabyteToByte,
} from '@/helpers/convert-storage'

export type Unit = 'b' | 'mb' | 'gb' | 'tb'

export function getUnit(value: number): Unit {
  if (value >= 1e12) {
    return 'tb'
  }
  if (value >= 1e9) {
    return 'gb'
  }
  if (value >= 1e6) {
    return 'mb'
  }
  return 'b'
}

export function convertFromByte(value: number, unit: Unit): number {
  if (unit === 'b') {
    return value
  }
  if (unit === 'mb') {
    return byteToMegabyte(value)
  }
  if (unit === 'gb') {
    return byteToGigabyte(value)
  }
  if (unit === 'tb') {
    return byteToTerabyte(value)
  }
  throw `Invalid unit ${unit}`
}

export function normalizeToByte(value: number, unit: Unit) {
  if (unit === 'b') {
    return value
  }
  if (unit === 'mb') {
    return megabyteToByte(value)
  }
  if (unit === 'gb') {
    return gigabyteToByte(value)
  }
  if (unit === 'tb') {
    return terabyteToByte(value)
  }
  throw `Invalid unit ${unit}`
}
