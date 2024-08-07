// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import fs from 'fs'
import Handlebars from 'handlebars'
import yaml from 'js-yaml'
import nodemailer from 'nodemailer'
import path from 'path'
import { getConfig } from '@/config/config'

type MessageParams = {
  subject: string
}

const config = getConfig().smtp

const transporter = nodemailer.createTransport({
  host: config.host,
  port: config.port,
  secure: config.secure,
  auth:
    config.username || config.password
      ? {
          user: config.username,
          pass: config.password,
        }
      : null,
})

export function sendTemplateMail(
  template: string,
  address: string,
  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  variables: Record<string, any>,
) {
  const params = yaml.load(
    fs.readFileSync(path.join('templates', template, 'params.yml'), 'utf8'),
  ) as MessageParams
  const html = Handlebars.compile(
    fs.readFileSync(path.join('templates', template, 'template.hbs'), 'utf8'),
  )(variables)
  const text = Handlebars.compile(
    fs.readFileSync(path.join('templates', template, 'template.txt'), 'utf8'),
  )(variables)
  return new Promise<void>((resolve, reject) => {
    transporter.sendMail(
      {
        from: `"${config.senderName}" <${config.senderAddress}>`,
        to: address,
        subject: params.subject,
        text,
        html,
      },
      (err) => {
        if (err) {
          reject(err)
        } else {
          resolve()
        }
      },
    )
  })
}
