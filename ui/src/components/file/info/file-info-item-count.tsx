import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { FileAPI, File, FileType } from '@/client/api/file'

export type FileInfoItemCountProps = {
  file: File
}

const FileInfoItemCount = ({ file }: FileInfoItemCountProps) => {
  const { data: count, error } = FileAPI.useGetCount(file.id)

  if (file.type !== FileType.Folder) {
    return null
  }

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
