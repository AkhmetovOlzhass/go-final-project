"use client"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import { BookOpen, Target, Clock, Zap, TrendingUp } from "lucide-react"
import type { Topic } from "@/lib/api"

interface ProblemSidebarProps {
  topics: Topic[]
  selectedTopic: string
  onTopicSelect: (topicId: string) => void
  stats?: {
    total: number
    solved: number
    easy: number
    medium: number
    hard: number
    extreme: number
  }
  isMobile?: boolean
}

export function ProblemSidebar({
  topics,
  selectedTopic,
  onTopicSelect,
  stats = { total: 0, solved: 0, easy: 0, medium: 0, hard: 0, extreme: 0 },
  isMobile = false,
}: ProblemSidebarProps) {
  const getDifficultyIcon = (difficulty: string) => {
    switch (difficulty) {
      case "easy":
        return <Target className="h-4 w-4 text-green-400" />
      case "medium":
        return <Clock className="h-4 w-4 text-orange-400" />
      case "hard":
        return <BookOpen className="h-4 w-4 text-red-400" />
      case "extreme":
        return <Zap className="h-4 w-4 text-purple-400" />
      default:
        return <Target className="h-4 w-4" />
    }
  }

  const progressPercentage = stats.total > 0 ? (stats.solved / stats.total) * 100 : 0

  return (
    <div
      className={`w-80 bg-gray-800/30 backdrop-blur-sm border-r border-gray-700 ${
        isMobile ? "h-full overflow-y-auto" : "hidden lg:block h-full overflow-y-auto"
      }`}
    >
      <div className="p-4 sm:p-6 space-y-4 sm:space-y-6">
        {/* Progress Stats */}
        <Card className="bg-gray-800/50 border-gray-700 shadow-xl hover:shadow-2xl transition-shadow duration-300">
          <CardHeader className="pb-4">
            <CardTitle className="text-lg flex items-center gap-2 text-white">
              <TrendingUp className="h-5 w-5 text-orange-400" />
              Progress Stats
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex justify-between items-center">
              <span className="text-2xl sm:text-3xl font-bold text-white">{stats.solved}</span>
              <span className="text-sm text-gray-400">/ {stats.total} Solved</span>
            </div>
            <div className="space-y-2">
              <Progress value={progressPercentage} className="h-3 bg-gray-700" />
              <div className="text-xs text-gray-400 text-center">{progressPercentage.toFixed(1)}% Complete</div>
            </div>

            <div className="space-y-3">
              <div className="flex justify-between items-center p-3 rounded-lg bg-gray-700/30 hover:bg-gray-700/50 transition-all duration-200 hover:scale-[1.02]">
                <div className="flex items-center gap-2">
                  {getDifficultyIcon("easy")}
                  <span className="text-sm text-gray-300">Easy</span>
                </div>
                <span className="text-sm font-semibold text-green-400">{stats.easy}</span>
              </div>

              <div className="flex justify-between items-center p-3 rounded-lg bg-gray-700/30 hover:bg-gray-700/50 transition-all duration-200 hover:scale-[1.02]">
                <div className="flex items-center gap-2">
                  {getDifficultyIcon("medium")}
                  <span className="text-sm text-gray-300">Medium</span>
                </div>
                <span className="text-sm font-semibold text-orange-400">{stats.medium}</span>
              </div>

              <div className="flex justify-between items-center p-3 rounded-lg bg-gray-700/30 hover:bg-gray-700/50 transition-all duration-200 hover:scale-[1.02]">
                <div className="flex items-center gap-2">
                  {getDifficultyIcon("hard")}
                  <span className="text-sm text-gray-300">Hard</span>
                </div>
                <span className="text-sm font-semibold text-red-400">{stats.hard}</span>
              </div>

              <div className="flex justify-between items-center p-3 rounded-lg bg-gray-700/30 hover:bg-gray-700/50 transition-all duration-200 hover:scale-[1.02]">
                <div className="flex items-center gap-2">
                  {getDifficultyIcon("extreme")}
                  <span className="text-sm text-gray-300">Extreme</span>
                </div>
                <span className="text-sm font-semibold text-purple-400">{stats.extreme}</span>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Topics Filter */}
        <Card className="bg-gray-800/50 border-gray-700 shadow-xl hover:shadow-2xl transition-shadow duration-300">
          <CardHeader className="pb-4">
            <CardTitle className="text-lg flex items-center gap-2 text-white">
              <BookOpen className="h-5 w-5 text-orange-400" />
              Topics
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <button
                onClick={() => onTopicSelect("all")}
                className={`w-full text-left px-4 py-3 rounded-lg text-sm transition-all duration-200 font-medium hover:scale-[1.02] ${
                  selectedTopic === "all"
                    ? "bg-gradient-to-r from-orange-500 to-orange-600 text-white shadow-lg shadow-orange-500/25"
                    : "text-gray-300 hover:bg-gray-700/50 hover:text-white"
                }`}
              >
                All Topics
              </button>

              {topics.map((topic) => (
                <button
                  key={topic.id}
                  onClick={() => onTopicSelect(topic.id)}
                  className={`w-full text-left px-4 py-3 rounded-lg text-sm transition-all duration-200 font-medium hover:scale-[1.02] ${
                    selectedTopic === topic.id
                      ? "bg-gradient-to-r from-orange-500 to-orange-600 text-white shadow-lg shadow-orange-500/25"
                      : "text-gray-300 hover:bg-gray-700/50 hover:text-white"
                  }`}
                >
                  {topic.title}
                </button>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
