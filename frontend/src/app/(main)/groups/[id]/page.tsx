"use client";

import React, { useEffect, useState, use } from "react";
import { Users, Calendar, MessageSquare, Plus } from "lucide-react";

interface Group {
  id: string;
  title: string;
  description: string;
  creator_id: string;
}

interface Event {
  id: string;
  title: string;
  description: string;
  event_date: string;
}

interface GroupDetailPageProps {
  params: Promise<{ id: string }>;
}

export default function GroupDetailPage({ params }: GroupDetailPageProps) {
  const { id } = use(params);
  const [group, setGroup] = useState<Group | null>(null);
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      fetch(`/api/groups/${id}`).then((res) => res.json()),
      fetch(`/api/groups/${id}/events`).then((res) => res.json())
    ])
      .then(([groupData, eventsData]) => {
        if (groupData && !groupData.error) {
          setGroup(groupData);
        }
        if (Array.isArray(eventsData)) {
          setEvents(eventsData);
        }
      })
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) {
    return <div className="p-8 text-center text-gray-500">Loading group details...</div>;
  }

  if (!group) {
    return <div className="p-8 text-center text-red-500 font-bold">Group not found</div>;
  }

  return (
    <div className="space-y-6">
      {/* Group Header */}
      <div className="bg-white rounded-xl border border-gray-100 p-6 shadow-sm">
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2 mb-2">
              <Users className="w-6 h-6 text-indigo-600" />
              {group.title}
            </h1>
            <p className="text-gray-600 mb-4">{group.description}</p>
          </div>
          <button className="bg-indigo-50 text-indigo-700 px-4 py-2 rounded-xl hover:bg-indigo-100 transition-colors font-medium text-sm">
            Join Group
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Main Content Area */}
        <div className="md:col-span-2 space-y-6">
          {/* Group Posts placeholder */}
          <div className="bg-white rounded-xl border border-gray-100 p-6 shadow-sm">
            <h2 className="font-bold text-gray-900 mb-4 flex items-center gap-2">
              <MessageSquare className="w-5 h-5 text-indigo-600" />
              Group Discussion
            </h2>
            <div className="flex gap-4">
              <div className="w-10 h-10 rounded-full bg-gray-200 shrink-0" />
              <input 
                type="text" 
                placeholder="Share something with the group..." 
                className="flex-1 bg-gray-50 border border-gray-200 rounded-xl px-4 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>
          </div>
        </div>

        {/* Sidebar / Events */}
        <div className="space-y-6">
          <div className="bg-white rounded-xl border border-gray-100 p-6 shadow-sm">
            <div className="flex items-center justify-between mb-4">
              <h2 className="font-bold text-gray-900 flex items-center gap-2">
                <Calendar className="w-5 h-5 text-indigo-600" />
                Events
              </h2>
              <button className="text-indigo-600 hover:text-indigo-800 p-1">
                <Plus className="w-4 h-4" />
              </button>
            </div>
            
            {events.length === 0 ? (
              <p className="text-sm text-gray-500">No upcoming events.</p>
            ) : (
              <div className="space-y-4">
                {events.map((e) => (
                  <div key={e.id} className="p-3 border border-gray-100 rounded-lg hover:bg-gray-50 transition-colors">
                    <h3 className="font-semibold text-gray-900 text-sm mb-1">{e.title}</h3>
                    <p className="text-xs text-gray-500 mb-2">{new Date(e.event_date).toLocaleString()}</p>
                    <button className="text-xs font-medium text-indigo-600 bg-indigo-50 px-2 py-1 rounded">RSVP</button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
