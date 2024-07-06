// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

export default function prettyBytes(value: number) {
  const UNITS = ['B', 'kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
  if (!Number.isFinite(value)) {
    throw new TypeError(
      `Expected a finite number, got ${typeof value}: ${value}`,
    )
  }
  const isNegative = value < 0
  if (isNegative) {
    value = -value
  }
  if (value < 1) {
    return (isNegative ? '-' : '') + value + ' B'
  }
  const exponent = Math.min(Math.floor(Math.log10(value) / 3), UNITS.length - 1)
  const number = Number((value / Math.pow(1000, exponent)).toPrecision(3))
  const unit = UNITS[exponent]
  return (isNegative ? '-' : '') + number + ' ' + unit
}
