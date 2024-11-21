// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

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
