// 密码规则配置
export const PASSWORD_RULES = [
  { text: '至少8个字符', test: (pwd: string) => pwd.length >= 8 },
  { text: '至少1个大写字母', test: (pwd: string) => /[A-Z]/.test(pwd) },
  { text: '至少1个小写字母', test: (pwd: string) => /[a-z]/.test(pwd) },
  { text: '至少1个数字', test: (pwd: string) => /\d/.test(pwd) },
  { text: '至少1个特殊字符', test: (pwd: string) => /[!@#$%^&*(),.?":{}|<>]/.test(pwd) },
];

// 密码规则提示文本
export const PASSWORD_RULES_TEXT = `密码规则：
• 至少8个字符
• 至少1个大写字母
• 至少1个小写字母
• 至少1个数字
• 至少1个特殊字符`;

// 验证密码是否符合所有规则
export const validatePassword = (password: string): { valid: boolean; failedRules: string[] } => {
  const failedRules: string[] = [];
  
  PASSWORD_RULES.forEach(rule => {
    if (!rule.test(password)) {
      failedRules.push(rule.text);
    }
  });

  return {
    valid: failedRules.length === 0,
    failedRules,
  };
};

// 获取密码规则状态（用于显示）
export const getPasswordRulesStatus = (password: string) => {
  return PASSWORD_RULES.map(rule => ({
    ...rule,
    passed: rule.test(password),
  }));
};
