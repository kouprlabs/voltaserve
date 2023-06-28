import {
  useCallback,
  useEffect,
  useState,
  ChangeEvent,
  KeyboardEvent,
} from 'react'
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
  InputGroup,
  InputLeftElement,
  Icon,
  InputRightElement,
  IconButton,
  Input,
  useColorModeValue,
} from '@chakra-ui/react'
import { SectionSpinner, variables } from '@koupr/ui'
import { IconClose, IconSearch } from '@koupr/ui'
import OrganizationAPI, {
  Organization,
  SortOrder,
} from '@/client/api/organization'
import { swrConfig } from '@/client/options'
import Pagination from '@/components/common/pagination'

type SearchProps = {
  query?: string
  onChange?: (value: string) => void
}

const Search = ({ query, onChange }: SearchProps) => {
  const [draft, setDraft] = useState('')
  const [text, setText] = useState('')
  const [isFocused, setIsFocused] = useState(false)

  useEffect(() => {
    if (query !== draft) {
      setDraft(query || '')
    }
  }, [query, draft])

  useEffect(() => {
    onChange?.(text)
  }, [text, onChange])

  const handleClear = useCallback(() => {
    setDraft('')
    setText('')
  }, [])

  const handleChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    setDraft(event.target.value || '')
  }, [])

  const handleSearch = useCallback((value: string) => {
    setText(value)
  }, [])

  const handleKeyDown = useCallback(
    (event: KeyboardEvent<HTMLInputElement>, value: string) => {
      if (event.key === 'Enter') {
        handleSearch(value)
      }
    },
    [handleSearch]
  )

  return (
    <HStack>
      <InputGroup>
        <InputLeftElement pointerEvents="none">
          <Icon as={IconSearch} color="gray.300" />
        </InputLeftElement>
        <Input
          value={draft}
          placeholder={draft || 'Search'}
          variant="filled"
          onKeyDown={(event) => handleKeyDown(event, draft)}
          onChange={handleChange}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
        />
        {draft && (
          <InputRightElement>
            <IconButton
              icon={<IconClose />}
              onClick={handleClear}
              size="xs"
              aria-label="Clear"
            />
          </InputRightElement>
        )}
      </InputGroup>
      {draft || (isFocused && draft) ? (
        <Button onClick={() => handleSearch(draft)} isDisabled={!draft}>
          Search
        </Button>
      ) : null}
    </HStack>
  )
}

type OrganizationSelectorProps = {
  onConfirm?: (organinzation: Organization) => void
}

const OrganizationSelector = ({ onConfirm }: OrganizationSelectorProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState('')
  const [selected, setSelected] = useState<Organization>()
  const [confirmed, setConfirmed] = useState<Organization>()
  const {
    data: list,
    error,
    mutate,
  } = OrganizationAPI.useList(
    { query, page, size: 5, sortOrder: SortOrder.Desc },
    swrConfig()
  )
  const selectionColor = useColorModeValue('gray.100', 'gray.600')
  const dimmedButtonLabelColor = useColorModeValue('gray.500', 'gray.500')
  const normalButtonLabelColor = useColorModeValue('black', 'white')

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
        onClick={onOpen}
        color={confirmed ? normalButtonLabelColor : dimmedButtonLabelColor}
      >
        {confirmed ? confirmed.name : 'Select Organization'}
      </Button>
      <Modal
        size="xl"
        isOpen={isOpen}
        onClose={onClose}
        closeOnOverlayClick={false}
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Select Organization</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Stack direction="column" spacing={variables.spacing}>
              <Search query={query} onChange={(value) => setQuery(value)} />
              {!list && error && (
                <Center h="300px">
                  <Text>Failed to load organizations.</Text>
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
                    <col style={{ width: '50px' }} />
                    <col style={{ width: 'auto' }} />
                  </colgroup>
                  <Tbody>
                    {list.data.map((o: Organization) => (
                      <Tr
                        key={o.id}
                        cursor="pointer"
                        bg={selected?.id === o.id ? selectionColor : 'auto'}
                        onClick={() => setSelected(o)}
                      >
                        <Td>
                          <Radio size="md" isChecked={selected?.id === o.id} />
                        </Td>
                        <Td>
                          <HStack spacing={variables.spacing}>
                            <Avatar
                              name={o.name}
                              size="sm"
                              width="40px"
                              height="40px"
                            />
                            <Text fontSize="14px">{o.name}</Text>
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

export default OrganizationSelector
