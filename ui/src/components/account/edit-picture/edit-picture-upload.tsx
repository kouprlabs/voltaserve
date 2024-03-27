import { ChangeEvent, useCallback, useRef, useState } from 'react'
import {
  IconButton,
  Image,
  Text,
  useColorModeValue,
  useToken,
} from '@chakra-ui/react'
import { IconEdit } from '@koupr/ui'
import cx from 'classnames'

export type EditPictureUploadProps = {
  name: string
  initialValue?: string
  disabled: boolean
  onChange: (event: any) => void
}

const EditPictureUpload = ({
  name,
  initialValue,
  onChange,
  disabled,
}: EditPictureUploadProps) => {
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
    <div className={cx('flex', 'flex-col', 'items-center', 'gap-1')}>
      <div
        className={cx(
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
          <div className={cx('relative', 'w-[400px]', 'h-[160px]')}>
            <Image
              src={src || initialValue || ''}
              className={cx(
                'rounded',
                'w-[400px]',
                'h-[160px]',
                'object-cover',
              )}
              alt="Account picture"
            />
            <IconButton
              icon={<IconEdit />}
              variant="solid-gray"
              className={cx('top-[10px', 'right-[5px]', 'absolute', 'z-[1000]')}
              aria-label=""
              disabled={disabled}
              onClick={handleEdit}
            />
          </div>
        ) : (
          <div
            className={cx(
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

export default EditPictureUpload
