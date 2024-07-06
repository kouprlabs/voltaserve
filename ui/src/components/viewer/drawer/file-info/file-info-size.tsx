// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import prettyBytes from '@/lib/helpers/pretty-bytes'

export type FileInfoSizeProps = {
  file: File
}

const FileInfoSize = ({ file }: FileInfoSizeProps) => {
  if (!file.snapshot?.original) {
    return null
  }
  return (
    <Stat>
      <StatLabel>File size</StatLabel>
      <StatNumber className={cx('text-base')}>
        {prettyBytes(file.snapshot?.original.size)}
      </StatNumber>
    </Stat>
  )
}

export default FileInfoSize
