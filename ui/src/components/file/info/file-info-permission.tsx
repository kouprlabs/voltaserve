// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Badge, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'

export type FileInfoPermissionProps = {
  file: File
}

const FileInfoPermission = ({ file }: FileInfoPermissionProps) => (
  <Stat>
    <StatLabel>Permission</StatLabel>
    <StatNumber className={cx('text-base')}>
      <Badge>{file.permission}</Badge>
    </StatNumber>
  </Stat>
)

export default FileInfoPermission
