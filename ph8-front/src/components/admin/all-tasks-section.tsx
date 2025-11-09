"use client"

import { useState, useEffect } from "react"
import { contentAPI, type Task, type Topic } from "@/lib/api"
import { TaskCard } from "@/components/admin/task-card"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Button } from "@/components/ui/button"

export function AllTasksSection() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [topics, setTopics] = useState<Topic[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")


  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [tasksData, topicsData] = await Promise.all([
        contentAPI.getAllTasks(),
        contentAPI.getAllTopics(),
      ])
      setTasks(tasksData)
      setTopics(topicsData)
    } catch (error) {
      setError("Failed to load data")
    } finally {
      setLoading(false)
    }
  }

  const handleDeleteTask = async (id: string) => {
    if (!confirm("Are you sure you want to delete this task?")) return
    try {
      await contentAPI.deleteTask(id)
      loadData()
    } catch (error) {
      setError("Failed to delete task")
    }
  }

  if (loading) return <div className="p-6 text-gray-300">Loading tasks...</div>

  return (
    <div className="p-6 bg-gray-950 min-h-screen">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">Tasks</h1>
          <p className="text-gray-400">Manage physics learning tasks</p>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {tasks.map((task) => (
        <TaskCard
          key={task.id}
          task={task}
          topics={topics}
          onDelete={handleDeleteTask}
          onUpdate={async (id, data) => {
            await contentAPI.updateTask(id, data)
            loadData()
          }}
        />
        ))}
      </div>

      {tasks.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-400">No tasks created yet. Create your first task to get started.</p>
        </div>
      )}
      
    </div>
  )
}
