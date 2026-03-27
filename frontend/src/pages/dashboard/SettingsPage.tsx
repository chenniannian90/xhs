import { useThemeStore } from '../../store/themeStore';
import { useAuthStore } from '../../store/authStore';
import Button from '../../components/ui/Button';
import Card from '../../components/ui/Card';

export default function SettingsPage() {
  const { theme, toggleTheme } = useThemeStore();
  const { user } = useAuthStore();

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">设置</h1>
        <p className="text-gray-600 dark:text-gray-400">管理你的账户设置</p>
      </div>

      <Card>
        <h2 className="text-lg font-semibold mb-4">主题设置</h2>
        <div className="flex items-center justify-between">
          <div>
            <p className="font-medium">当前主题</p>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              {theme === 'light' ? '浅色模式' : '深色模式'}
            </p>
          </div>
          <Button onClick={toggleTheme}>
            切换主题
          </Button>
        </div>
      </Card>

      <Card>
        <h2 className="text-lg font-semibold mb-4">账户信息</h2>
        <div className="space-y-4">
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">用户名</p>
            <p className="font-medium">{user?.username || '未设置'}</p>
          </div>
          <div>
            <p className="text-sm text-gray-600 dark:text-gray-400">邮箱</p>
            <p className="font-medium">{user?.email || '未设置'}</p>
          </div>
        </div>
      </Card>
    </div>
  );
}
