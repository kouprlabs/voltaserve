// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useMemo } from 'react'
import { Button, Card, CardBody, CardFooter, Text } from '@chakra-ui/react'
import { IconBolt, IconDelete } from '@koupr/ui'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import MosaicAPI from '@/client/api/mosaic'
import {
  geEditorPermission,
  geOwnerPermission,
  NONE_PERMISSION,
} from '@/client/api/permission'
import TaskAPI from '@/client/api/task'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/mosaic'

const MosaicOverviewSettings = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const mutateTaskCount = useAppSelector((state) => state.ui.tasks.mutateCount)
  const { data: info, mutate: mutateInfo } = MosaicAPI.useGetInfo(
    id,
    swrConfig(),
  )
  const { data: file, mutate: mutateFile } = FileAPI.useGet(id, swrConfig())

  const handleUpdate = useCallback(async () => {
    if (id) {
      await MosaicAPI.create(id)
      await mutateFile(await FileAPI.get(id))
      await mutateInfo(await MosaicAPI.getInfo(id))
      await mutateFiles?.()
      await mutateTaskCount?.(await TaskAPI.getCount())
      dispatch(modalDidClose())
    }
  }, [id, mutateFile, mutateFiles, mutateTaskCount, mutateInfo, dispatch])

  const handleDelete = useCallback(async () => {
    if (id) {
      await MosaicAPI.delete(id)
      await mutateFile(await FileAPI.get(id))
      await mutateInfo(await MosaicAPI.getInfo(id))
      await mutateFiles?.()
      await mutateTaskCount?.(await TaskAPI.getCount())
      dispatch(modalDidClose())
    }
  }, [id, mutateFiles, mutateTaskCount, mutateInfo, dispatch])

  const canCreate = useMemo(() => {
    return !!(
      !file?.snapshot?.task?.isPending &&
      info?.isOutdated &&
      geEditorPermission(file?.permission ?? NONE_PERMISSION)
    )
  }, [info, file])

  const canDelete = useMemo(() => {
    return (
      !file?.snapshot?.task?.isPending &&
      !info?.isOutdated &&
      geOwnerPermission(file?.permission ?? NONE_PERMISSION)
    )
  }, [info, file])

  if (!file || !info) {
    return null
  }

  return (
    <div className={cx('flex', 'flex-row', 'items-stretch', 'gap-1.5')}>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>Create a mosaic for the active snapshot.</Text>
        </CardBody>
        <CardFooter>
          <Button
            leftIcon={<IconBolt />}
            isDisabled={!canCreate}
            onClick={handleUpdate}
          >
            Create Mosaic
          </Button>
        </CardFooter>
      </Card>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>Delete mosaic from the active snapshot.</Text>
        </CardBody>
        <CardFooter>
          <Button
            colorScheme="red"
            leftIcon={<IconDelete />}
            isDisabled={!canDelete}
            onClick={handleDelete}
          >
            Delete Mosaic
          </Button>
        </CardFooter>
      </Card>
    </div>
  )
}

export default MosaicOverviewSettings
