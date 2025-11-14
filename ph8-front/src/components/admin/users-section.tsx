'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { usersAPI, type User } from '@/lib/api';

export function UsersSection() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadUsers();
  }, []);

  const loadUsers = async () => {
    try {
      const data = await usersAPI.getAllUsers();
      console.log(data);

      setUsers(data);
    } catch (error) {
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className='p-6 text-gray-300'>Loading users...</div>;
  }

  return (
    <div className='p-6 bg-gray-950 min-h-screen'>
      <div className='mb-6'>
        <h1 className='text-2xl font-bold text-white'>Users</h1>
        <p className='text-gray-400'>Manage platform users</p>
      </div>

      <div className='grid gap-4 md:grid-cols-2 lg:grid-cols-3'>
        {users.map((user) => (
          <Card
            key={user.id}
            className='bg-gray-900 border-gray-700 hover:bg-gray-800 transition-colors'
          >
            <CardHeader>
              <div className='flex justify-between items-start'>
                <div>
                  <CardTitle className='text-lg text-white'>{user.name}</CardTitle>
                  <CardDescription className='text-gray-400'>{user.email}</CardDescription>
                </div>
                <Badge
                  variant={user.role === 'Admin' ? 'default' : 'secondary'}
                  className={
                    user.role === 'Admin' ? 'bg-orange-600 text-white' : 'bg-gray-600 text-gray-300'
                  }
                >
                  {user.role}
                </Badge>
              </div>
            </CardHeader>
            <CardContent>
              <p className='text-sm text-gray-400'>
                Joined: {new Date(user.created_at).toLocaleDateString()}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      {users.length === 0 && (
        <div className='text-center py-12'>
          <p className='text-gray-400'>No users found.</p>
        </div>
      )}
    </div>
  );
}
