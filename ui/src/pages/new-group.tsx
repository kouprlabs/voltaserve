import { useCallback, useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
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
import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Heading,
  HStack,
  Input,
  Select,
  Stack,
} from '@chakra-ui/react'
import { Helmet } from 'react-helmet-async'
import GroupAPI from '@/api/group'
import OrganizationAPI, { Organization } from '@/api/organization'
import { geEditorPermission } from '@/api/permission'
import variables from '@/theme/variables'

type FormValues = {
  name: string
  organizationId: string
}

const NewGroupPage = () => {
  const navigate = useNavigate()
  const params = useParams()
  const orgId = params.org as string
  const { mutate } = useSWRConfig()
  const [orgs, setOrgs] = useState<Organization[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
    organizationId: Yup.string().required('Organization is required'),
  })

  useEffect(() => {
    ;(async () => {
      const result = await OrganizationAPI.getAll()
      setOrgs(result.filter((o) => geEditorPermission(o.permission)))
    })()
  }, [])

  const handleSubmit = useCallback(
    async (
      { name, organizationId }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
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
    [navigate, mutate]
  )

  return (
    <>
      <Helmet>
        <title>New Group</title>
      </Helmet>
      <Stack spacing={variables.spacing2Xl}>
        <Heading size="lg">New Group</Heading>
        <Formik
          enableReinitialize={true}
          initialValues={{ name: '', organizationId: orgId || '' }}
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
                        <Select {...field} placeholder=" ">
                          {orgs.map((o) => (
                            <option key={o.id} value={o.id}>
                              {o.name}
                            </option>
                          ))}
                        </Select>
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
