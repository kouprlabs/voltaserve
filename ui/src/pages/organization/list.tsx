import { useEffect } from 'react'
import {
  Link,
  useLocation,
  useNavigate,
  useSearchParams,
} from 'react-router-dom'
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
import {
  SectionSpinner,
  PagePagination,
  variables,
  usePagePagination,
} from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI, { SortOrder } from '@/client/api/organization'
import { swrConfig } from '@/client/options'
import { CreateOrganizationButton } from '@/components/top-bar/buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'
import { organizationPaginationStorage } from '@/infra/pagination'

const OrganizationListPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: organizationPaginationStorage(),
  })
  const {
    data: list,
    error,
    mutate,
  } = OrganizationAPI.useList(
    { query, page, size, sortOrder: SortOrder.Desc },
    swrConfig(),
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
        <Heading fontSize={variables.headingFontSize} pl={variables.spacingMd}>
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
          <PagePagination
            style={{ alignSelf: 'end' }}
            totalElements={list.totalElements}
            totalPages={list.totalPages}
            page={page}
            size={size}
            steps={steps}
            setPage={setPage}
            setSize={setSize}
          />
        )}
      </Stack>
    </>
  )
}

export default OrganizationListPage
