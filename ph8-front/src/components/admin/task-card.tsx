'use client';

import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Eye, Edit, Trash } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
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
import type { Task, Topic } from '@/lib/api';

interface UpdatePayload {
  title: string;
  bodyMd: string;
  difficulty: 'EASY' | 'MEDIUM' | 'HARD' | 'EXTREME';
  topicId: string;
  status: 'DRAFT' | 'PUBLISHED' | 'ARCHIVED';
  officialSolution: string;
  correctAnswer: string;
  answerType: 'TEXT' | 'NUMBER' | 'FORMULA';
  image?: File | null;
}

interface TaskCardProps {
  task: Task;
  topics: Topic[];
  onDelete?: (id: string) => void;
  onPublish?: (id: string) => void;
  /** image — опциональный файл, как и при создании */
  onUpdate?: (id: string, data: UpdatePayload, image?: File | null) => void;
}

export function TaskCard({ task, topics, onDelete, onPublish, onUpdate }: TaskCardProps) {
  const topicTitle = topics.find((t) => t.id === task.topic_id)?.title || 'Unknown';

  const [isEditOpen, setIsEditOpen] = useState(false);
  const [title, setTitle] = useState(task.title);
  const [bodyMd, setBodyMd] = useState(task.body_md || '');
  const [difficulty, setDifficulty] = useState<string>(task.difficulty || '');
  const [topicId, setTopicId] = useState<string>(task.topic_id || '');
  const [officialSolution, setOfficialSolution] = useState(task.official_solution || '');
  const [correctAnswer, setCorrectAnswer] = useState(task.correct_answer || '');
  const [answerType, setAnswerType] = useState<string>(task.answer_type || 'TEXT');

  // 🖼️ новое: файл картинки (опционально)
  const [image, setImage] = useState<File | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (onUpdate) {
      onUpdate(
        task.id,
        {
          title,
          bodyMd,
          difficulty: difficulty as 'EASY' | 'MEDIUM' | 'HARD' | 'EXTREME',
          topicId,
          status: task.status,
          officialSolution,
          correctAnswer,
          answerType: answerType as 'TEXT' | 'NUMBER' | 'FORMULA',
        },
        image, // 👈 передаём файл наверх (может быть null)
      );
    }
    setIsEditOpen(false);
    // по желанию можно очистить выбранный файл
    // setImage(null)
  };

  return (
    <>
      <Card className='bg-gray-900 border-gray-700 hover:bg-gray-800 transition-colors'>
        <CardHeader>
          <div className='flex justify-between items-start'>
            <div>
              <CardTitle className='text-lg text-white'>{task.title}</CardTitle>
              <CardDescription className='text-gray-400'>{task.description}</CardDescription>
            </div>
            <Badge
              variant={task.status === 'PUBLISHED' ? 'default' : 'secondary'}
              className={
                task.status === 'PUBLISHED'
                  ? 'bg-green-600 text-white'
                  : 'bg-gray-600 text-gray-300'
              }
            >
              {task.status}
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div className='space-y-2'>
            <p className='text-sm text-gray-400'>Topic: {topicTitle}</p>
            <p className='text-sm text-gray-400'>
              Created: {new Date(task.created_at).toLocaleDateString()}
            </p>
            <div className='flex gap-2'>
              <Button
                variant='outline'
                size='sm'
                className='border-gray-600 text-gray-300 hover:bg-gray-700 bg-transparent'
              >
                <Eye className='mr-2 h-4 w-4' />
                View
              </Button>

              <Button
                variant='outline'
                size='sm'
                className='border-gray-600 text-gray-300 hover:bg-gray-700 bg-transparent'
                onClick={() => setIsEditOpen(true)}
              >
                <Edit className='mr-2 h-4 w-4' />
                Edit
              </Button>

              {onDelete && (
                <Button
                  variant='outline'
                  size='sm'
                  className='border-gray-600 text-red-400 hover:bg-red-600 hover:text-white bg-transparent'
                  onClick={() => onDelete(task.id)}
                >
                  <Trash className='mr-2 h-4 w-4' />
                  Delete
                </Button>
              )}

              {onPublish && task.status === 'DRAFT' && (
                <Button
                  size='sm'
                  onClick={() => onPublish(task.id)}
                  className='bg-orange-600 hover:bg-orange-700 text-white'
                >
                  Publish
                </Button>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Edit Dialog */}
      <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
        <DialogContent className='max-w-2xl bg-gray-900 border-gray-700'>
          <DialogHeader>
            <DialogTitle className='text-white'>Edit Task</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit} className='space-y-4'>
            <div className='space-y-2'>
              <Label htmlFor='title' className='text-gray-300'>
                Title
              </Label>
              <Input
                id='title'
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className='bg-gray-800 border-gray-600 text-white'
                required
              />
            </div>

            <div className='space-y-2'>
              <Label htmlFor='bodyMd' className='text-gray-300'>
                Description (Markdown)
              </Label>
              <Textarea
                id='bodyMd'
                value={bodyMd}
                onChange={(e) => setBodyMd(e.target.value)}
                className='min-h-32 bg-gray-800 border-gray-600 text-white'
                required
              />
            </div>

            <div className='space-y-2'>
              <Label className='text-gray-300'>Topic</Label>
              <Select value={topicId} onValueChange={setTopicId} required>
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
              <Label className='text-gray-300'>Difficulty</Label>
              <Select value={difficulty} onValueChange={setDifficulty} required>
                <SelectTrigger className='bg-gray-800 border-gray-600 text-white'>
                  <SelectValue placeholder='Select difficulty' />
                </SelectTrigger>
                <SelectContent className='bg-gray-800 border-gray-600'>
                  <SelectItem value='EASY' className='text-green-400'>
                    Easy
                  </SelectItem>
                  <SelectItem value='MEDIUM' className='text-yellow-400'>
                    Medium
                  </SelectItem>
                  <SelectItem value='HARD' className='text-red-400'>
                    Hard
                  </SelectItem>
                  <SelectItem value='EXTREME' className='text-purple-400'>
                    Extreme
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Solution */}
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

            {/* Answer */}
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

            {/* Answer Type */}
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

            {/* 🖼️ Image (optional) */}
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
              {/* превью: новая картинка или текущая из задачи */}
              {(image || task.image_url) && (
                <div className='mt-2'>
                  <img
                    src={image ? URL.createObjectURL(image) : (task.image_url as string)}
                    alt='Preview'
                    className='max-h-48 rounded-md border border-gray-700'
                  />
                </div>
              )}
            </div>

            <Button type='submit' className='w-full bg-orange-600 hover:bg-orange-700 text-white'>
              Save Changes
            </Button>
          </form>
        </DialogContent>
      </Dialog>
    </>
  );
}
