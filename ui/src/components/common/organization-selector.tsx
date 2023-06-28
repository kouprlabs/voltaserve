import { useCallback, useEffect, useState, ChangeEvent } from 'react'
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
import { CreateOrganizationButton } from '@/components/top-bar/buttons'

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

  const handleClear = useCallback(() => {
    setQuery('')
  }, [])

  const handleChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    setQuery(event.target.value || '')
  }, [])

  const handleConfirm = useCallback(() => {
    if (selected) {
      setConfirmed(selected)
      onConfirm?.(selected)
      onClose()
    }
  }, [selected, onConfirm, onClose])

  return (
    <>
      <Button variant="outline" w="100%" onClick={onOpen}>
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
              <HStack>
                <InputGroup>
                  <InputLeftElement pointerEvents="none">
                    <Icon as={IconSearch} color="gray.300" />
                  </InputLeftElement>
                  <Input
                    value={query}
                    placeholder={query || 'Search'}
                    variant="filled"
                    onChange={handleChange}
                  />
                  {query && (
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
              </HStack>
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
                    <CreateOrganizationButton />
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
