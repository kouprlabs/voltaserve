import { ReactNode, useEffect } from 'react'
import { Center, Container, useToast } from '@chakra-ui/react'
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
    <Container
      h="100vh"
      position="relative"
      centerContent
      width={{ base: '100%', md: '400px' }}
    >
      <Center h="100%" w="100%">
        {children}
      </Center>
    </Container>
  )
}

export default FullLayout
