// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Card,
  CardBody,
  CircularProgress,
  IconButton,
  Text,
} from '@chakra-ui/react'
import { IconClose, IconSchedule, IconCheckCircle, IconError } from '@koupr/ui'
import cx from 'classnames'
import truncateMiddle from '@/lib/helpers/truncate-middle'
import { Upload, UploadDecorator, uploadRemoved } from '@/store/entities/uploads'
import { useAppDispatch } from '@/store/hook'

export type UploadsItemProps = {
  upload: Upload
}

const UploadItem = ({ upload: uploadProp }: UploadsItemProps) => {
  const dispatch = useAppDispatch()
  const upload = new UploadDecorator(uploadProp)

  return (
    <Card variant="outline">
      <CardBody>
        <div className={cx('flex', 'flex-col', 'gap-1')}>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1.5')}>
            {upload.isProgressing ? (
              <CircularProgress
                value={upload.progress}
                max={100}
                isIndeterminate={upload.progress === 100 && !upload.error}
                className={cx('text-black')}
                size="20px"
              />
            ) : null}
            {upload.isPending ? <IconSchedule className={cx('shrink-0', 'text-gray-500')} /> : null}
            {upload.isSucceeded ? <IconCheckCircle className={cx('shrink-0', 'text-green-500')} filled={true} /> : null}
            {upload.isFailed ? (
              <div className={cx('shrink-0', 'text-red-500')}>
                <IconError filled={true} />
              </div>
            ) : null}
            <div className={cx('flex', 'flex-col', 'grow')}>
              <span className={cx('font-semibold')}>{truncateMiddle(upload.blob.name, 40)}</span>
            </div>
            <IconButton
              icon={<IconClose />}
              size="xs"
              variant="outline"
              colorScheme={upload.isProgressing ? 'red' : 'gray'}
              title="Remove"
              aria-label="Remove"
              onClick={() => {
                upload.request?.abort()
                dispatch(uploadRemoved(upload.id))
              }}
            />
          </div>
          {upload.isFailed && upload.error ? (
            <Accordion allowMultiple>
              <AccordionItem className={cx('border-none')}>
                <AccordionButton className={cx('p-0.5')}>
                  <div className={cx('flex', 'flex-row', 'w-full')}>
                    <span className={cx('text-left', 'grow')}>Upload failed, click to show error</span>
                    <AccordionIcon />
                  </div>
                </AccordionButton>
                <AccordionPanel className={cx('p-0.5')}>
                  <Text
                    className={cx('text-red-500')}
                    dangerouslySetInnerHTML={{ __html: upload.error }}
                    noOfLines={5}
                  ></Text>
                </AccordionPanel>
              </AccordionItem>
            </Accordion>
          ) : null}
        </div>
      </CardBody>
    </Card>
  )
}

export default UploadItem
