// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useMemo } from 'react'
import { Tag } from '@chakra-ui/react'
import cx from 'classnames'
import parseEmailList from '@/lib/helpers/parse-email-list'

export type EmailTokenizerProps = {
  value: string
}

const EmailTokenizer = ({ value }: EmailTokenizerProps) => {
  const emails = useMemo(() => parseEmailList(value), [value])
  return (
    <>
      {emails.length > 0 ? (
        <div className={cx('flex', 'flex-wrap', 'gap-0.5')}>
          {emails.map((email, index) => (
            <Tag
              key={index}
              size="md"
              variant="solid"
              className={cx('rounded-full')}
            >
              {email}
            </Tag>
          ))}
        </div>
      ) : null}
    </>
  )
}

export default EmailTokenizer
