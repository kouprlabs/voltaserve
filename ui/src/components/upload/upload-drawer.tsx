// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useEffect } from 'react'
import { useDisclosure, Button } from '@chakra-ui/react'
import { AuxiliaryDrawer } from '@koupr/ui'
import cx from 'classnames'
import UploadList from '@/components/upload/upload-list'
import { IconClearAll, IconUpload } from '@/lib/components/icons'
import { completedUploadsCleared } from '@/store/entities/uploads'
import { useAppDispatch, useAppSelector } from '@/store/hook'

const UploadDrawer = () => {
  const dispatch = useAppDispatch()
  const hasPendingUploads = useAppSelector(
    (state) =>
      state.entities.uploads.items.filter((e) => !e.completed).length > 0,
  )
  const isDrawerOpen = useAppSelector((state) => state.ui.uploads.isDrawerOpen)
  const hasCompleted = useAppSelector(
    (state) =>
      state.entities.uploads.items.filter((e) => e.completed).length > 0,
  )
  const { isOpen, onOpen, onClose } = useDisclosure()

  useEffect(() => {
    if (isDrawerOpen) {
      onOpen()
    } else {
      onClose()
    }
  }, [isDrawerOpen, onOpen, onClose])

  const handleClearCompleted = useCallback(() => {
    dispatch(completedUploadsCleared())
  }, [dispatch])

  return (
    <AuxiliaryDrawer
      icon={<IconUpload />}
      isOpen={isOpen}
      onClose={onClose}
      onOpen={onOpen}
      hasBadge={hasPendingUploads}
      header="Uploads"
      body={<UploadList />}
      footer={
        <>
          {hasCompleted ? (
            <Button
              className={cx('w-full')}
              size="sm"
              leftIcon={<IconClearAll />}
              onClick={handleClearCompleted}
            >
              Clear Completed Items
            </Button>
          ) : null}
        </>
      }
    />
  )
}

export default UploadDrawer
