import { useCallback, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useSWRConfig } from 'swr'
import {
  Field,
  FieldAttributes,
  FieldProps,
  Form,
  Formik,
  FormikHelpers,
} from 'formik'
import * as Yup from 'yup'
import { Heading, HStack, Stack } from '@chakra-ui/react'
import { Button, FormControl, FormErrorMessage, Input } from '@chakra-ui/react'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI from '@/api/organization'
import variables from '@/theme/variables'

type FormValues = {
  name: string
}

const NewOrganizationPage = () => {
  const navigate = useNavigate()
  const { mutate } = useSWRConfig()
  const [isLoading, setIsLoading] = useState(false)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })

  const handleSubmit = useCallback(
    async (
      { name }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      setSubmitting(true)
      setIsLoading(true)
      try {
        const result = await OrganizationAPI.create({
          name,
        })
        mutate(`/organizations/${result.id}`, result)
        mutate(`/organizations`)
        setSubmitting(false)
        navigate(`/organization/${result.id}/member`)
      } catch {
        setIsLoading(false)
      } finally {
        setSubmitting(false)
      }
    },
    [navigate, mutate]
  )

  return (
    <>
      <Helmet>
        <title>New Organization</title>
      </Helmet>
      <Stack spacing={variables.spacing2Xl}>
        <Heading size="lg">New Organization</Heading>
        <Formik
          enableReinitialize={true}
          initialValues={{ name: '' }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <Stack spacing={variables.spacing2Xl}>
                <Stack spacing={variables.spacing}>
                  <Field name="name">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        maxW="400px"
                        isInvalid={errors.name && touched.name ? true : false}
                      >
                        <Input
                          {...field}
                          placeholder="Name"
                          disabled={isSubmitting}
                          autoFocus
                        />
                        <FormErrorMessage>{errors.name}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                </Stack>
                <HStack>
                  <Button
                    type="submit"
                    variant="solid"
                    colorScheme="blue"
                    isDisabled={isSubmitting || isLoading}
                    isLoading={isSubmitting || isLoading}
                  >
                    Create
                  </Button>
                </HStack>
              </Stack>
            </Form>
          )}
        </Formik>
      </Stack>
    </>
  )
}

export default NewOrganizationPage
