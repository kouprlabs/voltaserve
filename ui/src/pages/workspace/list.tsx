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
import WorkspaceAPI, { SortOrder } from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import PagePagination, {
  usePagePagination,
} from '@/components/common/page-pagination'
import { CreateWorkspaceButton } from '@/components/top-bar/buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'

const WorkspaceListPage = () => {
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, onPageChange, onSizeChange } = usePagePagination({
    localStoragePrefix: 'voltaserve',
    localStorageNamespace: 'workspace',
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

export default WorkspaceListPage
