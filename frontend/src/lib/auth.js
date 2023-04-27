import { initializeApp } from "firebase/app";
import { getAuth, setPersistence, browserSessionPersistence } from 'firebase/auth';

const firebaseConfig = {
    apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
    authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
    projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
    storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
    messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
    appId: import.meta.env.VITE_FIREBASE_APP_ID
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
export const auth = getAuth(app);
setPersistence(auth, browserSessionPersistence)

export const signInWithEmailAndPassword = async (email, password) => {
    console.log(auth);
    try {
        const { user } = await auth.signInWithEmailAndPassword(email, password);
        return user;
    } catch (error) {
        console.error(error);
        throw new Error(error.message);
    }
};

export const onAuthStateChanged = handler => {
    return auth.onAuthStateChanged(handler);
};
