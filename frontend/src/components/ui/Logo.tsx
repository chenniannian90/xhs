export default function Logo({ className = '' }: { className?: string }) {
  return (
    <div className={`flex items-center gap-2 ${className}`}>
      <div className="w-8 h-8 bg-gradient-to-br from-primary-500 to-primary-700 rounded-lg flex items-center justify-center">
        <span className="text-white text-xl font-bold">N</span>
      </div>
      <span className="text-xl font-bold text-gray-900 dark:text-white">NavHub</span>
    </div>
  );
}
