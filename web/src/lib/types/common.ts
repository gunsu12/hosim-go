// Common & Base Types for HOSIM EHR

export interface BaseEntity {
  id: string | number;
  createdAt?: string;
  updatedAt?: string;
  createdBy?: string;
  updatedBy?: string;
}

export interface ApiResponse<T = any> {
  code: number;
  status: string;
  message?: string;
  data: T;
}

export interface PaginatedResult<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}
