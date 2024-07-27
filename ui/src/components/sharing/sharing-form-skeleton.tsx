// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Skeleton } from '@chakra-ui/react'
import cx from 'classnames'

const SharingFormSkeleton = () => (
  <div className={cx('flex', 'flex-col', 'gap-1.5')}>
    <Skeleton className={cx('rounded-xl', 'w-[40px]')} />
    <Skeleton className={cx('rounded-xl', 'w-[40px]')} />
    <Skeleton className={cx('rounded-xl', 'w-[40px]')} />
  </div>
)

export default SharingFormSkeleton
