'use client';

import * as React from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { authAPI } from '@/lib/api';

export default function RegisterPage() {
  const router = useRouter();

  const [formData, setFormData] = React.useState({
    first_name: '',
    last_name: '',
    email: '',
    password: '',
    dob: '',
    nickname: '',
  });

  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState('');
  const [success, setSuccess] = React.useState('');

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  const handleSubmit = async (
    e: React.FormEvent<HTMLFormElement>
  ) => {
    e.preventDefault();

    setError('');
    setSuccess('');
    setIsLoading(true);

    try {
      await authAPI.register(formData);

      setSuccess('Registration successful! Redirecting...');

      setTimeout(() => {
        router.push('/feed');
      }, 1500);

    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : 'Registration failed.'
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-h-screen flex items-center justify-center bg-gray-30 px-0">
      <div className="w-full max-w-md bg-white rounded-xl shadow-sm border border-gray-200 p-36 sm:p-8">

        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-gray-900">
            Create Account
          </h1>

          <p className="text-gray-500 mt-2">
            Join the Social Network
          </p>
        </div>

        <form
          onSubmit={handleSubmit}
          className="space-y-5"
        >
          {/* Name row: horizontally aligned */}
          <div className="grid gap-4 sm:grid-cols-2">
            <Input
              label="First Name"
              name="first_name"
              placeholder="firstname"
              value={formData.first_name}
              onChange={handleChange}
              required
            />

            <Input
              label="Last Name"
              name="last_name"
              placeholder="lastname"
              value={formData.last_name}
              onChange={handleChange}
              required
            />
          </div>

          {/* Everything below: vertically stacked */}
          <Input
            label="Date of Birth"
            type="date"
            name="dob"
            value={formData.dob}
            onChange={handleChange}
            required
          />

          <Input
            label="Email"
            type="email"
            name="email"
            placeholder="joel@example.com"
            value={formData.email}
            onChange={handleChange}
            required
          />

          <Input
            label="Password"
            type="password"
            name="password"
            placeholder="Enter password"
            value={formData.password}
            onChange={handleChange}
            required
          />

          <Input
            label="Nickname"
            name="nickname"
            placeholder="nickname"
            value={formData.nickname}
            onChange={handleChange}
          />

          {error && (
            <div className="rounded-lg bg-red-50 border border-red-200 p-3 text-sm text-red-600">
              {error}
            </div>
          )}

          {success && (
            <div className="rounded-lg bg-green-50 border border-green-200 p-3 text-sm text-green-600">
              {success}
            </div>
          )}

          <Button
            type="submit"
            variant="primary"
            className="w-full"
            disabled={isLoading}
          >
            {isLoading
              ? 'Creating Account...'
              : 'Create Account'}
          </Button>
        </form>

        <p className="mt-6 text-center text-sm text-gray-500">
          Already have an account?{' '}
          <Link
            href="/login"
            className="text-indigo-600 hover:underline"
          >
            Sign in
          </Link>
        </p>

      </div>
    </div>
  );
}