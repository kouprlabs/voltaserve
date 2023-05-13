import { useEffect } from 'react'
import { Divider, Stack, Text } from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { uploadsDrawerClosed } from '@/store/ui/uploads-drawer'
import Item from './item'
import { queue } from './worker'

const List = () => {
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
      dispatch(uploadsDrawerClosed())
    }
  }, [items, dispatch])

  if (items.length === 0) {
    return <Text>There are no uploads.</Text>
  }

  return (
    <Stack spacing={variables.spacing}>
      {items.map((u, i) => (
        <Stack key={u.id} spacing={variables.spacing}>
          <Item upload={u} />
          {i !== items.length - 1 && <Divider />}
        </Stack>
      ))}
    </Stack>
  )
}

export default List
