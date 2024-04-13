import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import prettyBytes from '@/helpers/pretty-bytes'

export type FileInfoSizeProps = {
  file: File
}

const FileInfoSize = ({ file }: FileInfoSizeProps) => {
  if (!file.original) {
    return null
  }
  return (
    <Stat>
      <StatLabel>File size</StatLabel>
      <StatNumber className={cx('text-base')}>
        {prettyBytes(file.original.size)}
      </StatNumber>
    </Stat>
  )
}

export default FileInfoSize
