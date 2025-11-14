'use client';

import type React from 'react';
import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { Plus } from 'lucide-react';
import { contentAPI, type Task, type Topic } from '@/lib/api';
import { TaskCard } from '@/components/admin/task-card';

export function TasksSection() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState(true);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [title, setTitle] = useState('');
  const [bodyMd, setBodyMd] = useState('');
  const [difficulty, setDifficulty] = useState('');
  const [selectedTopicId, setSelectedTopicId] = useState('');
  const [officialSolution, setOfficialSolution] = useState('');
  const [correctAnswer, setCorrectAnswer] = useState('');
  const [answerType, setAnswerType] = useState('');
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState('');
  const [image, setImage] = useState<File | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [tasksData, topicsData] = await Promise.all([
        contentAPI.getDraftTasks(),
        contentAPI.getAllTopics(),
      ]);

      console.log(topicsData);

      setTasks(tasksData);
      setTopics(topicsData);
    } catch (error) {
      setError('Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateTask = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreating(true);
    setError('');
    try {
      await contentAPI.createTask({
        title,
        bodyMd,
        difficulty: difficulty as 'EASY' | 'MEDIUM' | 'HARD' | 'EXTREME',
        topicId: selectedTopicId,
        officialSolution,
        correctAnswer,
        answerType: answerType as 'TEXT' | 'NUMBER' | 'FORMULA',
        image,
      });
      setTitle('');
      setBodyMd('');
      setDifficulty('');
      setSelectedTopicId('');
      setOfficialSolution('');
      setCorrectAnswer('');
      setAnswerType('');
      setIsCreateOpen(false);
      setImage(null);
      loadData();
    } catch (error) {
      setError('Failed to create task');
    } finally {
      setCreating(false);
    }
  };

  const handlePublishTask = async (taskId: string) => {
    try {
      await contentAPI.publishTask(taskId);
      loadData();
    } catch (error) {
      setError('Failed to publish task');
    }
  };

  const handleDeleteTask = async (id: string) => {
    if (!confirm('Are you sure you want to delete this task?')) return;
    try {
      await contentAPI.deleteTask(id);
      loadData();
    } catch (error) {
      setError('Failed to delete task');
    }
  };

  if (loading) return <div className='p-6 text-gray-300'>Loading tasks...</div>;

  return (
    <div className='p-6 bg-gray-950 min-h-screen'>
      <div className='flex justify-between items-center mb-6'>
        <div>
          <h1 className='text-2xl font-bold text-white'>Tasks</h1>
          <p className='text-gray-400'>Manage physics learning tasks</p>
        </div>
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button className='bg-orange-600 hover:bg-orange-700 text-white'>
              <Plus className='mr-2 h-4 w-4' />
              Create Task
            </Button>
          </DialogTrigger>
          <DialogContent className='max-w-2xl bg-gray-900 border-gray-700'>
            <DialogHeader>
              <DialogTitle className='text-white'>Create New Task</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleCreateTask} className='space-y-4'>
              {error && (
                <div className='bg-red-900/50 border border-red-700 text-red-300 px-4 py-2 rounded-md text-sm'>
                  {error}
                </div>
              )}
              <div className='space-y-2'>
                <Label htmlFor='task-title' className='text-gray-300'>
                  Title
                </Label>
                <Input
                  id='task-title'
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder='Enter task title'
                  required
                  className='bg-gray-800 border-gray-600 text-white placeholder:text-gray-400'
                />
              </div>
              <div className='space-y-2'>
                <Label htmlFor='task-description' className='text-gray-300'>
                  Description (Markdown)
                </Label>
                <Textarea
                  id='task-description'
                  value={bodyMd}
                  onChange={(e) => setBodyMd(e.target.value)}
                  placeholder='Enter task description in Markdown format'
                  required
                  className='min-h-32 bg-gray-800 border-gray-600 text-white placeholder:text-gray-400'
                />
              </div>
              <div className='space-y-2'>
                <Label htmlFor='task-image' className='text-gray-300'>
                  Image (optional)
                </Label>
                <Input
                  id='task-image'
                  type='file'
                  accept='image/*'
                  onChange={(e) => setImage(e.target.files?.[0] || null)}
                  className='bg-gray-800 border-gray-600 text-white'
                />
              </div>
              <div className='space-y-2'>
                <Label htmlFor='topic' className='text-gray-300'>
                  Topic
                </Label>
                <Select value={selectedTopicId} onValueChange={setSelectedTopicId} required>
                  <SelectTrigger className='bg-gray-800 border-gray-600 text-white'>
                    <SelectValue placeholder='Select a topic' />
                  </SelectTrigger>
                  <SelectContent className='bg-gray-800 border-gray-600'>
                    {topics.map((topic) => (
                      <SelectItem
                        key={topic.id}
                        value={topic.id}
                        className='text-white hover:bg-gray-700'
                      >
                        {topic.title}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className='space-y-2'>
                <Label htmlFor='difficulty' className='text-gray-300'>
                  Difficulty
                </Label>
                <Select value={difficulty} onValueChange={setDifficulty} required>
                  <SelectTrigger className='bg-gray-800 border-gray-600 text-white'>
                    <SelectValue placeholder='Select difficulty level' />
                  </SelectTrigger>
                  <SelectContent className='bg-gray-800 border-gray-600'>
                    <SelectItem value='EASY' className='text-green-400 hover:bg-gray-700'>
                      Easy
                    </SelectItem>
                    <SelectItem value='MEDIUM' className='text-yellow-400 hover:bg-gray-700'>
                      Medium
                    </SelectItem>
                    <SelectItem value='HARD' className='text-red-400 hover:bg-gray-700'>
                      Hard
                    </SelectItem>
                    <SelectItem value='EXTREME' className='text-purple-400 hover:bg-gray-700'>
                      Extreme
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className='space-y-2'>
                <Label htmlFor='task-solution' className='text-gray-300'>
                  Solution (Markdown)
                </Label>
                <Textarea
                  id='task-solution'
                  value={officialSolution}
                  onChange={(e) => setOfficialSolution(e.target.value)}
                  placeholder='Enter task solution in Markdown format'
                  required
                  className='min-h-32 bg-gray-800 border-gray-600 text-white placeholder:text-gray-400'
                />
              </div>

              <div className='space-y-2'>
                <Label htmlFor='task-answer' className='text-gray-300'>
                  Answer
                </Label>
                <Input
                  id='task-answer'
                  value={correctAnswer}
                  onChange={(e) => setCorrectAnswer(e.target.value)}
                  placeholder='Enter answer'
                  required
                  className='bg-gray-800 border-gray-600 text-white placeholder:text-gray-400'
                />
              </div>
              <div className='space-y-2'>
                <Label htmlFor='answerType' className='text-gray-300'>
                  Answer type
                </Label>
                <Select value={answerType} onValueChange={setAnswerType} required>
                  <SelectTrigger className='bg-gray-800 border-gray-600 text-white'>
                    <SelectValue placeholder='Select answer type' />
                  </SelectTrigger>
                  <SelectContent className='bg-gray-800 border-gray-600'>
                    <SelectItem value='TEXT' className='text-green-400 hover:bg-gray-700'>
                      Text
                    </SelectItem>
                    <SelectItem value='NUMBER' className='text-green-400 hover:bg-gray-700'>
                      Number
                    </SelectItem>
                    <SelectItem value='FORMULA' className='text-green-400 hover:bg-gray-700'>
                      Formula
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <Button
                type='submit'
                disabled={creating}
                className='w-full bg-orange-600 hover:bg-orange-700 text-white'
              >
                {creating ? 'Creating...' : 'Create Task'}
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <div className='grid gap-4 md:grid-cols-2'>
        {tasks.map((task) => (
          <TaskCard
            key={task.id}
            task={task}
            topics={topics}
            onDelete={handleDeleteTask}
            onPublish={handlePublishTask}
            onUpdate={async (id, data) => {
              await contentAPI.updateTask(id, data);
              loadData();
            }}
          />
        ))}
      </div>

      {tasks.length === 0 && (
        <div className='text-center py-12'>
          <p className='text-gray-400'>
            No tasks created yet. Create your first task to get started.
          </p>
        </div>
      )}
    </div>
  );
}
