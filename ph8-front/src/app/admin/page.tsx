"use client"

import { useState, useEffect } from "react"
import { useAuth } from "@/lib/auth"
import { useRouter } from "next/navigation"
import { AdminSidebar } from "@/components/admin/admin-sidebar"
import { OverviewSection } from "@/components/admin/overview-section"
import { TopicsSection } from "@/components/admin/topics-section"
import { TasksSection } from "@/components/admin/tasks-section"
import { UsersSection } from "@/components/admin/users-section"
import { Button } from "@/components/ui/button"
import { AllTasksSection } from "@/components/admin/all-tasks-section"

export default function AdminPage() {
  const { user, loading, logout } = useAuth()
  const router = useRouter()
  const [activeSection, setActiveSection] = useState("overview")

  useEffect(() => {
    if (!loading && (!user || user.role !== "ADMIN")) {
      router.push("/")
    }
  }, [user, loading, router])

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-lg text-gray-100">Loading...</div>
      </div>
    )
  }

  if (!user || user.role !== "ADMIN") {
    return null
  }

  const renderSection = () => {
    switch (activeSection) {
      case "overview":
        return <OverviewSection />
      case "topics":
        return <TopicsSection />
      case "tasks":
        return <TasksSection />
      case "tasks_all":
        return <AllTasksSection />
      case "users":
        return <UsersSection />
      default:
        return <OverviewSection />
    }
  }

  return (
    <div className="min-h-screen bg-gray-950 flex">
      <AdminSidebar activeSection={activeSection} onSectionChange={setActiveSection} />
      <div className="flex-1 flex flex-col">
        <header className="bg-gradient-to-r from-gray-900 to-gray-800 border-b border-gray-700 px-6 py-4">
          <div className="flex justify-between items-center">
            <h1 className="text-xl font-semibold text-white">Physics Learning Platform - Admin</h1>
            <div className="flex items-center gap-4">
              <span className="text-sm text-gray-300">Welcome, {user.name}</span>
              <Button
                variant="outline"
                size="sm"
                onClick={() => router.push("/")}
                className="border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white bg-transparent"
              >
                Student View
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={logout}
                className="border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white bg-transparent"
              >
                Logout
              </Button>
            </div>
          </div>
        </header>
        <main className="flex-1 bg-gray-950">{renderSection()}</main>
      </div>
    </div>
  )
}
