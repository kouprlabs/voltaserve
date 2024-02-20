import { ChangeEvent, useCallback, useRef, useState } from 'react'
import {
  IconButton,
  Image,
  Text,
  useColorModeValue,
  useToken,
} from '@chakra-ui/react'
import { IconEdit } from '@koupr/ui'
import classNames from 'classnames'

type ImageUploadProps = {
  name: string
  initialValue?: string
  disabled: boolean
  onChange: (event: any) => void
}

const ImageUpload = ({
  name,
  initialValue,
  onChange,
  disabled,
}: ImageUploadProps) => {
  const [src, setSrc] = useState<string>()
  const hiddenInput = useRef<HTMLInputElement>(null)
  const blueColor = useToken(
    'colors',
    useColorModeValue('blue.600', 'blue.200'),
  )

  const handleFileChange = useCallback(
    (changeEvent: ChangeEvent<HTMLInputElement>) => {
      if (!changeEvent.target.files || changeEvent.target.files.length === 0) {
        return
      }
      const file = changeEvent.target.files.item(0)
      if (!file) {
        return
      }
      const reader = new FileReader()
      reader.onload = (readerEvent: ProgressEvent<FileReader>) => {
        if (
          readerEvent.target?.result &&
          typeof readerEvent.target.result === 'string'
        ) {
          setSrc(readerEvent.target.result)
        }
      }
      reader.readAsDataURL(file)
      onChange?.(changeEvent)
    },
    [onChange],
  )

  const handleEdit = useCallback(() => {
    if (!disabled && hiddenInput.current) {
      hiddenInput.current.click()
    }
  }, [disabled, hiddenInput])

  return (
    <div className={classNames('flex', 'flex-col', 'items-center', 'gap-1')}>
      <div
        className={classNames(
          'flex',
          'items-center',
          'justify-center',
          'rounded',
          'border',
          'border-dashed',
        )}
        style={{ borderColor: blueColor }}
      >
        {src || initialValue ? (
          <div className={classNames('relative', 'w-[400px]', 'h-[160px]')}>
            <Image
              w="400px"
              h="160px"
              objectFit="cover"
              src={src || initialValue || ''}
              className="rounded"
              alt="Account picture"
            />
            <IconButton
              icon={<IconEdit />}
              variant="solid-gray"
              top="10px"
              right="5px"
              position="absolute"
              zIndex={1000}
              aria-label=""
              disabled={disabled}
              onClick={handleEdit}
            />
          </div>
        ) : (
          <div
            className={classNames(
              'flex',
              'items-center',
              'justify-center',
              'cursor-pointer',
              'w-[400px]',
              'h-[160px]',
            )}
            onClick={handleEdit}
          >
            <Text color={blueColor}>Browse</Text>
          </div>
        )}
      </div>
      <input
        ref={hiddenInput}
        className="hidden"
        type="file"
        name={name}
        onChange={handleFileChange}
      />
    </div>
  )
}

export default ImageUpload
