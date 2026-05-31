export function SettingsPanelHeader(props: {
  title: string
  description?: string
  aside?: React.ReactNode
  badge?: React.ReactNode
}) {
  return (
    <div className="flex items-start justify-between gap-4">
      <div>
        <h2 className="flex items-center gap-2 text-2xl font-semibold text-foreground">
          {props.title}
          {props.badge}
        </h2>
        {props.description && <p className="mt-2 max-w-2xl text-sm text-muted-foreground">{props.description}</p>}
      </div>
      {props.aside}
    </div>
  )
}
