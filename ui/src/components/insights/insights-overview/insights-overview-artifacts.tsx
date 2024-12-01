// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Button } from '@chakra-ui/react'
import { IconOpenInNew, SectionError, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { useAppSelector } from '@/store/hook'

const InsightsOverviewArtifacts = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const {
    data: file,
    error: fileError,
    isLoading: fileIsLoading,
  } = FileAPI.useGet(id, swrConfig())
  const searchParams = new URLSearchParams({
    access_token: getAccessTokenOrRedirect(),
  })
  const fileIsReady = file && !fileError

  return (
    <>
      {fileIsLoading ? <SectionSpinner /> : null}
      {fileError ? <SectionError text={errorToString(fileError)} /> : null}
      {fileIsReady ? (
        <div
          className={cx(
            'flex',
            'flex-col',
            'items-center',
            'justify-center',
            'gap-1',
          )}
        >
          {file.snapshot?.text ? (
            <Button
              as="a"
              type="button"
              leftIcon={<IconOpenInNew />}
              href={`/proxy/api/v3/insights/${id}/text${file.snapshot?.text.extension}?${searchParams}`}
              target="_blank"
            >
              Open Plain Text File
            </Button>
          ) : null}
          {file.snapshot?.ocr ? (
            <Button
              as="a"
              type="button"
              leftIcon={<IconOpenInNew />}
              href={`/proxy/api/v3/insights/${id}/ocr${file.snapshot?.ocr.extension}?${searchParams}`}
              target="_blank"
            >
              Open Searchable PDF
            </Button>
          ) : null}
        </div>
      ) : null}
    </>
  )
}

export default InsightsOverviewArtifacts
