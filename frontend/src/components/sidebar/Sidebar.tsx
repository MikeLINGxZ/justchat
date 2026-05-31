import { SidebarHeader } from './SidebarHeader'
import { ConversationList } from './ConversationList'
import { SidebarFooter } from './SidebarFooter'

interface SidebarProps {
  onCollapse: () => void
}

export function Sidebar({ onCollapse }: SidebarProps) {
  return (
    <div className="flex flex-col h-full overflow-visible border-r border-border bg-background">
      <SidebarHeader onCollapse={onCollapse} />
      <ConversationList />
      <SidebarFooter />
    </div>
  )
}
