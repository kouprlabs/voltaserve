import { Badge, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'

export type FileInfoPermissionProps = {
  file: File
}

const FileInfoPermission = ({ file }: FileInfoPermissionProps) => (
  <Stat>
    <StatLabel>Permission</StatLabel>
    <StatNumber className={cx('text-base')}>
      <Badge>{file.permission}</Badge>
    </StatNumber>
  </Stat>
)

export default FileInfoPermission
