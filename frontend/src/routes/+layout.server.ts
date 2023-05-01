import { auth } from "$lib/auth.js";
import { redirect } from "@sveltejs/kit";
import type {ServerLoadEvent } from '@sveltejs/kit/types';

export async function load({ cookies }: ServerLoadEvent) {
    try {
        const user = auth.currentUser?.uid
        const token = cookies.get("token");
//        const user = token ? await auth.verifyIdToken(token) : null;
        return {
            uid: user
        };
    } catch {
        // The token is set but invalid or expired
        cookies.set("token", "", { maxAge: -1 });
        throw redirect(307, "/");
    }
}