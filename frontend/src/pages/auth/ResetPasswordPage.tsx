import { useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { CheckCircle, Eye, EyeOff, Lock } from 'lucide-react';
import { getPasswordRulesStatus } from '../../utils/passwordRules';
import api from '../../utils/api';
import Button from '../../components/ui/Button';
import Input from '../../components/ui/Input';
import PasswordInput from '../../components/ui/PasswordInput';
import Card from '../../components/ui/Card';

export default function ResetPasswordPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');

  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  // 检查密码规则
  const rulesStatus = getPasswordRulesStatus(password);
  const allRulesPassed = rulesStatus.every(rule => rule.passed);
  const passwordsMatch = password === confirmPassword && password !== '';

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    // 验证密码
    if (!allRulesPassed) {
      setError('密码不符合要求，请检查下方密码规则');
      return;
    }

    if (!passwordsMatch) {
      setError('两次输入的密码不一致');
      return;
    }

    setIsLoading(true);

    try {
      await api.post('/auth/reset-password', {
        token,
        password,
      });
      setSuccess(true);

      // 3秒后跳转到登录页
      setTimeout(() => {
        navigate('/login');
      }, 3000);
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || '重置失败，请稍后重试';

      // 提供更友好的错误提示
      if (errorMessage.includes('invalid') || errorMessage.includes('token')) {
        setError('重置链接无效或已过期，请重新申请密码重置');
      } else {
        setError(errorMessage);
      }
    } finally {
      setIsLoading(false);
    }
  };

  if (!token) {
    return (
      <div className="min-h-screen flex items-center justify-center px-4 bg-gray-50 dark:bg-gray-900">
        <Card className="max-w-md w-full p-8 text-center">
          <h1 className="text-xl font-bold text-gray-900 dark:text-white mb-4">
            无效的重置链接
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mb-6">
            重置链接不完整或已过期。请重新申请密码重置。
          </p>
          <Link to="/forgot-password">
            <Button>重新发送重置邮件</Button>
          </Link>
        </Card>
      </div>
    );
  }

  if (success) {
    return (
      <div className="min-h-screen flex items-center justify-center px-4 bg-gray-50 dark:bg-gray-900">
        <Card className="max-w-md w-full p-8 text-center">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-green-100 dark:bg-green-900 flex items-center justify-center">
            <CheckCircle size={32} className="text-green-600 dark:text-green-400" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
            密码重置成功！
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mb-6">
            您的密码已成功重置。现在可以使用新密码登录了。
          </p>
          <p className="text-sm text-gray-500 dark:text-gray-500">
            页面将在 3 秒后自动跳转到登录页面...
          </p>
          <Link to="/login">
            <Button>立即登录</Button>
          </Link>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center px-4 bg-gray-50 dark:bg-gray-900">
      <Card className="max-w-md w-full p-8">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
            <Lock size={32} className="text-primary-600 dark:text-primary-400" />
          </div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
            设置新密码
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            请输入您的新密码
          </p>
        </div>

        {error && (
          <div className="mb-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 px-4 py-3 rounded-lg">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4" autoComplete="on">
          <div>
            <PasswordInput
              label="新密码"
              name="new-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="•••••••••"
              required
              autoComplete="new-password"
            />

            {/* 密码规则提示 */}
            {password && (
              <div className="mt-2 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
                <p className="text-sm font-medium text-blue-800 dark:text-blue-300 mb-2">
                  密码要求：
                </p>
                <ul className="space-y-1">
                  {rulesStatus.map((rule, index) => (
                    <li key={index} className="flex items-center text-sm">
                      {rule.passed ? (
                        <svg className="w-4 h-4 text-green-500 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                        </svg>
                      ) : (
                        <svg className="w-4 h-4 text-red-400 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                          <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                        </svg>
                      )}
                      <span className={rule.passed ? 'text-green-700 dark:text-green-400 font-medium' : 'text-gray-600 dark:text-gray-400'}>
                        {rule.text}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>

          <div>
            <PasswordInput
              label="确认新密码"
              name="confirm-password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="•••••••••"
              required
              autoComplete="new-password"
              error={confirmPassword && !passwordsMatch ? '两次输入的密码不一致' : ''}
            />

            {/* 密码匹配状态指示 */}
            {confirmPassword && (
              <div className="mt-2 flex items-center text-sm">
                {passwordsMatch ? (
                  <div className="flex items-center text-green-600 dark:text-green-400">
                    <CheckCircle size={16} className="mr-1" />
                    密码一致
                  </div>
                ) : (
                  <div className="flex items-center text-red-600 dark:text-red-400">
                    ❌ 两次输入的密码不一致
                  </div>
                )}
              </div>
            )}
          </div>

          <Button
            type="submit"
            className="w-full"
            disabled={isLoading || !allRulesPassed || !passwordsMatch}
          >
            {isLoading ? '重置中...' : '重置密码'}
          </Button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-sm text-gray-600 dark:text-gray-400">
            链接有问题？
            {' '}
            <Link to="/forgot-password" className="text-primary-600 hover:text-primary-700 dark:text-primary-400 font-medium">
              重新发送邮件
            </Link>
          </p>
        </div>
      </Card>
    </div>
  );
}
