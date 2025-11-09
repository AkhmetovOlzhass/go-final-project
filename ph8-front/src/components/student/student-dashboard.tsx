"use client"

import { useState, useEffect } from "react"
import { ProblemSidebar } from "./problem-sidebar"
import { TaskTable } from "./task-table"
import { TaskViewer } from "./task-viewer"
import { contentAPI, type Task, type Topic } from "@/lib/api"
import { Menu, X } from "lucide-react"

export function StudentDashboard() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [topics, setTopics] = useState<Topic[]>([])
  const [filteredTasks, setFilteredTasks] = useState<Task[]>([])
  const [selectedTopicId, setSelectedTopicId] = useState<string>("all")
  const [searchQuery, setSearchQuery] = useState("")
  const [viewingTaskId, setViewingTaskId] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [userSolutions, setUserSolutions] = useState<string[]>([])
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const handleSolved = (taskId: string) => {
    setUserSolutions(prev => (prev.includes(taskId) ? prev : [...prev, taskId]))
  }

  useEffect(() => {
    loadData()
  }, [])

  useEffect(() => {
    filterTasks()
  }, [tasks, selectedTopicId, searchQuery])

  const loadData = async () => {
    try {
      const [topicsData] = await Promise.all([contentAPI.getAllTopics()])
      setTopics(topicsData)

      const allTasks: Task[] = []
      for (const topic of topicsData) {
        try {
          const topicTasks = await contentAPI.getTasksByTopic(topic.id)
          const publishedTasks = topicTasks.filter((task) => task.status === "PUBLISHED")
          allTasks.push(...publishedTasks)
        } catch (error) {
        }
      }

      const difficultyOrder: Record<string, number> = {
        EASY: 1,
        MEDIUM: 2,
        HARD: 3,
        EXTREME: 4,
      }

      allTasks.sort((a, b) => {
        const diffComp = difficultyOrder[a.difficulty] - difficultyOrder[b.difficulty]
        if (diffComp !== 0) return diffComp

        return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
      })

      setTasks(allTasks)

      const progress = await contentAPI.getUserProgress();
      const solvedTaskIds = progress
        .filter((p: any) => p.status === "SOLVED")
        .map((p: any) => p.taskId);

      setUserSolutions(solvedTaskIds);
    } catch (error) {
      console.error("Failed to load data:", error)
    } finally {
      setLoading(false)
    }
  }


  const filterTasks = () => {
    let filtered = tasks

    // Filter by topic
    if (selectedTopicId && selectedTopicId !== "all") {
      filtered = filtered.filter((task) => task.topicId === selectedTopicId)
    }

    // Filter by search query
    if (searchQuery) {
      const query = searchQuery.toLowerCase()
      filtered = filtered.filter(
        (task) =>
          task.title.toLowerCase().includes(query) ||
          task.description.toLowerCase().includes(query) ||
          task.content.toLowerCase().includes(query)
      )
    }

    setFilteredTasks(filtered)
  }

  const handleViewTask = (taskId: string) => {
    setViewingTaskId(taskId)
  }

  const handleBackToTasks = () => {
    setViewingTaskId(null)
  }

  if (viewingTaskId) {
    return <TaskViewer taskId={viewingTaskId} onBack={handleBackToTasks} onSolved={handleSolved} />
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="text-center space-y-4">
          <div className="relative">
            <div className="animate-spin rounded-full h-12 w-12 border-2 border-gray-700 mx-auto"></div>
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-orange-500 mx-auto absolute top-0 left-1/2 transform -translate-x-1/2"></div>
          </div>
          <div className="text-lg text-gray-300 animate-pulse">Loading physics challenges...</div>
        </div>
      </div>
    )
  }

  const difficultyStats = tasks.reduce(
    (acc, task) => {
      acc[task.difficulty.toLowerCase()] = (acc[task.difficulty.toLowerCase()] || 0) + 1
      return acc
    },
    {} as Record<string, number>,
  )

  const stats = {
    total: tasks.length,
    solved: userSolutions.length,
    easy: difficultyStats.easy || 0,
    medium: difficultyStats.medium || 0,
    hard: difficultyStats.hard || 0,
    extreme: difficultyStats.extreme || 0,
  }

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900 border-b border-gray-700">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 py-8 sm:py-12">
          <div className="text-center space-y-4 sm:space-y-6">
            <div className="flex items-center justify-between lg:hidden mb-4">
              <button
                onClick={() => setSidebarOpen(true)}
                className="p-2 text-gray-400 hover:text-white hover:bg-gray-800 rounded-lg transition-colors duration-200"
              >
                <Menu className="h-6 w-6" />
              </button>
              <div className="text-sm text-gray-400 bg-gray-800 px-3 py-1 rounded-full">
                {filteredTasks.length} challenges
              </div>
            </div>

            <h1 className="text-3xl sm:text-4xl md:text-6xl font-bold text-white text-balance">
              Physics{" "}
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-orange-400 to-orange-600">
                Challenges
              </span>
            </h1>
            <p className="text-lg sm:text-xl text-gray-300 max-w-2xl mx-auto text-balance leading-relaxed">
              Master physics through interactive problem-solving. Challenge yourself with real-world scenarios and build
              your understanding step by step.
            </p>
          </div>
        </div>
      </div>

      <div className="flex h-[calc(100vh-200px)]">
        {sidebarOpen && (
          <div className="fixed inset-0 z-50 lg:hidden">
            <div className="fixed inset-0 bg-black/50" onClick={() => setSidebarOpen(false)} />
            <div className="fixed left-0 top-0 h-full w-80 bg-gray-800 border-r border-gray-700 z-50">
              <div className="flex items-center justify-between p-4 border-b border-gray-700">
                <h2 className="text-lg font-semibold text-white">Menu</h2>
                <button
                  onClick={() => setSidebarOpen(false)}
                  className="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors duration-200"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>
              <ProblemSidebar
                topics={topics}
                selectedTopic={selectedTopicId}
                onTopicSelect={(topicId) => {
                  setSelectedTopicId(topicId)
                  setSidebarOpen(false)
                }}
                stats={stats}
                isMobile={true}
              />
            </div>
          </div>
        )}

        <ProblemSidebar
          topics={topics}
          selectedTopic={selectedTopicId}
          onTopicSelect={setSelectedTopicId}
          stats={stats}
        />

        <div className="flex-1 p-4 sm:p-6 lg:p-8 overflow-y-auto">
          <div className="mb-6 sm:mb-8">
            <div className="flex items-center justify-between mb-4 sm:mb-6">
              <h2 className="text-2xl sm:text-3xl font-bold text-white">All Challenges</h2>
              <div className="text-sm text-gray-400 bg-gray-800 px-3 py-1 rounded-full hidden lg:block">
                {filteredTasks.length} challenge{filteredTasks.length !== 1 ? "s" : ""} found
              </div>
            </div>
          </div>

          <TaskTable tasks={filteredTasks} topics={topics} onViewTask={handleViewTask} userSolutions={userSolutions} />
        </div>
      </div>
    </div>
  )
}
