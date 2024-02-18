import { useContext } from 'react'
import { DrawerContext, SwitchCard, IconInfoCircle } from '@koupr/ui'
import classNames from 'classnames'
import { File } from '@/client/api/file'
import DownloadButton from './download-button'
import FileInfo from './file-info'
import OpenNewTabButton from './open-new-tab-button'

type DrawerContentProps = {
  file: File
}

const DrawerContent = ({ file }: DrawerContentProps) => {
  const { isCollapsed } = useContext(DrawerContext)
  return (
    <div className={classNames('flex', 'flex-col', 'gap-1')}>
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
    </div>
  )
}

export default DrawerContent
