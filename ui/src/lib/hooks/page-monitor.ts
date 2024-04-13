export type UsePageMonitorMonitorOptions = {
  totalPages: number
  totalElements: number
  steps: number[]
}

export const usePageMonitor = ({
  totalPages,
  totalElements,
  steps,
}: UsePageMonitorMonitorOptions) => {
  const hasPageSwitcher = totalPages > 1
  const hasSizeSelector = totalElements > steps[0]

  return { hasPageSwitcher, hasSizeSelector }
}
