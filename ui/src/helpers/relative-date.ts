import dateFormat from 'dateformat'
import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'

let timeAgo: TimeAgo

export default function relativeDate(date: Date) {
  if (!timeAgo) {
    TimeAgo.addDefaultLocale(en)
    timeAgo = new TimeAgo('en-US')
  }
  const hoursDiff = Math.abs(new Date().getTime() - date.getTime()) / 3600000
  const isToday =
    new Date(date).setHours(0, 0, 0, 0) === new Date().setHours(0, 0, 0, 0)
  const isYesterday = new Date().getDate() - date.getDate() === 1
  const isThisYear = date.getFullYear() === new Date().getFullYear()
  if (hoursDiff <= 12 && isToday) {
    return timeAgo.format(date)
  } else if (isToday) {
    return 'Today'
  } else if (isYesterday) {
    return 'Yesterday'
  } else if (isThisYear) {
    return dateFormat(new Date(date), 'd mmm')
  } else {
    return dateFormat(new Date(date), 'd mmm yyyy')
  }
}
