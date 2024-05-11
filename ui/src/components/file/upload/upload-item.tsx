import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  CircularProgress,
  IconButton,
} from '@chakra-ui/react'
import cx from 'classnames'
import { IconClose, IconSchedule, IconCheckCircle, IconError } from '@/lib'
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
          <IconSchedule className={cx('shrink-0', 'text-gray-500')} />
        )}
        {upload.isSucceeded && (
          <IconCheckCircle
            className={cx('shrink-0', 'text-green-500')}
            filled={true}
          />
        )}
        {upload.isFailed && (
          <div className={cx('shrink-0', 'text-red-500')}>
            <IconError filled={true} />
          </div>
        )}
        <span
          className={cx(
            'grow',
            'text-ellipsis',
            'overflow-hidden',
            'whitespace-nowrap',
          )}
        >
          {upload.blob.name}
        </span>
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
          <AccordionItem className={cx('border-none')}>
            <AccordionButton className={cx('p-0.5', 'hover:bg-red-50')}>
              <div className={cx('flex', 'flex-row', 'w-full')}>
                <span className={cx('text-red-500', 'text-left', 'grow')}>
                  Upload failed. Click to expand.
                </span>
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
