import { z } from 'zod';

export const SessionSchema = z.object({
  id: z.uuidv7(),
  username: z.string(),
  email: z.email(),
  role: z.string(),
  avatar: z.url(),
});
export type Session = z.infer<typeof SessionSchema>;
