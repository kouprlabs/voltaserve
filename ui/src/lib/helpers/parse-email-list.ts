// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import * as Yup from 'yup'

export default function parseEmailList(value: string): string[] {
  return [...new Set(value.split(',').map((e: string) => e.trim()))].filter((e) => {
    if (e.length === 0) {
      return false
    }
    try {
      Yup.string()
        .email()
        .matches(/.+(\.[A-Za-z]{2,})$/, 'Email must end with a valid top-level domain')
        .validateSync(e)
      return true
    } catch {
      return false
    }
  })
}
