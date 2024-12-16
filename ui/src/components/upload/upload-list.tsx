// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useEffect } from 'react'
import { SectionPlaceholder } from '@koupr/ui'
import cx from 'classnames'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { drawerDidClose } from '@/store/ui/uploads'
import UploadMenu from '../common/upload-menu'
import UploadItem from './upload-item'
import { queue } from './upload-worker'

const UploadList = () => {
  const items = useAppSelector((state) => state.entities.uploads.items)
  const dispatch = useAppDispatch()

  useEffect(() => {
    for (const upload of items) {
      if (
        queue.findIndex((e) => e.id === upload.id) !== -1 ||
        upload.completed
      ) {
        continue
      }
      queue.push(upload)
    }
    if (items.length === 0) {
      dispatch(drawerDidClose())
    }
  }, [items, dispatch])

  return (
    <>
      {items.length === 0 ? (
        <SectionPlaceholder
          text="There are no uploads."
          content={<UploadMenu />}
        />
      ) : (
        <div className={cx('flex', 'flex-col', 'gap-1.5')}>
          {items.map((u) => (
            <div key={u.id} className={cx('flex', 'flex-col', 'gap-1.5')}>
              <UploadItem upload={u} />
            </div>
          ))}
        </div>
      )}
    </>
  )
}

export default UploadList
