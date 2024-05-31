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
  Circle,
  Icon,
  IconButton,
  Input,
  InputGroup,
  InputLeftElement,
  InputRightElement,
} from '@chakra-ui/react'
import cx from 'classnames'
import {
  decodeFileQuery,
  decodeQuery,
  encodeFileQuery,
  encodeQuery,
} from '@/helpers/query'
import { IconClose, IconSearch, IconTune } from '@/lib/components/icons'
import { useAppDispatch } from '@/store/hook'
import { modalDidOpen as searchFilterModalDidOpen } from '@/store/ui/search-filter'

const TopBarSearch = () => {
  const dispatch = useAppDispatch()
  const navigation = useNavigate()
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
  const query = isFiles
    ? decodeFileQuery(searchParams.get('q') as string)?.text
    : decodeQuery(searchParams.get('q') as string)
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
          navigation(
            `/workspace/${workspaceId}/file/${fileId}?q=${encodeFileQuery({ text: value })}`,
          )
        } else {
          navigation(`/workspace/${workspaceId}/file/${fileId}`)
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
          navigation(
            `/organization/${workspaceId}/member?q=${encodeQuery(value)}`,
          )
        } else {
          navigation(`/organization/${workspaceId}/member`)
        }
      } else if (isGroupMembers) {
        if (value) {
          navigation(`/group/${workspaceId}/member?q=${encodeQuery(value)}`)
        } else {
          navigation(`/group/${workspaceId}/member`)
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
      navigation,
    ],
  )

  const handleClear = useCallback(() => {
    setText('')
    if (isFiles) {
      navigation(`/workspace/${workspaceId}/file/${fileId}`)
    } else if (isWorkspaces) {
      navigation(`/workspace`)
    } else if (isGroups) {
      navigation(`/group`)
    } else if (isOrgs) {
      navigation(`/organization`)
    } else if (isOrgMembers) {
      navigation(`/organization/${workspaceId}/member`)
    } else if (isGroupMembers) {
      navigation(`/group/${workspaceId}/member`)
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
          value={text}
          placeholder={query || placeholder}
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
      <div className={cx('flex', 'items-center', 'justify-center', 'relative')}>
        <IconButton
          icon={<IconTune />}
          aria-label="Filters"
          onClick={handleFilterClick}
        />
        <Circle size="10px" bg="red" position="absolute" top={0} right={0} />
      </div>
      {text || (isFocused && text) ? (
        <Button
          leftIcon={<IconSearch />}
          onClick={() => handleSearch(text)}
          isDisabled={!text}
        >
          Search
        </Button>
      ) : null}
    </div>
  )
}

export default TopBarSearch
