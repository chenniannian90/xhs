export interface User {
  id: string;
  email: string;
  username: string;
  avatar_url?: string;
  email_verified: boolean;
  theme_preference: 'light' | 'dark';
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: string;
  user_id: string;
  name: string;
  description?: string;
  icon?: string;
  sort_order: number;
  is_public: boolean;
  share_token?: string;
  created_at: string;
  updated_at: string;
  sites?: Site[];
}

export interface Site {
  id: string;
  user_id: string;
  category_id: string;
  name: string;
  url: string;
  description?: string;
  icon?: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface ApiError {
  code: number;
  message: string;
}
