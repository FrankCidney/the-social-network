'use client';

import * as React from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';

export default function RegisterPage() {
  const [isLoading, setIsLoading] = React.useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    // Backend integration point: collect form values and send them to the register endpoint.
    // Example: POST /api/auth/register
    // TODO: replace this placeholder with the actual API call and success/error handling.

    setIsLoading(false);
  };

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h2 className="text-2xl font-bold tracking-tight text-text-main">Create an account</h2>
        <p className="text-sm text-gray-500">
          Join the community and start sharing
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          {/* Backend integration point: wire these fields to the register payload. */}
          <Input label="First Name" placeholder="John" required />
          <Input label="Last Name" placeholder="Doe" required />
        </div>
        <Input label="Email" placeholder="name@example.com" type="email" required />
        <Input label="Date of Birth" type="date" required />
        <Input label="Password" placeholder="••••••••" type="password" required />
        
        <div className="space-y-4 pt-2">
          <div className="text-xs text-gray-400 uppercase font-semibold tracking-wider">
            Optional Details
          </div>
          <Input label="Nickname" placeholder="johndoe" />
          <div className="space-y-1.5">
            {/* Backend integration point: send optional profile details here if supported by the API. */}
            <label className="text-sm font-medium text-gray-700 ml-1">About Me</label>
            <textarea 
              className="flex min-h-[80px] w-full rounded-bento border border-gray-200 bg-white px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              placeholder="Tell us about yourself..."
            />
          </div>
        </div>

        <Button 
          type="submit" 
          className="w-full" 
          disabled={isLoading}
        >
          {isLoading ? 'Creating account...' : 'Create Account'}
        </Button>
      </form>

      <div className="text-center text-sm">
        <span className="text-gray-500">Already have an account? </span>
        <Link href="/login" className="font-medium text-primary hover:underline">
          Sign In
        </Link>
      </div>
    </div>
  );
}
