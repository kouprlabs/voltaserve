import { useContext } from 'react'
import { DrawerContext, SwitchCard, IconInfoCircle } from '@koupr/ui'
import cx from 'classnames'
import { File } from '@/client/api/file'
import DrawerDownloadButton from './drawer-download-button'
import DrawerOpenNewTabButton from './drawer-open-new-tab-button'
import DrawerFileInfo from './file-info'

export type DrawerContentProps = {
  file: File
}

const DrawerContent = ({ file }: DrawerContentProps) => {
  const { isCollapsed } = useContext(DrawerContext)
  return (
    <div className={cx('flex', 'flex-col', 'gap-1')}>
      <DrawerDownloadButton file={file} isCollapsed={isCollapsed} />
      <DrawerOpenNewTabButton file={file} isCollapsed={isCollapsed} />
      <SwitchCard
        icon={<IconInfoCircle fontSize="18px" />}
        label="Show info"
        isCollapsed={isCollapsed}
        localStorageNamespace="file_info"
        expandedMinWidth="200px"
      >
        <DrawerFileInfo file={file} />
      </SwitchCard>
    </div>
  )
}

export default DrawerContent
