import { useCallback, useEffect, useState } from 'react'
import {
  Link,
  useSearchParams,
  useNavigate,
  useLocation,
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
  Select,
} from '@chakra-ui/react'
import { SectionSpinner, variables } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import WorkspaceAPI, { Workspace } from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import Pagination from '@/components/common/pagination'
import { CreateWorkspaceButton } from '@/components/top-bar/buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'

const WorkspaceListPage = () => {
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const navigate = useNavigate()
  const location = useLocation()
  const queryParams = new URLSearchParams(location.search)
  const currentPage = Number(queryParams.get('page')) || 1
  const [size, setSize] = useState(5)
  const {
    data: list,
    error,
    mutate,
  } = WorkspaceAPI.useList(
    {
      query,
      page: currentPage,
      size,
    },
    swrConfig()
  )

  useEffect(() => {
    mutate()
  }, [currentPage])

  useEffect(() => {
    if (!queryParams.has('page')) {
      queryParams.set('page', '1')
      navigate({ search: `?${queryParams.toString()}` })
    }
  }, [queryParams, navigate])

  const handlePageChange = useCallback(
    (page: number) => {
      queryParams.set('page', String(page))
      navigate({ search: `?${queryParams.toString()}` })
    },
    [queryParams, navigate]
  )

  useEffect(() => {
    mutate()
  }, [query, mutate])

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
        <Heading size="lg" pl={variables.spacingMd}>
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
              {list.data.map((w: Workspace) => (
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
                  <Td>{prettyDate(w.updateTime || w.createTime)}</Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        )}
        <HStack alignSelf="end">
          {list && list.totalPages > 1 ? (
            <Pagination
              page={list.page}
              totalPages={list.totalPages}
              onPageChange={handlePageChange}
            />
          ) : null}
          <Select
            defaultValue={size}
            onChange={(event) => {
              setSize(parseInt(event.target.value))
              mutate()
            }}
          >
            <option value="5">5 items</option>
            <option value="10">10 items</option>
            <option value="20">20 items</option>
            <option value="40">40 items</option>
            <option value="80">80 items</option>
            <option value="100">100 items</option>
          </Select>
        </HStack>
      </Stack>
    </>
  )
}

export default WorkspaceListPage
