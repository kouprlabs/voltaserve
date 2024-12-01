// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { RefObject, useEffect } from 'react'

export default function useFocusAndSelectAll(inputRef: RefObject<HTMLInputElement>, enable?: boolean) {
  useEffect(() => {
    setTimeout(() => {
      if (enable) {
        inputRef.current?.focus()
        setTimeout(() => inputRef.current?.select(), 100)
      }
    }, 100)
  }, [inputRef, enable])
}
