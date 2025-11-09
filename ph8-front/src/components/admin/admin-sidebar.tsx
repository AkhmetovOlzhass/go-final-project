"use client"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { BookOpen, Users, FileText, Settings } from "lucide-react"

interface AdminSidebarProps {
  activeSection: string
  onSectionChange: (section: string) => void
}

const sidebarItems = [
  {
    id: "overview",
    label: "Overview",
    icon: Settings,
  },
  {
    id: "topics",
    label: "Topics",
    icon: BookOpen,
  },
  {
    id: "tasks",
    label: "Draft Tasks",
    icon: FileText,
  },
  {
    id: "tasks_all",
    label: "All Tasks",
    icon: FileText,
  },
  {
    id: "users",
    label: "Users",
    icon: Users,
  },
]

export function AdminSidebar({ activeSection, onSectionChange }: AdminSidebarProps) {
  return (
    <div className="w-64 bg-gray-900 border-r border-gray-700">
      <div className="p-6">
        <h2 className="text-lg font-semibold text-white">Admin Panel</h2>
      </div>
      <ScrollArea className="flex-1">
        <div className="px-3 py-2">
          <div className="space-y-1">
            {sidebarItems.map((item) => {
              const Icon = item.icon
              return (
                <Button
                  key={item.id}
                  variant="ghost"
                  className={cn(
                    "w-full justify-start text-gray-300 hover:bg-gray-800 hover:text-white transition-all duration-200",
                    activeSection === item.id && "bg-orange-600 text-white hover:bg-orange-700",
                  )}
                  onClick={() => onSectionChange(item.id)}
                >
                  <Icon className="mr-2 h-4 w-4" />
                  {item.label}
                </Button>
              )
            })}
          </div>
        </div>
      </ScrollArea>
    </div>
  )
}
