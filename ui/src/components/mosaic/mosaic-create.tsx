// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback } from 'react'
import { Button, ModalBody, ModalFooter } from '@chakra-ui/react'
import cx from 'classnames'
import { MosaicAPI } from '@/client/api/mosaic'
import { TaskAPI } from '@/client/api/task'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/mosaic'

const MosaicCreate = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const mutateTasks = useAppSelector((state) => state.ui.tasks.mutateList)

  const handleCreate = useCallback(async () => {
    if (id) {
      await MosaicAPI.create(id, false)
      await mutateFiles?.()
      await mutateTasks?.(await TaskAPI.list())
      dispatch(modalDidClose())
    }
  }, [id, mutateFiles, mutateTasks, dispatch])

  if (!id) {
    return null
  }

  return (
    <>
      <ModalBody>
        <div
          className={cx(
            'flex',
            'flex-col',
            'items-center',
            'justify-center',
            'gap-1.5',
          )}
        >
          <p>
            Create a mosaic to enhance view performance of a large image by
            splitting it into smaller, manageable tiles. This makes browsing a
            high-resolution image faster and more efficient.
          </p>
        </div>
      </ModalBody>
      <ModalFooter>
        <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
          <Button
            type="button"
            variant="outline"
            colorScheme="blue"
            onClick={() => dispatch(modalDidClose())}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="solid"
            colorScheme="blue"
            onClick={handleCreate}
          >
            Create
          </Button>
        </div>
      </ModalFooter>
    </>
  )
}

export default MosaicCreate
