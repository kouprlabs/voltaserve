import cx from 'classnames'
import { Spinner } from './spinner'

interface SectionSpinnerProps {
  width?: string
  height?: string
}

const DEFAULT_WIDTH = '100%'
const DEFAULT_HEIGHT = '300px'

export const SectionSpinner = ({ width, height }: SectionSpinnerProps) => (
  <div
    className={cx('flex', 'items-center', 'justify-center')}
    style={{ width: width || DEFAULT_WIDTH, height: height || DEFAULT_HEIGHT }}
  >
    <Spinner />
  </div>
)
