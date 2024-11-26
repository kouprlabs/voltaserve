// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Stat, StatLabel, StatNumber, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import truncateMiddle from '@/lib/helpers/truncate-middle'

export type FileInfoNameProps = {
  file: File
}

const FileInfoName = ({ file }: FileInfoNameProps) => (
  <Stat>
    <StatLabel>Name</StatLabel>
    <StatNumber className={cx('text-base')}>
      {file.name.length > 40 ? (
        <Tooltip label={file.name}>
          <span>{truncateMiddle(file.name, 40)}</span>
        </Tooltip>
      ) : (
        <span>{truncateMiddle(file.name, 40)}</span>
      )}
    </StatNumber>
  </Stat>
)

export default FileInfoName
