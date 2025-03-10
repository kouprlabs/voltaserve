// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useMemo, useState } from 'react'
import { Button, ModalBody, ModalFooter } from '@chakra-ui/react'
import { SectionError, SectionSpinner, Select } from '@koupr/ui'
import { OptionBase, SingleValue } from 'chakra-react-select'
import cx from 'classnames'
import { EntityAPI } from '@/client/api/entity'
import { FileAPI } from '@/client/api/file'
import { SnapshotAPI, SnapshotLanguage } from '@/client/api/snapshot'
import { TaskAPI } from '@/client/api/task'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/insights'

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
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const mutateTasks = useAppSelector((state) => state.ui.tasks.mutateList)
  const [language, setLanguage] = useState<SnapshotLanguage>()
  const {
    data: languages,
    error: languagesError,
    isLoading: languagesIsLoading,
  } = SnapshotAPI.useGetLanguages(swrConfig())
  const {
    data: file,
    error: fileError,
    isLoading: fileIsLoading,
    mutate: mutateFile,
  } = FileAPI.useGet(id, swrConfig())
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
  const fileIsReady = file && !fileError
  const languagesIsReady = languages && !languagesError

  const handleCreate = useCallback(async () => {
    if (id && language) {
      await EntityAPI.create(id, { language: language.id }, false)
      await mutateFile(await FileAPI.get(id))
      await mutateFiles?.()
      await mutateTasks?.(await TaskAPI.list())
      dispatch(modalDidClose())
    }
  }, [language, id, mutateFile, mutateFiles, mutateTasks, dispatch])

  const handleLanguageChange = useCallback(
    (newValue: SingleValue<LanguageOption>) => {
      if (newValue?.value && languages) {
        setLanguage(languages.filter((e) => e.id === newValue.value)[0])
      }
    },
    [languages],
  )

  return (
    <>
      <ModalBody>
        {fileIsLoading ? <SectionSpinner /> : null}
        {fileError ? <SectionError text={errorToString(fileError)} /> : null}
        {fileIsReady ? (
          <>
            {languagesIsLoading ? <SectionSpinner /> : null}
            {languagesError ? (
              <SectionError text={errorToString(languagesError)} />
            ) : null}
            {languagesIsReady ? (
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
                  Select the language to use for collecting insights. During the
                  process, text will be extracted using OCR (optical character
                  recognition), and entities will be scanned using NER (named
                  entity recognition).
                </p>
                <Select<LanguageOption, false>
                  className={cx('w-full')}
                  defaultValue={existingLanguage}
                  options={languages.map((language) => ({
                    value: language.id,
                    label: language.name,
                  }))}
                  placeholder="Select Language"
                  selectedOptionStyle="check"
                  onChange={handleLanguageChange}
                />
              </div>
            ) : null}
          </>
        ) : null}
      </ModalBody>
      <ModalFooter>
        <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
          <Button
            type="button"
            variant="outline"
            colorScheme="blue"
            onClick={() => dispatch(modalDidClose())}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="solid"
            colorScheme="blue"
            isDisabled={!language}
            onClick={handleCreate}
          >
            Collect
          </Button>
        </div>
      </ModalFooter>
    </>
  )
}

export default InsightsCreate
