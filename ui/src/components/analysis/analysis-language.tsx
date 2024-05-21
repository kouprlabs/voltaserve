import { useEffect, useState } from 'react'
import { Select } from 'chakra-react-select'
import cx from 'classnames'
import AnalysisAPI from '@/client/api/analysis'
import reactSelectStyles from '@/styles/react-select'

export type AnalysisLanguageProps = {
  isReadOnly?: boolean
}

const AnalysisLanguage = ({}: AnalysisLanguageProps) => {
  const [value, setValue] = useState<string>()
  const { data: languages } = AnalysisAPI.useGetLanguages()

  useEffect(() => {
    console.log(value)
  }, [value])

  return (
    <div
      className={cx(
        'flex',
        'flex-col',
        'items-center',
        'justify-center',
        'gap-1.5',
      )}
    >
      <p>
        Specify the language to use for extracting text using OCR (optical
        character recognition), and for scanning entities using NER (named
        entity recognition).
      </p>
      {languages ? (
        <Select
          className={cx('w-full')}
          options={languages.map((language) => ({
            value: language.id,
            label: language.name,
          }))}
          placeholder="Select Language"
          selectedOptionStyle="check"
          chakraStyles={reactSelectStyles}
          onChange={(event) => {
            if (event) {
              setValue(event.value)
            }
          }}
        />
      ) : null}
    </div>
  )
}

export default AnalysisLanguage
