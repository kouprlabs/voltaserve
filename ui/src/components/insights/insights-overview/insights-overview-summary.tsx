// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useColorMode, useToken } from '@chakra-ui/react'
import { SectionError, SectionSpinner } from '@koupr/ui'
import { ResponsivePie } from '@nivo/pie'
import cx from 'classnames'
import { FileAPI } from '@/client'
import { EntityAPI, EntitySortBy, EntitySortOrder } from '@/client/api/entity'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { useAppSelector } from '@/store/hook'

const InsightsOverviewSummary = () => {
  const { colorMode } = useColorMode()
  const colors = useToken('colors', ['gray.200'])
  const colorsDark = useToken('colors', ['gray.500'])
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
  const fileIsReady = file && !fileError
  const {
    data: list,
    error: listError,
    isLoading: listIsLoading,
  } = EntityAPI.useList(
    file?.snapshot?.capabilities.entities ? id : null,
    {
      page: 1,
      size: 5,
      sortBy: EntitySortBy.Frequency,
      sortOrder: EntitySortOrder.Desc,
    },
    swrConfig(),
  )
  const listIsReady = list && !listError && list.totalElements >= 5

  return (
    <>
      {fileIsLoading ? <SectionSpinner /> : null}
      {fileError ? <SectionError text={errorToString(fileError)} /> : null}
      {fileIsReady ? (
        <div className={cx('flex', 'flex-col', 'gap-1.5')}>
          {file.snapshot?.capabilities.summary ? (
            <p>{file.snapshot?.summary}</p>
          ) : null}
          <>
            {listIsLoading ? <SectionSpinner /> : null}
            {listError ? (
              <SectionError text={errorToString(listError)} />
            ) : null}
            {listIsReady ? (
              <div
                className={cx(
                  'w-full',
                  'h-[300px]',
                  'flex',
                  'items-center',
                  'justify-center',
                )}
              >
                <ResponsivePie
                  data={list.data.map((entity) => ({
                    id: `${entity.text} (${entity.frequency})`,
                    label: `${entity.text} (${entity.frequency})`,
                    value: entity.frequency,
                  }))}
                  tooltip={() => null}
                  colors={colorMode === 'dark' ? colorsDark : colors}
                  margin={{ top: 40, right: 40, bottom: 40, left: 40 }}
                  innerRadius={0.65}
                  padAngle={3}
                  cornerRadius={4}
                  activeOuterRadiusOffset={8}
                  arcLabel={() => ''}
                  arcLinkLabelsSkipAngle={10}
                  arcLinkLabelsTextColor={
                    colorMode === 'light' ? 'rgb(26, 32, 44)' : 'white'
                  }
                  arcLinkLabelsThickness={2}
                  arcLinkLabelsColor={{ from: 'color' }}
                  arcLabelsSkipAngle={10}
                  arcLabelsTextColor={{
                    from: 'color',
                    modifiers: [['darker', 2]],
                  }}
                />
              </div>
            ) : null}
          </>
        </div>
      ) : null}
    </>
  )
}

export default InsightsOverviewSummary
