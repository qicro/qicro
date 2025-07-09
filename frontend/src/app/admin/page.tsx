import AdminDashboard from '@/components/admin/AdminDashboard';
import AuthGuard from '@/components/auth/AuthGuard';

export default function AdminPage() {
  return (
    <AuthGuard>
      <AdminDashboard />
    </AuthGuard>
  );
}