import { useEffect } from 'react'
import { Divider, Text } from '@chakra-ui/react'
import classNames from 'classnames'
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
    <div className={classNames('flex', 'flex-col', 'gap-1.5')}>
      {items.map((u, index) => (
        <div key={u.id} className={classNames('flex', 'flex-col', 'gap-1.5')}>
          <Item upload={u} />
          {index !== items.length - 1 && <Divider />}
        </div>
      ))}
    </div>
  )
}

export default List
