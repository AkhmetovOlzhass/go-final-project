"use client"

import { useRouter } from "next/navigation"
import { useEffect } from "react"

export default function HomePage() {
  const router = useRouter()

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center">
      <div className="text-center space-y-4">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto"></div>
        <div className="text-lg text-gray-600">Redirecting...</div>
      </div>
    </div>
  )
}
