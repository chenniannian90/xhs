import { Navigate } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

/**
 * 路由认证保护组件
 *
 * 功能：
 * - 检查用户是否已登录
 * - 未登录用户自动重定向到登录页
 * - 保留原始目标 URL，登录后可以返回
 *
 * 使用示例：
 * <Route path="/dashboard" element={
 *   <ProtectedRoute>
 *     <DashboardPage />
 *   </ProtectedRoute>
 * } />
 */
export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated } = useAuthStore();

  if (!isAuthenticated) {
    // 保存当前路径，登录后可以返回
    const currentPath = window.location.pathname + window.location.search;
    const loginPath = `/login?redirect=${encodeURIComponent(currentPath)}`;

    return <Navigate to={loginPath} replace />;
  }

  return <>{children}</>;
}
