import { page } from '$app/stores';
import { derived } from 'svelte/store';
//import { auth } from './auth';

type User = {
    uid: string;
    email?: string;
};

export const userAuth = derived<typeof page, User | null>(
    page,
        ($page, set) => {
            const { user } = $page.data;
            if (!user) {
                set(null);
                return;
            }

            set(user);
        },
        null
);