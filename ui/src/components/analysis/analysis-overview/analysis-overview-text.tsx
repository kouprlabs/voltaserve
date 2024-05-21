import AnalysisText from '../analysis-text'

export type AnalysisOverviewTextProps = {
  id: string
}

const AnalysisOverviewText = ({ id }: AnalysisOverviewTextProps) => {
  return <AnalysisText id={id} />
}

export default AnalysisOverviewText
