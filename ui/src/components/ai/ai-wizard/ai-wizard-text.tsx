import AIText from '../ai-text'

export type AIWizardTextProps = {
  id: string
}

const AIWizardText = ({ id }: AIWizardTextProps) => {
  return <AIText id={id} />
}

export default AIWizardText
