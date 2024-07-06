// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

export default function truncateMiddle(text: string, maxLength: number) {
  if (text.length <= maxLength) {
    return text
  }
  const half = Math.floor((maxLength - 3) / 2)
  return text.slice(0, half) + '…' + text.slice(text.length - half)
}
