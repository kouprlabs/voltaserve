import { useCallback } from 'react'
import {
  Button,
  ModalBody,
  ModalFooter,
  Step,
  StepDescription,
  StepIcon,
  StepIndicator,
  StepNumber,
  StepSeparator,
  StepStatus,
  StepTitle,
  Stepper,
  useSteps,
} from '@chakra-ui/react'
import cx from 'classnames'
import { useAppDispatch } from '@/store/hook'
import { modalDidClose, wizardDidComplete } from '@/store/ui/ai'
import AiWizardLanguage from './ai-wizard-language'
import AiWizardNamedEntities from './ai-wizard-named-entities'
import AiWizardText from './ai-wizard-text'

const steps = [
  { title: 'Text', description: 'Extract with OCR' },
  { title: 'Language', description: 'Detect with NLP' },
  { title: 'Named Entities', description: 'Scan with NER' },
]

const AiWizard = () => {
  const dispatch = useAppDispatch()
  const { activeStep, setActiveStep } = useSteps({
    index: 0,
    count: steps.length,
  })

  const handleNextStep = useCallback(() => {
    if (activeStep < steps.length) {
      setActiveStep(activeStep + 1)
    }
    if (activeStep === steps.length) {
      dispatch(wizardDidComplete(true))
    }
  }, [activeStep, setActiveStep, dispatch])

  return (
    <>
      <ModalBody>
        <div className={cx('flex', 'flex-col', 'gap-1.5')}>
          <Stepper index={activeStep}>
            {steps.map((step, index) => (
              <Step key={index}>
                <StepIndicator>
                  <StepStatus
                    complete={<StepIcon />}
                    incomplete={<StepNumber />}
                    active={<StepNumber />}
                  />
                </StepIndicator>
                <div className={cx('shrink-0')}>
                  <StepTitle>{step.title}</StepTitle>
                  <StepDescription>{step.description}</StepDescription>
                </div>
                <StepSeparator />
              </Step>
            ))}
          </Stepper>
          {activeStep === 0 ? <AiWizardText /> : null}
          {activeStep === 1 ? <AiWizardLanguage /> : null}
          {activeStep === 2 ? <AiWizardNamedEntities /> : null}
          {activeStep === steps.length ? (
            <div
              className={cx(
                'flex',
                'items-center',
                'justify-center',
                'h-[40px]',
              )}
            >
              Success!
            </div>
          ) : null}
        </div>
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
            type="submit"
            variant="solid"
            colorScheme={'blue'}
            onClick={handleNextStep}
          >
            {activeStep === steps.length ? 'Finish' : 'Next'}
          </Button>
        </div>
      </ModalFooter>
    </>
  )
}

export default AiWizard
