'use client';

import * as React from 'react';
import { motion } from 'framer-motion';
import { Bell, CheckCircle2, Clock3, MessageSquare, Sparkles, Users } from 'lucide-react';
import { useWebSocket } from '@/contexts/WebSocketContext';

type NotificationItem = {
  id: string;
  type: string;
  actor_id: string;
  created_at: string;
  is_read: boolean;
  message?: string;
};

function formatTimestamp(value: string) {
  const date = new Date(value);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60000);

  if (diffMins < 1) return 'Just now';
  if (diffMins < 60) return `${diffMins} minute${diffMins === 1 ? '' : 's'} ago`;

  const diffHours = Math.floor(diffMins / 60);
  if (diffHours < 24) return `${diffHours} hour${diffHours === 1 ? '' : 's'} ago`;

  const diffDays = Math.floor(diffHours / 24);
  return `${diffDays} day${diffDays === 1 ? '' : 's'} ago`;
}

function getNotificationTitle(type: string) {
  switch (type) {
    case 'follow_request':
      return 'Follow request';
    case 'group_invite':
      return 'Group invite';
    case 'group_event':
      return 'Group update';
    case 'notification':
      return 'New activity';
    default:
      return 'Notification';
  }
}

function getNotificationDescription(type: string, actorId: string) {
  switch (type) {
    case 'follow_request':
      return `User ${actorId} wants to follow you.`;
    case 'group_invite':
      return `User ${actorId} invited you to join a group.`;
    case 'group_event':
      return `User ${actorId} shared a new group update.`;
    case 'notification':
      return `User ${actorId} triggered a new notification.`;
    default:
      return `User ${actorId} sent an update.`;
  }
}

function getNotificationIcon(type: string) {
  switch (type) {
    case 'follow_request':
      return <Users className="w-5 h-5 text-indigo-600" />;
    case 'group_invite':
    case 'group_event':
      return <Sparkles className="w-5 h-5 text-amber-600" />;
    case 'notification':
      return <MessageSquare className="w-5 h-5 text-emerald-600" />;
    default:
      return <Bell className="w-5 h-5 text-slate-600" />;
  }
}
