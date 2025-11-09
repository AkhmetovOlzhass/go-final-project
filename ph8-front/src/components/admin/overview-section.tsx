"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { BookOpen, FileText, Users, Eye } from "lucide-react"
import { contentAPI, usersAPI } from "@/lib/api"

export function OverviewSection() {
  const [stats, setStats] = useState({
    topics: 0,
    tasks: 0,
    users: 0,
    loading: true,
  })

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    try {
      const [topics, tasks, users] = await Promise.all([
        contentAPI.getAllTopics(),
        contentAPI.getDraftTasks(),
        usersAPI.getAllUsers(),
      ])
      console.log(topics)

      setStats({
        topics: topics.length,
        tasks: tasks.length,
        users: users.length,
        loading: false,
      })
    } catch (error) {
      setStats((prev) => ({ ...prev, loading: false }))
    }
  }

  if (stats.loading) {
    return <div className="p-6 text-gray-300">Loading overview...</div>
  }

  return (
    <div className="p-6 bg-gray-950 min-h-screen">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-white">Admin Overview</h1>
        <p className="text-gray-400">Physics Learning Platform dashboard</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card className="bg-gray-900 border-gray-700 hover:bg-gray-800 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-gray-300">Total Topics</CardTitle>
            <BookOpen className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stats.topics}</div>
            <p className="text-xs text-gray-400">Physics learning topics</p>
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-700 hover:bg-gray-800 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-gray-300">Draft Tasks</CardTitle>
            <FileText className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stats.tasks}</div>
            <p className="text-xs text-gray-400">Unpublished tasks</p>
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-700 hover:bg-gray-800 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-gray-300">Total Users</CardTitle>
            <Users className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stats.users}</div>
            <p className="text-xs text-gray-400">Registered students</p>
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-700 hover:bg-gray-800 transition-colors">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-gray-300">Platform Status</CardTitle>
            <Eye className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-400">Active</div>
            <p className="text-xs text-gray-400">System operational</p>
          </CardContent>
        </Card>
      </div>

      <div className="mt-8">
        <Card className="bg-gray-900 border-gray-700">
          <CardHeader>
            <CardTitle className="text-white">Quick Actions</CardTitle>
            <CardDescription className="text-gray-400">Common administrative tasks</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <p className="text-sm text-gray-300">• Create new physics topics for students to explore</p>
              <p className="text-sm text-gray-300">• Add interactive tasks and problems to existing topics</p>
              <p className="text-sm text-gray-300">• Review and publish draft content</p>
              <p className="text-sm text-gray-300">• Monitor user engagement and progress</p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
