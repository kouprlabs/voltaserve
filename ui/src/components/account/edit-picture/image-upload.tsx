import React, { ChangeEvent, useCallback, useRef, useState } from 'react'
import {
  Box,
  Center,
  IconButton,
  Image,
  Text,
  useColorModeValue,
  VStack,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
import { IconEdit } from '@koupr/ui'

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
  const blueColor = useColorModeValue('blue.600', 'blue.200')

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
    [onChange]
  )

  const handleEdit = useCallback(() => {
    if (!disabled && hiddenInput.current) {
      hiddenInput.current.click()
    }
  }, [disabled, hiddenInput])

  return (
    <VStack spacing={variables.spacingSm}>
      <Center className="rounded border border-dashed" borderColor={blueColor}>
        {src || initialValue ? (
          <Box position="relative" width="400px" height="160px">
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
          </Box>
        ) : (
          <Center
            className="cursor-pointer"
            width="400px"
            height="160px"
            onClick={handleEdit}
          >
            <Text color={blueColor}>Browse</Text>
          </Center>
        )}
      </Center>
      <input
        ref={hiddenInput}
        className="hidden"
        type="file"
        name={name}
        onChange={handleFileChange}
      />
    </VStack>
  )
}

export default ImageUpload
