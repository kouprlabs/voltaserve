import {
  useCallback,
  useEffect,
  useState,
  ChangeEvent,
  KeyboardEvent,
} from 'react'
import {
  Button,
  HStack,
  InputGroup,
  InputLeftElement,
  Icon,
  InputRightElement,
  IconButton,
  Input,
} from '@chakra-ui/react'
import { IconClose, IconSearch } from '@koupr/ui'

type SearchInputProps = {
  query?: string
  onChange?: (value: string) => void
}

const SearchInput = ({ query, onChange }: SearchInputProps) => {
  const [draft, setDraft] = useState('')
  const [text, setText] = useState('')
  const [isFocused, setIsFocused] = useState(false)

  useEffect(() => {
    if (query !== draft) {
      setDraft(query || '')
    }
  }, [query, draft])

  useEffect(() => {
    onChange?.(text)
  }, [text, onChange])

  const handleClear = useCallback(() => {
    setDraft('')
    setText('')
  }, [])

  const handleChange = useCallback((event: ChangeEvent<HTMLInputElement>) => {
    setDraft(event.target.value || '')
  }, [])

  const handleSearch = useCallback((value: string) => {
    setText(value)
  }, [])

  const handleKeyDown = useCallback(
    (event: KeyboardEvent<HTMLInputElement>, value: string) => {
      if (event.key === 'Enter') {
        handleSearch(value)
      }
    },
    [handleSearch]
  )

  return (
    <HStack>
      <InputGroup>
        <InputLeftElement pointerEvents="none">
          <Icon as={IconSearch} color="gray.300" />
        </InputLeftElement>
        <Input
          value={draft}
          placeholder={draft || 'Search'}
          variant="filled"
          onKeyDown={(event) => handleKeyDown(event, draft)}
          onChange={handleChange}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
        />
        {draft && (
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
      {draft || (isFocused && draft) ? (
        <Button onClick={() => handleSearch(draft)} isDisabled={!draft}>
          Search
        </Button>
      ) : null}
    </HStack>
  )
}

export default SearchInput
