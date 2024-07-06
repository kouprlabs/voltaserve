// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { useCallback, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Heading } from '@chakra-ui/react'
import { Button, FormControl, FormErrorMessage, Input } from '@chakra-ui/react'
import {
  Field,
  FieldAttributes,
  FieldProps,
  Form,
  Formik,
  FormikHelpers,
} from 'formik'
import * as Yup from 'yup'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI from '@/client/api/organization'
import { useAppSelector } from '@/store/hook'

type FormValues = {
  name: string
}

const NewOrganizationPage = () => {
  const navigate = useNavigate()
  const mutate = useAppSelector((state) => state.ui.organizations.mutate)
  const [isLoading, setIsLoading] = useState(false)
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })

  const handleSubmit = useCallback(
    async (
      { name }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>,
    ) => {
      setSubmitting(true)
      setIsLoading(true)
      try {
        const result = await OrganizationAPI.create({
          name,
        })
        mutate?.()
        setSubmitting(false)
        navigate(`/organization/${result.id}/member`)
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
        <title>New Organization</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5')}>
        <Heading className={cx('text-heading')}>New Organization</Heading>
        <Formik
          enableReinitialize={true}
          initialValues={{ name: '' }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <div className={cx('flex', 'flex-col', 'gap-3.5')}>
                <div className={cx('flex', 'flex-col', 'gap-1.5')}>
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
                </div>
                <div
                  className={cx('flex', 'flex-row', 'items-center', 'gap-1')}
                >
                  <Button
                    type="submit"
                    variant="solid"
                    colorScheme="blue"
                    isDisabled={isSubmitting || isLoading}
                    isLoading={isSubmitting || isLoading}
                  >
                    Create
                  </Button>
                  <Button as={Link} to="/organization" variant="solid">
                    Cancel
                  </Button>
                </div>
              </div>
            </Form>
          )}
        </Formik>
      </div>
    </>
  )
}

export default NewOrganizationPage
