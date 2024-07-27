// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Progress, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import StorageAPI from '@/client/api/storage'
import prettyBytes from '@/lib/helpers/pretty-bytes'

export type FileInfoStorageUsageProps = {
  file: File
}

const FileInfoStorageUsage = ({ file }: FileInfoStorageUsageProps) => {
  const { data, error } = StorageAPI.useGetFileUsage(file.id)
  return (
    <Stat>
      <StatLabel>Storage usage</StatLabel>
      <StatNumber className={cx('text-base')}>
        <div className={cx('flex', 'flex-col', 'gap-0.5')}>
          {error ? (
            <span className={cx('text-red-500')}>Failed to load.</span>
          ) : null}
          {data && !error ? (
            <>
              <span>
                {prettyBytes(data.bytes)} of {prettyBytes(data.maxBytes)} used
              </span>
              <Progress size="sm" value={data.percentage} hasStripe />
            </>
          ) : null}
          {!data && !error ? (
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

export default FileInfoStorageUsage
