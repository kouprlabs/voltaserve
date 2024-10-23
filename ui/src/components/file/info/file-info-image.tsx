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

export type FileInfoImageProps = {
  file: File
}

const FileInfoImage = ({ file }: FileInfoImageProps) => {
  if (!file.snapshot?.original?.image) {
    return null
  }
  return (
    <Stat>
      <StatLabel>Image dimensions</StatLabel>
      <StatNumber className={cx('text-base')}>
        {file.snapshot?.original.image.width}x
        {file.snapshot?.original.image.height}
      </StatNumber>
    </Stat>
  )
}

export default FileInfoImage
