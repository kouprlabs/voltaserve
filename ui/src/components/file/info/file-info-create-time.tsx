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
