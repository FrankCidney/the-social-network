'use client';

import * as React from 'react';
import { motion } from 'framer-motion';
import { Users, Plus, Search, Globe, Loader } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { groupAPI } from '@/lib/api';

type Group = {
  id: string;
  creator_id: string;
  title: string;
  description: string;
  created_at: string;
};

export default function GroupsPage() {
  const [searchQuery, setSearchQuery] = React.useState('');
  const [groups, setGroups] = React.useState<Group[]>([]);
  const [isLoading, setIsLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    const fetchGroups = async () => {
      try {
        // Backend integration point: fetch the list of groups from the API here.
        const response = await groupAPI.getGroups();
        const groupList = Array.isArray(response)
          ? response
          : response.groups ?? [];

        setGroups(groupList);
      } catch (err) {
        // Backend integration point: surface any API/network errors here.
        setError(err instanceof Error ? err.message : 'Failed to load groups');
      } finally {
        setIsLoading(false);
      }
    };

    fetchGroups();
  }, []);

  const filteredGroups = groups.filter((group) =>
    group.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    group.description.toLowerCase().includes(searchQuery.toLowerCase())
  );

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-96">
        <Loader className="w-8 h-8 text-indigo-600 animate-spin" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-xl p-6 text-red-600">
        <p className="font-bold">Error loading groups</p>
        <p className="text-sm mt-1">{error}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-xl border border-gray-100 p-6">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
              <div className="p-2 bg-indigo-100 rounded-lg">
                <Users className="w-6 h-6 text-indigo-600" />
              </div>
              Groups
            </h1>
            <p className="text-gray-500 text-sm mt-2">
              Discover communities and connect with people who share your interests.
            </p>
          </div>
          <Button variant="primary" size="lg" className="flex items-center gap-2">
            <Plus className="w-5 h-5" />
            {/* Backend integration point: connect this action to the create-group endpoint. */}
            Create Group
          </Button>
        </div>

        <div className="grid gap-4 sm:grid-cols-2 items-center">
          <div>
            <p className="text-sm text-gray-500">Browse available groups</p>
            <p className="text-xl font-semibold text-gray-900">{groups.length.toLocaleString()} groups available</p>
          </div>
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              type="text"
              placeholder="Search groups by name or interest..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full bg-gray-50 border border-gray-200 rounded-lg pl-10 pr-4 py-3 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
            />
          </div>
        </div>
      </div>

      <div className="space-y-4">
        {filteredGroups.length > 0 ? (
          filteredGroups.map((group, index) => (
            <motion.div
              key={group.id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.05 }}
              className="bg-white rounded-xl border border-gray-100 overflow-hidden hover:border-gray-200 transition-colors"
            >
              <div className="p-6">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center gap-4 flex-1">
                    <div className="w-16 h-16 rounded-lg bg-gradient-to-br from-indigo-400 to-indigo-600 flex items-center justify-center flex-shrink-0">
                      <span className="text-white font-bold text-xl">{group.title.charAt(0).toUpperCase()}</span>
                    </div>

                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <h3 className="text-lg font-bold text-gray-900">{group.title}</h3>
                        <span title="Public community">
                          <Globe className="w-4 h-4 text-gray-400" />
                        </span>
                      </div>
                      <p className="text-gray-600 text-sm mb-3">{group.description}</p>
                      <div className="flex items-center gap-4 text-xs text-gray-500">
                        <span className="flex items-center gap-1">
                          <Users className="w-4 h-4" />
                          Created {new Date(group.created_at).toLocaleDateString('en-US', {
                            year: 'numeric',
                            month: 'short',
                            day: 'numeric',
                          })}
                        </span>
                        <span className="rounded-full bg-gray-100 px-2 py-1">Creator: {group.creator_id}</span>
                      </div>
                    </div>
                  </div>

                  <Button variant="secondary" size="md">
                    {/* Backend integration point: route this action to the group details or join flow. */}
                    View Group
                  </Button>
                </div>
              </div>
            </motion.div>
          ))
        ) : (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="bg-white rounded-xl border border-gray-100 p-12 text-center"
          >
            <Users className="w-12 h-12 text-gray-300 mx-auto mb-4" />
            <h3 className="text-lg font-bold text-gray-900 mb-2">No groups found</h3>
            <p className="text-gray-500 text-sm">Try adjusting your search term or check back later.</p>
          </motion.div>
        )}
      </div>
    </div>
  );
}
