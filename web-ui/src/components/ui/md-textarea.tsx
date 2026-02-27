import { Textarea } from '@/components/ui/textarea'

interface MdTextareaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  value: string
  onChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => void
  rows?: number
}

const HINTS = [
  { code: '**жирный**',     title: 'жирный' },
  { code: '*курсив*',       title: 'курсив' },
  { code: '~~зачёркнутый~~', title: 'зачёркнутый' },
  { code: '[текст](url)',   title: 'ссылка' },
  { code: '`код`',          title: 'код' },
]

/**
 * Textarea for Markdown text.
 * Text is stored in standard Markdown; the backend converts to Telegram MarkdownV2 on send.
 */
export function MdTextarea({ value, onChange, rows = 4, ...props }: MdTextareaProps) {
  return (
    <div className="space-y-1.5">
      <div className="flex flex-wrap items-center gap-x-1.5 gap-y-1">
        <span className="text-xs text-muted-foreground">Markdown:</span>
        {HINTS.map(h => (
          <code
            key={h.code}
            className="text-xs bg-muted text-muted-foreground px-1.5 py-0.5 rounded font-mono leading-none"
            title={h.title}
          >
            {h.code}
          </code>
        ))}
      </div>
      <Textarea value={value} onChange={onChange} rows={rows} {...props} />
    </div>
  )
}
