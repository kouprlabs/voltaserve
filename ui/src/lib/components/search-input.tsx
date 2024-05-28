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
  InputRightElement,
  IconButton,
  Input,
} from '@chakra-ui/react'
import cx from 'classnames'
import { IconClose, IconSearch } from './icons'

type SearchInputProps = {
  placeholder?: string
  query?: string
  onChange?: (value: string) => void
}

export const SearchInput = ({
  placeholder,
  query,
  onChange,
}: SearchInputProps) => {
  const [draft, setDraft] = useState('')
  const [text, setText] = useState('')
  const [isFocused, setIsFocused] = useState(false)

  useEffect(() => {
    setDraft(query || '')
  }, [query])

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
    [handleSearch],
  )

  return (
    <HStack>
      <InputGroup>
        <InputLeftElement className={cx('pointer-events-none')}>
          <IconSearch className={cx('text-gray-300')} />
        </InputLeftElement>
        <Input
          value={draft}
          placeholder={draft || placeholder || 'Search'}
          variant="filled"
          onKeyDown={(event) => handleKeyDown(event, draft)}
          onChange={handleChange}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
        />
        {draft ? (
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
      {draft || (isFocused && draft) ? (
        <Button onClick={() => handleSearch(draft)} isDisabled={!draft}>
          Search
        </Button>
      ) : null}
    </HStack>
  )
}
