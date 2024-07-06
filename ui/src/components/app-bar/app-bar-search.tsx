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
import cx from 'classnames'
import { Query as FileQuery } from '@/client/api/file'
import { IconClose, IconSearch, IconTune } from '@/lib/components/icons'
import NotificationBadge from '@/lib/components/notification-badge'
import {
  decodeFileQuery,
  decodeQuery,
  encodeFileQuery,
  encodeQuery,
} from '@/lib/helpers/query'
import { useAppDispatch } from '@/store/hook'
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
      isGroupMembers,
    [isWorkspaces, isFiles, isGroups, isOrgs, isOrgMembers, isGroupMembers],
  )
  const hasFileQuery = useMemo(() => {
    const fileQuery = query as FileQuery
    return isFiles &&
      fileQuery &&
      (fileQuery.type ||
        fileQuery.createTimeAfter ||
        fileQuery.createTimeBefore ||
        fileQuery.updateTimeAfter ||
        fileQuery.updateTimeBefore)
      ? true
      : false
  }, [isFiles, query])
  const placeholder = useMemo(() => {
    if (isWorkspaces) {
      return 'Search Workspaces'
    } else if (isFiles) {
      return 'Search Files'
    } else if (isGroups) {
      return 'Search Groups'
    } else if (isOrgs) {
      return 'Search Organizations'
    } else if (isOrgMembers) {
      return 'Search Organization Members'
    } else if (isGroupMembers) {
      return 'Search Group Members'
    }
  }, [isWorkspaces, isFiles, isGroups, isOrgs, isOrgMembers, isGroupMembers])
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
      navigate,
    ],
  )

  const handleClear = useCallback(() => {
    setBuffer('')
    if (isFiles) {
      navigate(`/workspace/${workspaceId}/file/${fileId}`)
    } else if (isWorkspaces) {
      navigate(`/workspace`)
    } else if (isGroups) {
      navigate(`/group`)
    } else if (isOrgs) {
      navigate(`/organization`)
    } else if (isOrgMembers) {
      navigate(`/organization/${workspaceId}/member`)
    } else if (isGroupMembers) {
      navigate(`/group/${workspaceId}/member`)
    }
  }, [
    workspaceId,
    fileId,
    isFiles,
    isWorkspaces,
    isGroups,
    isOrgs,
    isOrgMembers,
    isGroupMembers,
    navigate,
  ])

  const handleKeyDown = useCallback(
    (event: KeyboardEvent<HTMLInputElement>) => {
      if (event.key === 'Enter') {
        handleSearch(buffer)
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
    <div className={cx('flex', 'flex-row', 'gap-0.5')}>
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
      <NotificationBadge hasBadge={hasFileQuery}>
        <IconButton
          icon={<IconTune />}
          aria-label="Filters"
          onClick={handleFilterClick}
        />
      </NotificationBadge>
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
