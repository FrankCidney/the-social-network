'use client';

import * as React from 'react';
import { motion } from 'framer-motion';
import { Bell, CheckCircle2, Clock3, MessageSquare, Sparkles, Users } from 'lucide-react';
import { useWebSocket } from '@/contexts/WebSocketContext';
import { notificationsAPI, NotificationItem, NotificationListResponse } from '@/lib/api';

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

// The endpoint may return a bare array or a wrapped object — handle both.
function unwrapNotifications(
  data: NotificationListResponse | NotificationItem[]
): NotificationItem[] {
  return Array.isArray(data) ? data : data.notifications;
}

export default function NotificationsPage() {
  const { socket, isConnected } = useWebSocket();

  const [notifications, setNotifications] = React.useState<NotificationItem[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [markingAll, setMarkingAll] = React.useState(false);

  const loadNotifications = React.useCallback(async () => {
    try {
      setLoading(true);
      const data = await notificationsAPI.getNotifications();
      setNotifications(unwrapNotifications(data));
      setError(null);
    } catch (err) {
      console.error('Failed to load notifications', err);
      setError('Could not load your notifications.');
    } finally {
      setLoading(false);
    }
  }, []);

  React.useEffect(() => {
    loadNotifications();
  }, [loadNotifications]);

  // Live notifications pushed over the shared WebSocket connection.
  React.useEffect(() => {
    if (!socket) return;

    const handleMessage = (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);
        if (
          data.type === 'notification' ||
          data.type === 'follow_request' ||
          data.type === 'group_invite' ||
          data.type === 'group_event'
        ) {
          const payload = data.payload as NotificationItem;
          setNotifications((prev) => {
            if (prev.some((item) => item.id === payload.id)) return prev;
            return [
              {
                ...payload,
                message: payload.message ?? getNotificationDescription(payload.type, payload.actor_id),
              },
              ...prev,
            ];
          });
        }
      } catch (err) {
        console.error('Failed to parse notification payload', err);
      }
    };

    socket.addEventListener('message', handleMessage);
    return () => socket.removeEventListener('message', handleMessage);
  }, [socket]);

  const unreadCount = notifications.filter((item) => !item.is_read).length;

  const handleMarkAsRead = async (notificationId: string) => {
    const previous = notifications;
    setNotifications((prev) =>
      prev.map((item) => (item.id === notificationId ? { ...item, is_read: true } : item))
    );
    try {
      await notificationsAPI.markAsRead(notificationId);
    } catch (err) {
      console.error('Failed to mark notification as read', err);
      setNotifications(previous); // roll back on failure
    }
  };

  const markAllAsRead = async () => {
    if (unreadCount === 0) return;
    const previous = notifications;
    setNotifications((prev) => prev.map((item) => ({ ...item, is_read: true })));
    try {
      setMarkingAll(true);
      await notificationsAPI.markAllAsRead();
    } catch (err) {
      console.error('Failed to mark all as read', err);
      setNotifications(previous);
      setError('Could not mark everything as read. Please try again.');
    } finally {
      setMarkingAll(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-xl border border-gray-100 p-6">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
              <div className="p-2 bg-amber-100 rounded-lg">
                <Bell className="w-6 h-6 text-amber-600" />
              </div>
              Notifications
            </h1>
            <p className="text-gray-500 text-sm mt-2">
              Keep track of follow requests, group invites, and other recent activity in one place.
            </p>
          </div>

          <button
            onClick={markAllAsRead}
            disabled={markingAll || unreadCount === 0}
            className="inline-flex items-center gap-2 rounded-full border border-gray-200 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <CheckCircle2 className="w-4 h-4" />
            {markingAll ? 'Marking…' : 'Mark all read'}
          </button>
        </div>

        <div className="mt-6 grid gap-4 sm:grid-cols-3">
          <div className="rounded-xl border border-gray-100 bg-gray-50 p-4">
            <p className="text-sm text-gray-500">Unread</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">{unreadCount}</p>
          </div>
          <div className="rounded-xl border border-gray-100 bg-gray-50 p-4">
            <p className="text-sm text-gray-500">Live updates</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">{isConnected ? 'On' : 'Connecting'}</p>
          </div>
          <div className="rounded-xl border border-gray-100 bg-gray-50 p-4">
            <p className="text-sm text-gray-500">Total</p>
            <p className="mt-2 text-2xl font-semibold text-gray-900">{notifications.length}</p>
          </div>
        </div>
      </div>

      <div className="space-y-4">
        {loading && (
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-white rounded-xl border border-gray-100 p-4 animate-pulse">
                <div className="flex items-start gap-3">
                  <div className="w-9 h-9 rounded-full bg-gray-100 shrink-0" />
                  <div className="flex-1 space-y-2">
                    <div className="w-40 h-3 bg-gray-100 rounded" />
                    <div className="w-64 h-3 bg-gray-100 rounded" />
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {!loading && error && (
          <div className="bg-white rounded-xl border border-gray-100 p-8 text-center">
            <p className="text-sm text-red-600 mb-3">{error}</p>
            <button
              onClick={loadNotifications}
              className="text-sm text-indigo-600 hover:underline"
            >
              Retry
            </button>
          </div>
        )}

        {!loading && !error && notifications.length === 0 && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="bg-white rounded-xl border border-gray-100 p-12 text-center"
          >
            <Bell className="w-12 h-12 text-gray-300 mx-auto mb-4" />
            <h3 className="text-lg font-bold text-gray-900 mb-2">No notifications yet</h3>
            <p className="text-gray-500 text-sm">
              New activity will appear here as soon as the app sends it.
            </p>
          </motion.div>
        )}

        {!loading &&
          !error &&
          notifications.map((notification, index) => (
            <motion.div
              key={notification.id}
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.04 }}
              onClick={() => !notification.is_read && handleMarkAsRead(notification.id)}
              className={`rounded-xl border p-4 transition-colors ${
                notification.is_read
                  ? 'border-gray-100 bg-white'
                  : 'border-indigo-100 bg-indigo-50/70 cursor-pointer hover:bg-indigo-50'
              }`}
            >
              <div className="flex items-start gap-3">
                <div className="mt-0.5 rounded-full bg-white p-2 shadow-sm">
                  {getNotificationIcon(notification.type)}
                </div>
                <div className="flex-1">
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <div>
                      <p className="font-semibold text-gray-900">{getNotificationTitle(notification.type)}</p>
                      <p className="text-sm text-gray-600 mt-1">
                        {notification.message ?? getNotificationDescription(notification.type, notification.actor_id)}
                      </p>
                    </div>
                    <div className="flex items-center gap-2 text-xs text-gray-500">
                      <Clock3 className="w-3.5 h-3.5" />
                      {formatTimestamp(notification.created_at)}
                    </div>
                  </div>
                </div>
              </div>
            </motion.div>
          ))}
      </div>
    </div>
  );
}
