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
} from '@chakra-ui/react'
import { SectionSpinner, variables } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import { swrConfig } from '@/api/options'
import OrganizationAPI, { Organization } from '@/api/organization'
import { CreateOrganizationButton } from '@/components/top-bar/buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'

const OrganizationListPage = () => {
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { data: orgs, error } = OrganizationAPI.useGetAllOrSearch(
    query ? { search: { text: query } } : undefined,
    swrConfig()
  )
  return (
    <>
      <Helmet>
        <title>Organizations</title>
      </Helmet>
      <Stack direction="column" spacing={variables.spacing2Xl}>
        <Heading size="lg" pl={variables.spacingMd}>
          Organizations
        </Heading>
        {!orgs && error && (
          <Center h="300px">
            <Text>Failed to load organizations.</Text>
          </Center>
        )}
        {!orgs && !error && <SectionSpinner />}
        {orgs && orgs.length === 0 && (
          <Center h="300px">
            <VStack spacing={variables.spacing}>
              <Text>There are no organizations.</Text>
              <CreateOrganizationButton />
            </VStack>
          </Center>
        )}
        {orgs && orgs.length > 0 && (
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Name</Th>
                <Th>Date</Th>
              </Tr>
            </Thead>
            <Tbody>
              {orgs.map((o: Organization) => (
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
                  <Td>{prettyDate(o.updateTime || o.createTime)}</Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        )}
      </Stack>
    </>
  )
}

export default OrganizationListPage
