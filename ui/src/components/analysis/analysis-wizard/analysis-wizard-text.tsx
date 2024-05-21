import AnalysisText from '../analysis-text'

export type AnalysisWizardTextProps = {
  id: string
}

const AnalysisWizardText = ({ id }: AnalysisWizardTextProps) => {
  return <AnalysisText id={id} />
}

export default AnalysisWizardText
