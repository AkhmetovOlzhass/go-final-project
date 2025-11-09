"use client"

import type React from "react"
import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Plus, Edit, Trash } from "lucide-react"
import { contentAPI, type Topic } from "@/lib/api"

export function TopicsSection() {
  const [topics, setTopics] = useState<Topic[]>([])
  const [loading, setLoading] = useState(true)

  // Create
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [schoolClass, setSchoolClass] = useState("TEN")
  const [editSchoolClass, setEditSchoolClass] = useState("TEN")
  const [title, setTitle] = useState("")
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState("")

  // Edit
  const [isEditOpen, setIsEditOpen] = useState(false)
  const [editTitle, setEditTitle] = useState("")
  const [selectedTopic, setSelectedTopic] = useState<Topic | null>(null)
  const [updating, setUpdating] = useState(false)

  useEffect(() => {
    loadTopics()
  }, [])

  const loadTopics = async () => {
    try {
      const data = await contentAPI.getAllTopics()
      setTopics(data)
    } catch (error) {
      setError("Failed to load topics")
    } finally {
      setLoading(false)
    }
  }

  const handleCreateTopic = async (e: React.FormEvent) => {
    e.preventDefault()
    setCreating(true)
    setError("")
    try {
      await contentAPI.createTopic({ title, schoolClass })
      setTitle("")
      setIsCreateOpen(false)
      loadTopics()
    } catch (error) {
      setError("Failed to create topic")
    } finally {
      setCreating(false)
    }
  }

  const handleEditClick = (topic: Topic) => {
    setSelectedTopic(topic)
    setEditTitle(topic.title)
    setEditSchoolClass(topic.schoolClass)
    setIsEditOpen(true)
  }

  const handleUpdateTopic = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedTopic) return
    setUpdating(true)
    try {
      await contentAPI.updateTopic(selectedTopic.id, { title: editTitle, schoolClass: editSchoolClass })
      setIsEditOpen(false)
      setSelectedTopic(null)
      loadTopics()
    } catch (error) {
      setError("Failed to update topic")
    } finally {
      setUpdating(false)
    }
  }

  const handleDeleteTopic = async (id: string) => {
    if (!confirm("Are you sure you want to delete this topic?")) return
    try {
      await contentAPI.deleteTopic(id)
      loadTopics()
    } catch (error) {
      setError("Failed to delete topic")
    }
  }

  if (loading) {
    return <div className="p-6 text-gray-300">Loading topics...</div>
  }

  return (
    <div className="p-6 bg-gray-950 min-h-screen">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">Topics</h1>
          <p className="text-gray-400">Manage physics learning topics</p>
        </div>

        {/* Create Topic */}
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button className="bg-orange-600 hover:bg-orange-700 text-white">
              <Plus className="mr-2 h-4 w-4" />
              Create Topic
            </Button>
          </DialogTrigger>
          <DialogContent className="bg-gray-900 border-gray-700">
            <DialogHeader>
              <DialogTitle className="text-white">Create New Topic</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleCreateTopic} className="space-y-4">
              {error && (
                <div className="bg-red-900/50 border border-red-700 text-red-300 px-4 py-2 rounded-md text-sm">
                  {error}
                </div>
              )}
              <div className="space-y-2">
                <Label htmlFor="title" className="text-gray-300">Title</Label>
                <Input
                  id="title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="Enter topic title"
                  required
                  className="bg-gray-800 border-gray-600 text-white placeholder:text-gray-400"
                />
              </div>
              <div className="space-y-2">
  <Label htmlFor="schoolClass" className="text-gray-300">Class</Label>
  <select
    id="schoolClass"
    value={schoolClass}
    onChange={(e) => setSchoolClass(e.target.value)}
    className="w-full rounded-md bg-gray-800 border border-gray-600 text-white px-3 py-2"
  >
    <option value="SEVEN">7 класс</option>
    <option value="EIGHT">8 класс</option>
    <option value="NINE">9 класс</option>
    <option value="TEN">10 класс</option>
    <option value="ELEVEN">11 класс</option>
  </select>
</div>
              <Button type="submit" disabled={creating} className="w-full bg-orange-600 hover:bg-orange-700 text-white">
                {creating ? "Creating..." : "Create Topic"}
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {/* Topics grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {topics.map((topic) => (
          <Card key={topic.id} className="bg-gray-900 border-gray-700 hover:bg-gray-800 transition-colors">
            <CardHeader>
              <CardTitle className="text-lg text-white">{topic.title}</CardTitle>
              <CardDescription className="text-gray-400">{topic.description}</CardDescription>
              <CardDescription className="text-gray-400">
                {topic.schoolClass ? `Class: ${topic.schoolClass}` : "No class"}
              </CardDescription>
            </CardHeader>
            <CardContent className="flex justify-between items-center">
              <p className="text-sm text-gray-500">Created: {new Date(topic.createdAt).toLocaleDateString()}</p>
              <div className="flex gap-2">
                <Button variant="ghost" size="sm" onClick={() => handleEditClick(topic)} className="text-gray-300 hover:text-white">
                  <Edit className="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="sm" onClick={() => handleDeleteTopic(topic.id)} className="text-red-400 hover:text-red-600">
                  <Trash className="h-4 w-4" />
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Edit Dialog */}
      <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
        <DialogContent className="bg-gray-900 border-gray-700">
          <DialogHeader>
            <DialogTitle className="text-white">Edit Topic</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleUpdateTopic} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="editTitle" className="text-gray-300">Title</Label>
              <Input
                id="editTitle"
                value={editTitle}
                onChange={(e) => setEditTitle(e.target.value)}
                placeholder="Enter topic title"
                required
                className="bg-gray-800 border-gray-600 text-white placeholder:text-gray-400"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="editSchoolClass" className="text-gray-300">Class</Label>
              <select
                id="editSchoolClass"
                value={editSchoolClass}
                onChange={(e) => setEditSchoolClass(e.target.value)}
                className="w-full rounded-md bg-gray-800 border border-gray-600 text-white px-3 py-2"
              >
                <option value="SEVEN">7 класс</option>
                <option value="EIGHT">8 класс</option>
                <option value="NINE">9 класс</option>
                <option value="TEN">10 класс</option>
                <option value="ELEVEN">11 класс</option>
              </select>
            </div>
            <Button type="submit" disabled={updating} className="w-full bg-orange-600 hover:bg-orange-700 text-white">
              {updating ? "Updating..." : "Update Topic"}
            </Button>
          </form>
        </DialogContent>
      </Dialog>

      {topics.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-400">No topics created yet. Create your first topic to get started.</p>
        </div>
      )}
    </div>
  )
}
