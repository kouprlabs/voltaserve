// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useMemo } from 'react'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'
import variables from '@/lib/variables'

export type ViewerVideoProps = {
  file: File
}

const ViewerVideo = ({ file }: ViewerVideoProps) => {
  const download = useMemo(() => file.snapshot?.original, [file])
  const url = useMemo(() => {
    if (!download || !download.extension) {
      return ''
    }
    return `/proxy/api/v3/files/${file.id}/preview${
      download.extension
    }?${new URLSearchParams({
      access_token: getAccessTokenOrRedirect(),
    })}`
  }, [file, download])

  if (!download) {
    return null
  }

  return (
    <video
      controls
      autoPlay
      style={{
        maxWidth: '100%',
        maxHeight: '100%',
        borderRadius: variables.borderRadius,
      }}
    >
      <source src={url} />
    </video>
  )
}

export default ViewerVideo
