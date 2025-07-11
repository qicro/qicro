// 时区转换测试
import { formatTime, formatMessageTime, formatDate } from './utils';

// 测试用例
const testCases = [
  // UTC时间字符串（数据库存储格式）
  '2024-07-11T06:30:00Z',        // UTC上午6:30
  '2024-07-11T14:30:00.000Z',    // UTC下午2:30
  '2024-07-10T16:00:00+00:00',   // UTC昨天下午4:00
  '2024-01-15T10:00:00Z',        // UTC今年1月15日上午10:00
  '2023-12-25T12:00:00Z',        // UTC去年12月25日中午12:00
  
  // 当前时间测试
  new Date().toISOString(),      // 当前时间
  new Date(Date.now() - 5 * 60 * 1000).toISOString(),  // 5分钟前
  new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(), // 2小时前
  new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString(), // 1天前
];

console.log('=== 时区转换测试 ===');
console.log('当前时区:', Intl.DateTimeFormat().resolvedOptions().timeZone);
console.log('');

testCases.forEach((dateString, index) => {
  console.log(`测试案例 ${index + 1}:`);
  console.log(`输入: ${dateString}`);
  console.log(`formatTime: ${formatTime(dateString)}`);
  console.log(`formatMessageTime: ${formatMessageTime(dateString)}`);
  console.log(`formatDate: ${formatDate(dateString)}`);
  console.log('---');
});

// 验证JavaScript内置的时区转换
console.log('=== JavaScript时区转换验证 ===');
const utcDate = new Date('2024-07-11T06:30:00Z');
console.log('UTC时间:', utcDate.toISOString());
console.log('本地时间:', utcDate.toLocaleString('zh-CN'));
console.log('时区偏移量:', utcDate.getTimezoneOffset() / 60, '小时');

export {};