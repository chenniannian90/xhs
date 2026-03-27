import { Link } from 'react-router-dom';
import { Home, ArrowLeft } from 'lucide-react';

export default function NotFoundPage() {
  return (
    <div className="min-h-screen flex items-center justify-center px-4 bg-gray-50 dark:bg-gray-900">
      <div className="text-center">
        <div className="inline-block mb-4">
          <h1 className="text-9xl font-bold text-primary-600 dark:text-primary-400">404</h1>
        </div>
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-2">
          页面未找到
        </h2>
        <p className="text-gray-600 dark:text-gray-400 mb-8 max-w-md">
          抱歉，您访问的页面不存在。可能是链接错误或页面已被移除。
        </p>
        <div className="flex gap-4 justify-center">
          <Link
            to="/dashboard"
            className="inline-flex items-center gap-2 px-6 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
          >
            <Home size={18} />
            返回首页
          </Link>
          <button
            onClick={() => window.history.back()}
            className="inline-flex items-center gap-2 px-6 py-3 bg-gray-200 dark:bg-gray-800 text-gray-900 dark:text-white rounded-lg hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors"
          >
            <ArrowLeft size={18} />
            返回上一页
          </button>
        </div>
      </div>
    </div>
  );
}
