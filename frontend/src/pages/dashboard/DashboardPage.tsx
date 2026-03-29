import { Link } from 'react-router-dom';
import { FolderKanban, Globe, TrendingUp, Clock } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import api from '../../utils/api';
import Card from '../../components/ui/Card';

export default function DashboardPage() {
  // 获取统计数据
  const { data: categories, isLoading: categoriesLoading } = useQuery({
    queryKey: ['categories'],
    queryFn: () => api.get('/categories').then(res => res.data.data),
  });

  const { data: sites, isLoading: sitesLoading } = useQuery({
    queryKey: ['sites'],
    queryFn: () => api.get('/sites').then(res => res.data.data),
  });

  // 计算统计���据
  const categoryCount = categories?.length || 0;
  const siteCount = sites?.length || 0;
  const totalCategories = categories || [];

  // 最近访问的站点（取前 4 个）
  const recentSites = sites?.slice(0, 4) || [];

  if (categoriesLoading || sitesLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          <p className="mt-2 text-gray-600 dark:text-gray-400">加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="p-6">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-primary-100 dark:bg-primary-900 rounded-lg">
              <FolderKanban className="text-primary-600 dark:text-primary-400" size={24} />
            </div>
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-400">分类总数</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">{categoryCount}</p>
            </div>
          </div>
        </Card>

        <Card className="p-6">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-green-100 dark:bg-green-900 rounded-lg">
              <Globe className="text-green-600 dark:text-green-400" size={24} />
            </div>
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-400">站点总数</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">{siteCount}</p>
            </div>
          </div>
        </Card>

        <Card className="p-6">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-purple-100 dark:bg-purple-900 rounded-lg">
              <TrendingUp className="text-purple-600 dark:text-purple-400" size={24} />
            </div>
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-400">本月访问</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">--</p>
            </div>
          </div>
        </Card>

        <Card className="p-6">
          <div className="flex items-center gap-4">
            <div className="p-3 bg-orange-100 dark:bg-orange-900 rounded-lg">
              <Clock className="text-orange-600 dark:text-orange-400" size={24} />
            </div>
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-400">最近更新</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white">刚刚</p>
            </div>
          </div>
        </Card>
      </div>

      {/* Categories */}
      <Card className="p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold text-gray-900 dark:text-white">我的分类</h2>
          <Link
            to="/dashboard/categories"
            className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400"
          >
            查看全部 →
          </Link>
        </div>
        {totalCategories.length === 0 ? (
          <div className="text-center py-8 text-gray-500 dark:text-gray-400">
            还没有分类，<Link to="/dashboard/categories/new" className="text-primary-600 hover:underline">创建第一个分类</Link>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {totalCategories.slice(0, 4).map((category: any) => (
              <Link
                key={category.id}
                to={`/dashboard/categories/${category.id}`}
                className="group p-4 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-primary-500 dark:hover:border-primary-500 transition-all hover:shadow-md"
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h3 className="font-semibold text-gray-900 dark:text-white group-hover:text-primary-600 dark:group-hover:text-primary-400">
                      {category.name}
                    </h3>
                    {category.description && (
                      <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                        {category.description}
                      </p>
                    )}
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </Card>

      {/* Recent Sites */}
      <Card className="p-6">
        <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">最近添加</h2>
        {recentSites.length === 0 ? (
          <div className="text-center py-8 text-gray-500 dark:text-gray-400">
            还没有站点，<Link to="/dashboard/categories/new" className="text-primary-600 hover:underline">添加第一个站点</Link>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {recentSites.map((site: any) => (
              <a
                key={site.id}
                href={site.url}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-3 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
              >
                <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-gray-100 to-gray-200 dark:from-gray-700 dark:to-gray-600 flex items-center justify-center">
                  <span className="text-lg font-bold text-gray-600 dark:text-gray-300">
                    {site.name[0]}
                  </span>
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-medium text-gray-900 dark:text-white truncate">{site.name}</p>
                  <p className="text-sm text-gray-500 dark:text-gray-400 truncate">{site.url}</p>
                </div>
              </a>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
}
