import { API_BASE_URL } from './auth';

export interface Topic {
  id: string;
  title: string;
  slug: string;
  parent_id?: string;
  created_at: string;
  updated_at: string;
  school_class: string;
}

export interface Task {
  id: string;
  title: string;
  description: string;
  content: string;
  body_md: string;
  difficulty: 'EASY' | 'MEDIUM' | 'HARD' | 'EXTREME';
  status: 'DRAFT' | 'PUBLISHED' | 'ARCHIVED';
  official_solution: string;
  correct_answer: string;
  answer_type: 'TEXT' | 'NUMBER' | 'FORMULA';
  topic_id: string;
  author_id: string;
  created_at: string;
  updated_at: string;
  image_url: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  role: 'Admin' | 'Teacher' | 'Student';
  created_at: string;
}

export interface Solution {
  id: string;
  body_md: string;
  task_id: string;
  user_id: string;
  created_at: string;
  updated_at: string;
}

export async function refreshAccessToken() {
  const refreshToken = localStorage.getItem('refresh_token');
  if (!refreshToken) throw new Error('No refresh token');

  const response = await fetch(`${API_BASE_URL}/api/v1/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });

  if (!response.ok) throw new Error('Failed to refresh token');

  const data = await response.json();
  localStorage.setItem('access_token', data.access_token);
  return data.access_token;
}

export async function apiFetch(input: RequestInfo, init?: RequestInit): Promise<Response> {
  let token = localStorage.getItem('access_token');

  let headers = {
    ...(init?.headers || {}),
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
  };

  let response = await fetch(input, { ...init, headers });

  if (response.status === 401) {
    try {
      token = await refreshAccessToken();
      headers = {
        ...(init?.headers || {}),
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      };
      response = await fetch(input, { ...init, headers });
    } catch (e) {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      throw new Error('Session expired, please login again');
    }
  }

  return response;
}

function getAuthHeaders() {
  const token = localStorage.getItem('access_token');
  return {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
  };
}

export const contentAPI = {
  // Topics
  async createTopic(data: {
    title: string;
    slug: string;
    school_class: string;
    parent_id: string | null;
  }) {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/topics`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) throw new Error('Failed to create topic');
    return response.json();
  },

  async updateTopic(
    id: string,
    data: { title: string; slug: string; school_class: string; parent_id: string | null },
  ) {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/topics/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error('Failed to update topic');
    return response.json();
  },

  async deleteTopic(id: string) {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/topics/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to delete topic');
    return true;
  },

  async getAllTopics(): Promise<Topic[]> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/topics`);
    if (!response.ok) throw new Error('Failed to fetch topics');

    return response.json();
  },

  async getTopicById(id: string): Promise<Topic> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/topics/${id}`);
    if (!response.ok) throw new Error('Failed to fetch topic');
    return response.json();
  },

  // Tasks
  async getAllTasks(): Promise<Task[]> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/tasks`);
    if (!response.ok) throw new Error('Failed to fetch topics');

    return response.json();
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
    formData.append('title', data.title);
    formData.append('body_md', data.bodyMd);
    formData.append('difficulty', data.difficulty);
    formData.append('status', 'DRAFT');
    formData.append('topic_id', data.topicId);
    formData.append('official_solution', data.officialSolution);
    formData.append('correct_answer', data.correctAnswer);
    formData.append('answer_type', data.answerType);

    if (data.image) {
      formData.append('image_url', data.image);
    }

    const token = localStorage.getItem('access_token');

    const response = await fetch(`${API_BASE_URL}/api/v1/tasks`, {
      method: 'POST',
      headers: {
        ...(token && { Authorization: `Bearer ${token}` }),
      },
      body: formData,
    });

    if (!response.ok) throw new Error('Failed to create task');
    return response.json();
  },

  async updateTask(
    id: string,
    data: {
      title: string;
      bodyMd: string;
      difficulty: string;
      topicId: string;
      status: string;
      officialSolution?: string;
      correctAnswer?: string;
      answerType: string;
      image?: File | null;
    },
  ) {
    const formData = new FormData();

    formData.append('title', data.title);
    formData.append('body_md', data.bodyMd);
    formData.append('difficulty', data.difficulty);
    formData.append('topic_id', data.topicId);
    formData.append('status', data.status);
    formData.append('answer_type', data.answerType);

    if (data.officialSolution) formData.append('official_solution', data.officialSolution);

    if (data.correctAnswer) formData.append('correct_answer', data.correctAnswer);

    if (data.image) formData.append('image_url', data.image);

    const token = localStorage.getItem('access_token');

    const response = await fetch(`${API_BASE_URL}/api/v1/tasks/${id}`, {
      method: 'PUT',
      headers: {
        ...(token && { Authorization: `Bearer ${token}` }),
      },
      body: formData,
    });

    if (!response.ok) throw new Error('Failed to update task');

    return response.json();
  },

  async deleteTask(id: string) {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/tasks/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to create task');
    return response.json();
  },

  async publishTask(taskId: string) {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/tasks/${taskId}/publish`, {
      method: 'POST',
      headers: getAuthHeaders(),
    });

    if (!response.ok) throw new Error('Failed to publish task');
    return response.json();
  },

  async getDraftTasks(): Promise<Task[]> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/tasks/drafts`, {
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to fetch draft tasks');
    return response.json();
  },

  async getTask(id: string): Promise<Task> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/tasks/${id}`);
    if (!response.ok) throw new Error('Failed to fetch task');
    return response.json();
  },

  async submitAnswer(taskId: string, answer: string) {
    const response = await fetch(`${API_BASE_URL}/api/v1/content/tasks/${taskId}/submit`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({ answer }),
    });

    if (!response.ok) throw new Error('Failed to submit answer');
    return response.json();
  },

  async getTasksByTopic(topicId: string): Promise<Task[]> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/tasks/topic/${topicId}`);
    if (!response.ok) throw new Error('Failed to fetch tasks');
    return response.json();
  },

  async getUserProgress() {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/content/tasks/progress`, {
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to fetch user progress');
    return response.json();
  },
};

export const usersAPI = {
  async getAllUsers(): Promise<User[]> {
    const response = await apiFetch(`${API_BASE_URL}/api/v1/user/all`, {
      headers: getAuthHeaders(),
    });
    if (!response.ok) throw new Error('Failed to fetch users');
    return response.json();
  },
};
