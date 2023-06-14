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
import { swrConfig } from '@/api/options'
import WorkspaceAPI, { Workspace } from '@/api/workspace'
import { CreateWorkspaceButton } from '@/components/top-bar/buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'

const WorkspaceListPage = () => {
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { data: workspaces, error } = WorkspaceAPI.useGetAllOrSearch(
    query ? { search: { text: query } } : undefined,
    swrConfig()
  )
  return (
    <>
      <Helmet>
        <title>Workspaces</title>
      </Helmet>
      <Stack direction="column" spacing={variables.spacing2Xl}>
        <Heading size="lg" pl={variables.spacingMd}>
          Workspaces
        </Heading>
        {!workspaces && error && (
          <Center h="300px">
            <Text>Failed to load workspaces.</Text>
          </Center>
        )}
        {!workspaces && !error && <SectionSpinner />}
        {workspaces && workspaces.length === 0 && !error ? (
          <Center h="300px">
            <VStack spacing={variables.spacing}>
              <Text>There are no workspaces.</Text>
              <CreateWorkspaceButton />
            </VStack>
          </Center>
        ) : null}
        {workspaces && workspaces.length > 0 && (
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
              {workspaces.map((w: Workspace) => (
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
      </Stack>
    </>
  )
}

export default WorkspaceListPage
