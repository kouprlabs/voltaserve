import { useEffect, useState } from 'react'
import { Select } from 'chakra-react-select'
import cx from 'classnames'
import reactSelectStyles from '@/styles/react-select'

export type AILanguageProps = {
  isReadOnly?: boolean
}

const AILanguage = ({}: AILanguageProps) => {
  const [value, setValue] = useState<string>()

  useEffect(() => {
    console.log(value)
  }, [value])

  return (
    <div className={cx('flex', 'items-center', 'justify-center')}>
      <Select
        className={cx('w-full')}
        options={[
          { value: 'ara', label: 'Arabic' },
          { value: 'zho', label: 'Chinese Simplified' },
          { value: 'zho', label: 'Chinese Traditional' },
          { value: 'deu', label: 'German' },
          { value: 'eng', label: 'English' },
          { value: 'fra', label: 'French' },
          { value: 'hin', label: 'Hindi' },
          { value: 'ita', label: 'Italian' },
          { value: 'jpn', label: 'Japanese' },
          { value: 'nld', label: 'Dutch' },
          { value: 'por', label: 'Portuguese' },
          { value: 'rus', label: 'Russian' },
          { value: 'spa', label: 'Spanish' },
          { value: 'swe', label: 'Swedish' },
        ]}
        placeholder="Select Language"
        selectedOptionStyle="check"
        chakraStyles={reactSelectStyles}
        onChange={(event) => {
          if (event) {
            setValue(event.value)
          }
        }}
      />
    </div>
  )
}

export default AILanguage
