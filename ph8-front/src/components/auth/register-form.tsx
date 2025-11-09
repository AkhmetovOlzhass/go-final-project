"use client"

import type React from "react"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useAuth } from "@/lib/auth"
import { AlertCircle } from "lucide-react"

interface RegisterFormProps {
  onToggleMode: () => void
}

export function RegisterForm({ onToggleMode }: RegisterFormProps) {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [name, setName] = useState("")
  const [error, setError] = useState("")
  const { register, loading } = useAuth()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    try {
      await register(email, password, name)
    } catch (error) {
      setError(error instanceof Error ? error.message : "Registration failed. Please try again.")
    }
  }

  return (
    <Card className="w-full max-w-md bg-gray-900/80 backdrop-blur-xl border-gray-700/50 shadow-2xl hover:shadow-orange-500/10 transition-all duration-300">
      <CardHeader className="text-center space-y-2">
        <CardTitle className="text-2xl font-bold text-white">Create Account</CardTitle>
        <CardDescription className="text-gray-300">Join our Physics Learning platform</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-3 flex items-center gap-2 animate-in slide-in-from-top-2 duration-300">
              <AlertCircle className="h-4 w-4 text-red-400 flex-shrink-0" />
              <p className="text-red-300 text-sm">{error}</p>
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="name" className="text-gray-200">
              Full Name
            </Label>
            <Input
              id="name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              placeholder="Enter your full name"
              className="bg-gray-800/50 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500 focus:ring-orange-500/20 transition-all duration-200"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="email" className="text-gray-200">
              Email
            </Label>
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              placeholder="Enter your email"
              className="bg-gray-800/50 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500 focus:ring-orange-500/20 transition-all duration-200"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password" className="text-gray-200">
              Password
            </Label>
            <Input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              placeholder="Create a password"
              className="bg-gray-800/50 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500 focus:ring-orange-500/20 transition-all duration-200"
            />
          </div>
          <Button
            type="submit"
            className="w-full bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-medium py-2.5 transition-all duration-200 transform hover:scale-[1.02] disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
            disabled={loading}
          >
            {loading ? "Creating account..." : "Create Account"}
          </Button>
        </form>
        <div className="mt-4 text-center">
          <Button
            variant="link"
            onClick={onToggleMode}
            className="text-gray-300 hover:text-orange-400 transition-colors duration-200"
          >
            Already have an account? Sign in
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}
