// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { TaskStatus } from '@/client'
import { File } from '@/client/api/file'
import IconBadgeError from './icon-badge-error'
import IconBadgeInsights from './icon-badge-insights'
import IconBadgeMosaic from './icon-badge-mosaic'
import IconBadgeProcessing from './icon-badge-processing'
import IconBadgeShared from './icon-badge-shared'
import IconBadgeWaiting from './icon-badge-waiting'

export type IconBadgeProps = {
  file: File
  isLoading?: boolean
}

const IconBadge = ({ file, isLoading }: IconBadgeProps) => {
  return (
    <>
      {file.type === 'file' ? (
        <>
          {file.snapshot?.task?.status === TaskStatus.Waiting ? (
            <IconBadgeWaiting />
          ) : null}
          {file.snapshot?.task?.status === TaskStatus.Running || isLoading ? (
            <IconBadgeProcessing />
          ) : null}
          {file.snapshot?.task?.status === TaskStatus.Error ? (
            <IconBadgeError />
          ) : null}
          {file.isShared ? <IconBadgeShared /> : null}
          {file.snapshot?.capabilities.entities ||
          file.snapshot?.capabilities.summary ? (
            <IconBadgeInsights />
          ) : null}
          {file.snapshot?.capabilities.mosaic ? <IconBadgeMosaic /> : null}
        </>
      ) : null}
      {file.type === 'folder' ? (
        <>
          {file.isShared ? <IconBadgeShared /> : null}
          {isLoading ? <IconBadgeProcessing /> : null}
        </>
      ) : null}
    </>
  )
}

export default IconBadge
