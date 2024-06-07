import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import prettyDate from '@/lib/helpers/pretty-date'

export type FileInfoUpdateTimeProps = {
  file: File
}

const FileInfoUpdateTime = ({ file }: FileInfoUpdateTimeProps) => {
  if (
    !file.updateTime ||
    (file.updateTime &&
      file.createTime.includes(file.updateTime.replaceAll('Z', '')))
  ) {
    return null
  }
  return (
    <Stat>
      <StatLabel>Update time</StatLabel>
      <StatNumber className={cx('text-base')}>
        {prettyDate(file.updateTime)}
      </StatNumber>
    </Stat>
  )
}

export default FileInfoUpdateTime
