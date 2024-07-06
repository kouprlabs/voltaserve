// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { Button } from '@chakra-ui/react'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import { swrConfig } from '@/client/options'
import { IconOpenInNew } from '@/lib/components/icons'
import { useAppSelector } from '@/store/hook'

const MosaicOverviewArtifacts = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const { data: file } = FileAPI.useGet(id, swrConfig())

  if (!file) {
    return null
  }

  return (
    <div
      className={cx(
        'flex',
        'flex-col',
        'items-center',
        'justify-center',
        'gap-1',
      )}
    >
      <Button
        as="a"
        type="button"
        leftIcon={<IconOpenInNew />}
        target="_blank"
        href={`/file/${file.id}/mosaic`}
      >
        View Mosaic
      </Button>
    </div>
  )
}

export default MosaicOverviewArtifacts
