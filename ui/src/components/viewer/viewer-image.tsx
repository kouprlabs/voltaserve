// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useMemo, useState } from 'react'
import { SectionSpinner, variables } from '@koupr/ui'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

export type ViewerImageProps = {
  file: File
}

const ViewerImage = ({ file }: ViewerImageProps) => {
  const [isLoading, setIsLoading] = useState(true)
  const url = useMemo(() => {
    if (file.snapshot?.preview && file.snapshot?.preview.extension) {
      return `/proxy/api/v3/files/${file.id}/preview${
        file.snapshot?.preview.extension
      }?${new URLSearchParams({
        access_token: getAccessTokenOrRedirect(),
      })}`
    }
  }, [file])

  if (!file.snapshot?.preview) {
    return null
  }

  return (
    <div className={cx('flex', 'flex-col', 'w-full', 'h-full', 'gap-1.5')}>
      <div
        className={cx(
          'relative',
          'flex',
          'items-center',
          'justify-center',
          'grow',
          'w-full',
          'h-full',
          'overflow-scroll',
        )}
      >
        {isLoading ? <SectionSpinner /> : null}
        <img
          src={url}
          style={{
            objectFit: 'contain',
            width: isLoading ? 0 : 'auto',
            height: isLoading ? 0 : '90%',
            visibility: isLoading ? 'hidden' : 'visible',
            borderRadius: variables.borderRadius,
          }}
          onLoad={() => setIsLoading(false)}
          alt={file.name}
        />
      </div>
    </div>
  )
}

export default ViewerImage
