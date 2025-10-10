import { readable } from 'svelte/store'

interface TimerOptions {
  interval?: number
  format?: Intl.DateTimeFormatOptions
  locale?: string
}

export default function(options: TimerOptions = {}) {
  const {locale = 'en-US', format = {
       month: 'short',
       day: 'numeric',
       year: 'numeric',
       hour: '2-digit',
       minute: '2-digit',
       second: '2-digit'
  }} = options

  return readable<string>(new Intl.DateTimeFormat(locale, format).format(new Date()), set => {
    const update = () => set(new Intl.DateTimeFormat(locale, format).format(new Date()))
    update()
    const intervalId = setInterval(update, options.interval || 500)

    return () => clearInterval(intervalId)
  })
}

