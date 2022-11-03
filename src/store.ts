import { writable } from 'svelte/store';

export const auth = writable<any>();
export const name = writable<string>();
export const profileImgUrl = writable<string>();
export const handle = writable<string>();

// TODO: Type this
export function setAuth(value: any) {
  auth.set(value)
  name.set(value.name)
  handle.set(value.username)
  profileImgUrl.set(value.profile_image_url)
}