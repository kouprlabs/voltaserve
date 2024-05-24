import { ChangeEvent, useCallback, useRef, useState } from 'react'
import { IconButton, Image } from '@chakra-ui/react'
import cx from 'classnames'
import { IconEdit } from '@/lib'

export type EditPictureUploadProps = {
  name: string
  initialValue?: string
  disabled: boolean
  onChange: (event: ChangeEvent<HTMLInputElement>) => void
}

const EditPictureUpload = ({
  name,
  initialValue,
  onChange,
  disabled,
}: EditPictureUploadProps) => {
  const [src, setSrc] = useState<string>()
  const hiddenInput = useRef<HTMLInputElement>(null)

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
          'border-blue-600',
          'dark:border-blue-200',
        )}
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
              className={cx(
                'top-[10px]',
                'right-[5px]',
                'absolute',
                'z-[1000]',
              )}
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
            <span className={cx('text-blue-600', 'dark:text-blue-200')}>
              Browse
            </span>
          </div>
        )}
      </div>
      <input
        ref={hiddenInput}
        className={cx('hidden')}
        type="file"
        name={name}
        onChange={handleFileChange}
      />
    </div>
  )
}

export default EditPictureUpload
