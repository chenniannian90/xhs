import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import Button from '../../components/ui/Button';
import Input from '../../components/ui/Input';
import PasswordInput from '../../components/ui/PasswordInput';
import { getPasswordRulesStatus } from '../../utils/passwordRules';

export default function RegisterPage() {
  const navigate = useNavigate();
  const { register } = useAuthStore();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [errorDetails, setErrorDetails] = useState('');
  const [suggestions, setSuggestions] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);

  // 检查密码是否符合规则
  const rulesStatus = getPasswordRulesStatus(password);
  const allRulesPassed = rulesStatus.every(rule => rule.passed);

  // 检查两次密码是否一致
  const passwordsMatch = password === confirmPassword && password !== '';
  const showMatchWarning = confirmPassword !== '' && !passwordsMatch;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setErrorDetails('');
    setSuggestions([]);

    // 前端验证密码规则
    if (!allRulesPassed) {
      setError('密码不符合要求，请检查下方密码规则');
      return;
    }

    // 验证两次密码是否一致
    if (!passwordsMatch) {
      setError('两次输入的密码不一致');
      return;
    }

    setLoading(true);

    try {
      await register(username, email, password);
      navigate('/dashboard');
    } catch (err: any) {
      const errorMessage = err.message || '注册失败';

      if (errorMessage.toLowerCase().includes('password') ||
          errorMessage.toLowerCase().includes('密码')) {
        setError(errorMessage + '，请检查下方密码规则');
      } else {
        setError(errorMessage);
        // Set additional details and suggestions if available
        if (err.details) {
          setErrorDetails(err.details);
        }
        if (err.suggestions && Array.isArray(err.suggestions)) {
          setSuggestions(err.suggestions);
        }
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center px-4 py-8">
      <div className="max-w-md w-full">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">创建账号</h1>
          <p className="text-gray-600 dark:text-gray-400">加入 NavHub</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4" autoComplete="on">
          <Input
            label="用户名"
            type="text"
            name="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="yourname"
            required
            autoComplete="username"
          />

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
            <PasswordInput
              label="密码"
              name="new-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
              autoComplete="new-password"
            />

            {/* 密码规则提示 - 始终显示 */}
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
                    ) : password ? (
                      <svg className="w-4 h-4 text-red-400 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                      </svg>
                    ) : (
                      <svg className="w-4 h-4 text-gray-300 dark:text-gray-600 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                      </svg>
                    )}
                    <span className={
                      rule.passed
                        ? 'text-green-700 dark:text-green-400 font-medium'
                        : password
                          ? 'text-red-600 dark:text-red-400'
                          : 'text-gray-600 dark:text-gray-400'
                    }>
                      {rule.text}
                    </span>
                  </li>
                ))}
              </ul>
              {!password && (
                <p className="text-xs text-blue-600 dark:text-blue-400 mt-2">
                  请输入密码，满足所有要求后即可注册
                </p>
              )}
            </div>
          </div>

          <div>
            <PasswordInput
              label="确认密码"
              name="confirm-password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="••••••••"
              required
              autoComplete="new-password"
              error={showMatchWarning ? '两次输入的密码不一致' : ''}
            />
            
            {/* 密码匹配状态指示 */}
            {confirmPassword && (
              <div className="mt-2 flex items-center text-sm">
                {passwordsMatch ? (
                  <div className="flex items-center text-green-600 dark:text-green-400">
                    <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    密码一致
                  </div>
                ) : (
                  <div className="flex items-center text-red-600 dark:text-red-400">
                    <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                    </svg>
                    两次输入的密码不一致
                  </div>
                )}
              </div>
            )}
          </div>

          {error && (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 px-4 py-3 rounded-lg">
              <div className="flex items-start">
                <svg className="w-5 h-5 mr-2 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                </svg>
                <div className="flex-1">
                  <p className="font-medium mb-1">{error}</p>
                  {errorDetails && (
                    <p className="text-sm text-red-600 dark:text-red-400 mb-2">{errorDetails}</p>
                  )}
                  {suggestions.length > 0 && (
                    <ul className="mt-2 space-y-2">
                      {suggestions.map((suggestion, index) => {
                        const isLogin = suggestion.includes('直接登录') || suggestion.includes('前往登录');
                        const isReset = suggestion.includes('忘记密码') || suggestion.includes('重置');

                        return (
                          <li key={index} className="flex items-start text-sm">
                            <svg className="w-4 h-4 mr-2 flex-shrink-0 mt-0.5 text-red-500" fill="currentColor" viewBox="0 0 20 20">
                              <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
                            </svg>
                            <span className="text-red-600 dark:text-red-400">
                              {suggestion.split('：')[0]}：
                              {isLogin ? (
                                <Link to="/login" className="ml-1 text-primary-600 hover:text-primary-700 font-medium underline">
                                  {suggestion.split('：')[1]}
                                </Link>
                              ) : isReset ? (
                                <Link to="/forgot-password" className="ml-1 text-primary-600 hover:text-primary-700 font-medium underline">
                                  {suggestion.split('：')[1]}
                                </Link>
                              ) : (
                                <span className="ml-1">{suggestion.split('：')[1]}</span>
                              )}
                            </span>
                          </li>
                        );
                      })}
                    </ul>
                  )}
                </div>
              </div>
            </div>
          )}

          <Button 
            type="submit" 
            className="w-full" 
            disabled={loading || !allRulesPassed || !passwordsMatch}
          >
            {loading ? '注册中...' : '注册'}
          </Button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-gray-600 dark:text-gray-400">
            已有账号？{' '}
            <Link to="/login" className="text-primary-600 hover:text-primary-700 font-medium">
              登录
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
