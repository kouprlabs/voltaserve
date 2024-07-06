// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { File } from '@/client/api/file'

export enum FileViewType {
  Grid = 'grid',
  List = 'list',
}

export type FileCommonProps = {
  file: File
  scale: number
  viewType: FileViewType
  isPresentational?: boolean
  isDragging?: boolean
  isLoading?: boolean
  isSelectionMode?: boolean
}
