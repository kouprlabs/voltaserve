import { useContext } from 'react'
import { Stack } from '@chakra-ui/react'
import { DrawerContext, variables } from '@koupr/ui'
import { IconInfoCircle } from '@koupr/ui'
import { File } from '@/api/file'
import FileInfo from '@/components/common/file-info'
import SwitchCard from '@/components/common/switch-card'
import DownloadButton from './download-button'
import OpenNewTabButton from './open-new-tab-button'

type DrawerContentProps = {
  file: File
}

const DrawerContent = ({ file }: DrawerContentProps) => {
  const { isCollapsed } = useContext(DrawerContext)
  return (
    <Stack spacing={variables.spacing}>
      <Stack spacing={variables.spacingSm}>
        <DownloadButton file={file} isCollapsed={isCollapsed} />
        <OpenNewTabButton file={file} isCollapsed={isCollapsed} />
        <SwitchCard
          icon={<IconInfoCircle fontSize="18px" />}
          label="Show info"
          isCollapsed={isCollapsed}
          localStorageNamespace="file_info"
          expandedMinWidth="200px"
        >
          <FileInfo file={file} />
        </SwitchCard>
      </Stack>
    </Stack>
  )
}

export default DrawerContent
