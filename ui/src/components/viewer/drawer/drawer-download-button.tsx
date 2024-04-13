import { Button, IconButton } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import downloadFile from '@/helpers/download-file'
import { IconDownload } from '@/lib'

export type DrawerDownloadButtonProps = {
  file: File
  isCollapsed?: boolean
}

const DrawerDownloadButton = ({
  file,
  isCollapsed,
}: DrawerDownloadButtonProps) => {
  if (isCollapsed) {
    return (
      <IconButton
        icon={<IconDownload fontSize="18px" />}
        variant="solid"
        colorScheme="blue"
        aria-label="Download"
        title="Download"
        className={cx('h-[50px]', 'w-[50px]', 'p-1.5', 'rounded-md')}
        onClick={() => downloadFile(file)}
      />
    )
  } else {
    return (
      <Button
        leftIcon={<IconDownload fontSize="18px" />}
        variant="solid"
        colorScheme="blue"
        className={cx('h-[50px]', 'w-full', 'p-1.5', 'rounded-md')}
        onClick={() => downloadFile(file)}
      >
        Download
      </Button>
    )
  }
}

export default DrawerDownloadButton
