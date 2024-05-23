import { useMemo } from 'react'
import { ResponsivePie } from '@nivo/pie'
import cx from 'classnames'
import InsightsAPI, { SortBy, SortOrder } from '@/client/api/insights'
import { swrConfig } from '@/client/options'
import { useAppSelector } from '@/store/hook'

const InsightsOverviewChart = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const { data: entities } = InsightsAPI.useListEntities(
    id,
    { page: 1, size: 5, sortBy: SortBy.Frequency, sortOrder: SortOrder.Desc },
    swrConfig(),
  )
  const data = useMemo(() => {
    if (entities && entities.size >= 5) {
      return entities.data.map((entity) => ({
        id: entity.text,
        label: entity.text,
        value: entity.frequency,
      }))
    }
  }, [entities])

  return (
    <div
      className={cx(
        'w-full',
        'h-[320px]',
        'flex',
        'items-center',
        'justify-center',
      )}
    >
      {data ? (
        <ResponsivePie
          data={data}
          colors={{ scheme: 'greys' }}
          margin={{ top: 40, right: 40, bottom: 40, left: 40 }}
          innerRadius={0.6}
          padAngle={3}
          cornerRadius={3}
          activeOuterRadiusOffset={8}
          borderWidth={2}
          borderColor={{
            from: 'color',
            modifiers: [['darker', 0.2]],
          }}
          arcLinkLabelsSkipAngle={10}
          arcLinkLabelsTextColor="rgb(26, 32, 44)"
          arcLinkLabelsThickness={2}
          arcLinkLabelsColor={{
            from: 'color',
            modifiers: [['darker', 0.2]],
          }}
          arcLabelsSkipAngle={10}
          arcLabelsTextColor={{
            from: 'color',
            modifiers: [['darker', 2]],
          }}
        />
      ) : (
        <p>Not enough data to render the chart.</p>
      )}
    </div>
  )
}

export default InsightsOverviewChart
