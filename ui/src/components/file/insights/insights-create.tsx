import { useCallback, useMemo, useState } from 'react'
import { Button, ModalBody, ModalFooter } from '@chakra-ui/react'
import { OptionBase, Select, SingleValue } from 'chakra-react-select'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import InsightsAPI, { Language } from '@/client/api/insights'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/insights'
import { reactSelectStyles } from '@/styles/react-select'

interface LanguageOption extends OptionBase {
  label: string
  value: string
}

const InsightsCreate = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFile = useAppSelector((state) => state.ui.insights.mutateFile)
  const [language, setLanguage] = useState<Language>()
  const [isLoading, setIsLoading] = useState(false)
  const { data: languages } = InsightsAPI.useGetLanguages(swrConfig())
  const { data: summary } = InsightsAPI.useGetMetadata(id, swrConfig())
  const { data: file } = FileAPI.useGet(id, swrConfig())
  const existingLanguage = useMemo<LanguageOption | undefined>(() => {
    if (file && languages && file.snapshot?.language) {
      const value = file.snapshot.language
        ? languages.filter((e) => e.id === file.snapshot?.language)[0]
        : undefined
      if (value) {
        setLanguage(value)
        return { value: value.id, label: value.name }
      }
    }
  }, [file, languages])

  const handleCreate = useCallback(async () => {
    if (id && language) {
      try {
        setIsLoading(true)
        await InsightsAPI.create(id, { languageId: language.id })
        mutateFile?.()
        setIsLoading(false)
      } catch (error) {
        setIsLoading(false)
      } finally {
        setIsLoading(false)
      }
    }
  }, [language, id, mutateFile])

  const handleLanguageChange = useCallback(
    (value: SingleValue<LanguageOption>) => {
      if (value?.value && languages) {
        setLanguage(languages.filter((e) => e.id === value.value)[0])
      }
    },
    [languages],
  )

  if (!id || !file || !summary || !languages) {
    return null
  }

  return (
    <>
      <ModalBody>
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
            Select the language to use for creating insights.
            <br />
            During the process, text will be extracted using OCR (optical
            character recognition), and entities will be scanned using NER
            (named entity recognition).
          </p>
          {languages ? (
            <Select<LanguageOption, false>
              className={cx('w-full')}
              defaultValue={existingLanguage}
              options={languages.map((language) => ({
                value: language.id,
                label: language.name,
              }))}
              placeholder="Select Language"
              selectedOptionStyle="check"
              chakraStyles={reactSelectStyles()}
              isDisabled={isLoading}
              onChange={handleLanguageChange}
            />
          ) : null}
        </div>
      </ModalBody>
      <ModalFooter>
        <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
          <Button
            type="button"
            variant="outline"
            colorScheme="blue"
            isDisabled={isLoading}
            onClick={() => dispatch(modalDidClose())}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="solid"
            colorScheme="blue"
            isLoading={isLoading}
            isDisabled={!language}
            onClick={handleCreate}
          >
            Create
          </Button>
        </div>
      </ModalFooter>
    </>
  )
}

export default InsightsCreate
