// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { FileAPI, File } from '@/client/api/file'

export type FileInfoItemCountProps = {
  file: File
}

const FileInfoItemCount = ({ file }: FileInfoItemCountProps) => {
  const { data: count, error } = FileAPI.useGetCount(file.id)

  return (
    <Stat>
      <StatLabel>Item count</StatLabel>
      <StatNumber className={cx('text-base')}>
        <div className={cx('flex', 'flex-col', 'gap-0.5')}>
          {error ? (
            <span className={cx('text-red-500')}>Failed to load.</span>
          ) : null}
          {count != null && !error ? <span>{count}</span> : null}
          {count == null && !error ? <span>Calculating…</span> : null}
        </div>
      </StatNumber>
    </Stat>
  )
}

export default FileInfoItemCount
