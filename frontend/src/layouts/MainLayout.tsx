import * as React from 'react';
import { Outlet, Link, useLocation } from 'react-router-dom';
import { 
  Home, 
  Users, 
  UserCircle, 
  MessageSquare, 
  Bell, 
  Search,
  LogOut,
  PlusSquare
} from 'lucide-react';
import { cn } from '../lib/utils';

const MainLayout = () => {
  const location = useLocation();

  const navItems = [
    { icon: Home, label: 'Feed', path: '/feed' },
    { icon: Users, label: 'Groups', path: '/groups' },
    { icon: MessageSquare, label: 'Messages', path: '/messages' },
    { icon: Bell, label: 'Notifications', path: '/notifications' },
    { icon: UserCircle, label: 'Profile', path: '/profile' },
  ];

  return (
    <div className="min-h-screen bg-canvas text-text-main flex justify-center">
      <div className="max-w-[1280px] w-full flex gap-6 p-6">
        
        {/* Left Sidebar (Desktop) */}
        <aside className="hidden lg:flex flex-col w-64 gap-6 sticky top-6 h-fit">
          <div className="bg-surface p-6 rounded-bento border border-gray-100 space-y-8">
            <h1 className="text-xl font-bold text-primary px-2">SocialNet</h1>
            
            <nav className="space-y-1">
              {navItems.map((item) => (
                <Link
                  key={item.path}
                  to={item.path}
                  className={cn(
                    "flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors",
                    location.pathname === item.path 
                      ? "bg-indigo-50 text-primary" 
                      : "text-gray-600 hover:bg-gray-50"
                  )}
                >
                  <item.icon className="w-5 h-5" />
                  {item.label}
                </Link>
              ))}
            </nav>

            <button className="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-privacy-private hover:bg-red-50 w-full transition-colors">
              <LogOut className="w-5 h-5" />
              Sign Out
            </button>
          </div>

          <button className="bg-primary text-white p-4 rounded-bento shadow-lg shadow-indigo-100 flex items-center justify-center gap-2 font-bold hover:bg-opacity-90 transition-all">
            <PlusSquare className="w-5 h-5" />
            Create Post
          </button>
        </aside>

        {/* Main Content Area */}
        <main className="flex-1 max-w-2xl space-y-6">
          <Outlet />
        </main>

        {/* Right Sidebar (Desktop) */}
        <aside className="hidden xl:flex flex-col w-80 gap-6 sticky top-6 h-fit">
          <div className="bg-surface p-4 rounded-bento border border-gray-100">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input 
                type="text" 
                placeholder="Search groups..." 
                className="w-full bg-gray-50 border-none rounded-lg py-2 pl-10 pr-4 text-sm focus:ring-2 focus:ring-primary"
              />
            </div>
          </div>

          <div className="bg-surface p-6 rounded-bento border border-gray-100 space-y-4">
            <h3 className="font-bold text-sm text-gray-400 uppercase tracking-wider px-1">Suggested Groups</h3>
            <div className="space-y-4">
              {[1, 2, 3].map((i) => (
                <div key={i} className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-gray-100" />
                  <div className="flex-1">
                    <p className="text-sm font-bold">Tech Enthusiasts</p>
                    <p className="text-xs text-gray-500">1.2k members</p>
                  </div>
                  <button className="text-xs font-bold text-primary hover:underline">Join</button>
                </div>
              ))}
            </div>
          </div>
        </aside>

        {/* Mobile Nav (Bottom) */}
        <nav className="lg:hidden fixed bottom-0 left-0 right-0 bg-surface border-t border-gray-100 flex justify-around p-3 z-50">
          {navItems.map((item) => (
            <Link key={item.path} to={item.path} className="p-2">
              <item.icon className={cn("w-6 h-6", location.pathname === item.path ? "text-primary" : "text-gray-400")} />
            </Link>
          ))}
        </nav>

      </div>
    </div>
  );
};

export default MainLayout;
