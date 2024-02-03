import { useCallback, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Heading,
  HStack,
  Input,
  Stack,
} from '@chakra-ui/react'
import { variables } from '@koupr/ui'
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
import { Helmet } from 'react-helmet-async'
import GroupAPI from '@/client/api/group'
import OrganizationSelector from '@/components/common/organization-selector'

type FormValues = {
  name: string
  organizationId: string
}

const NewGroupPage = () => {
  const navigate = useNavigate()
  const { org } = useParams()
  const { mutate } = useSWRConfig()
  const [isLoading, setIsLoading] = useState(false)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
    organizationId: Yup.string().required('Organization is required'),
  })

  const handleSubmit = useCallback(
    async (
      { name, organizationId }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      setIsLoading(true)
      try {
        const result = await GroupAPI.create({
          name,
          organizationId,
        })
        mutate(`/groups/${result.id}`, result)
        mutate(`/groups`)
        setSubmitting(false)
        navigate(`/group/${result.id}/member`)
      } catch {
        setIsLoading(false)
      } finally {
        setSubmitting(false)
      }
    },
    [navigate, mutate],
  )

  return (
    <>
      <Helmet>
        <title>New Group</title>
      </Helmet>
      <Stack spacing={variables.spacing2Xl}>
        <Heading fontSize={variables.headingFontSize}>New Group</Heading>
        <Formik
          enableReinitialize={true}
          initialValues={{ name: '', organizationId: org || '' }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting, setFieldValue }) => (
            <Form>
              <Stack spacing={variables.spacing2Xl}>
                <Stack spacing={variables.spacing}>
                  <Field name="name">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        maxW="400px"
                        isInvalid={errors.name && touched.name ? true : false}
                      >
                        <FormLabel>Name</FormLabel>
                        <Input {...field} disabled={isSubmitting} autoFocus />
                        <FormErrorMessage>{errors.name}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <Field name="organizationId">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        maxW="400px"
                        isInvalid={
                          errors.organizationId && touched.organizationId
                            ? true
                            : false
                        }
                      >
                        <FormLabel>Organization</FormLabel>
                        <OrganizationSelector
                          onConfirm={(value) =>
                            setFieldValue(field.name, value.id)
                          }
                        />
                        <FormErrorMessage>
                          {errors.organizationId}
                        </FormErrorMessage>
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
                  <Button as={Link} to="/group" variant="solid">
                    Cancel
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

export default NewGroupPage
