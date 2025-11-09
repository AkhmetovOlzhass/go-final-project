"use client"

import type React from "react"

import { useState, useEffect, useRef } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { ArrowLeft, Camera, Loader2, User } from "lucide-react"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8081"

export default function EditProfilePage() {
  const router = useRouter()
  const fileInputRef = useRef<HTMLInputElement>(null)

  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [profile, setProfile] = useState({
    email: "",
    displayName: "",
    avatar: "",
  })
  const [avatarFile, setAvatarFile] = useState<File | null>(null)
  const [avatarPreview, setAvatarPreview] = useState<string>("")

  useEffect(() => {
    loadProfile()
    
  }, [])


  const loadProfile = async () => {
    try {
      const token = localStorage.getItem("access_token")
      if (!token) {
        router.push("/")
        return
      }

      const response = await fetch(`${API_BASE_URL}/api/v1/user/profile`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error("Failed to load profile")
      }

      const data = await response.json()
      
      setProfile({
        email: data.email || "",
        displayName: data.displayName || "",
        avatar: data.avatarUrl || "",
      })
      setAvatarPreview(data.avatarUrl || "")
    } catch (error) {
      console.error("Failed to load profile:", error)
    } finally {
      setLoading(false)
    }
  }

  const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setAvatarFile(file)
      const reader = new FileReader()
      reader.onloadend = () => {
        setAvatarPreview(reader.result as string)
      }
      reader.readAsDataURL(file)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)

    try {
      const token = localStorage.getItem("access_token")
      if (!token) {
        router.push("/")
        return
      }

      const formData = new FormData()
      formData.append("email", profile.email)
      formData.append("displayName", profile.displayName)

      if (avatarFile) {
        formData.append("avatar", avatarFile)
      }

      const response = await fetch(`${API_BASE_URL}/api/v1/user/profile`, {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      })

      if (!response.ok) {
        throw new Error("Failed to update profile")
      }

    router.push("/learn")
    } catch (error) {
      console.error("Failed to update profile:", error)
    } finally {
      setSaving(false)
    }
  }

  const getInitials = (name: string) => {
    return name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2)
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="text-center space-y-4">
          <div className="relative">
            <div className="animate-spin rounded-full h-12 w-12 border-2 border-gray-700 mx-auto"></div>
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-orange-500 mx-auto absolute top-0 left-1/2 transform -translate-x-1/2"></div>
          </div>
          <div className="text-lg text-gray-300 animate-pulse">Loading profile...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-900">
      <div className="bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900 border-b border-gray-700">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 py-6">
          <button
            onClick={() => router.back()}
            className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors duration-200 mb-4"
          >
            <ArrowLeft className="h-5 w-5" />
            <span>Back</span>
          </button>
          <div className="space-y-2">
            <h1 className="text-3xl sm:text-4xl font-bold text-white">
              Edit{" "}
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-orange-400 to-orange-600">
                Profile
              </span>
            </h1>
            <p className="text-lg text-gray-300">Update your personal information and avatar</p>
          </div>
        </div>
      </div>

      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8">
        <Card className="bg-gray-800 border-gray-700">
          <CardHeader>
            <CardTitle className="text-white">Profile Information</CardTitle>
            <CardDescription className="text-gray-400">Update your display name and profile picture</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              <div className="flex flex-col items-center gap-4 pb-6 border-b border-gray-700">
                <div className="relative group">
                  <Avatar className="h-32 w-32 border-4 border-gray-700">
                    <AvatarImage src={avatarPreview || "/placeholder.svg"} alt={profile.displayName} />
                    <AvatarFallback className="bg-gradient-to-br from-orange-400 to-orange-600 text-white text-3xl font-semibold">
                      {profile.displayName ? getInitials(profile.displayName) : <User className="h-12 w-12" />}
                    </AvatarFallback>
                  </Avatar>
                  <button
                    type="button"
                    onClick={() => fileInputRef.current?.click()}
                    className="absolute bottom-0 right-0 bg-orange-500 hover:bg-orange-600 text-white p-3 rounded-full transition-colors duration-200 shadow-lg"
                  >
                    <Camera className="h-5 w-5" />
                  </button>
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/*"
                    onChange={handleAvatarChange}
                    className="hidden"
                  />
                </div>
                <p className="text-sm text-gray-400">Click the camera icon to upload a new avatar</p>
              </div>

              <div className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="email" className="text-gray-300">
                    Email
                  </Label>
                  <Input
                    id="email"
                    type="email"
                    value={profile.email}
                    onChange={(e) => setProfile({ ...profile, email: e.target.value })}
                    className="bg-gray-900 border-gray-700 text-white focus:border-orange-500 focus:ring-orange-500"
                    disabled
                  />
                  <p className="text-xs text-gray-500">Email cannot be changed</p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="displayName" className="text-gray-300">
                    Display Name
                  </Label>
                  <Input
                    id="displayName"
                    type="text"
                    value={profile.displayName}
                    onChange={(e) => setProfile({ ...profile, displayName: e.target.value })}
                    placeholder="Enter your display name"
                    className="bg-gray-900 border-gray-700 text-white placeholder:text-gray-500 focus:border-orange-500 focus:ring-orange-500"
                    required
                  />
                </div>
              </div>

              <div className="flex gap-4 pt-4">
                <Button
                  type="submit"
                  disabled={saving}
                  className="flex-1 bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold"
                >
                  {saving ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Saving...
                    </>
                  ) : (
                    "Save Changes"
                  )}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => router.back()}
                  className="flex-1 bg-transparent border-gray-700 text-gray-300 hover:bg-gray-800 hover:text-white"
                >
                  Cancel
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
