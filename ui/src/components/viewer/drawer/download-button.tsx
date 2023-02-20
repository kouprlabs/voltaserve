import { Button, IconButton } from '@chakra-ui/react'
import { File } from '@/api/file'
import { IconDownload } from '@/components/common/icon'
import variables from '@/theme/variables'
import downloadFile from '@/helpers/download-file'

type DownloadButtonProps = {
  file: File
  isCollapsed?: boolean
}

const DownloadButton = ({ file, isCollapsed }: DownloadButtonProps) => {
  if (isCollapsed) {
    return (
      <IconButton
        icon={<IconDownload fontSize="18px" />}
        variant="solid"
        colorScheme="blue"
        w="50px"
        h="50px"
        p={variables.spacing}
        borderRadius={variables.borderRadiusSm}
        aria-label="Download"
        title="Download"
        onClick={() => downloadFile(file)}
      />
    )
  } else {
    return (
      <Button
        leftIcon={<IconDownload fontSize="18px" />}
        variant="solid"
        colorScheme="blue"
        w="100%"
        h="50px"
        p={variables.spacing}
        borderRadius={variables.borderRadiusSm}
        onClick={() => downloadFile(file)}
      >
        Download
      </Button>
    )
  }
}

export default DownloadButton
