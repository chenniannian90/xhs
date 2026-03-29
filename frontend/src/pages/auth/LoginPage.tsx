import { useState } from 'react';
import { useNavigate, Link, useSearchParams } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import { PASSWORD_RULES } from '../../utils/passwordRules';
import Button from '../../components/ui/Button';
import Input from '../../components/ui/Input';
import PasswordInput from '../../components/ui/PasswordInput';

export default function LoginPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { login } = useAuthStore();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // 获取登录后要重定向的路径
  const redirectTo = searchParams.get('redirect') || '/dashboard';
  const [showPasswordRules, setShowPasswordRules] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login(email, password);
      navigate(redirectTo, { replace: true });
    } catch (err: any) {
      const errorMessage = err.message || '登录失败';

      // 提供更友好的错误提示
      if (errorMessage.toLowerCase().includes('invalid') ||
          errorMessage.toLowerCase().includes('邮箱或密码错误') ||
          errorMessage.toLowerCase().includes('invalid email or password')) {
        setError('用户名或密码错误');
      } else if (errorMessage.toLowerCase().includes('account') ||
                 errorMessage.toLowerCase().includes('deactivated')) {
        setError('账号已被停用，请联系管理员。');
      } else if (errorMessage.toLowerCase().includes('not found')) {
        setError('账号不存在，请先注册。');
      } else {
        setError(errorMessage);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center px-4">
      <div className="max-w-md w-full">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">欢迎回来</h1>
          <p className="text-gray-600 dark:text-gray-400">登录到 NavHub</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4" autoComplete="on">
          <Input
            label="邮箱"
            type="email"
            name="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="your@email.com"
            required
            autoComplete="email"
          />

          <div>
            <div className="flex items-center justify-between mb-1">
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                密码
              </label>
              <div className="flex gap-2">
                <Link
                  to="/forgot-password"
                  className="text-xs text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
                >
                  忘记密码？
                </Link>
                <button
                  type="button"
                  onClick={() => setShowPasswordRules(!showPasswordRules)}
                  className="text-xs text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
                >
                  {showPasswordRules ? '隐藏' : '显示'}密码规则
                </button>
              </div>
            </div>
            <PasswordInput
              label=""
              name="current-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
              autoComplete="current-password"
            />

            {/* 密码规则提示（可折叠） */}
            {showPasswordRules && (
              <div className="mt-2 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
                <p className="text-sm font-medium text-blue-800 dark:text-blue-300 mb-2">
                  密码规则：
                </p>
                <ul className="space-y-1">
                  {PASSWORD_RULES.map((rule, index) => (
                    <li key={index} className="flex items-center text-sm text-gray-700 dark:text-gray-300">
                      <svg className="w-4 h-4 text-blue-500 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
                      </svg>
                      {rule.text}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>

          {error && (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 px-4 py-3 rounded-lg">
              <div className="flex items-start gap-2">
                <span className="text-xl">⚠️</span>
                <div className="flex-1">
                  <pre className="whitespace-pre-wrap text-sm">{error}</pre>
                  {error.includes('忘记密码') && (
                    <div className="mt-2">
                      <Link to="/forgot-password" className="text-sm underline hover:text-red-700 dark:hover:text-red-300">
                        前往重置密码 →
                      </Link>
                    </div>
                  )}
                </div>
              </div>
            </div>
          )}

          <Button type="submit" className="w-full" disabled={loading}>
            {loading ? '登录中...' : '登录'}
          </Button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-gray-600 dark:text-gray-400">
            还没有账号？{' '}
            <Link to="/register" className="text-primary-600 hover:text-primary-700 font-medium">
              注册
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
