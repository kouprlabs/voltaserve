import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box,
  CircularProgress,
  HStack,
  IconButton,
  Stack,
  Text,
} from '@chakra-ui/react'
import {
  Upload,
  UploadDecorator,
  uploadRemoved,
} from '@/store/entities/uploads'
import { useAppDispatch } from '@/store/hook'
import {
  IconClose,
  IconTime,
  IconCheckCircleFill,
  IconAlertCircleFill,
} from '@/components/common/icon'
import variables from '@/theme/variables'

type ItemProps = {
  upload: Upload
}

const FileUploadItem = ({ upload: uploadProp }: ItemProps) => {
  const dispatch = useAppDispatch()
  const upload = new UploadDecorator(uploadProp)
  return (
    <Stack spacing={variables.spacingSm}>
      <HStack justifyContent="space-between" h="25px">
        {upload.isProgressing && (
          <CircularProgress
            value={upload.progress}
            max={100}
            isIndeterminate={upload.progress === 100 && !upload.error}
            size="20px"
          />
        )}
        {upload.isPending && (
          <Box color="gray.500" flexShrink={0}>
            <IconTime fontSize="21px" />
          </Box>
        )}
        {upload.isSucceeded && (
          <Box color="green" flexShrink={0}>
            <IconCheckCircleFill fontSize="22px" />
          </Box>
        )}
        {upload.isFailed && (
          <Box color="red" flexShrink={0}>
            <IconAlertCircleFill fontSize="22px" />
          </Box>
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
      </HStack>
      {upload.isFailed && (
        <Accordion allowMultiple>
          <AccordionItem border="none">
            <AccordionButton p={variables.spacingXs} _hover={{ bg: 'red.50' }}>
              <Stack direction="row" w="100%">
                <Text color="red" flexGrow={1} textAlign="left">
                  Upload failed. Click to expand.
                </Text>
                <AccordionIcon color="red" />
              </Stack>
            </AccordionButton>
            <AccordionPanel p={variables.spacingXs}>
              {upload.error}
            </AccordionPanel>
          </AccordionItem>
        </Accordion>
      )}
    </Stack>
  )
}

export default FileUploadItem
