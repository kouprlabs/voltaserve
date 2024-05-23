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
    const colors = [
      'hsl(56, 70%, 50%)',
      'hsl(33, 70%, 50%)',
      'hsl(216, 70%, 50%)',
      'hsl(110, 70%, 50%)',
      'hsl(78, 70%, 50%)',
    ]
    if (entities && entities.size >= 5) {
      return entities.data.map((entity, index) => ({
        id: entity.text,
        label: entity.text,
        value: entity.frequency,
        color: colors[index],
      }))
    }
  }, [entities])

  if (!data) {
    return null
  }

  return (
    <div className={cx('w-full', 'h-[320px]')}>
      <ResponsivePie
        data={data}
        margin={{ top: 40, right: 40, bottom: 40, left: 40 }}
        innerRadius={0.5}
        padAngle={0.7}
        cornerRadius={3}
        activeOuterRadiusOffset={8}
        borderWidth={1}
        borderColor={{
          from: 'color',
          modifiers: [['darker', 0.2]],
        }}
        arcLinkLabelsSkipAngle={10}
        arcLinkLabelsTextColor="#333333"
        arcLinkLabelsThickness={2}
        arcLinkLabelsColor={{ from: 'color' }}
        arcLabelsSkipAngle={10}
        arcLabelsTextColor={{
          from: 'color',
          modifiers: [['darker', 2]],
        }}
      />
    </div>
  )
}

export default InsightsOverviewChart
