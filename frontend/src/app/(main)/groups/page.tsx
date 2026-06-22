"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import { Users, Plus } from "lucide-react";

interface Group {
  id: string;
  title: string;
  description: string;
  creator_id: string;
}

export default function GroupsPage() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("/api/groups")
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data)) {
          setGroups(data);
        }
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="p-8 text-center text-gray-500">Loading groups...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
          <Users className="w-6 h-6 text-indigo-600" />
          Groups
        </h1>
        <button className="flex items-center gap-2 bg-indigo-600 text-white px-4 py-2 rounded-xl hover:bg-indigo-700 transition-colors font-medium">
          <Plus className="w-4 h-4" />
          Create Group
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {groups.length === 0 ? (
          <div className="col-span-full p-8 text-center text-gray-500 bg-white rounded-xl border border-gray-100">
            No groups found. Be the first to create one!
          </div>
        ) : (
          groups.map((g) => (
            <Link 
              key={g.id} 
              href={`/groups/${g.id}`}
              className="bg-white p-5 rounded-xl border border-gray-100 shadow-sm hover:shadow-md transition-shadow cursor-pointer block"
            >
              <h3 className="font-bold text-lg text-gray-900 mb-2">{g.title}</h3>
              <p className="text-gray-600 text-sm mb-4 line-clamp-2">{g.description}</p>
              <div className="text-xs text-indigo-600 font-medium">View Group &rarr;</div>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}
