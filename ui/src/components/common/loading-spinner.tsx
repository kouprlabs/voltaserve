import { Center, Spinner } from '@chakra-ui/react'

type LoadingSpinnerProps = {
  width?: string
  height?: string
}

const DEFAULT_WIDTH = '100%'
const DEFAULT_HEIGHT = '300px'

const LoadingSpinner = ({ width, height }: LoadingSpinnerProps) => (
  <Center w={width || DEFAULT_WIDTH} h={height || DEFAULT_HEIGHT}>
    <Spinner size="sm" thickness="4px" />
  </Center>
)

export default LoadingSpinner
