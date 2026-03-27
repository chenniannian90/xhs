import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import api from '../../utils/api';
import Card from '../../components/ui/Card';
import Button from '../../components/ui/Button';
import { ChevronLeft } from 'lucide-react';

export default function CategoryDetailPage() {
  const { id: categoryId } = useParams<{ id: string }>();

  const { data: category, isLoading } = useQuery({
    queryKey: ['category', categoryId],
    queryFn: () => api.get(`/categories/${categoryId}`).then(res => res.data),
    enabled: !!categoryId,
  });

  if (isLoading) {
    return <div className="flex items-center justify-center h-64">加载中...</div>;
  }

  if (!category) {
    return <div className="flex items-center justify-center h-64">分类未找到</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="sm" onClick={() => window.history.back()}>
          <ChevronLeft size={20} />
        </Button>
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{category.name}</h1>
          <p className="text-gray-600 dark:text-gray-400">{category.description}</p>
        </div>
      </div>

      <Card>
        <h2 className="text-lg font-semibold mb-4">站点列表</h2>
        <p className="text-gray-600 dark:text-gray-400">暂无站点</p>
      </Card>
    </div>
  );
}
