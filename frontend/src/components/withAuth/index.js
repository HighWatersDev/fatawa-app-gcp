import React from 'react';
import { useAuth } from '../../hooks/withAuth';
import { useNavigate } from 'react-router-dom';

export const withAuth = (Component) => {
    return (props) => {
        let navigate = useNavigate();
        const {user, isLoading} = useAuth();

        React.useEffect(() => {
            if (!isLoading && !user) {
                navigate('/login');
            }
        }, [user, isLoading, navigate]);

        if (isLoading) {
            return <p>Loading...</p>;
        }

        return <Component {...props} user={user}/>;
    };
};

//export default withAuth
