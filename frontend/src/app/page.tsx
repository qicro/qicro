'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/store/auth';
import { Loader2 } from 'lucide-react';
import { ThemeToggle } from '@/components/theme-toggle';

export default function Home() {
  const router = useRouter();
  const { user, token } = useAuthStore();

  useEffect(() => {
    // 如果用户已登录，重定向到聊天页面
    if (token && user) {
      router.push('/chat');
    } else {
      // 否则重定向到认证页面
      router.push('/auth');
    }
  }, [token, user, router]);

  // 显示加载状态
  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <div className="absolute top-4 left-4">
        <h1 className="text-2xl font-bold text-foreground">Qicro</h1>
      </div>
      <div className="absolute top-4 right-4">
        <ThemeToggle />
      </div>
      <div className="text-center">
        <Loader2 className="mx-auto h-8 w-8 animate-spin" />
        <p className="mt-4 text-muted-foreground">Loading Qicro...</p>
      </div>
    </div>
  );
}