"use client";

import React, { useState, useEffect } from 'react';
import { Bell } from 'lucide-react';
import { useWebSocket } from '@/contexts/WebSocketContext';

interface Notification {
  id: string;
  type: string;
  actor_id: string;
  created_at: string;
  is_read: boolean;
}

export function NotificationDropdown() {
  const { socket, isConnected } = useWebSocket();
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    if (!socket) return;

    const handleMessage = (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "notification" || data.type === "follow_request" || data.type === "group_invite" || data.type === "group_event") {
          setNotifications((prev) => [data.payload, ...prev]);
        }
      } catch (err) {
        console.error("Failed to parse websocket message", err);
      }
    };

    socket.addEventListener('message', handleMessage);

    return () => {
      socket.removeEventListener('message', handleMessage);
    };
  }, [socket]);

  const unreadCount = notifications.filter(n => !n.is_read).length;

  return (
    <div className="relative">
      <button 
        onClick={() => setIsOpen(!isOpen)}
        className="p-2 rounded-full hover:bg-gray-100 relative"
      >
        <Bell className="w-5 h-5 text-gray-600" />
        {unreadCount > 0 && (
          <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full border-2 border-white"></span>
        )}
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 bg-white rounded-xl shadow-lg border border-gray-100 overflow-hidden z-50">
          <div className="p-3 border-b border-gray-100 bg-gray-50 flex justify-between items-center">
            <h3 className="font-semibold text-gray-900">Notifications</h3>
            <span className="text-xs text-indigo-600 font-medium">
              {isConnected ? "Live" : "Connecting..."}
            </span>
          </div>
          <div className="max-h-96 overflow-y-auto">
            {notifications.length === 0 ? (
              <div className="p-4 text-center text-sm text-gray-500">
                No new notifications
              </div>
            ) : (
              notifications.map((notif) => (
                <div key={notif.id} className="p-3 hover:bg-gray-50 border-b border-gray-50 last:border-0 cursor-pointer transition-colors">
                  <p className="text-sm text-gray-800">
                    <span className="font-bold">User {notif.actor_id}</span> triggered a {notif.type}
                  </p>
                  <p className="text-xs text-gray-400 mt-1">
                    {new Date(notif.created_at).toLocaleTimeString()}
                  </p>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
}
