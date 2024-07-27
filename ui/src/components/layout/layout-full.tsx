// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { ReactNode, useEffect } from 'react'
import { useToast } from '@chakra-ui/react'
import cx from 'classnames'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { errorCleared } from '@/store/ui/error'

export type LayoutFullProps = {
  children?: ReactNode
}

const LayoutFull = ({ children }: LayoutFullProps) => {
  const toast = useToast()
  const error = useAppSelector((state) => state.ui.error.value)
  const dispatch = useAppDispatch()

  useEffect(() => {
    if (error) {
      toast({
        title: error,
        status: 'error',
        isClosable: true,
      })
      dispatch(errorCleared())
    }
  }, [error, toast, dispatch])

  return (
    <div
      className={cx(
        'relative',
        'flex',
        'flex-col',
        'items-center',
        'w-full',
        'h-[100vh]',
      )}
    >
      <div
        className={cx(
          'flex',
          'items-center',
          'justify-center',
          'w-full',
          'md:w-[400px]',
          'h-full',
          'p-2',
          'md:p-0',
        )}
      >
        {children}
      </div>
    </div>
  )
}

export default LayoutFull
