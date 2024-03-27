import cx from 'classnames'
import { File } from '@/client/api/file'
import FileInfoCreateTime from './file-info-create-time'
import FileInfoExtension from './file-info-extension'
import FileInfoImage from './file-info-image'
import FileInfoPermission from './file-info-permission'
import FileInfoSize from './file-info-size'
import FileInfoStorageUsage from './file-info-storage-usage'
import FileInfoUpdateTime from './file-info-update-time'

export type DrawerFileInfoProps = {
  file: File
}

const DrawerFileInfo = ({ file }: DrawerFileInfoProps) => {
  if (!file.original) {
    return null
  }
  return (
    <div className={cx('flex', 'flex-col', 'gap-1')}>
      <FileInfoImage file={file} />
      <FileInfoSize file={file} />
      <FileInfoExtension file={file} />
      <FileInfoStorageUsage file={file} />
      <FileInfoPermission file={file} />
      <FileInfoCreateTime file={file} />
      <FileInfoUpdateTime file={file} />
    </div>
  )
}

export default DrawerFileInfo
