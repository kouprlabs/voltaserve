import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  CircularProgress,
  IconButton,
  Text,
  useToken,
} from '@chakra-ui/react'
import {
  IconClose,
  IconTime,
  IconCheckCircleFill,
  IconAlertCircleFill,
} from '@koupr/ui'
import cx from 'classnames'
import {
  Upload,
  UploadDecorator,
  uploadRemoved,
} from '@/store/entities/uploads'
import { useAppDispatch } from '@/store/hook'

export type UploadItemProps = {
  upload: Upload
}

const UploadItem = ({ upload: uploadProp }: UploadItemProps) => {
  const dispatch = useAppDispatch()
  const upload = new UploadDecorator(uploadProp)
  const grayColor = useToken('colors', 'gray.500')
  const greenColor = useToken('colors', 'green')
  const redColor = useToken('colors', 'red')

  return (
    <div className={cx('flex', 'flex-col', 'gap-1')}>
      <div
        className={cx(
          'flex',
          'flex-row',
          'items-center',
          'gap-0.5',
          'justify-between',
          'h-2.5',
        )}
      >
        {upload.isProgressing && (
          <CircularProgress
            value={upload.progress}
            max={100}
            isIndeterminate={upload.progress === 100 && !upload.error}
            className={cx('text-black')}
            size="20px"
          />
        )}
        {upload.isPending && (
          <div className={cx('shrink-0')} style={{ color: grayColor }}>
            <IconTime fontSize="21px" />
          </div>
        )}
        {upload.isSucceeded && (
          <div className={cx('shrink-0')} style={{ color: greenColor }}>
            <IconCheckCircleFill fontSize="22px" />
          </div>
        )}
        {upload.isFailed && (
          <div className={cx('shrink-0')} style={{ color: redColor }}>
            <IconAlertCircleFill fontSize="22px" />
          </div>
        )}
        <Text
          className={cx(
            'grow',
            'text-ellipsis',
            'overflow-hidden',
            'whitespace-nowrap',
          )}
        >
          {upload.file.name}
        </Text>
        <IconButton
          icon={<IconClose />}
          size="xs"
          variant="outline"
          colorScheme={upload.isProgressing ? 'red' : 'gray'}
          aria-label=""
          onClick={() => {
            upload.request?.abort()
            dispatch(uploadRemoved(upload.id))
          }}
        />
      </div>
      {upload.isFailed && (
        <Accordion allowMultiple>
          <AccordionItem border="none">
            <AccordionButton className={cx('p-0.5')} _hover={{ bg: 'red.50' }}>
              <div className={cx('flex', 'flex-row', 'w-full')}>
                <Text className={cx('text-red-500', 'text-left', 'grow')}>
                  Upload failed. Click to expand.
                </Text>
                <AccordionIcon className={cx('text-red-500')} />
              </div>
            </AccordionButton>
            <AccordionPanel className={cx('p-0.5')}>
              {upload.error}
            </AccordionPanel>
          </AccordionItem>
        </Accordion>
      )}
    </div>
  )
}

export default UploadItem
