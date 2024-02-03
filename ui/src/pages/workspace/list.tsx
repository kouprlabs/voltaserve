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
import WorkspaceAPI, { SortOrder } from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import { CreateWorkspaceButton } from '@/components/top-bar/buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'
import { workspacePaginationStorage } from '@/infra/pagination'

const WorkspaceListPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: workspacePaginationStorage(),
  })
  const {
    data: list,
    error,
    mutate,
  } = WorkspaceAPI.useList(
    { query, page, size, sortOrder: SortOrder.Desc },
    swrConfig(),
  )

  useEffect(() => {
    mutate()
  }, [query, page, size, mutate])

  return (
    <>
      <Helmet>
        <title>Workspaces</title>
      </Helmet>
      <Stack
        direction="column"
        spacing={variables.spacing2Xl}
        pb={variables.spacing2Xl}
      >
        <Heading fontSize={variables.headingFontSize} pl={variables.spacingMd}>
          Workspaces
        </Heading>
        {!list && error && (
          <Center h="300px">
            <Text>Failed to load workspaces.</Text>
          </Center>
        )}
        {!list && !error && <SectionSpinner />}
        {list && list.data.length === 0 && !error ? (
          <Center h="300px">
            <VStack spacing={variables.spacing}>
              <Text>There are no workspaces.</Text>
              <CreateWorkspaceButton />
            </VStack>
          </Center>
        ) : null}
        {list && list.data.length > 0 && (
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Name</Th>
                <Th>Organization</Th>
                <Th>Permission</Th>
                <Th>Date</Th>
              </Tr>
            </Thead>
            <Tbody>
              {list.data.map((w) => (
                <Tr key={w.id}>
                  <Td>
                    <HStack spacing={variables.spacing}>
                      <Avatar
                        name={w.name}
                        size="sm"
                        width="40px"
                        height="40px"
                      />
                      <ChakraLink
                        as={Link}
                        to={`/workspace/${w.id}/file/${w.rootId}`}
                        textDecoration="none"
                      >
                        <Text>{w.name}</Text>
                      </ChakraLink>
                    </HStack>
                  </Td>
                  <Td>
                    <ChakraLink
                      as={Link}
                      to={`/organization/${w.organization.id}/member`}
                      textDecoration="none"
                    >
                      {w.organization.name}
                    </ChakraLink>
                  </Td>
                  <Td>
                    <Badge>{w.permission}</Badge>
                  </Td>
                  <Td>{prettyDate(w.createTime)}</Td>
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

export default WorkspaceListPage
