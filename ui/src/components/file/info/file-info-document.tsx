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
import cx from 'classnames'
import { File } from '@/client/api/file'

export type FileInfoDocumentProps = {
  file: File
}

const FileInfoDocument = ({ file }: FileInfoDocumentProps) => {
  const document =
    file.snapshot?.preview?.document || file.snapshot?.ocr?.document
  if (!document) {
    return null
  }
  return (
    <Stat>
      <StatLabel>Document pages</StatLabel>
      <StatNumber className={cx('text-base')}>
        {document.page?.count}
      </StatNumber>
    </Stat>
  )
}

export default FileInfoDocument
