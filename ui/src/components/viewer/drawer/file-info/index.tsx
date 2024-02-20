import classNames from 'classnames'
import { File } from '@/client/api/file'
import CreateTimeInfo from './create-time-info'
import ExtensionInfo from './extension-info'
import ImageInfo from './image-info'
import PermissionInfo from './permission-info'
import SizeInfo from './size-info'
import StorageUsageInfo from './storage-usage-info'
import UpdateTimeInfo from './update-time-info'

type FileInfoProps = {
  file: File
}

const FileInfo = ({ file }: FileInfoProps) => {
  if (!file.original) {
    return null
  }
  return (
    <div className={classNames('flex', 'flex-col', 'gap-1')}>
      <ImageInfo file={file} />
      <SizeInfo file={file} />
      <ExtensionInfo file={file} />
      <StorageUsageInfo file={file} />
      <PermissionInfo file={file} />
      <CreateTimeInfo file={file} />
      <UpdateTimeInfo file={file} />
    </div>
  )
}

export default FileInfo
