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
  variables,
} from '@koupr/ui'
import classNames from 'classnames'
import {
  Upload,
  UploadDecorator,
  uploadRemoved,
} from '@/store/entities/uploads'
import { useAppDispatch } from '@/store/hook'

type ItemProps = {
  upload: Upload
}

const Item = ({ upload: uploadProp }: ItemProps) => {
  const dispatch = useAppDispatch()
  const upload = new UploadDecorator(uploadProp)
  const grayColor = useToken('colors', 'gray.500')
  const greenColor = useToken('colors', 'green')
  const redColor = useToken('colors', 'red')

  return (
    <div className={classNames('flex', 'flex-col', 'gap-1')}>
      <div
        className={classNames(
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
            color="black"
            size="20px"
          />
        )}
        {upload.isPending && (
          <div className={classNames(`text-[${grayColor}]`, 'shrink-0')}>
            <IconTime fontSize="21px" />
          </div>
        )}
        {upload.isSucceeded && (
          <div className={classNames(`text-[${greenColor}]`, 'shrink-0')}>
            <IconCheckCircleFill fontSize="22px" />
          </div>
        )}
        {upload.isFailed && (
          <div className={classNames(`text-[${redColor}]`, 'shrink-0')}>
            <IconAlertCircleFill fontSize="22px" />
          </div>
        )}
        <Text
          flexGrow={1}
          textOverflow="ellipsis"
          overflow="hidden"
          whiteSpace="nowrap"
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
            <AccordionButton p={variables.spacingXs} _hover={{ bg: 'red.50' }}>
              <div className={classNames('flex', 'flex-row', 'w-full')}>
                <Text color="red" flexGrow={1} textAlign="left">
                  Upload failed. Click to expand.
                </Text>
                <AccordionIcon color="red" />
              </div>
            </AccordionButton>
            <AccordionPanel p={variables.spacingXs}>
              {upload.error}
            </AccordionPanel>
          </AccordionItem>
        </Accordion>
      )}
    </div>
  )
}

export default Item
