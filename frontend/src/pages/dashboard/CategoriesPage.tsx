import { useState } from 'react';
import { Link } from 'react-router-dom';
import { Plus, MoreVertical, Edit, Trash2, Share, Share2, Eye } from 'lucide-react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../../utils/api';
import Button from '../../components/ui/Button';
import Card from '../../components/ui/Card';
import ConfirmDialog from '../../components/ui/ConfirmDialog';

export default function CategoriesPage() {
  const queryClient = useQueryClient();
  const [showMenu, setShowMenu] = useState<string | null>(null);
  const [confirmDialog, setConfirmDialog] = useState<{
    isOpen: boolean;
    title: string;
    message: string;
    onConfirm: () => void;
  }>({
    isOpen: false,
    title: '',
    message: '',
    onConfirm: () => {},
  });

  // 获取分类列表
  const { data: categories = [], isLoading, error } = useQuery({
    queryKey: ['categories'],
    queryFn: async () => {
      const response = await api.get('/categories');
      return response.data.data;
    },
  });

  // 删除分类的 mutation
  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`/categories/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      setShowMenu(null);
    },
  });

  // 切换公开状态的 mutation
  const togglePublicMutation = useMutation({
    mutationFn: async ({ id, isPublic }: { id: string; isPublic: boolean }) => {
      if (isPublic) {
        await api.delete(`/categories/${id}/share`);
      } else {
        await api.post(`/categories/${id}/share`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      setShowMenu(null);
    },
  });

  const handleDelete = (id: string) => {
    setConfirmDialog({
      isOpen: true,
      title: '确认删除',
      message: '确定要删除这个分类吗？分类下的所有站点也会被删除。此操作无法撤销。',
      onConfirm: () => deleteMutation.mutate(id),
    });
  };

  const handleTogglePublic = (id: string, isPublic: boolean) => {
    togglePublicMutation.mutate({ id, isPublic });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          <p className="mt-2 text-gray-600 dark:text-gray-400">加载中...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <Card className="p-8 bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800">
          <p className="text-red-600 dark:text-red-400">加载失败，请刷新页面重试</p>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">分类管理</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            管理您的导航分类，共 {categories.length} 个分类
          </p>
        </div>
        <Link to="/dashboard/categories/new">
          <Button>
            <Plus size={18} className="mr-1" />
            新建分类
          </Button>
        </Link>
      </div>

      {/* Categories Grid */}
      {categories.length === 0 ? (
        <Card className="p-12 text-center">
          <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
            <Plus size={32} className="text-gray-400" />
          </div>
          <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">
            还没有分类
          </h3>
          <p className="text-gray-600 dark:text-gray-400 mb-4">
            创建您的第一个分类来开始整理网站
          </p>
          <Link to="/dashboard/categories/new">
            <Button>
              <Plus size={18} className="mr-1" />
              创建分类
            </Button>
          </Link>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {categories.map((category: any) => (
            <Card
              key={category.id}
              className="p-5 hover:shadow-lg transition-shadow group"
            >
              {/* Card Header */}
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-3">
                  {category.icon && (
                    <div className="w-12 h-12 rounded-lg bg-gradient-to-br from-primary-100 to-primary-200 dark:from-primary-900 dark:to-primary-800 flex items-center justify-center text-2xl">
                      {category.icon}
                    </div>
                  )}
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-900 dark:text-white truncate">
                      {category.name}
                    </h3>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {category.sites?.length || 0} 个站点
                    </p>
                  </div>
                </div>
                <div className="relative">
                  <button
                    onClick={() => setShowMenu(showMenu === category.id ? null : category.id)}
                    className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700"
                  >
                    <MoreVertical size={18} className="text-gray-500" />
                  </button>

                  {/* Dropdown Menu */}
                  {showMenu === category.id && (
                    <div className="absolute right-0 top-8 z-10 w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 py-1">
                      <Link
                        to={`/dashboard/categories/${category.id}`}
                        className="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                        onClick={() => setShowMenu(null)}
                      >
                        <Eye size={16} />
                        查看站点
                      </Link>
                      <button
                        className="w-full flex items-center gap-2 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                        onClick={() => {
                          setShowMenu(null);
                          // 编辑功能将在后续版本实现
                        }}
                      >
                        <Edit size={16} />
                        编辑
                      </button>
                      <button
                        className="w-full flex items-center gap-2 px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                        onClick={() => {
                          handleTogglePublic(category.id, !!category.is_public);
                        }}
                      >
                        {category.is_public ? (
                          <>
                            <Share2 size={16} />
                            取消分享
                          </>
                        ) : (
                          <>
                            <Share size={16} />
                            公开分享
                          </>
                        )}
                      </button>
                      <button
                        className="w-full flex items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20"
                        onClick={() => handleDelete(category.id)}
                      >
                        <Trash2 size={16} />
                        删除
                      </button>
                    </div>
                  )}
                </div>
              </div>

              {/* Description */}
              {category.description && (
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-3 line-clamp-2">
                  {category.description}
                </p>
              )}

              {/* Public Badge */}
              {category.is_public && (
                <div className="flex items-center gap-2 text-xs text-primary-600 dark:text-primary-400">
                  <Share size={14} />
                  <span>已公开分享</span>
                  {category.share_token && (
                    <a
                      href={`${window.location.origin}/shared/${category.share_token}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="hover:underline"
                    >
                      查看链接
                    </a>
                  )}
                </div>
              )}
            </Card>
          ))}
        </div>
      )}

      {/* Confirm Dialog */}
      <ConfirmDialog
        isOpen={confirmDialog.isOpen}
        title={confirmDialog.title}
        message={confirmDialog.message}
        confirmText="删除"
        cancelText="取消"
        type="danger"
        onConfirm={confirmDialog.onConfirm}
        onCancel={() => setConfirmDialog((prev) => ({ ...prev, isOpen: false }))}
      />
    </div>
  );
}
