'use client';

import * as React from 'react';
import { CheckCircle2, Pencil, Sparkles, User, Users } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { profileAPI, ProfileResponse, PublicUser, resolveAssetUrl } from '@/lib/api';

function getDisplayName(user: PublicUser) {
  return user.nickname?.trim() ? user.nickname : `${user.first_name} ${user.last_name}`;
}

function getInitials(user: PublicUser) {
  return `${user.first_name?.[0] ?? ''}${user.last_name?.[0] ?? ''}`.toUpperCase();
}

function formatDate(value?: string) {
  if (!value) return 'Not provided';
  const date = new Date(value);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

export default function ProfilePage() {
  const [profile, setProfile] = React.useState<ProfileResponse | null>(null);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    const loadProfile = async () => {
      setLoading(true);
      setError(null);

      try {
        const data = await profileAPI.getMyProfile();
        setProfile(data);
      } catch (err) {
        console.error('Failed to load profile', err);
        setError('Unable to load your profile right now.');
      } finally {
        setLoading(false);
      }
    };

    void loadProfile();
  }, []);

  if (loading) {
    return (
      <div className="space-y-4">
        <div className="animate-pulse rounded-xl border border-gray-100 bg-white p-6">
          <div className="h-8 w-64 rounded bg-gray-100 mb-4" />
          <div className="grid gap-4 sm:grid-cols-3">
            <div className="h-24 rounded-xl bg-gray-100" />
            <div className="h-24 rounded-xl bg-gray-100" />
            <div className="h-24 rounded-xl bg-gray-100" />
          </div>
        </div>
        <div className="animate-pulse rounded-xl border border-gray-100 bg-white p-6 h-72" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-xl p-6 text-red-700">
        <h1 className="text-xl font-semibold">Profile unavailable</h1>
        <p className="mt-2 text-sm">{error}</p>
      </div>
    );
  }

  if (!profile) {
    return null;
  }

  const avatarUrl = resolveAssetUrl(profile.user.avatar_path);
  const name = getDisplayName(profile.user);
  const handle = `@${profile.user.first_name.toLowerCase()}${profile.user.last_name.toLowerCase()}`;

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-xl border border-gray-100 p-6">
        <div className="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
          <div className="flex items-center gap-5">
            <div className="relative flex h-20 w-20 items-center justify-center rounded-3xl bg-indigo-100 text-3xl font-bold text-indigo-700 overflow-hidden">
              {avatarUrl ? (
                <img src={avatarUrl} alt={name} className="h-full w-full object-cover" />
              ) : (
                getInitials(profile.user)
              )}
            </div>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">{name}</h1>
              <p className="text-sm text-gray-500 mt-1">{handle}</p>
              <div className="mt-3 inline-flex items-center gap-2 rounded-full border border-gray-200 bg-gray-50 px-3 py-1 text-sm text-gray-600">
                <CheckCircle2 className="w-4 h-4 text-emerald-500" />
                {profile.user.is_public ? 'Public profile' : 'Private profile'}
              </div>
            </div>
          </div>

          <Button variant="secondary" size="lg" className="flex items-center gap-2">
            <Pencil className="w-4 h-4" />
            Edit profile
          </Button>
        </div>

        <div className="mt-8 grid gap-4 sm:grid-cols-3">
          <div className="rounded-3xl border border-gray-100 bg-gray-50 p-5">
            <p className="text-sm text-gray-500">Followers</p>
            <p className="mt-2 text-3xl font-semibold text-gray-900">{profile.follower_count.toLocaleString()}</p>
          </div>
          <div className="rounded-3xl border border-gray-100 bg-gray-50 p-5">
            <p className="text-sm text-gray-500">Following</p>
            <p className="mt-2 text-3xl font-semibold text-gray-900">{profile.following_count.toLocaleString()}</p>
          </div>
          <div className="rounded-3xl border border-gray-100 bg-gray-50 p-5">
            <p className="text-sm text-gray-500">Posts</p>
            <p className="mt-2 text-3xl font-semibold text-gray-900">{profile.post_count.toLocaleString()}</p>
          </div>
        </div>
      </div>

      <div className="grid gap-4 lg:grid-cols-3">
        <div className="lg:col-span-2 space-y-4">
          <div className="rounded-xl border border-gray-100 bg-white p-6">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-xl font-semibold text-gray-900">About you</h2>
                <p className="text-sm text-gray-500 mt-1">A quick summary of your profile info.</p>
              </div>
              <Sparkles className="w-5 h-5 text-indigo-500" />
            </div>
            <p className="mt-6 text-gray-600 leading-7">
              {profile.about_me?.trim() || 'You have not added an about section yet. Add a short bio to let people know what you care about.'}
            </p>
          </div>

          <div className="rounded-xl border border-gray-100 bg-white p-6">
            <h2 className="text-xl font-semibold text-gray-900">Profile details</h2>
            <div className="mt-6 space-y-4">
              <div className="flex items-center justify-between rounded-3xl border border-gray-100 bg-gray-50 p-4">
                <span className="text-sm text-gray-500">Full name</span>
                <span className="font-medium text-gray-900">{profile.user.first_name} {profile.user.last_name}</span>
              </div>
              <div className="flex items-center justify-between rounded-3xl border border-gray-100 bg-gray-50 p-4">
                <span className="text-sm text-gray-500">Nickname</span>
                <span className="font-medium text-gray-900">{profile.user.nickname || 'Not set'}</span>
              </div>
              <div className="flex items-center justify-between rounded-3xl border border-gray-100 bg-gray-50 p-4">
                <span className="text-sm text-gray-500">Birthday</span>
                <span className="font-medium text-gray-900">{formatDate(profile.dob)}</span>
              </div>
            </div>
          </div>
        </div>

        <div className="space-y-4">
          <div className="rounded-xl border border-gray-100 bg-white p-6">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-xl font-semibold text-gray-900">Profile summary</h2>
                <p className="text-sm text-gray-500 mt-1">Your profile activity at a glance.</p>
              </div>
              <Users className="w-5 h-5 text-indigo-500" />
            </div>

            <div className="mt-6 grid gap-4">
              <div className="rounded-3xl border border-gray-100 bg-gray-50 p-4">
                <p className="text-sm text-gray-500">Visibility</p>
                <p className="mt-2 font-semibold text-gray-900">{profile.user.is_public ? 'Open to everyone' : 'Only visible to your followers'}</p>
              </div>
              <div className="rounded-3xl border border-gray-100 bg-gray-50 p-4">
                <p className="text-sm text-gray-500">Member since</p>
                <p className="mt-2 font-semibold text-gray-900">{profile.dob ? formatDate(profile.dob) : 'Date not available'}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
