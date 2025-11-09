"use client"

import { useAuth } from "@/lib/auth"
import { useRouter } from "next/navigation"
import { StudentDashboard } from "@/components/student/student-dashboard"
import { Button } from "@/components/ui/button"
import { LogOut, Settings, UserCog } from "lucide-react"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

export default function LearnPage() {
  const { user, loading, logout } = useAuth()
  const router = useRouter()

  if (loading || !user) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-center space-y-4">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto"></div>
          <div className="text-lg text-gray-300">Loading...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-950">
      <header className="bg-gradient-to-r from-gray-900 to-gray-800 border-b border-gray-700 sticky top-0 z-50 backdrop-blur-sm">
        <div className="flex justify-between items-center max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-gradient-to-r from-orange-500 to-orange-600 rounded-lg flex items-center justify-center shadow-lg">
              <span className="text-white font-bold text-sm">P</span>
            </div>
            <h1 className="text-xl font-bold text-white">Physics Learning Platform</h1>
          </div>

          <div className="flex items-center gap-4">
            {user.role === "Admin" && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => router.push("/admin")}
                className="text-gray-300 hover:text-white hover:bg-gray-700 transition-all duration-200 cursor-pointer"
              >
                <Settings className="mr-2 h-4 w-4" />
                Admin Panel
              </Button>
            )}

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <button className="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-700/50 transition-all duration-200 cursor-pointer group">
                  <Avatar className="h-9 w-9 ring-2 ring-orange-500/50 group-hover:ring-orange-500 transition-all">
                    <AvatarImage src={user.avatarUrl || "/placeholder.svg"} alt={user.displayName} />
                    <AvatarFallback className="bg-gradient-to-br from-orange-500 to-orange-600 text-white font-semibold">
                      {user.displayName.charAt(0).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex flex-col items-start">
                    <span className="text-sm font-medium text-white group-hover:text-orange-400 transition-colors">
                      {user.displayName}
                    </span>
                    <span className="text-xs text-gray-400">{user.email}</span>
                  </div>
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-56 bg-gray-800 border-gray-700">
                <DropdownMenuLabel className="text-gray-300">My Account</DropdownMenuLabel>
                <DropdownMenuSeparator className="bg-gray-700" />
                <DropdownMenuItem
                  onClick={() => router.push("/profile")}
                  className="text-gray-300 hover:text-white hover:bg-gray-700 cursor-pointer focus:bg-gray-700 focus:text-white"
                >
                  <UserCog className="mr-2 h-4 w-4" />
                  Edit Profile
                </DropdownMenuItem>
                <DropdownMenuSeparator className="bg-gray-700" />
                <DropdownMenuItem
                  onClick={async () => {
                    await logout()
                    router.push("/auth")
                  }}
                  className="text-red-400 hover:text-red-300 hover:bg-red-500/10 cursor-pointer focus:bg-red-500/10 focus:text-red-300"
                  variant="destructive"
                >
                  <LogOut className="mr-2 h-4 w-4" />
                  Logout
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>

      <main>
        <StudentDashboard />
      </main>
    </div>
  )
}
