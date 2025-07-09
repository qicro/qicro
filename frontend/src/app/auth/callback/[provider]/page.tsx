'use client';

import { useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuthStore } from '@/store/auth';
import { Loader2 } from 'lucide-react';

interface OAuthCallbackProps {
  params: Promise<{
    provider: string;
  }>;
}

export default function OAuthCallback({ params }: OAuthCallbackProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { handleOAuthCallback, error } = useAuthStore();

  useEffect(() => {
    const handleCallback = async () => {
      const resolvedParams = await params;
      const code = searchParams.get('code');
      const error_param = searchParams.get('error');

      if (error_param) {
        router.push('/auth?error=' + encodeURIComponent(error_param));
        return;
      }

      if (code) {
        try {
          await handleOAuthCallback(resolvedParams.provider, code);
          router.push('/chat');
        } catch {
          router.push('/auth?error=' + encodeURIComponent('OAuth login failed'));
        }
      } else {
        router.push('/auth?error=' + encodeURIComponent('No authorization code received'));
      }
    };

    handleCallback();
  }, [searchParams, params, handleOAuthCallback, router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <Loader2 className="mx-auto h-8 w-8 animate-spin" />
        <p className="mt-4 text-gray-600">Completing authentication...</p>
        {error && (
          <p className="mt-2 text-red-600">{error}</p>
        )}
      </div>
    </div>
  );
}