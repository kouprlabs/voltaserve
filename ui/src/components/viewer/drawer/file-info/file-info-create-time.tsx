import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import prettyDate from '@/helpers/pretty-date'

export type FileInfoCreateTimeProps = {
  file: File
}

const FileInfoCreateTime = ({ file }: FileInfoCreateTimeProps) => (
  <Stat>
    <StatLabel>Create time</StatLabel>
    <StatNumber className={cx('text-base')}>
      {prettyDate(file.createTime)}
    </StatNumber>
  </Stat>
)

export default FileInfoCreateTime
