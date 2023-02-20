import { Stack } from '@chakra-ui/react'
import { File } from '@/api/file'
import variables from '@/theme/variables'
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
    <Stack spacing={variables.spacingSm}>
      <ImageInfo file={file} />
      <SizeInfo file={file} />
      <ExtensionInfo file={file} />
      <StorageUsageInfo file={file} />
      <PermissionInfo file={file} />
      <CreateTimeInfo file={file} />
      <UpdateTimeInfo file={file} />
    </Stack>
  )
}

export default FileInfo
