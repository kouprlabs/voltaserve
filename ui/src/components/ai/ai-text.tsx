import { Button } from '@chakra-ui/react'
import cx from 'classnames'
import { IconDownload } from '@/lib'

const AIText = () => {
  return (
    <div
      className={cx(
        'flex',
        'flex-col',
        'items-center',
        'justify-center',
        'gap-1',
      )}
    >
      <Button type="button" leftIcon={<IconDownload />}>
        Download Text
      </Button>
      <Button type="button" leftIcon={<IconDownload />}>
        Download PDF
      </Button>
    </div>
  )
}

export default AIText
