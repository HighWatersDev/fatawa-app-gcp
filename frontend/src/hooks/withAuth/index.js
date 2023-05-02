import {useEffect, useState} from 'react';
import {auth} from '../../firebase';

export const useAuth = () => {
    const [user, setUser] = useState(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        return auth.onAuthStateChanged((authUser) => {
            if (authUser) {
                setUser(authUser);
            } else {
                setUser(null);
            }
            setIsLoading(false);
        });
    }, []);

    return { user, isLoading };
};