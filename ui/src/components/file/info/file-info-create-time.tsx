// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import { RelativeDate } from '@koupr/ui'
import cx from 'classnames'
import { File } from '@/client/api/file'

export type FileInfoCreateTimeProps = {
  file: File
}

const FileInfoCreateTime = ({ file }: FileInfoCreateTimeProps) => (
  <Stat>
    <StatLabel>Created</StatLabel>
    <StatNumber className={cx('text-base')}>
      <RelativeDate date={new Date(file.createTime)} />
    </StatNumber>
  </Stat>
)

export default FileInfoCreateTime
