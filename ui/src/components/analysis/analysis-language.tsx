import { useCallback, useMemo, useState } from 'react'
import { OptionBase, Select, SingleValue } from 'chakra-react-select'
import cx from 'classnames'
import AnalysisAPI, { Language } from '@/client/api/analysis'
import FileAPI from '@/client/api/file'
import { swrConfig } from '@/client/options'
import { useAppSelector } from '@/store/hook'
import reactSelectStyles from '@/styles/react-select'

export type AnalysisLanguageProps = {
  isReadOnly?: boolean
}

interface LanguageOption extends OptionBase {
  label: string
  value: string
}

const AnalysisLanguage = ({}: AnalysisLanguageProps) => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const [_, setValue] = useState<Language>()
  const { data: languages } = AnalysisAPI.useGetLanguages(swrConfig())
  const { data: summary } = AnalysisAPI.useGetSummary(id, swrConfig())
  const { data: file } = FileAPI.useGet(id, swrConfig())
  const defaultValue = useMemo<LanguageOption | undefined>(() => {
    if (file && summary && languages && file.snapshot?.language) {
      const value = summary.hasLanguage
        ? languages.filter((e) => e.id === file.snapshot?.language)[0]
        : undefined
      if (value) {
        return { value: value.id, label: value.name }
      }
    }
  }, [file, summary, languages])

  const handleChange = useCallback(
    () => (value: SingleValue<LanguageOption>) => {
      if (value && languages) {
        setValue(languages.filter((e) => e.id === value.value)[0])
      }
    },
    [languages],
  )

  if (!id || !file || !summary || !languages) {
    return null
  }

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
        Select the language to use for extracting text using OCR (optical
        character recognition), and for scanning entities using NER (named
        entity recognition).
      </p>
      {languages ? (
        <Select<LanguageOption, false>
          className={cx('w-full')}
          defaultValue={defaultValue}
          options={languages.map((language) => ({
            value: language.id,
            label: language.name,
          }))}
          placeholder="Select Language"
          selectedOptionStyle="check"
          chakraStyles={reactSelectStyles}
          onChange={handleChange}
        />
      ) : null}
    </div>
  )
}

export default AnalysisLanguage
