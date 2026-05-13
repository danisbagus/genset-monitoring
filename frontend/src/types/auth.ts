export interface LoginResponse {
  token: string;
  user: User;
}

export interface User {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'operator' | 'viewer';
}

export interface TelemetryPreview {
  label: string;
  value: string | number;
  unit: string;
  trend?: 'up' | 'down' | 'stable';
  status: 'normal' | 'warning' | 'critical';
}
