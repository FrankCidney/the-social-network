'use client';

import * as React from 'react';
import { CheckCircle2, Pencil, Sparkles, Users, X } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import {
  profileAPI,
  ProfileResponse,
  PublicUser,
  resolveAssetUrl,
  UpdateProfilePayload,
} from '@/lib/api';

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

// Converts a dob string to yyyy-mm-dd for the <input type="date"> value.
function toDateInputValue(value?: string) {
  if (!value) return '';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '';
  return date.toISOString().slice(0, 10);
}

type EditFormState = {
  first_name: string;
  last_name: string;
  nickname: string;
  about_me: string;
  dob: string;
  is_public: boolean;
};

export default function ProfilePage() {
  const [profile, setProfile] = React.useState<ProfileResponse | null>(null);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const [isEditing, setIsEditing] = React.useState(false);
  const [form, setForm] = React.useState<EditFormState | null>(null);
  const [saving, setSaving] = React.useState(false);
  const [saveError, setSaveError] = React.useState<string | null>(null);

  const [avatarFile, setAvatarFile] = React.useState<File | null>(null);
  const [avatarPreview, setAvatarPreview] = React.useState<string | null>(null);
  const [uploadingAvatar, setUploadingAvatar] = React.useState(false);
  const avatarInputRef = React.useRef<HTMLInputElement>(null);

  const loadProfile = React.useCallback(async () => {
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
  }, []);

  React.useEffect(() => {
    void loadProfile();
  }, [loadProfile]);

  const startEditing = () => {
    if (!profile) return;
    setForm({
      first_name: profile.user.first_name,
      last_name: profile.user.last_name,
      nickname: profile.user.nickname ?? '',
      about_me: profile.about_me ?? '',
      dob: toDateInputValue(profile.dob),
      is_public: profile.user.is_public,
    });
    setSaveError(null);
    setIsEditing(true);
  };

  const cancelEditing = () => {
    setIsEditing(false);
    setForm(null);
    setSaveError(null);
    setAvatarFile(null);
    setAvatarPreview(null);
    if (avatarInputRef.current) avatarInputRef.current.value = '';
  };

  const handleAvatarSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setAvatarFile(file);
    setAvatarPreview(URL.createObjectURL(file));
  };

  const handleSave = async () => {
    if (!form || !profile) return;
    setSaving(true);
    setSaveError(null);

    try {
      // Upload the new avatar first, if one was chosen, so it's reflected
      // in the same profile refresh as the rest of the field updates.
      if (avatarFile) {
        setUploadingAvatar(true);
        await profileAPI.uploadAvatar(avatarFile);
        setUploadingAvatar(false);
      }

      const payload: UpdateProfilePayload = {
        first_name: form.first_name.trim(),
        last_name: form.last_name.trim(),
        nickname: form.nickname.trim() || undefined,
        about_me: form.about_me.trim(),
        dob: form.dob || undefined,
        is_public: form.is_public,
      };

      const updated = await profileAPI.updateProfile(payload);
      setProfile(updated);
      setIsEditing(false);
      setForm(null);
      setAvatarFile(null);
      setAvatarPreview(null);
      if (avatarInputRef.current) avatarInputRef.current.value = '';
    } catch (err) {
      console.error('Failed to update profile', err);
      setSaveError('Could not save your changes. Please try again.');
    } finally {
      setSaving(false);
      setUploadingAvatar(false);
    }
  };

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
        <button onClick={loadProfile} className="mt-3 text-sm font-medium underline">
          Try again
        </button>
      </div>
    );
  }

  if (!profile) {
    return null;
  }

  const avatarUrl = avatarPreview || resolveAssetUrl(profile.user.avatar_path);
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
              {isEditing && (
                <button
                  type="button"
                  onClick={() => avatarInputRef.current?.click()}
                  className="absolute inset-0 flex items-center justify-center bg-black/40 text-white text-xs font-medium opacity-0 hover:opacity-100 transition-opacity"
                >
                  Change
                </button>
              )}
              <input
                ref={avatarInputRef}
                type="file"
                accept="image/*"
                className="hidden"
                onChange={handleAvatarSelect}
              />
            </div>
            <div>
              {isEditing && form ? (
                <div className="space-y-2">
                  <div className="flex gap-2">
                    <input
                      value={form.first_name}
                      onChange={(e) => setForm({ ...form, first_name: e.target.value })}
                      placeholder="First name"
                      className="w-32 rounded-lg border border-gray-200 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                    <input
                      value={form.last_name}
                      onChange={(e) => setForm({ ...form, last_name: e.target.value })}
                      placeholder="Last name"
                      className="w-32 rounded-lg border border-gray-200 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                  </div>
                  <input
                    value={form.nickname}
                    onChange={(e) => setForm({ ...form, nickname: e.target.value })}
                    placeholder="Nickname (optional)"
                    className="w-full rounded-lg border border-gray-200 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  />
                </div>
              ) : (
                <>
                  <h1 className="text-3xl font-bold text-gray-900">{name}</h1>
                  <p className="text-sm text-gray-500 mt-1">{handle}</p>
                </>
              )}

              <div className="mt-3 inline-flex items-center gap-2 rounded-full border border-gray-200 bg-gray-50 px-3 py-1 text-sm text-gray-600">
                <CheckCircle2 className="w-4 h-4 text-emerald-500" />
                {isEditing && form ? (
                  <label className="flex items-center gap-2 cursor-pointer">
                    <input
                      type="checkbox"
                      checked={form.is_public}
                      onChange={(e) => setForm({ ...form, is_public: e.target.checked })}
                    />
                    {form.is_public ? 'Public profile' : 'Private profile'}
                  </label>
                ) : (
                  <>{profile.user.is_public ? 'Public profile' : 'Private profile'}</>
                )}
              </div>
            </div>
          </div>

          {isEditing ? (
            <div className="flex items-center gap-2">
              <Button
                variant="ghost"
                size="lg"
                onClick={cancelEditing}
                disabled={saving}
                className="flex items-center gap-2"
              >
                <X className="w-4 h-4" />
                Cancel
              </Button>
              <Button
                variant="primary"
                size="lg"
                onClick={handleSave}
                disabled={saving || !form?.first_name.trim() || !form?.last_name.trim()}
                className="flex items-center gap-2"
              >
                {saving ? (uploadingAvatar ? 'Uploading photo…' : 'Saving…') : 'Save changes'}
              </Button>
            </div>
          ) : (
            <Button variant="secondary" size="lg" onClick={startEditing} className="flex items-center gap-2">
              <Pencil className="w-4 h-4" />
              Edit profile
            </Button>
          )}
        </div>

        {saveError && (
          <p className="mt-4 text-sm text-red-600">{saveError}</p>
        )}

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
            {isEditing && form ? (
              <textarea
                value={form.about_me}
                onChange={(e) => setForm({ ...form, about_me: e.target.value })}
                rows={4}
                placeholder="Tell people a bit about yourself..."
                className="mt-6 w-full rounded-xl border border-gray-200 px-4 py-3 text-sm text-gray-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 resize-none"
              />
            ) : (
              <p className="mt-6 text-gray-600 leading-7">
                {profile.about_me?.trim() || 'You have not added an about section yet. Add a short bio to let people know what you care about.'}
              </p>
            )}
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
                {isEditing && form ? (
                  <input
                    type="date"
                    value={form.dob}
                    onChange={(e) => setForm({ ...form, dob: e.target.value })}
                    className="rounded-lg border border-gray-200 px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  />
                ) : (
                  <span className="font-medium text-gray-900">{formatDate(profile.dob)}</span>
                )}
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
                <p className="text-sm text-gray-500">Birthday on file</p>
                <p className="mt-2 font-semibold text-gray-900">{profile.dob ? formatDate(profile.dob) : 'Date not available'}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
