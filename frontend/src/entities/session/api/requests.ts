import { api } from '@/shared/api';

import { type Session, SessionSchema } from '../model/schema';
import type { LoginDTO, RegisterDTO } from './types';

export async function fetchMe(): Promise<Session> {
  const response = await api.get('/auth/me');
  return SessionSchema.parse(response.data);
}

export async function logout() {
  await api.post('/auth/logout');
}

export async function login(dto: LoginDTO) {
  await api.post('/auth/login', dto);
}

export async function register(dto: RegisterDTO) {
  await api.post('/auth/register', dto);
}
