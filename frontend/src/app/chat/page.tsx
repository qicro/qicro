import Chat from '@/components/chat/Chat';
import AuthGuard from '@/components/auth/AuthGuard';

export default function ChatPage() {
  return (
    <AuthGuard>
      <Chat />
    </AuthGuard>
  );
}