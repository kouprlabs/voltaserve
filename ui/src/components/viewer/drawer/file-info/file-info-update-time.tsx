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
import prettyDate from '@/lib/helpers/pretty-date'

export type FileInfoUpdateTimeProps = {
  file: File
}

const FileInfoUpdateTime = ({ file }: FileInfoUpdateTimeProps) => {
  if (
    !file.updateTime ||
    (file.updateTime &&
      file.createTime.includes(file.updateTime.replaceAll('Z', '')))
  ) {
    return null
  }
  return (
    <Stat>
      <StatLabel>Update time</StatLabel>
      <StatNumber className={cx('text-base')}>
        {prettyDate(file.updateTime)}
      </StatNumber>
    </Stat>
  )
}

export default FileInfoUpdateTime
