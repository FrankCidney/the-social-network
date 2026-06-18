'use client';

import React from 'react';
import Link from 'next/link';
import { Home, Users, MessageSquare, Bell, User, LogOut, Search } from 'lucide-react';

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-gray-50">
      {/* Navigation Bar */}
      <nav className="sticky top-0 z-50 bg-white border-b border-gray-100">
        <div className="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
          <div className="flex items-center gap-8">
            <Link href="/feed" className="text-2xl font-bold text-indigo-600">Social</Link>
            <div className="hidden md:flex relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input 
                type="text" 
                placeholder="Search social..." 
                className="bg-gray-100 rounded-full py-2 pl-10 pr-4 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 w-64 transition-all"
              />
            </div>
          </div>
          
          <div className="flex items-center gap-2">
            <button className="p-2 rounded-full hover:bg-gray-100 relative">
              <Bell className="w-5 h-5 text-gray-600" />
              <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full border-2 border-white"></span>
            </button>
            <Link href="/profile" className="flex items-center gap-2 p-1 pl-3 rounded-full hover:bg-gray-100 transition-colors">
              <span className="text-sm font-medium hidden sm:inline">My Profile</span>
              <div className="w-8 h-8 rounded-full bg-indigo-100" />
            </Link>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 py-6 grid grid-cols-1 md:grid-cols-12 gap-6">
        {/* Sidebar */}
        <aside className="hidden md:block md:col-span-3 space-y-2">
          <SidebarItem icon={<Home className="w-5 h-5" />} label="Home Feed" href="/feed" active />
          <SidebarItem icon={<Users className="w-5 h-5" />} label="Groups" href="/groups" />
          <SidebarItem icon={<MessageSquare className="w-5 h-5" />} label="Messages" href="/messages" />
          <SidebarItem icon={<Bell className="w-5 h-5" />} label="Notifications" href="/notifications" />
          <SidebarItem icon={<User className="w-5 h-5" />} label="Profile" href="/profile" />
          <hr className="my-4 border-gray-100" />
          <SidebarItem icon={<LogOut className="w-5 h-5" />} label="Logout" href="/login" />
        </aside>

        {/* Main Content */}
        <main className="md:col-span-6">
          {children}
        </main>

        {/* Right Sidebar (Suggestions/Trends) */}
        <aside className="hidden lg:block lg:col-span-3 space-y-6">
          <div className="bg-white rounded-xl border border-gray-100 p-4">
            <h3 className="font-bold text-gray-900 mb-4">Who to follow</h3>
            <div className="space-y-4">
              {[1, 2, 3].map(i => (
                <div key={i} className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-full bg-gray-200" />
                    <div className="text-sm">
                      <p className="font-bold">User {i}</p>
                      <p className="text-gray-500 text-xs">@user_{i}</p>
                    </div>
                  </div>
                  <button className="text-xs font-bold text-indigo-600 hover:text-indigo-700">Follow</button>
                </div>
              ))}
            </div>
          </div>
        </aside>
      </div>
    </div>
  );
}

function SidebarItem({ icon, label, href, active = false }: { icon: React.ReactNode, label: string, href: string, active?: boolean }) {
  return (
    <Link 
      href={href}
      className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-colors font-medium ${
        active 
          ? 'bg-indigo-50 text-indigo-700' 
          : 'text-gray-600 hover:bg-gray-100'
      }`}
    >
      {icon}
      <span>{label}</span>
    </Link>
  );
}
