// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import cx from 'classnames'
import { File } from '@/client/api/file'
import FileInfoItemCount from '@/components/file/info/file-info-item-count'
import FileInfoCreateTime from './file-info-create-time'
import FileInfoExtension from './file-info-extension'
import FileInfoImage from './file-info-image'
import FileInfoPermission from './file-info-permission'
import FileInfoSize from './file-info-size'
import FileInfoStorageUsage from './file-info-storage-usage'
import FileInfoUpdateTime from './file-info-update-time'

export type FileInfoEmbedProps = {
  file: File
}

const FileInfoEmbed = ({ file }: FileInfoEmbedProps) => (
  <div className={cx('flex', 'flex-col', 'gap-1')}>
    <FileInfoImage file={file} />
    <FileInfoSize file={file} />
    <FileInfoExtension file={file} />
    <FileInfoItemCount file={file} />
    <FileInfoStorageUsage file={file} />
    <FileInfoPermission file={file} />
    <FileInfoCreateTime file={file} />
    <FileInfoUpdateTime file={file} />
  </div>
)

export default FileInfoEmbed
