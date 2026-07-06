'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { Home, Users, MessageSquare, Bell, User, LogOut, Search } from 'lucide-react';
import { PeopleDiscoveryPanel } from '@/components/discovery/PeopleDiscoveryPanel';
import { WebSocketProvider } from '@/contexts/WebSocketContext';
import { NotificationDropdown } from '@/components/notifications/NotificationDropdown';
import { authAPI, isAuthenticationError, profileAPI } from '@/lib/api';

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const router = useRouter();
  const searchParams = useSearchParams();
  const isGroupsRoute = pathname.startsWith('/groups');
  const [authChecked, setAuthChecked] = React.useState(false);
  const [authError, setAuthError] = React.useState<string | null>(null);
  const [isLoggingOut, setIsLoggingOut] = React.useState(false);

  React.useEffect(() => {
    let cancelled = false;

    const verifySession = async () => {
      setAuthChecked(false);
      setAuthError(null);

      try {
        await profileAPI.getMyProfile();
        if (!cancelled) {
          setAuthChecked(true);
        }
      } catch (error) {
        if (cancelled) {
          return;
        }

        if (isAuthenticationError(error)) {
          const queryString = searchParams.toString();
          const nextPath = `${pathname}${queryString ? `?${queryString}` : ''}`;
          router.replace(`/login?next=${encodeURIComponent(nextPath)}`);
          return;
        }

        setAuthError(
          error instanceof Error ? error.message : 'Unable to verify your session.'
        );
        setAuthChecked(true);
      }
    };

    void verifySession();

    return () => {
      cancelled = true;
    };
  }, [pathname, router, searchParams]);

  const handleLogout = React.useCallback(async () => {
    setIsLoggingOut(true);
    setAuthError(null);

    try {
      await authAPI.logout();
    } catch (error) {
      if (!isAuthenticationError(error)) {
        setAuthError(
          error instanceof Error ? error.message : 'Unable to log out right now.'
        );
        setIsLoggingOut(false);
        return;
      }
    }

    router.replace('/login');
    router.refresh();
  }, [router]);

  if (!authChecked) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="mx-auto flex min-h-screen max-w-7xl items-center justify-center px-4">
          <div className="rounded-xl border border-gray-100 bg-white px-6 py-4 text-sm text-gray-500 shadow-sm">
            Checking your session...
          </div>
        </div>
      </div>
    );
  }

  if (authError) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="mx-auto flex min-h-screen max-w-7xl items-center justify-center px-4">
          <div className="rounded-xl border border-red-200 bg-white px-6 py-5 text-center">
            <p className="text-sm text-red-600">{authError}</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <WebSocketProvider>
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
              <NotificationDropdown />
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
          <SidebarItem icon={<Home className="w-5 h-5" />} label="Home Feed" href="/feed" active={pathname === '/feed'} />
          <SidebarItem icon={<Users className="w-5 h-5" />} label="Groups" href="/groups" active={pathname.startsWith('/groups')} />
          <SidebarItem icon={<MessageSquare className="w-5 h-5" />} label="Messages" href="/messages" active={pathname === '/messages'} />
          <SidebarItem icon={<Bell className="w-5 h-5" />} label="Notifications" href="/notifications" active={pathname === '/notifications'} />
          <SidebarItem icon={<User className="w-5 h-5" />} label="Profile" href="/profile" active={pathname === '/profile'} />
          <hr className="my-4 border-gray-100" />
          <SidebarAction
            icon={<LogOut className="w-5 h-5" />}
            label={isLoggingOut ? 'Logging out...' : 'Logout'}
            onClick={handleLogout}
            disabled={isLoggingOut}
          />
        </aside>

          {/* Main Content */}
          <main className={isGroupsRoute ? 'md:col-span-9' : 'md:col-span-6'}>
            {children}
          </main>

          {/* Right Sidebar (Suggestions/Trends) */}
          {!isGroupsRoute && (
          <aside className="hidden lg:block lg:col-span-3 space-y-6">
            <PeopleDiscoveryPanel />
          </aside>
          )}
        </div>
      </div>
    </WebSocketProvider>
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

function SidebarAction({
  icon,
  label,
  onClick,
  disabled = false,
}: {
  icon: React.ReactNode;
  label: string;
  onClick: () => void;
  disabled?: boolean;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      className="flex w-full items-center gap-3 rounded-xl px-4 py-3 text-left font-medium text-gray-600 transition-colors hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-60"
    >
      {icon}
      <span>{label}</span>
    </button>
  );
}
