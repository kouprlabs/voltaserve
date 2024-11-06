// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { ChangeEvent, KeyboardEvent, useEffect, useMemo, useState } from 'react'
import { useCallback } from 'react'
import {
  useLocation,
  useNavigate,
  useParams,
  useSearchParams,
} from 'react-router-dom'
import {
  Button,
  Icon,
  IconButton,
  Input,
  InputGroup,
  InputLeftElement,
  InputRightElement,
} from '@chakra-ui/react'
import { IconClose, IconSearch, IconTune, NotificationBadge } from '@koupr/ui'
import cx from 'classnames'
import { Query as FileQuery } from '@/client/api/file'
import {
  decodeFileQuery,
  decodeQuery,
  encodeFileQuery,
  encodeQuery,
} from '@/lib/helpers/query'
import store from '@/store/configure-store'
import { useAppDispatch } from '@/store/hook'
import { errorOccurred } from '@/store/ui/error'
import { modalDidOpen as searchFilterModalDidOpen } from '@/store/ui/search-filter'

const AppBarSearch = () => {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const location = useLocation()
  const { id: workspaceId, fileId } = useParams()
  const [searchParams] = useSearchParams()
  const isWorkspaces = useMemo(
    () => location.pathname === '/workspace',
    [location],
  )
  const isFiles = useMemo(
    () =>
      location.pathname.includes('/workspace/') &&
      location.pathname.includes('/file/'),
    [location],
  )
  const isGroups = useMemo(() => location.pathname === '/group', [location])
  const isOrgs = useMemo(
    () => location.pathname === '/organization',
    [location],
  )
  const isOrgMembers = useMemo(
    () =>
      location.pathname.includes('/organization/') &&
      location.pathname.includes('/member'),
    [location],
  )
  const isGroupMembers = useMemo(
    () =>
      location.pathname.includes('/group/') &&
      location.pathname.includes('/member'),
    [location],
  )
  const isConsoleUsers = useMemo(
    () =>
      location.pathname.includes('/console/') &&
      location.pathname.includes('/users'),
    [location],
  )
  const isConsoleGroups = useMemo(
    () =>
      location.pathname.includes('/console/') &&
      location.pathname.includes('/groups'),
    [location],
  )
  const isConsoleWorkspaces = useMemo(
    () =>
      location.pathname.includes('/console/') &&
      location.pathname.includes('/workspaces'),
    [location],
  )
  const isConsoleOrganizations = useMemo(
    () =>
      location.pathname.includes('/console/') &&
      location.pathname.includes('/organizations'),
    [location],
  )
  const query: string | FileQuery | undefined = isFiles
    ? decodeFileQuery(searchParams.get('q') as string)
    : decodeQuery(searchParams.get('q') as string)
  const parsedQuery = useMemo(
    () =>
      (isFiles && query ? (query as FileQuery).text : (query as string)) || '',
    [isFiles, query],
  )
  const isAvailable = useMemo(
    () =>
      isWorkspaces ||
      isFiles ||
      isGroups ||
      isOrgs ||
      isOrgMembers ||
      isGroupMembers ||
      isConsoleUsers ||
      isConsoleGroups ||
      isConsoleOrganizations ||
      isConsoleWorkspaces,
    [
      isWorkspaces,
      isFiles,
      isGroups,
      isOrgs,
      isOrgMembers,
      isGroupMembers,
      isConsoleUsers,
      isConsoleGroups,
      isConsoleWorkspaces,
      isConsoleOrganizations,
    ],
  )
  const hasFileQuery = useMemo(() => {
    const fileQuery = query as FileQuery
    return Boolean(
      isFiles &&
        fileQuery &&
        (fileQuery.type ||
          fileQuery.createTimeAfter ||
          fileQuery.createTimeBefore ||
          fileQuery.updateTimeAfter ||
          fileQuery.updateTimeBefore),
    )
  }, [isFiles, query])
  const placeholder = useMemo(() => {
    if (isWorkspaces || isConsoleWorkspaces) {
      return 'Search Workspaces'
    } else if (isFiles) {
      return 'Search Files'
    } else if (isGroups || isConsoleGroups) {
      return 'Search Groups'
    } else if (isOrgs || isConsoleOrganizations) {
      return 'Search Organizations'
    } else if (isOrgMembers) {
      return 'Search Organization Members'
    } else if (isGroupMembers) {
      return 'Search Group Members'
    } else if (isConsoleUsers) {
      return 'Search Users'
    }
  }, [
    isWorkspaces,
    isFiles,
    isGroups,
    isOrgs,
    isOrgMembers,
    isGroupMembers,
    isConsoleUsers,
    isConsoleGroups,
    isConsoleWorkspaces,
    isConsoleOrganizations,
  ])
  const [buffer, setBuffer] = useState(parsedQuery)
  const [isFocused, setIsFocused] = useState(false)

  useEffect(() => {
    if (query) {
      setBuffer(parsedQuery)
    } else {
      setBuffer('')
    }
  }, [parsedQuery, isFiles])

  const handleSearch = useCallback(
    (value: string) => {
      if (isFiles) {
        if (value) {
          const encodedQuery = encodeFileQuery({
            ...(query as FileQuery),
            text: value,
          })
          navigate(`/workspace/${workspaceId}/file/${fileId}?q=${encodedQuery}`)
        } else {
          navigate(`/workspace/${workspaceId}/file/${fileId}`)
        }
      } else if (isWorkspaces) {
        if (value) {
          navigate(`/workspace?q=${encodeQuery(value)}`)
        } else {
          navigate(`/workspace`)
        }
      } else if (isGroups) {
        if (value) {
          navigate(`/group?q=${encodeQuery(value)}`)
        } else {
          navigate(`/group`)
        }
      } else if (isOrgs) {
        if (value) {
          navigate(`/organization?q=${encodeQuery(value)}`)
        } else {
          navigate(`/organization`)
        }
      } else if (isOrgMembers) {
        if (value) {
          navigate(
            `/organization/${workspaceId}/member?q=${encodeQuery(value)}`,
          )
        } else {
          navigate(`/organization/${workspaceId}/member`)
        }
      } else if (isGroupMembers) {
        if (value) {
          navigate(`/group/${workspaceId}/member?q=${encodeQuery(value)}`)
        } else {
          navigate(`/group/${workspaceId}/member`)
        }
      } else if (isConsoleUsers) {
        if (value) {
          navigate(`/console/users?q=${encodeQuery(value)}`)
        } else {
          navigate(`/console/users`)
        }
      } else if (isConsoleWorkspaces) {
        if (value) {
          navigate(`/console/workspaces?q=${encodeQuery(value)}`)
        } else {
          navigate(`/console/workspaces`)
        }
      } else if (isConsoleOrganizations) {
        if (value) {
          navigate(`/console/organizations?q=${encodeQuery(value)}`)
        } else {
          navigate(`/console/organizations`)
        }
      } else if (isConsoleGroups) {
        if (value) {
          navigate(`/console/groups?q=${encodeQuery(value)}`)
        } else {
          navigate(`/console/groups`)
        }
      }
    },
    [
      workspaceId,
      fileId,
      isFiles,
      isWorkspaces,
      isGroups,
      isOrgs,
      isOrgMembers,
      isGroupMembers,
      isConsoleUsers,
      isConsoleOrganizations,
      isConsoleGroups,
      isConsoleWorkspaces,
      navigate,
    ],
  )

  const handleClear = useCallback(() => {
    setBuffer('')
    handleSearch('')
  }, [
    workspaceId,
    fileId,
    isFiles,
    isWorkspaces,
    isGroups,
    isOrgs,
    isOrgMembers,
    isGroupMembers,
    isConsoleUsers,
    isConsoleOrganizations,
    isConsoleGroups,
    isConsoleWorkspaces,
    navigate,
  ])

  const handleKeyDown = useCallback(
    (event: KeyboardEvent<HTMLInputElement>) => {
      if (event.key === 'Enter' && buffer.trim().length >= 3) {
        handleSearch(buffer.trim())
      } else if (
        event.key === 'Enter' &&
        buffer.trim().length > 0 &&
        buffer.length < 3 &&
        !isFiles
      ) {
        store.dispatch(
          errorOccurred('Search query needs at least 3 characters'),
        )
      } else if (event.key === 'Enter' && buffer.trim().length === 0) {
        handleClear()
      }
    },
    [buffer, handleSearch],
  )

  const handleChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    setBuffer(event.target.value || '')
  }, [])

  const handleFilterClick = useCallback(() => {
    dispatch(searchFilterModalDidOpen())
  }, [dispatch])

  if (!isAvailable) {
    return null
  }

  return (
    <div className={cx('flex', 'flex-row', 'gap-0.5', 'grow')}>
      <InputGroup>
        <InputLeftElement pointerEvents="none">
          <Icon as={IconSearch} className={cx('text-gray-300')} />
        </InputLeftElement>
        <Input
          value={buffer}
          placeholder={parsedQuery || placeholder}
          variant="filled"
          onKeyDown={handleKeyDown}
          onChange={handleChange}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
        />
        {query ? (
          <InputRightElement>
            <IconButton
              icon={<IconClose />}
              onClick={handleClear}
              size="xs"
              aria-label="Clear"
            />
          </InputRightElement>
        ) : null}
      </InputGroup>
      {isFiles ? (
        <NotificationBadge hasBadge={hasFileQuery}>
          <IconButton
            icon={<IconTune />}
            aria-label="Filters"
            onClick={handleFilterClick}
          />
        </NotificationBadge>
      ) : null}
      {buffer || (isFocused && buffer) ? (
        <Button
          leftIcon={<IconSearch />}
          onClick={() => handleSearch(buffer)}
          isDisabled={!buffer}
        >
          Search
        </Button>
      ) : null}
    </div>
  )
}

export default AppBarSearch
