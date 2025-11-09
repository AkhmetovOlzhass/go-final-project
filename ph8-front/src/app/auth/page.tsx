"use client"

import { useEffect, useState } from "react"
import { LoginForm } from "@/components/auth/login-form"
import { useRouter } from "next/navigation"
import { RegisterForm } from "@/components/auth/register-form"
import { useAuth } from "@/lib/auth"

export default function AuthPage() {
  const [isLogin, setIsLogin] = useState(true)
  const { user, loading, logout } = useAuth()

  const router = useRouter()

  useEffect(() => {
    if (user) {
      router.push("/learn")
    }
  }, [user, loading, router])

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-800 to-black p-4 relative overflow-hidden">
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-orange-900/20 via-transparent to-transparent"></div>
      <div className="absolute inset-0 bg-grid-white/[0.02] bg-[size:50px_50px]"></div>

      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-orange-500/10 rounded-full blur-3xl animate-pulse"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-blue-500/10 rounded-full blur-3xl animate-pulse delay-1000"></div>
      </div>

      <div className="w-full max-w-md relative z-10 animate-in fade-in-0 slide-in-from-bottom-4 duration-700">
        {isLogin ? (
          <LoginForm onToggleMode={() => setIsLogin(false)} />
        ) : (
          <RegisterForm onToggleMode={() => setIsLogin(true)} />
        )}
      </div>
    </div>
  )
}
