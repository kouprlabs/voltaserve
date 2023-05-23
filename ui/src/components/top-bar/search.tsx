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
  HStack,
  Icon,
  IconButton,
  Input,
  InputGroup,
  InputLeftElement,
  InputRightElement,
} from '@chakra-ui/react'
import { IconClose, IconSearch } from '@koupr/ui'
import { decodeQuery, encodeQuery } from '@/helpers/query'

const Search = () => {
  const navigation = useNavigate()
  const location = useLocation()
  const params = useParams()
  const workspaceId = params.id as string
  const fileId = params.fileId as string
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const isWorkspaces = useMemo(
    () => location.pathname === '/workspace',
    [location]
  )
  const isFiles = useMemo(
    () =>
      location.pathname.includes('/workspace/') &&
      location.pathname.includes('/file/'),
    [location]
  )
  const isGroups = useMemo(() => location.pathname === '/group', [location])
  const isOrgs = useMemo(
    () => location.pathname === '/organization',
    [location]
  )
  const isAvailable = useMemo(
    () => isWorkspaces || isFiles || isGroups || isOrgs,
    [isWorkspaces, isFiles, isGroups, isOrgs]
  )
  const placeholder = useMemo(() => {
    if (isWorkspaces) {
      return `Search Workspaces`
    }
    if (isFiles) {
      return `Search Files`
    }
    if (isGroups) {
      return `Search Groups`
    }
    if (isOrgs) {
      return `Search Organizations`
    }
  }, [isWorkspaces, isFiles, isGroups, isOrgs])
  const [text, setText] = useState(query || '')
  const [isFocused, setIsButtonVisible] = useState(false)

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
            `/workspace/${workspaceId}/file/${fileId}?q=${encodeQuery(value)}`
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
      }
    },
    [workspaceId, fileId, isFiles, isWorkspaces, isGroups, isOrgs, navigation]
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
    }
  }, [workspaceId, fileId, isFiles, isWorkspaces, isGroups, isOrgs, navigation])

  const handleKeyDown = useCallback(
    (event: KeyboardEvent<HTMLInputElement>) => {
      if (event.key === 'Enter') {
        handleSearch(text)
      }
    },
    [text, handleSearch]
  )

  const handleChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    setText(event.target.value || '')
  }, [])

  if (!isAvailable) {
    return null
  }

  return (
    <HStack>
      <InputGroup>
        <InputLeftElement pointerEvents="none">
          <Icon as={IconSearch} color="gray.300" />
        </InputLeftElement>
        <Input
          value={text}
          placeholder={query || placeholder}
          variant="filled"
          onKeyDown={handleKeyDown}
          onChange={handleChange}
          onFocus={() => setIsButtonVisible(true)}
          onBlur={() => setIsButtonVisible(false)}
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
    </HStack>
  )
}

export default Search
