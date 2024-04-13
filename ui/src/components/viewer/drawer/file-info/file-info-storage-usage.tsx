import { Progress, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import StorageAPI from '@/client/api/storage'
import prettyBytes from '@/helpers/pretty-bytes'

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
          {error && <span className={cx('text-red-500')}>Failed to load.</span>}
          {data && !error && (
            <>
              <span>
                {prettyBytes(data.bytes)} of {prettyBytes(data.maxBytes)} used
              </span>
              <Progress size="sm" value={data.percentage} hasStripe />
            </>
          )}
          {!data && !error && (
            <>
              <span>Calculating…</span>
              <Progress size="sm" value={0} hasStripe />
            </>
          )}
        </div>
      </StatNumber>
    </Stat>
  )
}

export default FileInfoStorageUsage
