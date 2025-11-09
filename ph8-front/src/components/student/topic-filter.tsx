"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { BookOpen, Filter, ChevronDown } from "lucide-react"
import { contentAPI, type Topic } from "@/lib/api"

interface TopicFilterProps {
  selectedTopicId: string | null
  onTopicSelect: (topicId: string | null) => void
}

export function TopicFilter({ selectedTopicId, onTopicSelect }: TopicFilterProps) {
  const [topics, setTopics] = useState<Topic[]>([])
  const [loading, setLoading] = useState(true)
  const [isExpanded, setIsExpanded] = useState(true)

  useEffect(() => {
    loadTopics()
  }, [])

  const loadTopics = async () => {
    try {
      const data = await contentAPI.getAllTopics()
      setTopics(data)
    } catch (error) {
      console.error("Failed to load topics:", error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <Card className="bg-slate-900/80 backdrop-blur-sm border-slate-700/50">
        <CardContent className="p-4">
          <div className="flex items-center gap-2 mb-4">
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-orange-500"></div>
            <span className="text-sm text-slate-400">Loading topics...</span>
          </div>
        </CardContent>
      </Card>
    )
  }

  const selectedTopic = topics.find((t) => t.id === selectedTopicId)

  return (
    <Card className="bg-slate-900/80 backdrop-blur-sm border-slate-700/50 overflow-hidden">
      <CardContent className="p-0">
        {/* Header */}
        <div
          className="flex items-center justify-between p-4 cursor-pointer hover:bg-slate-800/50 transition-colors duration-200"
          onClick={() => setIsExpanded(!isExpanded)}
        >
          <div className="flex items-center gap-2">
            <Filter className="h-4 w-4 text-orange-400" />
            <h3 className="font-semibold text-white">Topics</h3>
          </div>
          <ChevronDown
            className={`h-4 w-4 text-slate-400 transition-transform duration-200 ${isExpanded ? "rotate-180" : ""}`}
          />
        </div>

        {/* Content */}
        <div
          className={`transition-all duration-300 overflow-hidden ${isExpanded ? "max-h-96 opacity-100" : "max-h-0 opacity-0"}`}
        >
          <div className="px-4 pb-4 space-y-2">
            {/* All Topics Button */}
            <Button
              variant={selectedTopicId === null ? "default" : "ghost"}
              size="sm"
              className={`w-full justify-start transition-all duration-200 ${
                selectedTopicId === null
                  ? "bg-gradient-to-r from-orange-600 to-orange-500 text-white hover:from-orange-500 hover:to-orange-400"
                  : "text-slate-300 hover:text-white hover:bg-slate-800/60"
              }`}
              onClick={() => onTopicSelect(null)}
            >
              <BookOpen className="mr-2 h-4 w-4" />
              All Topics
            </Button>

            {/* Individual Topics */}
            {topics.map((topic) => (
              <Button
                key={topic.id}
                variant={selectedTopicId === topic.id ? "default" : "ghost"}
                size="sm"
                className={`w-full justify-start transition-all duration-200 ${
                  selectedTopicId === topic.id
                    ? "bg-gradient-to-r from-orange-600 to-orange-500 text-white hover:from-orange-500 hover:to-orange-400"
                    : "text-slate-300 hover:text-white hover:bg-slate-800/60"
                }`}
                onClick={() => onTopicSelect(topic.id)}
              >
                <BookOpen className="mr-2 h-4 w-4" />
                <span className="truncate">{topic.title}</span>
              </Button>
            ))}
          </div>
        </div>

        {/* Selected Topic Info */}
        {selectedTopic && (
          <div className="border-t border-slate-700/50 p-4 bg-slate-800/30">
            <Badge variant="secondary" className="mb-2 bg-orange-500/10 text-orange-400 border-orange-500/20">
              Active Filter
            </Badge>
            <p className="font-medium text-white text-sm">{selectedTopic.title}</p>
            <p className="text-slate-400 text-xs mt-1 leading-relaxed">{selectedTopic.description}</p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
