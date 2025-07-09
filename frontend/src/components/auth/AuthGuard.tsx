'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth';
import { Loader2 } from 'lucide-react';

interface AuthGuardProps {
  children: React.ReactNode;
}

export default function AuthGuard({ children }: AuthGuardProps) {
  const router = useRouter();
  const { user, token, isLoading, getProfile } = useAuthStore();

  useEffect(() => {
    // 如果没有token，重定向到登录页
    if (!token) {
      router.push('/auth');
      return;
    }

    // 如果有token但没有用户信息，尝试获取用户信息
    if (token && !user) {
      getProfile();
    }
  }, [token, user, router, getProfile]);

  // 正在加载或没有认证信息时显示加载状态
  if (isLoading || !token || !user) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <Loader2 className="mx-auto h-8 w-8 animate-spin" />
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return <>{children}</>;
}