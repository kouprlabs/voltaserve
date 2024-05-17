import { RefObject, useEffect } from 'react'

export default function useFocusAndSelectAll(
  inputRef: RefObject<HTMLInputElement>,
  enable?: boolean,
) {
  useEffect(() => {
    setTimeout(() => {
      if (enable) {
        inputRef.current?.focus()
        setTimeout(() => inputRef.current?.select(), 100)
      }
    }, 100)
  }, [inputRef, enable])
}
