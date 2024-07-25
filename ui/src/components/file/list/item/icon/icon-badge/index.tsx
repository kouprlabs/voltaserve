// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { File } from '@/client/api/file'
import { Status } from '@/client/api/snapshot'
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
          {file.snapshot?.status === Status.Waiting ? (
            <IconBadgeWaiting />
          ) : null}
          {file.snapshot?.status === Status.Processing || isLoading ? (
            <IconBadgeProcessing />
          ) : null}
          {file.snapshot?.status === Status.Error ? <IconBadgeError /> : null}
          {file.isShared ? <IconBadgeShared /> : null}
          {file.snapshot?.entities ? <IconBadgeInsights /> : null}
          {file.snapshot?.mosaic ? <IconBadgeMosaic /> : null}
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
