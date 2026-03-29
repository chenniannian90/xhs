import { useState } from 'react';
import { Search as SearchIcon, Globe, Folder, ExternalLink } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import api from '../../utils/api';
import Input from '../../components/ui/Input';
import Card from '../../components/ui/Card';
import Button from '../../components/ui/Button';

export default function SearchPage() {
  const [query, setQuery] = useState('');
  const [searchType, setSearchType] = useState<'all' | 'sites' | 'categories'>('all');

  // 获取所有数据用于搜索
  const { data: categories = [], isLoading: categoriesLoading } = useQuery({
    queryKey: ['categories'],
    queryFn: async () => {
      const response = await api.get('/categories');
      return response.data.data;
    },
  });

  const { data: allSites = [], isLoading: sitesLoading } = useQuery({
    queryKey: ['sites'],
    queryFn: async () => {
      const response = await api.get('/sites');
      return response.data.data;
    },
  });

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

  // 客户端过滤
  const filteredSites = query
    ? allSites.filter((site: any) =>
        site.name.toLowerCase().includes(query.toLowerCase()) ||
        (site.description && site.description.toLowerCase().includes(query.toLowerCase())) ||
        site.url.toLowerCase().includes(query.toLowerCase())
      )
    : [];

  const filteredCategories = query
    ? categories.filter((cat: any) =>
        cat.name.toLowerCase().includes(query.toLowerCase()) ||
        (cat.description && cat.description.toLowerCase().includes(query.toLowerCase()))
      )
    : [];

  const showSites = searchType === 'all' || searchType === 'sites';
  const showCategories = searchType === 'all' || searchType === 'categories';

  const getFavicon = (url: string) => {
    try {
      const domain = new URL(url).hostname;
      return `https://www.google.com/s2/favicons?domain=${domain}&sz=64`;
    } catch {
      return '';
    }
  };

  return (
    <div className="space-y-6 max-w-4xl mx-auto">
      {/* Search Header */}
      <div className="text-center py-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
          搜索导航
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          搜索您的站点和分类
        </p>
      </div>

      {/* Search Box */}
      <Card className="p-6">
        <div className="flex gap-4 mb-4">
          <div className="flex-1">
            <Input
              type="text"
              placeholder="搜索站点名称、描述或 URL..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              autoFocus
            />
          </div>
          <Button>
            <SearchIcon size={18} className="mr-2" />
            搜索
          </Button>
        </div>

        {/* Search Type Filters */}
        <div className="flex gap-2">
          <Button
            variant={searchType === 'all' ? 'primary' : 'ghost'}
            size="sm"
            onClick={() => setSearchType('all')}
          >
            全部
          </Button>
          <Button
            variant={searchType === 'sites' ? 'primary' : 'ghost'}
            size="sm"
            onClick={() => setSearchType('sites')}
          >
            <Globe size={16} className="mr-1" />
            站点
          </Button>
          <Button
            variant={searchType === 'categories' ? 'primary' : 'ghost'}
            size="sm"
            onClick={() => setSearchType('categories')}
          >
            <Folder size={16} className="mr-1" />
            分类
          </Button>
        </div>
      </Card>

      {/* Search Results */}
      {query && (
        <div className="space-y-6">
          {/* Sites Results */}
          {showSites && (
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">
                站点 ({filteredSites.length})
              </h2>
              {filteredSites.length === 0 ? (
                <Card className="p-8 text-center text-gray-500 dark:text-gray-400">
                  没有找到匹配的站点
                </Card>
              ) : (
                <div className="space-y-2">
                  {filteredSites.map((site: any) => (
                    <Card
                      key={site.id}
                      className="p-4 hover:shadow-md transition-shadow cursor-pointer"
                    >
                      <a
                        href={site.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center gap-4"
                      >
                        <img
                          src={getFavicon(site.url)}
                          alt=""
                          className="w-10 h-10 rounded"
                          onError={(e) => {
                            e.currentTarget.src = 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><text y="18" font-size="16">🔗</text></svg>';
                          }}
                        />
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2">
                            <span className="font-medium text-gray-900 dark:text-white">
                              {site.name}
                            </span>
                            <ExternalLink size={14} className="text-gray-400" />
                          </div>
                          {site.description && (
                            <p className="text-sm text-gray-600 dark:text-gray-400">
                              {site.description}
                            </p>
                          )}
                          <p className="text-xs text-gray-400 dark:text-gray-500">
                            {site.url}
                          </p>
                        </div>
                      </a>
                    </Card>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* Categories Results */}
          {showCategories && (
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">
                分类 ({filteredCategories.length})
              </h2>
              {filteredCategories.length === 0 ? (
                <Card className="p-8 text-center text-gray-500 dark:text-gray-400">
                  没有找到匹配的分类
                </Card>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  {filteredCategories.map((cat: any) => (
                    <Card
                      key={cat.id}
                      className="p-4 hover:shadow-md transition-shadow cursor-pointer"
                      onClick={() => (window.location.href = `/dashboard/categories/${cat.id}`)}
                    >
                      <div className="flex items-center gap-3">
                        {cat.icon && (
                          <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-primary-100 to-primary-200 dark:from-primary-900 dark:to-primary-800 flex items-center justify-center text-xl">
                            {cat.icon}
                          </div>
                        )}
                        <div>
                          <h3 className="font-medium text-gray-900 dark:text-white">
                            {cat.name}
                          </h3>
                          <p className="text-sm text-gray-500 dark:text-gray-400">
                            {cat.sites?.length || 0} 个站点
                          </p>
                        </div>
                      </div>
                    </Card>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      )}

      {/* Empty State */}
      {!query && (
        <Card className="p-12 text-center">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
            <SearchIcon size={32} className="text-gray-400" />
          </div>
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
            开始搜索
          </h3>
          <p className="text-gray-600 dark:text-gray-400">
            输入关键词搜索您的站点和分类
          </p>
        </Card>
      )}
    </div>
  );
}
