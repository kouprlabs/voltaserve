// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useEffect, useState } from 'react'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import {
  Avatar,
  Button,
  Heading,
  IconButton,
  Tab,
  TabList,
  Tabs,
} from '@chakra-ui/react'
import { IconEdit, NumberTag, SectionError, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import { InvitationAPI } from '@/client/api/invitation'
import { errorToString } from '@/client/error'
import { AuthUserAPI } from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import AccountEditPicture from '@/components/account/edit-picture'
import { getPictureUrl } from '@/lib/helpers/picture'
import { truncateEnd } from '@/lib/helpers/truncate-end'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/account'
import { AccountExtensions } from '@/types/extensibility'

export type AccountLayoutProps = {
  extensions?: AccountExtensions
}

const AccountLayout = ({ extensions }: AccountLayoutProps) => {
  const location = useLocation()
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const [isImageModalOpen, setIsImageModalOpen] = useState(false)
  const {
    data: user,
    isLoading: userIsLoading,
    error: userError,
    mutate,
  } = AuthUserAPI.useGet(swrConfig())
  const { data: invitationCount } =
    InvitationAPI.useGetIncomingCount(swrConfig())
  const [tabIndex, setTabIndex] = useState(0)
  const userIsReady = user && !userError

  useEffect(() => {
    const segments = location.pathname.split('/')
    const segment = segments[segments.length - 1]
    if (segment === 'settings') {
      setTabIndex(0)
    } else if (segment === 'invitation') {
      setTabIndex(1)
    } else {
      const index = extensions?.pages
        ?.filter((page) => page.tab)
        .findIndex((page) => page.path === location.pathname)
      if (index !== undefined && index !== -1) {
        setTabIndex(index + 2)
      }
    }
  }, [location])

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate, dispatch])

  return (
    <>
      {userIsLoading ? <SectionSpinner /> : null}
      {userError ? <SectionError text={errorToString(userError)} /> : null}
      {userIsReady ? (
        <>
          <Helmet>
            <title>{user.fullName}</title>
          </Helmet>
          <div className={cx('flex', 'flex-row', 'gap-2.5')}>
            <div
              className={cx(
                'flex',
                'flex-col',
                'gap-2',
                'items-center',
                'w-[250px]',
              )}
            >
              <div className={cx('flex', 'flex-col', 'gap-2', 'items-center')}>
                <div className={cx('relative', 'shrink-0')}>
                  <Avatar
                    name={user.fullName}
                    src={user.picture ? getPictureUrl(user.picture) : undefined}
                    size="2xl"
                    className={cx(
                      'w-[165px]',
                      'h-[165px]',
                      'border',
                      'border-gray-300',
                      'dark:border-gray-700',
                    )}
                  />
                  <IconButton
                    icon={<IconEdit />}
                    variant="solid-gray"
                    right="5px"
                    bottom="10px"
                    position="absolute"
                    zIndex={1000}
                    title="Edit picture"
                    aria-label="Edit picture"
                    onClick={() => setIsImageModalOpen(true)}
                  />
                </div>
                <Heading className={cx('text-center', 'text-heading')}>
                  {truncateEnd(user.fullName, 50)}
                </Heading>
              </div>
              <div className={cx('w-full', 'gap-1')}>
                <Button
                  variant="outline"
                  colorScheme="red"
                  type="submit"
                  className={cx('w-full')}
                  onClick={() => navigate('/sign-out')}
                >
                  Sign Out
                </Button>
              </div>
            </div>
            <div className={cx('w-full', 'pb-1.5')}>
              <Tabs
                variant="solid-rounded"
                colorScheme="gray"
                index={tabIndex}
                className={cx('pb-2.5')}
              >
                <TabList>
                  <Tab onClick={() => navigate('/account/settings')}>
                    Settings
                  </Tab>
                  <Tab onClick={() => navigate('/account/invitation')}>
                    <div
                      className={cx(
                        'flex',
                        'flex-row',
                        'items-center',
                        'gap-0.5',
                      )}
                    >
                      <span>Invitations</span>
                      {invitationCount && invitationCount > 0 ? (
                        <NumberTag isActive={tabIndex === 1}>
                          {invitationCount}
                        </NumberTag>
                      ) : null}
                    </div>
                  </Tab>
                  {extensions?.pages
                    ?.filter((page) => page.tab)
                    .map((page, index) => (
                      <Tab key={index} onClick={() => navigate(page.path)}>
                        {page.tab!.label}
                      </Tab>
                    ))}
                </TabList>
              </Tabs>
              <Outlet />
            </div>
          </div>
          <AccountEditPicture
            open={isImageModalOpen}
            user={user}
            onClose={() => setIsImageModalOpen(false)}
          />
        </>
      ) : null}
    </>
  )
}

export default AccountLayout
