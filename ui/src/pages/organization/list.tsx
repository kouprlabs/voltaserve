import { useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import {
  Heading,
  Link as ChakraLink,
  Stack,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  Text,
  HStack,
  Center,
  VStack,
  Avatar,
  Badge,
} from '@chakra-ui/react'
import { SectionSpinner, variables } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI, { SortOrder } from '@/client/api/organization'
import { swrConfig } from '@/client/options'
import PagePagination, {
  usePagePagination,
} from '@/components/common/page-pagination'
import { CreateOrganizationButton } from '@/components/top-bar/buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'

const OrganizationListPage = () => {
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, onPageChange, onSizeChange } = usePagePagination({
    localStoragePrefix: 'voltaserve',
    localStorageNamespace: 'organization',
  })
  const {
    data: list,
    error,
    mutate,
  } = OrganizationAPI.useList(
    { query, page, size, sortOrder: SortOrder.Desc },
    swrConfig()
  )

  useEffect(() => {
    mutate()
  }, [query, page, size, mutate])

  return (
    <>
      <Helmet>
        <title>Organizations</title>
      </Helmet>
      <Stack
        direction="column"
        spacing={variables.spacing2Xl}
        pb={variables.spacing2Xl}
      >
        <Heading size="lg" pl={variables.spacingMd}>
          Organizations
        </Heading>
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
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Name</Th>
                <Th>Permission</Th>
                <Th>Date</Th>
              </Tr>
            </Thead>
            <Tbody>
              {list.data.map((o) => (
                <Tr key={o.id}>
                  <Td>
                    <HStack spacing={variables.spacing}>
                      <Avatar
                        name={o.name}
                        size="sm"
                        width="40px"
                        height="40px"
                      />
                      <ChakraLink
                        as={Link}
                        to={`/organization/${o.id}/member`}
                        textDecoration="none"
                      >
                        <Text>{o.name}</Text>
                      </ChakraLink>
                    </HStack>
                  </Td>
                  <Td>
                    <Badge>{o.permission}</Badge>
                  </Td>
                  <Td>{prettyDate(o.createTime)}</Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        )}
        {list && (
          <HStack alignSelf="end">
            <PagePagination
              totalPages={list.totalPages}
              page={page}
              size={size}
              onPageChange={onPageChange}
              onSizeChange={onSizeChange}
            />
          </HStack>
        )}
      </Stack>
    </>
  )
}

export default OrganizationListPage
