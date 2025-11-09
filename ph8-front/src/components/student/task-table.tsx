"use client"

import { useState } from "react"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Check, Search } from "lucide-react"
import type { Task, Topic } from "@/lib/api"

interface TaskTableProps {
  tasks: Task[]
  topics: Topic[]
  onViewTask: (taskId: string) => void
  userSolutions?: string[]
}

export function TaskTable({ tasks, topics, onViewTask, userSolutions = [] }: TaskTableProps) {
  const [searchQuery, setSearchQuery] = useState("")
  const [selectedDifficulty, setSelectedDifficulty] = useState<string>("all")
  const [selectedTopic, setSelectedTopic] = useState<string>("all")
  const [classSort, setClassSort] = useState<"none" | "asc" | "desc">("asc") // ← селект сортировки

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case "EASY":
        return "text-green-400 bg-green-900/30 border-green-700 hover:bg-green-900/50"
      case "MEDIUM":
        return "text-orange-400 bg-orange-900/30 border-orange-700 hover:bg-orange-900/50"
      case "HARD":
        return "text-red-400 bg-red-900/30 border-red-700 hover:bg-red-900/50"
      case "EXTREME":
        return "text-purple-400 bg-purple-900/30 border-purple-700 hover:bg-purple-900/50"
      default:
        return "text-gray-400 bg-gray-800 border-gray-600"
    }
  }

  const filteredTasks = tasks.filter((task) => {
    const matchesSearch = task.title.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesDifficulty = selectedDifficulty === "all" || task.difficulty === selectedDifficulty
    const matchesTopic = selectedTopic === "all" || task.topicId === selectedTopic
    return matchesSearch && matchesDifficulty && matchesTopic
  })

  const getTopicTitle = (topicId: string) => {
    return topics.find((topic) => topic.id === topicId)?.title || "Unknown"
  }

  const classMap: Record<string, string> = {
    SEVEN: "7",
    EIGHT: "8",
    NINE: "9",
    TEN: "10",
    ELEVEN: "11",
  }

  const getTopicClass = (topicId: string) => {
    const schoolClass = topics.find((topic) => topic.id === topicId)?.schoolClass
    return schoolClass ? classMap[schoolClass] : "Unknown"
  }

  const getClassNumber = (topicId: string) => {
    const c = parseInt(getTopicClass(topicId), 10)
    return Number.isNaN(c) ? 999 : c
  }

  const isTaskSolved = (taskId: string) => {
    return userSolutions.includes(taskId)
  }

  // Сортируем уже отфильтрованные задачи по выбранному направлению классов
  const tasksToRender =
    classSort === "none"
      ? filteredTasks
      : [...filteredTasks].sort((a, b) => {
          const ca = getClassNumber(a.topicId)
          const cb = getClassNumber(b.topicId)
          return classSort === "asc" ? ca - cb : cb - ca
        })

  // Calculate acceptance rate (mock data for now)
  const getAcceptanceRate = (taskId: string) => {
    const rates = ["45.2%", "52.1%", "38.7%", "61.3%", "29.4%", "71.8%"]
    return rates[taskId.length % rates.length]
  }

  return (
    <div className="bg-gray-800/50 backdrop-blur-sm rounded-xl border border-gray-700 overflow-hidden shadow-2xl">
      {/* Header with search and filters */}
      <div className="p-6 border-b border-gray-700 bg-gray-800/30">
        <div className="flex flex-col lg:flex-row gap-4 items-start lg:items-center justify-between">
          <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
            <Input
              placeholder="Search problems..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10 bg-gray-900/50 border-gray-600 text-white placeholder-gray-400 focus:border-orange-500 focus:ring-orange-500/20 transition-all duration-200"
            />
          </div>

          <div className="flex gap-3">
            <select
              value={selectedDifficulty}
              onChange={(e) => setSelectedDifficulty(e.target.value)}
              className="px-4 py-2 bg-gray-800/80 border border-gray-600 rounded-lg text-sm text-white focus:border-orange-500 focus:ring-2 focus:ring-orange-500/20 transition-all duration-200 hover:bg-gray-700/80 hover:border-gray-500 cursor-pointer appearance-none bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTIiIGhlaWdodD0iOCIgdmlld0JveD0iMCAwIDEyIDgiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxwYXRoIGQ9Ik0xIDFMNiA2TDExIDEiIHN0cm9rZT0iIzlDQTNBRiIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiLz4KPC9zdmc+')] bg-no-repeat bg-[right_12px_center] pr-10"
            >
              <option value="all" className="bg-gray-800 text-white">
                All Difficulties
              </option>
              <option value="EASY" className="bg-gray-800 text-white">
                Easy
              </option>
              <option value="MEDIUM" className="bg-gray-800 text-white">
                Medium
              </option>
              <option value="HARD" className="bg-gray-800 text-white">
                Hard
              </option>
              <option value="EXTREME" className="bg-gray-800 text-white">
                Extreme
              </option>
            </select>

            <select
              value={selectedTopic}
              onChange={(e) => setSelectedTopic(e.target.value)}
              className="px-4 py-2 bg-gray-800/80 border border-gray-600 rounded-lg text-sm text-white focus:border-orange-500 focus:ring-2 focus:ring-orange-500/20 transition-all duration-200 hover:bg-gray-700/80 hover:border-gray-500 cursor-pointer appearance-none bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTIiIGhlaWdodD0iOCIgdmlld0JveD0iMCAwIDEyIDgiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxwYXRoIGQ9Ik0xIDFMNiA2TDExIDEiIHN0cm9rZT0iIzlDQTNBRiIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiLz4KPC9zdmc+')] bg-no-repeat bg-[right_12px_center] pr-10"
            >
              <option value="all" className="bg-gray-800 text-white">
                All Topics
              </option>
              {topics.map((topic) => (
                <option key={topic.id} value={topic.id} className="bg-gray-800 text-white">
                  {topic.title}
                </option>
              ))}
            </select>

            {/* ▼ Новый селект сортировки по классам */}
            <select
              value={classSort}
              onChange={(e) => setClassSort(e.target.value as "none" | "asc" | "desc")}
              className="px-4 py-2 bg-gray-800/80 border border-gray-600 rounded-lg text-sm text-white focus:border-orange-500 focus:ring-2 focus:ring-orange-500/20 transition-all duration-200 hover:bg-gray-700/80 hover:border-gray-500 cursor-pointer appearance-none bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTIiIGhlaWdodD0iOCIgdmlld0JveD0iMCAwIDEyIDgiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjxwYXRoIGQ9Ik0xIDFMNiA2TDExIDEiIHN0cm9rZT0iIzlDQTNBRiIgc3Ryb2tlLXdpZHRoPSIyIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiLz4KPC9zdmc+')] bg-no-repeat bg-[right_12px_center] pr-10"
              title="Class Sort"
              aria-label="Sort by class"
            >
              <option value="asc" className="bg-gray-800 text-white">
                Class ↑ (7→11)
              </option>
              <option value="desc" className="bg-gray-800 text-white">
                Class ↓ (11→7)
              </option>
            </select>
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-800/50 border-b border-gray-700">
            <tr>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-300 uppercase tracking-wider w-16">
                Status
              </th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-300 uppercase tracking-wider">
                Title
              </th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-300 uppercase tracking-wider w-28 hidden sm:table-cell">
                Class
              </th>
              {/* <th className="px-6 py-4 text-left text-xs font-semibold text-gray-300 uppercase tracking-wider w-28 hidden sm:table-cell">
                Acceptance
              </th> */}
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-300 uppercase tracking-wider w-28">
                Difficulty
              </th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-gray-300 uppercase tracking-wider w-36 hidden md:table-cell">
                Topic
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-700">
            {tasksToRender.map((task, index) => (
              <tr
                key={task.id}
                className="hover:bg-gray-700/30 cursor-pointer transition-all duration-200 group hover:shadow-lg hover:shadow-orange-500/5"
                onClick={() => onViewTask(task.id)}
              >
                <td className="px-6 py-5 whitespace-nowrap">
                  <div className="flex items-center justify-center">
                    {isTaskSolved(task.id) ? (
                      <div className="w-6 h-6 bg-gradient-to-r from-green-500 to-green-400 rounded-full flex items-center justify-center shadow-lg shadow-green-500/25 group-hover:scale-110 transition-transform duration-200">
                        <Check className="w-3.5 h-3.5 text-white" />
                      </div>
                    ) : (
                      <div className="w-6 h-6 border-2 border-gray-500 rounded-full group-hover:border-orange-400 transition-colors duration-200"></div>
                    )}
                  </div>
                </td>

                <td className="px-6 py-5">
                  <div className="flex flex-col">
                    <div className="text-sm font-semibold text-white group-hover:text-orange-400 transition-colors duration-200 text-balance">
                      {index + 1}. {task.title}
                    </div>
                  </div>
                </td>
                <td className="px-6 py-5 whitespace-nowrap text-sm text-gray-400 hidden md:table-cell font-medium">
                  {getTopicClass(task.topicId)}
                </td>
                {/* <td className="px-6 py-5 whitespace-nowrap text-sm text-gray-400 hidden sm:table-cell font-medium">
                  {getAcceptanceRate(task.id)}
                </td> */}

                <td className="px-6 py-5 whitespace-nowrap">
                  <Badge
                    className={`${getDifficultyColor(task.difficulty)} text-xs font-semibold px-3 py-1.5 transition-all duration-200 border`}
                  >
                    {task.difficulty}
                  </Badge>
                </td>

                <td className="px-6 py-5 whitespace-nowrap text-sm text-gray-400 hidden md:table-cell font-medium">
                  {getTopicTitle(task.topicId)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {tasksToRender.length === 0 && (
        <div className="text-center py-16">
          <div className="text-gray-400 text-sm">No problems found matching your criteria.</div>
          <div className="text-gray-500 text-xs mt-2">Try adjusting your search or filters.</div>
        </div>
      )}
    </div>
  )
}
