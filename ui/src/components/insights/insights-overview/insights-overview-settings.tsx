// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useMemo } from 'react'
import { Button, Card, CardBody, CardFooter, Text } from '@chakra-ui/react'
import { IconBolt, IconDelete, SectionError, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import { EntityAPI } from '@/client/api/entity'
import { FileAPI } from '@/client/api/file'
import {
  geEditorPermission,
  geOwnerPermission,
  NONE_PERMISSION,
} from '@/client/api/permission'
import { TaskAPI, TaskStatus } from '@/client/api/task'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/insights'

const InsightsOverviewSettings = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const mutateTaskCount = useAppSelector((state) => state.ui.tasks.mutateCount)
  const {
    data: file,
    error: fileError,
    isLoading: fileIsLoading,
    mutate: mutateFile,
  } = FileAPI.useGet(id, swrConfig())
  const canCollect = useMemo(() => {
    return (
      file?.snapshot?.task?.status !== TaskStatus.Running &&
      geEditorPermission(file?.permission ?? NONE_PERMISSION)
    )
  }, [file])
  const canDelete = useMemo(() => {
    return (
      file?.snapshot?.task?.status !== TaskStatus.Running &&
      geOwnerPermission(file?.permission ?? NONE_PERMISSION)
    )
  }, [file])
  const fileIsReady = file && !fileError

  const handleCreate = useCallback(async () => {
    if (file && file.snapshot?.language) {
      await EntityAPI.create(file.id, { language: file.snapshot?.language })
      await mutateFile(await FileAPI.get(file.id))
      await mutateFiles?.()
      await mutateTaskCount?.(await TaskAPI.getCount())
      dispatch(modalDidClose())
    }
  }, [file, mutateFiles, mutateTaskCount, dispatch])

  const handleDelete = useCallback(async () => {
    if (id) {
      await EntityAPI.delete(id)
      await mutateFile(await FileAPI.get(id))
      await mutateFiles?.()
      await mutateTaskCount?.(await TaskAPI.getCount())
      dispatch(modalDidClose())
    }
  }, [id, mutateFile, mutateFiles, mutateTaskCount, dispatch])

  return (
    <>
      {fileIsLoading ? <SectionSpinner /> : null}
      {fileError ? <SectionError text={errorToString(fileError)} /> : null}
      {fileIsReady ? (
        <div className={cx('flex', 'flex-col', 'items-stretch', 'gap-1.5')}>
          {file.snapshot?.capabilities.entities ? (
            <Card size="md" variant="outline">
              <CardBody>
                <Text>Delete entities from the active snapshot.</Text>
              </CardBody>
              <CardFooter>
                <Button
                  className={cx('w-full')}
                  colorScheme="red"
                  leftIcon={<IconDelete />}
                  isDisabled={!canDelete}
                  onClick={handleDelete}
                >
                  Delete Entities
                </Button>
              </CardFooter>
            </Card>
          ) : (
            <Card size="md" variant="outline">
              <CardBody>
                <Text>Collect entities for the active snapshot.</Text>
              </CardBody>
              <CardFooter>
                <Button
                  className={cx('w-full')}
                  leftIcon={<IconBolt />}
                  isDisabled={!canCollect}
                  onClick={handleCreate}
                >
                  Collect Entities
                </Button>
              </CardFooter>
            </Card>
          )}
        </div>
      ) : null}
    </>
  )
}

export default InsightsOverviewSettings
