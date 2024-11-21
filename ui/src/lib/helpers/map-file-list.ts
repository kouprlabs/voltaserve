// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

export default function mapFileList(fileList: FileList | null): File[] {
  if (!fileList || fileList.length === 0) {
    return []
  }
  const files = []
  for (let i = 0; i < fileList.length; i++) {
    files.push(fileList[i])
  }
  return files
}
