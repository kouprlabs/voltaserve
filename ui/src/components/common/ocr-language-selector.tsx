import { useCallback, useEffect, useState } from 'react'
import {
  Button,
  Text,
  Center,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  useDisclosure,
  VStack,
  Stack,
  Table,
  Tr,
  Tbody,
  Td,
  HStack,
  Avatar,
  Radio,
  Box,
  useColorModeValue,
} from '@chakra-ui/react'
import { SectionSpinner, variables } from '@koupr/ui'
import OcrLanguageAPI, {
  OcrLanguage,
  SortOrder,
} from '@/client/api/ocr-language'
import { swrConfig } from '@/client/options'
import Pagination from '@/components/common/pagination'
import SearchInput from '@/components/common/search-input'

type OcrLanguageSelectorProps = {
  valueId?: string
  isDisabled?: boolean
  buttonLabel?: string
  onConfirm?: (ocrLanguage: OcrLanguage) => void
}

const OcrLanguageSelector = ({
  valueId,
  isDisabled,
  buttonLabel,
  onConfirm,
}: OcrLanguageSelectorProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState('')
  const [selected, setSelected] = useState<OcrLanguage>()
  const [confirmed, setConfirmed] = useState<OcrLanguage>()
  const {
    data: list,
    error,
    mutate,
  } = OcrLanguageAPI.useList(
    { query, page, size: 5, sortOrder: SortOrder.Desc },
    swrConfig(),
  )
  const selectionColor = useColorModeValue('gray.100', 'gray.600')
  const dimmedButtonLabelColor = useColorModeValue('gray.500', 'gray.500')
  const normalButtonLabelColor = useColorModeValue('black', 'white')
  const [isFetching, setIsFetching] = useState(false)

  useEffect(() => {
    mutate()
  }, [page, query, mutate])

  useEffect(() => {
    if (!isOpen) {
      setPage(1)
      setSelected(undefined)
      setQuery('')
    }
  }, [isOpen])

  useEffect(() => {
    if (!list || !valueId) {
      return
    }
    async function fetch() {
      setIsFetching(true)
      try {
        const result = await OcrLanguageAPI.list({ query: valueId })
        if (result.size > 0) {
          setConfirmed(result.data[0])
        }
      } finally {
        setIsFetching(false)
      }
    }
    fetch()
  }, [list, valueId])

  const handleConfirm = useCallback(() => {
    if (selected) {
      setConfirmed(selected)
      onConfirm?.(selected)
      onClose()
    }
  }, [selected, onConfirm, onClose])

  return (
    <>
      <Button
        variant="outline"
        w="100%"
        isDisabled={isDisabled}
        isLoading={isFetching}
        onClick={onOpen}
        color={confirmed ? normalButtonLabelColor : dimmedButtonLabelColor}
      >
        {confirmed ? confirmed.name : buttonLabel ?? 'Select OCR Language'}
      </Button>
      <Modal
        size="xl"
        isOpen={isOpen}
        onClose={onClose}
        closeOnOverlayClick={false}
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>{buttonLabel ?? 'Select OCR Language'}</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Stack direction="column" spacing={variables.spacing}>
              <SearchInput
                query={query}
                onChange={(value) => setQuery(value)}
              />
              {!list && error && (
                <Center h="300px">
                  <Text>Failed to load OCR languages.</Text>
                </Center>
              )}
              {!list && !error && <SectionSpinner />}
              {list && list.data.length === 0 && (
                <Center h="300px">
                  <VStack spacing={variables.spacing}>
                    <Text>There are no organizations.</Text>
                  </VStack>
                </Center>
              )}
              {list && list.data.length > 0 && (
                <Table variant="simple" size="sm">
                  <colgroup>
                    <col style={{ width: '40px' }} />
                    <col style={{ width: 'auto' }} />
                  </colgroup>
                  <Tbody>
                    {list.data.map((ol) => (
                      <Tr
                        key={ol.id}
                        cursor="pointer"
                        bg={selected?.id === ol.id ? selectionColor : 'auto'}
                        onClick={() => setSelected(ol)}
                      >
                        <Td px={variables.spacingXs} textAlign="center">
                          <Radio size="md" isChecked={selected?.id === ol.id} />
                        </Td>
                        <Td px={variables.spacingXs}>
                          <HStack spacing={variables.spacing}>
                            <Avatar
                              name={ol.name}
                              size="sm"
                              width="40px"
                              height="40px"
                            />
                            <Text fontSize={variables.bodyFontSize}>
                              {ol.name}
                            </Text>
                          </HStack>
                        </Td>
                      </Tr>
                    ))}
                  </Tbody>
                </Table>
              )}
              {list && (
                <Box alignSelf="end">
                  {list.totalPages > 1 ? (
                    <Pagination
                      size="md"
                      maxButtons={3}
                      page={page}
                      totalPages={list.totalPages}
                      onPageChange={(value) => setPage(value)}
                    />
                  ) : null}
                </Box>
              )}
            </Stack>
          </ModalBody>
          <ModalFooter>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              mr={variables.spacingSm}
              onClick={onClose}
            >
              Cancel
            </Button>
            <Button
              variant="solid"
              colorScheme="blue"
              isDisabled={!selected}
              onClick={handleConfirm}
            >
              Confirm
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  )
}

export default OcrLanguageSelector
