import type { LoadEvent } from '@sveltejs/kit/types';

export async function load({ parent }: LoadEvent) {
    const { uid } = await parent();
    return { uid };
}