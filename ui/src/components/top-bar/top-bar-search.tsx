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
import { decodeQuery, encodeQuery } from '@/helpers/query'
import { IconClose, IconSearch } from '@/lib'

const TopBarSearch = () => {
  const navigation = useNavigate()
  const location = useLocation()
  const { id, fileId } = useParams()
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
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
  const [text, setText] = useState(query || '')
  const [isFocused, setIsFocused] = useState(false)

  useEffect(() => {
    if (query) {
      setText(query || '')
    } else {
      setText('')
    }
  }, [query])

  const handleSearch = useCallback(
    (value: string) => {
      if (isFiles) {
        if (value) {
          navigation(`/workspace/${id}/file/${fileId}?q=${encodeQuery(value)}`)
        } else {
          navigation(`/workspace/${id}/file/${fileId}`)
        }
      } else if (isWorkspaces) {
        if (value) {
          navigation(`/workspace?q=${encodeQuery(value)}`)
        } else {
          navigation(`/workspace`)
        }
      } else if (isGroups) {
        if (value) {
          navigation(`/group?q=${encodeQuery(value)}`)
        } else {
          navigation(`/group`)
        }
      } else if (isOrgs) {
        if (value) {
          navigation(`/organization?q=${encodeQuery(value)}`)
        } else {
          navigation(`/organization`)
        }
      } else if (isOrgMembers) {
        if (value) {
          navigation(`/organization/${id}/member?q=${encodeQuery(value)}`)
        } else {
          navigation(`/organization/${id}/member`)
        }
      } else if (isGroupMembers) {
        if (value) {
          navigation(`/group/${id}/member?q=${encodeQuery(value)}`)
        } else {
          navigation(`/group/${id}/member`)
        }
      }
    },
    [
      id,
      fileId,
      isFiles,
      isWorkspaces,
      isGroups,
      isOrgs,
      isOrgMembers,
      isGroupMembers,
      navigation,
    ],
  )

  const handleClear = useCallback(() => {
    setText('')
    if (isFiles) {
      navigation(`/workspace/${id}/file/${fileId}`)
    } else if (isWorkspaces) {
      navigation(`/workspace`)
    } else if (isGroups) {
      navigation(`/group`)
    } else if (isOrgs) {
      navigation(`/organization`)
    } else if (isOrgMembers) {
      navigation(`/organization/${id}/member`)
    } else if (isGroupMembers) {
      navigation(`/group/${id}/member`)
    }
  }, [
    id,
    fileId,
    isFiles,
    isWorkspaces,
    isGroups,
    isOrgs,
    isOrgMembers,
    isGroupMembers,
    navigation,
  ])

  const handleKeyDown = useCallback(
    (event: KeyboardEvent<HTMLInputElement>) => {
      if (event.key === 'Enter') {
        handleSearch(text)
      }
    },
    [text, handleSearch],
  )

  const handleChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    setText(event.target.value || '')
  }, [])

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
          value={text}
          placeholder={query || placeholder}
          variant="filled"
          onKeyDown={handleKeyDown}
          onChange={handleChange}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
        />
        {query && (
          <InputRightElement>
            <IconButton
              icon={<IconClose />}
              onClick={handleClear}
              size="xs"
              aria-label="Clear"
            />
          </InputRightElement>
        )}
      </InputGroup>
      {text || (isFocused && text) ? (
        <Button onClick={() => handleSearch(text)} isDisabled={!text}>
          Search
        </Button>
      ) : null}
    </div>
  )
}

export default TopBarSearch
