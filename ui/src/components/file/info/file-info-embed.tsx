// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import cx from 'classnames'
import { File, FileType } from '@/client/api/file'
import FileInfoDocument from '@/components/file/info/file-info-document'
import FileInfoItemCount from '@/components/file/info/file-info-item-count'
import FileInfoCreateTime from './file-info-create-time'
import FileInfoExtension from './file-info-extension'
import FileInfoImage from './file-info-image'
import FileInfoName from './file-info-name'
import FileInfoPermission from './file-info-permission'
import FileInfoStorage from './file-info-storage'
import FileInfoUpdateTime from './file-info-update-time'

export type FileInfoEmbedProps = {
  file: File
}

const FileInfoEmbed = ({ file }: FileInfoEmbedProps) => (
  <div className={cx('flex', 'flex-col', 'gap-1')}>
    <FileInfoName file={file} />
    <FileInfoImage file={file} />
    <FileInfoDocument file={file} />
    <FileInfoExtension file={file} />
    {file.type === FileType.Folder ? <FileInfoItemCount file={file} /> : null}
    <FileInfoStorage file={file} />
    <FileInfoPermission file={file} />
    <FileInfoCreateTime file={file} />
    <FileInfoUpdateTime file={file} />
  </div>
)

export default FileInfoEmbed
