// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import React, { ReactNode, useEffect, useRef } from 'react'
import cx from 'classnames'

interface TextProps extends React.HTMLAttributes<HTMLSpanElement> {
  children?: ReactNode
  noOfLines?: number
  maxCharacters?: number
}

const Text: React.FC<TextProps> = ({
  children,
  noOfLines,
  maxCharacters,
  className,
  ...props
}: TextProps) => {
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const text = children?.toString() || ''
    if (ref.current && maxCharacters && text.length > maxCharacters) {
      ref.current.textContent = text.slice(0, maxCharacters).trim() + '…'
    }
  }, [children, maxCharacters])

  return (
    <span
      {...props}
      ref={ref}
      style={{
        display: noOfLines !== undefined ? '-webkit-box' : undefined,
        WebkitBoxOrient: noOfLines !== undefined ? 'vertical' : undefined,
        WebkitLineClamp: noOfLines,
      }}
      className={cx(
        { 'whitespace-nowrap': maxCharacters !== undefined },
        {
          'overflow-hidden':
            noOfLines !== undefined || maxCharacters !== undefined,
        },
        className,
      )}
    >
      {children}
    </span>
  )
}

export default Text
