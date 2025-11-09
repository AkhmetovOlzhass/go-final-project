import { API_BASE_URL } from "./auth"

export interface Topic {
  id: string
  title: string
  slug: string
  description: string
  parentId?: string
  createdAt: string
  updatedAt: string
  schoolClass: string
}

export interface Task {
  id: string
  title: string
  description: string
  content: string
  bodyMd: string
  difficulty: "EASY" | "MEDIUM" | "HARD" | "EXTREME"
  status: "DRAFT" | "PUBLISHED" | "ARCHIVED"
  officialSolution: string
  correctAnswer: string
  answerType: "TEXT" | "NUMBER" | "FORMULA"
  topicId: string
  authorId: string
  createdAt: string
  updatedAt: string
  imageUrl: string
}

export interface User {
  id: string
  email: string
  name: string
  role: "ADMIN" | "USER"
  createdAt: string
}

export interface Solution {
  id: string
  bodyMd: string
  taskId: string
  userId: string
  createdAt: string
  updatedAt: string
}

export async function refreshAccessToken() {
  const refreshToken = localStorage.getItem("refresh_token");
  if (!refreshToken) throw new Error("No refresh token");

  const response = await fetch(`${API_BASE_URL}/api/v1/auth/refresh`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refreshToken }),
  });

  if (!response.ok) throw new Error("Failed to refresh token");

  const data = await response.json();
  localStorage.setItem("access_token", data.accessToken);
  return data.accessToken;
}

export async function apiFetch(input: RequestInfo, init?: RequestInit): Promise<Response> {
  let token = localStorage.getItem("access_token");

  let headers = {
    ...(init?.headers || {}),
    "Content-Type": "application/json",
    ...(token && { Authorization: `Bearer ${token}` }),
  };

  let response = await fetch(input, { ...init, headers });

  if (response.status === 401) {
    try {
      token = await refreshAccessToken();
      headers = {
        ...(init?.headers || {}),
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      };
      response = await fetch(input, { ...init, headers });
    } catch (e) {
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
      throw new Error("Session expired, please login again");
    }
  }

  return response;
}


function getAuthHeaders() {
  const token = localStorage.getItem("access_token")
  return {
    "Content-Type": "application/json",
    ...(token && { Authorization: `Bearer ${token}` }),
  }
}

export const contentAPI = {
  // Topics
  async createTopic(data: { title: string; schoolClass: string }) {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/topics`, {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) throw new Error("Failed to create topic");
    return response.json();
  },

  async updateTopic(id: string, data: {title: string, schoolClass: string}){
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/topics/${id}`, {
      method: "PUT",
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    })
    if (!response.ok) throw new Error("Failed to update topic")
    return response.json()
  },

  async deleteTopic(id: string){
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/topics/${id}`, {
      method: "DELETE",
      headers: getAuthHeaders(),
    })
    if (!response.ok) throw new Error("Failed to delete topic")
    return true
  },

  async getAllTopics(): Promise<Topic[]> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/topics`)
    if (!response.ok) throw new Error("Failed to fetch topics")
      
    return response.json()
  },

  async getTopicById(id: string): Promise<Topic> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/topics/${id}`)
    if (!response.ok) throw new Error("Failed to fetch topic")
    return response.json()
  },

  // Tasks
  async getAllTasks(): Promise<Task[]> {    
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/tasks`)
    if (!response.ok) throw new Error("Failed to fetch topics")
      
    return response.json()
  },

  async createTask(data: {
    title: string;
    bodyMd: string;
    difficulty: string;
    topicId: string;
    officialSolution: string;
    correctAnswer: string;
    answerType: string;
    image?: File | null;
  }) {
    const formData = new FormData();
    formData.append("title", data.title);
    formData.append("bodyMd", data.bodyMd);
    formData.append("difficulty", data.difficulty);
    formData.append("topicId", data.topicId);
    formData.append("officialSolution", data.officialSolution);
    formData.append("correctAnswer", data.correctAnswer);
    formData.append("answerType", data.answerType);
  
    if (data.image) {
      formData.append("image", data.image); // 👈 файл
    }
  
    const token = localStorage.getItem("access_token");
  
    const response = await fetch(`${API_BASE_URL}/api/v1/content/tasks`, {
      method: "POST",
      headers: {
        ...(token && { Authorization: `Bearer ${token}` }),
      },
      body: formData,
    });
  
    if (!response.ok) throw new Error("Failed to create task");
    return response.json();
  },

  async updateTask(id: string, data: { title: string; bodyMd: string; difficulty: string; topicId: string }) {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/tasks/${id}`, {
      method: "PUT",
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    })
    if (!response.ok) throw new Error("Failed to create task")
    return response.json()
  },

  async deleteTask(id: string) {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/tasks/${id}`, {
      method: "DELETE",
      headers: getAuthHeaders()
    })
    if (!response.ok) throw new Error("Failed to create task")
    return response.json()
  },

  async publishTask(taskId: string) {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/tasks/publish`, {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({ taskId }),
    })
    if (!response.ok) throw new Error("Failed to publish task")
    return response.json()
  },

  async getDraftTasks(): Promise<Task[]> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/tasks/drafts`, {
      headers: getAuthHeaders(),
    })
    if (!response.ok) throw new Error("Failed to fetch draft tasks")
    return response.json()
  },

  async getTask(id: string): Promise<Task> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/tasks/${id}`)
    if (!response.ok) throw new Error("Failed to fetch task")
    return response.json()
  },

  async submitAnswer(taskId: string, answer: string) {
    const response = await fetch(`${API_BASE_URL}/api/v1/content/tasks/${taskId}/submit`, {
      method: "POST",
      headers: getAuthHeaders(),
      body: JSON.stringify({ answer }),
    })

    if (!response.ok) throw new Error("Failed to submit answer")
    return response.json()
  },

  async getTasksByTopic(topicId: string): Promise<Task[]> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/topics/${topicId}/tasks`)
    if (!response.ok) throw new Error("Failed to fetch tasks")
    return response.json()
  },

  async getUserProgress() {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/tasks/progress`, {
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error("Failed to fetch user progress");
    return response.json();
  }
}

export const usersAPI = {
  async getAllUsers(): Promise<User[]> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/users`, {
      headers: getAuthHeaders(),
    })
    if (!response.ok) throw new Error("Failed to fetch users")
    return response.json()
  },
}
