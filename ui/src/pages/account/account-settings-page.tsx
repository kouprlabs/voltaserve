// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useState } from 'react'
import {
  IconButton,
  IconButtonProps,
  Progress,
  Tooltip,
} from '@chakra-ui/react'
import {
  IconEdit,
  IconDelete,
  IconWarning,
  SectionSpinner,
  Form,
  SectionError,
} from '@koupr/ui'
import cx from 'classnames'
import StorageAPI from '@/client/api/storage'
import { errorToString } from '@/client/error'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import AccountChangePassword from '@/components/account/account-change-password'
import AccountDelete from '@/components/account/account-delete'
import AccountEditEmail from '@/components/account/account-edit-email'
import AccountEditFullName from '@/components/account/account-edit-full-name'
import AccountThemeSwitcher from '@/components/account/account-theme-switcher'
import prettyBytes from '@/lib/helpers/pretty-bytes'
import { truncateEnd } from '@/lib/helpers/truncate-end'
import truncateMiddle from '@/lib/helpers/truncate-middle'
import { AccountExtensions } from '@/types/extensibility'

const EditButton = (props: IconButtonProps) => (
  <IconButton
    icon={<IconEdit />}
    className={cx('h-[40px]', 'w-[40px]')}
    {...props}
  />
)

export type AccountSettingsPageProps = {
  extensions?: AccountExtensions
}

const AccountSettingsPage = ({ extensions }: AccountSettingsPageProps) => {
  const {
    data: user,
    error: userError,
    isLoading: userIsLoading,
  } = UserAPI.useGet()
  const {
    data: storageUsage,
    error: storageUsageError,
    isLoading: storageUsageIsLoading,
  } = StorageAPI.useGetAccountUsage(swrConfig())
  const [isFullNameModalOpen, setIsFullNameModalOpen] = useState(false)
  const [isEmailModalOpen, setIsEmailModalOpen] = useState(false)
  const [isPasswordModalOpen, setIsPasswordModalOpen] = useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const userIsReady = user && !userError
  const storageUsageIsReady = storageUsage && !storageUsageError

  return (
    <>
      {userIsLoading ? <SectionSpinner /> : null}
      {userError ? <SectionError text={errorToString(userError)} /> : null}
      {userIsReady ? (
        <>
          <Form
            sections={[
              ...(extensions?.settings?.sections || []),
              {
                title: 'Storage',
                content: (
                  <>
                    {storageUsageError ? (
                      <SectionError
                        text={errorToString(storageUsageError)}
                        height="auto"
                      />
                    ) : null}
                    {storageUsageIsReady ? (
                      <>
                        <span>
                          {prettyBytes(storageUsage.bytes)} of{' '}
                          {prettyBytes(storageUsage.maxBytes)} used
                        </span>
                        <Progress value={storageUsage.percentage} hasStripe />
                      </>
                    ) : null}
                    {storageUsageIsLoading ? (
                      <>
                        <span>Calculating…</span>
                        <Progress value={0} hasStripe />
                      </>
                    ) : null}
                  </>
                ),
              },
              {
                title: 'Basics',
                rows: [
                  {
                    label: 'Full name',
                    content: (
                      <>
                        <span>{truncateEnd(user.fullName, 50)}</span>
                        <EditButton
                          title="Edit name"
                          aria-label="Edit name"
                          onClick={() => setIsFullNameModalOpen(true)}
                        />
                      </>
                    ),
                  },
                ],
              },
              {
                title: 'Credentials',
                rows: [
                  {
                    label: 'Email',
                    content: (
                      <>
                        {user.pendingEmail ? (
                          <div
                            className={cx(
                              'flex',
                              'flex-row',
                              'gap-0.5',
                              'items-center',
                            )}
                          >
                            <Tooltip label="Please check your inbox to confirm your email.">
                              <div
                                className={cx(
                                  'flex',
                                  'items-center',
                                  'justify-center',
                                  'cursor-default',
                                )}
                              >
                                <IconWarning
                                  className={cx('text-yellow-400')}
                                />
                              </div>
                            </Tooltip>
                            <span>{truncateMiddle(user.pendingEmail, 50)}</span>
                          </div>
                        ) : null}
                        {!user.pendingEmail ? (
                          <span>
                            {truncateMiddle(
                              user.pendingEmail || user.email,
                              50,
                            )}
                          </span>
                        ) : null}
                        <EditButton
                          title="Edit email"
                          aria-label="Edit email"
                          onClick={() => setIsEmailModalOpen(true)}
                        />
                      </>
                    ),
                  },
                  {
                    label: 'Password',
                    content: (
                      <EditButton
                        title="Change password"
                        aria-label="Change password"
                        onClick={() => setIsPasswordModalOpen(true)}
                      />
                    ),
                  },
                ],
              },
              {
                title: 'Theme',
                rows: [
                  {
                    label: 'Dark mode',
                    content: <AccountThemeSwitcher />,
                  },
                ],
              },
              {
                title: 'Advanced',
                rows: [
                  {
                    label: 'Delete account',
                    content: (
                      <IconButton
                        icon={<IconDelete />}
                        variant="solid"
                        colorScheme="red"
                        title="Delete account"
                        aria-label="Delete account"
                        onClick={() => setIsDeleteModalOpen(true)}
                      />
                    ),
                  },
                ],
              },
            ]}
          />
          <AccountEditFullName
            open={isFullNameModalOpen}
            user={user}
            onClose={() => setIsFullNameModalOpen(false)}
          />
          <AccountEditEmail
            open={isEmailModalOpen}
            user={user}
            onClose={() => setIsEmailModalOpen(false)}
          />
          <AccountChangePassword
            open={isPasswordModalOpen}
            user={user}
            onClose={() => setIsPasswordModalOpen(false)}
          />
          <AccountDelete
            open={isDeleteModalOpen}
            onClose={() => setIsDeleteModalOpen(false)}
          />
        </>
      ) : null}
    </>
  )
}

export default AccountSettingsPage
