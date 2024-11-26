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
import { SectionError, SectionPlaceholder, SectionSpinner } from '@koupr/ui'
import { ResponsivePie } from '@nivo/pie'
import cx from 'classnames'
import InsightsAPI, { SortBy, SortOrder } from '@/client/api/insights'
import { swrConfig } from '@/client/options'
import { useAppSelector } from '@/store/hook'

const InsightsOverviewChart = () => {
  const { colorMode } = useColorMode()
  const colors = useToken('colors', ['gray.200'])
  const colorsDark = useToken('colors', ['gray.500'])
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const {
    data: list,
    error: listError,
    isLoading: isListLoading,
  } = InsightsAPI.useListEntities(
    id,
    { page: 1, size: 5, sortBy: SortBy.Frequency, sortOrder: SortOrder.Desc },
    swrConfig(),
  )
  const isListError = !list && listError
  const isListEmpty = list && !listError && list.totalElements < 5
  const isListReady = list && !listError && list.totalElements >= 5

  return (
    <>
      {isListLoading ? <SectionSpinner /> : null}
      {isListError ? <SectionError text="Failed to load chart." /> : null}
      {isListEmpty ? (
        <SectionPlaceholder text="Not enough data to render the chart." />
      ) : null}
      {isListReady ? (
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
  )
}

export default InsightsOverviewChart
