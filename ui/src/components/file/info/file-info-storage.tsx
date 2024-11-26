// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Progress, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import StorageAPI from '@/client/api/storage'
import prettyBytes from '@/lib/helpers/pretty-bytes'

export type FileInfoStorageProps = {
  file: File
}

const FileInfoStorage = ({ file }: FileInfoStorageProps) => {
  const { data: usage, error } = StorageAPI.useGetFileUsage(file.id)
  return (
    <Stat>
      <StatLabel>Storage</StatLabel>
      <StatNumber className={cx('text-base')}>
        <div className={cx('flex', 'flex-col', 'gap-0.5')}>
          {error ? (
            <span className={cx('text-red-500')}>Failed to load.</span>
          ) : null}
          {usage && !error ? (
            <>
              <span>
                {prettyBytes(usage.bytes)} of {prettyBytes(usage.maxBytes)} used
              </span>
              <Progress size="sm" value={usage.percentage} hasStripe />
            </>
          ) : null}
          {!usage && !error ? (
            <>
              <span>Calculating…</span>
              <Progress size="sm" value={0} hasStripe />
            </>
          ) : null}
        </div>
      </StatNumber>
    </Stat>
  )
}

export default FileInfoStorage
