import { ReactNode, useEffect, useMemo, useState } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import {
  Avatar,
  Box,
  Button,
  Heading,
  HStack,
  IconButton,
  Stack,
  Tab,
  TabList,
  Tabs,
  Tag,
  Text,
  VStack,
} from '@chakra-ui/react'
import NotificationAPI from '@/api/notification'
import { swrConfig } from '@/api/options'
import UserAPI from '@/api/user'
import AccountEditPicture from '@/components/account/edit-picture'
import { IconEdit } from '@/components/common/icon'
import variables from '@/theme/variables'

const AccountLayout = () => {
  const location = useLocation()
  const navigate = useNavigate()
  const [isImageModalOpen, setIsImageModalOpen] = useState(false)
  const { data: user } = UserAPI.useGet(swrConfig())
  const { data: notfications } = NotificationAPI.useGetAll(swrConfig())
  const invitationCount = useMemo(
    () => notfications?.filter((e) => e.type === 'new_invitation').length,
    [notfications]
  )
  const [tabIndex, setTabIndex] = useState(0)

  useEffect(() => {
    const segments = location.pathname.split('/')
    const segment = segments[segments.length - 1]
    if (segment === 'settings') {
      setTabIndex(0)
    } else if (segment === 'invitation') {
      setTabIndex(1)
    }
  }, [location])

  if (!user) {
    return null
  }

  return (
    <Stack direction="row" spacing={variables.spacingLg}>
      <VStack spacing={variables.spacingMd} width="250px">
        <VStack spacing={variables.spacingMd}>
          <Box position="relative" flexShrink={0}>
            <Avatar
              name={user.fullName}
              src={user.picture}
              width="165px"
              height="165px"
              size="2xl"
            />
            <IconButton
              icon={<IconEdit />}
              variant="solid-gray"
              right="5px"
              bottom="10px"
              position="absolute"
              zIndex={1000}
              aria-label=""
              onClick={() => setIsImageModalOpen(true)}
            />
          </Box>
          <Heading fontSize="16px" textAlign="center">
            {user.fullName}
          </Heading>
        </VStack>
        <VStack width="100%" spacing={variables.spacingSm}>
          <Button
            variant="outline"
            colorScheme="red"
            width="100%"
            type="submit"
            onClick={() => navigate('/sign-out')}
          >
            Sign out
          </Button>
        </VStack>
      </VStack>
      <Stack w="100%" pb={variables.spacing}>
        <Tabs
          variant="solid-rounded"
          colorScheme="gray"
          pb={variables.spacingLg}
          index={tabIndex}
        >
          <TabList>
            <Tab onClick={() => navigate('/account/settings')}>Settings</Tab>
            <Tab onClick={() => navigate('/account/invitation')}>
              <HStack>
                <Text>Invitations</Text>
                {invitationCount && invitationCount > 0 ? (
                  <Tag borderRadius="full">{invitationCount}</Tag>
                ) : null}
              </HStack>
            </Tab>
          </TabList>
        </Tabs>
        <Outlet />
      </Stack>
      <AccountEditPicture
        open={isImageModalOpen}
        user={user}
        onClose={() => setIsImageModalOpen(false)}
      />
    </Stack>
  )
}

export default AccountLayout
