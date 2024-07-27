// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import {
  byteToGigabyte,
  byteToMegabyte,
  byteToTerabyte,
  gigabyteToByte,
  megabyteToByte,
  terabyteToByte,
} from '@/lib/helpers/convert-storage'

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
  throw new Error(`Invalid unit: ${unit}`)
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
  throw new Error(`Invalid unit: ${unit}`)
}
