// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useContext } from 'react'
import { SwitchCard, IconInfo, SidenavContext } from '@koupr/ui'
import cx from 'classnames'
import { File } from '@/client/api/file'
import FileInfoEmbed from '@/components/file/info/file-info-embed'
import DrawerDownloadButton from './drawer-download-button'
import DrawerOpenNewTabButton from './drawer-open-new-tab-button'

export type DrawerContentProps = {
  file: File
}

const DrawerContent = ({ file }: DrawerContentProps) => {
  const { isCollapsed } = useContext(SidenavContext)
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
        <FileInfoEmbed file={file} />
      </SwitchCard>
    </div>
  )
}

export default DrawerContent
