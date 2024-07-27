// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useContext } from 'react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { DrawerContext } from '@/lib/components/drawer'
import { IconInfo } from '@/lib/components/icons'
import SwitchCard from '@/lib/components/switch-card'
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
        icon={<IconInfo />}
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
