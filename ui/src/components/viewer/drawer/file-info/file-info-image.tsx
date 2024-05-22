import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'

export type FileInfoImageProps = {
  file: File
}

const FileInfoImage = ({ file }: FileInfoImageProps) => {
  if (!file.snapshot?.original?.image) {
    return null
  }
  return (
    <Stat>
      <StatLabel>Image dimensions</StatLabel>
      <StatNumber className={cx('text-base')}>
        {file.snapshot?.original.image.width}x
        {file.snapshot?.original.image.height}
      </StatNumber>
    </Stat>
  )
}

export default FileInfoImage
