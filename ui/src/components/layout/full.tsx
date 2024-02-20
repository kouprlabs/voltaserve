import { ReactNode, useEffect } from 'react'
import { useToast } from '@chakra-ui/react'
import classNames from 'classnames'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { errorCleared } from '@/store/ui/error'

type FullLayoutProps = {
  children?: ReactNode
}

const FullLayout = ({ children }: FullLayoutProps) => {
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
      className={classNames(
        'relative',
        'flex',
        'flex-col',
        'items-center',
        'w-full',
        'h-[100vh]',
      )}
    >
      <div
        className={classNames(
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

export default FullLayout
