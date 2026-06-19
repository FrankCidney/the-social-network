import * as React from 'react';
import { motion } from 'framer-motion';
import { Image, MessageCircle, MoreHorizontal, Share2 } from 'lucide-react';
import { Button } from '../components/ui/Button';

const Feed = () => {
  return (
    <div className="space-y-6">
      {/* Composer (Summary view) */}
      <div className="bg-surface p-4 rounded-bento border border-gray-100 flex gap-4">
        <div className="w-10 h-10 rounded-full bg-indigo-100 flex-shrink-0" />
        <button className="flex-1 bg-gray-50 rounded-full px-5 text-left text-gray-500 text-sm hover:bg-gray-100 transition-colors">
          What's on your mind?
        </button>
        <div className="flex gap-2">
          <Button variant="ghost" size="sm" className="rounded-full w-10 p-0">
            <Image className="w-5 h-5 text-gray-400" />
          </Button>
        </div>
      </div>

      {/* Posts */}
      {[1, 2].map((post) => (
        <motion.div 
          key={post}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-surface rounded-bento border border-gray-100 overflow-hidden"
        >
          {/* Post Header */}
          <div className="p-4 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-gray-200" />
              <div>
                <p className="text-sm font-bold">Alex Johnson</p>
                <p className="text-xs text-gray-500">2 hours ago • Public</p>
              </div>
            </div>
            <button className="text-gray-400 hover:text-gray-600">
              <MoreHorizontal className="w-5 h-5" />
            </button>
          </div>

          {/* Post Content */}
          <div className="px-4 pb-4 space-y-3">
            <p className="text-sm leading-relaxed">
              Just finished setting up the new office space! Really happy with how the "Bento Box" layout turned out for the project dashboard. 🍱✨
            </p>
            <div className="aspect-video bg-gray-50 rounded-xl border border-gray-100 flex items-center justify-center">
              <span className="text-gray-300 italic text-sm">Post Image Placeholder</span>
            </div>
          </div>

          {/* Post Actions */}
          <div className="px-4 py-3 bg-gray-50/50 border-t border-gray-100 flex items-center gap-6">
            <button className="flex items-center gap-2 text-xs font-bold text-gray-500 hover:text-primary transition-colors">
              <MessageCircle className="w-4 h-4" />
              12 Comments
            </button>
            <button className="flex items-center gap-2 text-xs font-bold text-gray-500 hover:text-primary transition-colors">
              <Share2 className="w-4 h-4" />
              Share
            </button>
          </div>
        </motion.div>
      ))}
    </div>
  );
};

export default Feed;
