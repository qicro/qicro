import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// 时间格式化工具函数
export function formatTime(dateString: string | Date): string {
  // 确保输入的UTC时间正确转换为本地时间
  const date = typeof dateString === 'string' ? new Date(dateString) : dateString;
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '时间无效';
  }
  
  const now = new Date();
  // 使用本地时间戳进行比较
  const diff = now.getTime() - date.getTime();
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  
  if (days > 0) {
    return `${days}天前`;
  } else if (hours > 0) {
    return `${hours}小时前`;
  } else if (minutes > 0) {
    return `${minutes}分钟前`;
  } else if (seconds > 0) {
    return `${seconds}秒前`;
  } else {
    return '刚刚';
  }
}

export function formatMessageTime(dateString: string | Date): string {
  const date = typeof dateString === 'string' ? new Date(dateString) : dateString;
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '时间无效';
  }
  
  const now = new Date();
  
  // 使用日期字符串进行比较（避免时间影响）
  const dateStr = date.toDateString();
  const nowStr = now.toDateString();
  
  // 如果是今天，只显示时间
  if (dateStr === nowStr) {
    return date.toLocaleTimeString('zh-CN', { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  }
  
  // 如果是昨天
  const yesterday = new Date(now);
  yesterday.setDate(yesterday.getDate() - 1);
  if (dateStr === yesterday.toDateString()) {
    return `昨天 ${date.toLocaleTimeString('zh-CN', { 
      hour: '2-digit', 
      minute: '2-digit' 
    })}`;
  }
  
  // 如果是今年
  if (date.getFullYear() === now.getFullYear()) {
    return date.toLocaleDateString('zh-CN', { 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
  
  // 完整日期
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
}

// 简单的日期格式化，用于列表显示
export function formatDate(dateString: string | Date): string {
  const date = typeof dateString === 'string' ? new Date(dateString) : dateString;
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '日期无效';
  }
  
  // 使用本地时间格式化
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
}